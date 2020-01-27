// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package element

import (
	"errors"
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type element struct {
	handle rpi.DXElement
	origin rpi.DXPoint
	size   rpi.DXSize
	bitmap bitmap.Bitmap

	sync.Mutex
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Config) Name() string { return "gopi/element" }

func (config Config) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(element)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION Element

func (this *element) Init(config Config) error {
	// Set size, and for any zeros try and set from bitmap
	size := rpi.DXSize{uint32(config.Size.W), uint32(config.Size.H)}
	if size.W == 0 && config.Bitmap != nil {
		size.W = uint32(config.Bitmap.Size().W)
	}
	if size.H == 0 && config.Bitmap != nil {
		size.H = uint32(config.Bitmap.Size().H)
	}
	// Check size, neither dimension can be zero
	if size.W == 0 || size.H == 0 {
		return gopi.ErrBadParameter.WithPrefix("size")
	} else {
		this.origin = rpi.DXPoint{int32(config.Origin.X), int32(config.Origin.Y)}
		this.size = size
	}
	// Check update
	if config.Update == rpi.DX_NO_HANDLE {
		return gopi.ErrBadParameter.WithPrefix("update")
	}
	// Check display
	if config.Display == rpi.DX_NO_HANDLE {
		return gopi.ErrBadParameter.WithPrefix("display")
	}
	// Check opacity
	if config.Opacity < 0.0 || config.Opacity > 1.0 {
		return gopi.ErrBadParameter.WithPrefix("opacity")
	}
	// If no bitmap, then create one
	if config.Bitmap == nil {
		if bm, err := gopi.New(bitmap.Config{Size: this.Size(), Mode: config.Flags.Config()}, this.Log.Clone(bitmap.Config{}.Name())); err != nil {
			return err
		} else {
			this.bitmap = bm.(bitmap.Bitmap)
		}
	} else {
		this.bitmap = config.Bitmap
	}

	// Retain bitmap
	this.bitmap.Retain()

	// Set element properties
	src_size := this.bitmap.DXSize()
	dest_rect := rpi.DXNewRect(this.origin.X, this.origin.Y, this.size.W, this.size.H)
	alpha := rpi.DXAlpha{
		Opacity: uint32(opacity_from_float(config.Opacity)),
	}
	clamp := rpi.DXClamp{}
	transform := rpi.DX_TRANSFORM_NONE
	protection := rpi.DX_PROTECTION_NONE

	// Handle alpha flags
	if config.Flags.Mod()&gopi.SURFACE_FLAG_ALPHA_FROM_SOURCE != 0 {
		alpha.Flags |= rpi.DX_ALPHA_FLAG_FROM_SOURCE
	} else {
		alpha.Flags |= rpi.DX_ALPHA_FLAG_FIXED_ALL_PIXELS
	}

	// Add element
	if handle, err := rpi.DXElementAdd(config.Update, config.Display, config.Layer, dest_rect, this.bitmap.DXHandle(), src_size, protection, alpha, clamp, transform); err != nil {
		this.releaseBitmap()
		return err
	} else {
		this.handle = handle
	}

	// Success
	return nil
}

func (this *element) RemoveElement(update rpi.DXUpdate) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check parameters
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}
	// Remove element
	if this.handle != rpi.DX_NO_HANDLE {
		if err := rpi.DXElementRemove(update, this.handle); err != nil {
			return err
		}
	}

	// Release bitmap
	this.releaseBitmap()

	// Release resources
	this.handle = 0

	// Return success
	return nil
}

func (this *element) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Can't close until element removed
	if this.handle != rpi.DX_NO_HANDLE {
		return errors.New("Call RemoveElement before Close")
	}

	// Return sucess
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// CHANGE ELEMENT PROPERTIES

/*
DX_CHANGE_FLAG_LAYER     DXChangeFlags = (1 << 0)
	DX_CHANGE_FLAG_OPACITY   DXChangeFlags = (1 << 1)
	 DXChangeFlags = (1 << 2)
	DX_CHANGE_FLAG_SRC_RECT  DXChangeFlags = (1 << 3)
	DX_CHANGE_FLAG_MASK      DXChangeFlags = (1 << 4)
	DX_CHANGE_FLAG_TRANSFORM DXChangeFlags = (1 << 5)
*/

/*

func DXElementChangeAttributes(update DXUpdate, element DXElement, flags DXChangeFlags, layer uint16, opacity uint8, dest_rect, src_rect DXRect, transform DXTransform) error {
	if C.vc_dispmanx_element_change_attributes(
		C.DISPMANX_UPDATE_HANDLE_T(update),
		C.DISPMANX_ELEMENT_HANDLE_T(element),
		C.uint32_t(flags),
		C.int32_t(layer),
		C.uint8_t(opacity),
		dest_rect, src_rect, 0, C.DISPMANX_TRANSFORM_T(transform)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}
*/

