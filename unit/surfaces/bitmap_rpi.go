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

type bitmap struct {
	flags           gopi.SurfaceFlags
	size            rpi.DXSize
	handle          rpi.DXResource
	stride          uint32
	dxtype          rpi.DXImageType
	bytes_per_pixel uint32
}

////////////////////////////////////////////////////////////////////////////////
// NEW / DESTROY

func NewBitmap(flags gopi.SurfaceFlags, size gopi.Size) (*bitmap, error) {
	// Check parameters
	if flags.Type() != gopi.SURFACE_FLAG_BITMAP {
		return nil, gopi.ErrBadParameter.WithPrefix("flags")
	} else if size.W <= 0.0 || size.H <= 0.0 {
		return nil, gopi.ErrBadParameter.WithPrefix("size")
	}

	// Create bitmap
	b := &bitmap{
		size:  rpi.DXSize{uint32(size.W), uint32(size.H)},
		flags: gopi.SURFACE_FLAG_BITMAP | flags.Config(),
	}
	switch flags.Config() {
	case gopi.SURFACE_FLAG_RGBA32:
		b.dxtype = rpi.DX_IMAGE_TYPE_RGBA32
		b.bytes_per_pixel = 4
	case gopi.SURFACE_FLAG_RGB888:
		b.dxtype = rpi.DX_IMAGE_TYPE_RGB888
		b.bytes_per_pixel = 3
	case gopi.SURFACE_FLAG_RGB565:
		b.dxtype = rpi.DX_IMAGE_TYPE_RGB565
		b.bytes_per_pixel = 2
	default:
		return nil, gopi.ErrNotImplemented
	}

	// Create resource
	if handle, err := rpi.DXResourceCreate(b.dxtype, b.size); err != nil {
		return nil, err
	} else {
		b.handle = handle
		b.stride = rpi.DXAlignUp(b.size.W, 16) * b.bytes_per_pixel
		return b, nil
	}
}

func (this *bitmap) Destroy() error {
	if this.handle == 0 {
		return nil
	}
	if err := rpi.DXResourceDelete(this.handle); err != nil {
		return err
	}
	// Release resourcfes
	this.handle = 0
	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Bitmap

func (this *bitmap) Type() gopi.SurfaceFlags {
	return this.flags.Config()
}

func (this *bitmap) Size() gopi.Size {
	return gopi.Size{float32(this.size.W), float32(this.size.H)}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *bitmap) String() string {
	return fmt.Sprintf("<graphics.bitmap id=0x%08X type=%v size=%v stride=%v>", this.handle, this.flags.ConfigString(), this.size, this.stride)
}
