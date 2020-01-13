// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces

import (
	"fmt"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type nativesurface struct {
	handle rpi.DXElement
	size   rpi.DXSize
	origin rpi.DXPoint
}

type surface struct {
	flags   gopi.SurfaceFlags
	opacity float32
	layer   uint16
	native  *nativesurface
	bitmap  *bitmap
}

////////////////////////////////////////////////////////////////////////////////
// NEW AND DESTROY

func NewSurface(flags gopi.SurfaceFlags, opacity float32, layer uint16, native *nativesurface) *surface {
	return &surface{flags, opacity, layer, native, nil}
}

func NewNativeSurface(update rpi.DXUpdate, bitmap *bitmap, display rpi.DXDisplayHandle, flags gopi.SurfaceFlags, opacity float32, layer uint16, origin gopi.Point, size gopi.Size) (*nativesurface, error) {
	// Check update
	if update == 0 {
		return nil, gopi.ErrOutOfOrder.WithPrefix("update")
	}
	// Set alpha
	alpha := rpi.DXAlpha{
		Opacity: uint32(opacity_from_float(opacity)),
	}
	if flags.Mod()&gopi.SURFACE_FLAG_ALPHA_FROM_SOURCE != 0 {
		alpha.Flags |= rpi.DX_ALPHA_FLAG_FROM_SOURCE
	} else {
		alpha.Flags |= rpi.DX_ALPHA_FLAG_FIXED_ALL_PIXELS
	}

	// Clamp, transform and protection
	clamp := rpi.DXClamp{}
	transform := rpi.DX_TRANSFORM_NONE
	protection := rpi.DX_PROTECTION_NONE

	// If there is a bitmap, then the source rectangle is set from that
	dest_rect := rpi.DXNewRect(int32(origin.X), int32(origin.Y), uint32(size.W), uint32(size.H))
	src_size := rpi.DXRectSize(dest_rect)
	dest_size := rpi.DXRectSize(dest_rect)
	dest_origin := rpi.DXRectOrigin(dest_rect)

	// Check size - uint16
	if src_size.W > 0xFFFF || src_size.H > 0xFFFF {
		return nil, gopi.ErrBadParameter.WithPrefix("size")
	}
	if dest_size.W > 0xFFFF || dest_size.H > 0xFFFF {
		return nil, gopi.ErrBadParameter.WithPrefix("size")
	}

	// Adjust size for source
	src_size.W = src_size.W << 16
	src_size.H = src_size.H << 16

	// Get source resource
	src_resource := rpi.DXResource(0)
	if bitmap != nil {
		src_resource = bitmap.handle
	}

	// Create the element
	if handle, err := rpi.DXElementAdd(update, display, layer, dest_rect, src_resource, src_size, protection, alpha, clamp, transform); err != nil {
		return nil, err
	} else {
		return &nativesurface{handle, dest_size, dest_origin}, nil
	}
}

func (this *surface) Destroy(update rpi.DXUpdate) error {
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}

	if this.native != nil {
		if err := this.native.Destroy(update); err != nil {
			return err
		}
	}

	// Release resources
	this.native = nil

	// Return success
	return nil
}

func (this *nativesurface) Destroy(update rpi.DXUpdate) error {
	if update == 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}

	// Remove element
	if this.handle != 0 {
		if err := rpi.DXElementRemove(update, this.handle); err != nil {
			return err
		}
	}

	// Release resouces
	this.handle = 0

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Surface

func (this *surface) Type() gopi.SurfaceFlags {
	return this.flags.Type()
}

func (this *surface) Size() gopi.Size {
	return gopi.Size{float32(this.native.size.W), float32(this.native.size.H)}
}

func (this *surface) Origin() gopi.Point {
	return gopi.Point{float32(this.native.origin.X), float32(this.native.origin.Y)}
}

func (this *surface) Opacity() float32 {
	return this.opacity
}

func (this *surface) Layer() uint16 {
	return this.layer
}

func (this *surface) Bitmap() gopi.Bitmap {
	return this.bitmap
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *surface) String() string {
	return fmt.Sprintf("<graphics.surface id=0x%08X flags=%v size=%v origin=%v opacity=%v layer=%v>", this.native.handle, this.flags, this.native.size, this.native.origin, this.opacity, this.layer)
}