func (this *element) SetSize(update rpi.DXUpdate, size gopi.Size) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Determine new size
	dxsize := rpi.DXSize{uint32(size.W), uint32(size.H)}

	// Check parameters
	if dxsize.W == 0 || dxsize.H == 0 {
		return gopi.ErrBadParameter.WithPrefix("size")
	}
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}
	if this.handle == rpi.DX_NO_HANDLE {
		return nil
	}

	// Do change
	dest_rect := rpi.DXNewRect(this.origin.X, this.origin.Y, dxsize.W, dxsize.H)
	if err := rpi.DXElementChangeAttributes(update, this.handle, rpi.DX_CHANGE_FLAG_DEST_RECT, 0, 0, dest_rect, nil, 0); err != nil {
		return err
	} else {
		this.size = dxsize
	}

	// Success
	return nil
}

func (this *element) SetOrigin(update rpi.DXUpdate, origin gopi.Point) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Determine new origin
	dxorigin := rpi.DXPoint{int32(origin.X), int32(origin.Y)}

	// Check parameters
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}
	if this.handle == rpi.DX_NO_HANDLE {
		return nil
	}

	// Do change
	dest_rect := rpi.DXNewRect(dxorigin.X, dxorigin.Y, this.size.W, this.size.H)
	if err := rpi.DXElementChangeAttributes(update, this.handle, rpi.DX_CHANGE_FLAG_DEST_RECT, 0, 0, dest_rect, nil, 0); err != nil {
		return err
	} else {
		this.origin = dxorigin
	}

	// Success
	return nil
}

func (this *element) SetLayer(update rpi.DXUpdate, layer uint16) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check parameters
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}
	if this.handle == rpi.DX_NO_HANDLE {
		return nil
	}

	// Do change
	if err := rpi.DXElementChangeAttributes(update, this.handle, rpi.DX_CHANGE_FLAG_LAYER, layer, 0, nil, nil, 0); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *element) SetOpacity(update rpi.DXUpdate, opacity float32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check parameters
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}
	if this.handle == rpi.DX_NO_HANDLE {
		return nil
	}
	if opacity < 0.0 || opacity > 1.0 {
		return gopi.ErrBadParameter.WithPrefix("opacity")
	}

	// Do change
	if err := rpi.DXElementChangeAttributes(update, this.handle, rpi.DX_CHANGE_FLAG_OPACITY, 0, opacity_from_float(opacity), nil, nil, 0); err != nil {
		return err
	}

	// Success
	return nil
}

func (this *element) SetBitmap(update rpi.DXUpdate, bm bitmap.Bitmap) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check parameters
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}
	if this.handle == rpi.DX_NO_HANDLE {
		return nil
	}
	if bm == nil {
		return gopi.ErrBadParameter.WithPrefix("bitmap")
	}

	// Retain the bitmap
	bm.Retain()

	// Change the src_rect and the source
	if err := rpi.DXElementChangeAttributes(update, this.handle, rpi.DX_CHANGE_FLAG_SRC_RECT, 0, 0, nil, bm.DXRect(), 0); err != nil {
		bm.Release()
		return err
	} else if err := rpi.DXElementChangeSource(update, this.handle, bm.DXHandle()); err != nil {
		bm.Release()
		return err
	} else if err := this.releaseBitmap(); err != nil {
		bm.Release()
		return err
	} else {
		this.bitmap = bm
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *element) releaseBitmap() error {
	if this.bitmap == nil {
		return nil
	}
	if this.bitmap.Release() {
		if err := this.bitmap.Close(); err != nil {
			return err
		}
	}
	this.bitmap = nil
	return nil
}

func opacity_from_float(opacity float32) uint8 {
	if opacity < 0.0 {
		opacity = 0.0
	} else if opacity > 1.0 {
		opacity = 1.0
	}
	// Opacity is between 0 (fully transparent) and 255 (fully opaque)
	return uint8(opacity * float32(0xFF))
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *element) String() string {
	if this.handle == 0 {
		return "<" + Config{}.Name() + " handle=nil" + ">"
	} else {
		return "<" + Config{}.Name() +
			" handle=" + fmt.Sprint(this.handle) +
			" origin=" + fmt.Sprint(this.origin) +
			" size=" + fmt.Sprint(this.size) +
			" bitmap=" + fmt.Sprint(this.bitmap) +
			">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *element) Size() gopi.Size {
	return gopi.Size{float32(this.size.W), float32(this.size.H)}
}

func (this *element) Origin() gopi.Point {
	return gopi.Point{float32(this.origin.X), float32(this.origin.Y)}
}

func (this *element) Bitmap() bitmap.Bitmap {
	return this.bitmap
}
