// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package dispmanx

import (
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Update interface {
	// Create a new bitmap
	NewBitmap(ImageType, uint32, uint32) (Bitmap, error)

	// Release a bitmap
	ReleaseBitmap(Bitmap) error

	// Create a new element
	AddElement(rect rpi.DXRect, resource Bitmap, layer uint16, opacity uint8) (Element, error)

	// Remove an element
	RemoveElement(Element) error

	// Perform AddElement, RemoveElement and bitmap operations within Do
	Do(int32, func() error) error

	// Close
	Close() error
}

type update struct {
	display  rpi.DXDisplayHandle
	bitmaps  map[Bitmap]*bitmap
	elements map[Element]*element
	update   rpi.DXUpdate

	sync.RWMutex
}

////////////////////////////////////////////////////////////////////////////////
// NEW AND CLOSE

func NewUpdate(display rpi.DXDisplayHandle) (Update, error) {
	this := new(update)
	this.display = display
	this.bitmaps = make(map[Bitmap]*bitmap)
	this.elements = make(map[Element]*element)
	return this, nil
}

func (this *update) Close() error {
	this.Lock()
	defer this.Unlock()

	// Remove elements
	if update, err := rpi.DXUpdateStart(0); err != nil {
		return err
	} else {
		for _,element := range this.elements {
			if _, err := element.Close(update); err != nil {
				return err
			}
		}
		if err := rpi.DXUpdateSubmitSync(update); err != nil {
			return err
		}
	}

	// Remove bitmaps
	for _,bitmap := range this.bitmaps {
		if err := bitmap.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.elements = nil
	this.bitmaps = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// DO

func (this *update) Do(priority int32, cb func() error) error {
	this.Lock()
	defer this.Unlock()

	if update, err := rpi.DXUpdateStart(priority); err != nil {
		return err
	} else {
		defer func() {
			rpi.DXUpdateSubmitSync(update)
			this.update = 0
		}()
		this.update = update
		return cb()
	}
}

func (this *update) Update() rpi.DXUpdate {
	return this.update
}

////////////////////////////////////////////////////////////////////////////////
// ADD AND REMOVE BITMAPS

func (this *update) NewBitmap(imageType ImageType, w, h uint32) (Bitmap, error) {
	if bm, err := NewBitmap(imageType, w, h); err != nil {
		return nil, err
	} else if bm_,ok := bm.(*bitmap); ok == false {		
		return nil,gopi.ErrInternalAppError.WithPrefix("NewBitmap")
	} else {
		bm_.Retain()
		this.bitmaps[bm] = bm_
		return bm, nil
	}
}

func (this *update) ReleaseBitmap(bm Bitmap) error {
	if bm_, exists := this.bitmaps[bm]; exists == false {
		return gopi.ErrNotFound.WithPrefix("ReleaseBitmap")
	} else if bm_.Release() == false {
		return nil
	} else {
		delete(this.bitmaps, bm)
		return bm_.Close()
	}
}

func (this *update) AddElement(rect rpi.DXRect, resource Bitmap, layer uint16, opacity uint8) (Element, error) {
	update := this.Update()
	if update == 0 {
		return nil, gopi.ErrInternalAppError.WithPrefix("AddElement")
	}
	if elem, err := NewElement(update, this.display, rect, resource, layer, opacity); err != nil {
		return nil, err
	} else {
		this.elements[elem] = elem.(*element)
		return elem, nil
	}
}

func (this *update) RemoveElement(element Element) error {
	update := this.Update()
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("RemoveElement")
	}
	if _,exists := this.elements[element]; exists == false {
		return gopi.ErrNotFound.WithPrefix("RemoveElement")
	}
	// Remove the element and the bitmap
	delete(this.elements,element)
	if release, err := element.Close(update); err != nil {
		return err
	} else if release {
		if bm,exists := this.bitmaps[element.Bitmap()]; exists == false {
			return gopi.ErrNotFound.WithPrefix("RemoveElement")
		} else {
			delete(this.bitmaps,element.Bitmap())
			return bm.Close()
		}
	}
	// Return success
	return nil
}


