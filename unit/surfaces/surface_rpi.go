// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces

import (
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

////////////////////////////////////////////////////////////////////////////////
// NEW AND DESTROY

func NewSurface(update rpi.DXUpdate, bitmap *bitmap, display rpi.DXDisplayHandle, flags gopi.SurfaceFlags, opacity float32, layer uint16, origin gopi.Point, size gopi.Size) (*nativesurface, error) {
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
