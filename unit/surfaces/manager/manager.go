// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package manager

import (
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
	element "github.com/djthorpe/gopi/v2/unit/surfaces/element"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type manager struct {
	display rpi.DXDisplayHandle
	bitmap  map[bitmap.Bitmap]bool
	element map[element.Element]bool

	sync.Mutex
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Config) Name() string { return "gopi/surfaces/manager" }

func (config Config) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(manager)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION Manager

func (this *manager) Init(config Config) error {
	// Check display
	if config.Display == rpi.DX_NO_HANDLE {
		return gopi.ErrBadParameter.WithPrefix("display")
	} else {
		this.display = config.Display
	}

	// Set bitmaps and elements
	this.bitmap = make(map[bitmap.Bitmap]bool)
	this.element = make(map[element.Element]bool)

	return nil
}

func (this *manager) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Remove all elements and bitmaps
	if err := this.removeAllElements(); err != nil {
		return err
	}
	if err := this.removeAllBitmaps(); err != nil {
		return err
	}

	// Release resources
	this.element = nil
	this.bitmap = nil
	this.display = rpi.DX_NO_HANDLE

	// Return sucess
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// BITMAPS

func (this *manager) NewBitmap(size gopi.Size, mode gopi.SurfaceFlags) (bitmap.Bitmap, error) {
	if bm, err := gopi.New(bitmap.Config{size, mode}, this.Log.Clone(bitmap.Config{}.Name())); err != nil {
		return nil, err
	} else {
		bm_ := bm.(bitmap.Bitmap)
		bm_.Retain()
		this.bitmap[bm_] = true
		return bm_, nil
	}
}

func (this *manager) ReleaseBitmap(bm bitmap.Bitmap) error {
	if bm == nil {
		return gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if _,exists := this.bitmap[bm]; exists == false {
		return gopi.ErrNotFound.WithPrefix("bitmap")
	} else {
		delete(this.bitmap,bm)
		if bm.Release() {
			return bm.Close()
		} else {
			return nil
		}
	}
}

func (this *manager) AddElementWithSize(origin gopi.Point,size gopi.Size,layer uint16,opacity float32,flags) (element.Element, error) {
	// TODO
	
	if em, err := gopi.New(element.Config{
		Origin:  origin,
		Flags:   flags,
		Opacity: opacity,
		Update:  update,
		Display: this.display,
	}, nil); err != nil {
		return nil, err
	} else {
		return em.(element.Element), nil
	}
}


////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *manager) String() string {
	if this.display == rpi.DX_NO_HANDLE {
		return "<" + this.Log.Name() +
			" display=nil" +
			">"
	} else {
		fmt.Println("display=", this.display)
		return "<" + this.Log.Name() +
			" display=" + fmt.Sprint(this.display) +
			">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *manager) removeAllElements() error {
	// Ignore if no elements
	if len(this.element) == 0 {
		return nil
	}
	// Start an update
	update, err := rpi.DXUpdateStart(0)
	if err != nil {
		return err
	}
	// Remove elements
	errs := gopi.NewCompoundError()
	for em := range this.element {
		errs.Add(em.RemoveElement(update))
	}
	// Update GPU
	errs.Add(rpi.DXUpdateSubmitSync(update))
	// Delete elements from map
	for em := range this.element {
		delete(this.element, em)
	}
	// Return any errors
	return errs.ErrorOrSelf()
}

func (this *manager) removeAllBitmaps() error {
	// Ignore if no bitmaps
	if len(this.bitmap) == 0 {
		return nil
	}
	// Release bitmaps
	errs := gopi.NewCompoundError()
	for bm := range this.bitmap {
		errs.Add(bm.Close())
	}
	// Delete bitmaps from map
	for bm := range this.bitmap {
		delete(this.bitmap, bm)
	}
	// Return any errors
	return errs.ErrorOrSelf()
}
