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

type Element interface {
	Bitmap() Bitmap
	Close(update rpi.DXUpdate) (bool,error)
}

type element struct {
	handle rpi.DXElement
	bitmap *bitmap

	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// NEW AND CLOSE

func NewElement(update rpi.DXUpdate, display rpi.DXDisplayHandle, rect rpi.DXRect, resource Bitmap, layer uint16, opacity uint8) (Element, error) {
	// Make new element
	this := new(element)

	// Check incoming parameters
	if update == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("update")
	}
	if display == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("display")
	}
	if resource == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("bitmap")
	}
	if bitmap_, ok := resource.(*bitmap); ok == false || bitmap_.handle == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("bitmap")
	} else {
		this.bitmap = bitmap_
	}

	// Determine size and destination
	src_size := resource.Size()
	dest_origin := rpi.DXRectOrigin(rect)
	dest_size := rpi.DXRectSize(rect)
	if dest_size.W == 0 {
		dest_size.W = src_size.W
	}
	if dest_size.H == 0 {
		dest_size.H = src_size.H
	}
	dest_rect := rpi.DXNewRect(dest_origin.X, dest_origin.Y, dest_size.W, dest_size.H)
	alpha := rpi.DXAlpha{
		Opacity: uint32(opacity),
	}
	clamp := rpi.DXClamp{}
	transform := rpi.DX_TRANSFORM_NONE
	protection := rpi.DX_PROTECTION_NONE

	if update == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("update")
	} else if handle, err := rpi.DXElementAdd(update, display, layer, dest_rect,this.bitmap.Retain(), src_size, protection, alpha, clamp, transform); err != nil {
		this.bitmap.Release()
		return nil, err
	} else {
		this.handle = handle
	}

	// Return success
	return this, nil
}

func (this *element) Close(update rpi.DXUpdate) (bool,error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if update == 0 {
		return false,gopi.ErrInternalAppError.WithPrefix("Close")
	} else if this.handle == rpi.DX_NO_HANDLE {
		return false,gopi.ErrInternalAppError.WithPrefix("Close")
	} else {
		err := rpi.DXElementRemove(update,this.handle)
		release := this.bitmap.Release()

		// Release resources
		this.handle = rpi.DX_NO_HANDLE
		this.bitmap = nil

		// Return success
		return release, err
	}
}

func (this *element) Bitmap() Bitmap {
	return this.bitmap
}



