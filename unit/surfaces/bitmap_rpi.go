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
	"image"
	"image/color"

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
	dxrow           *rpi.DXData
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

	// Set stride
	stride := rpi.DXAlignUp(b.size.W, 16) * b.bytes_per_pixel

	// Create resource
	if handle, err := rpi.DXResourceCreate(b.dxtype, b.size); err != nil {
		return nil, err
	} else if dxrow := rpi.DXNewData(uint(stride)); dxrow == nil {
		rpi.DXResourceDelete(handle)
		return nil, gopi.ErrInternalAppError.WithPrefix("dxrow")
	} else {
		b.handle = handle
		b.stride = stride
		b.dxrow = dxrow
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

	// Release row
	this.dxrow.Free()

	// Release resources
	this.handle = 0
	this.dxrow = nil

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
// CLEAR TO COLOR

func (this *bitmap) ClearToColor(c gopi.Color) {
	// Create a strip of color
	src := color_to_bytes(c, this.dxtype)
	row := this.dxrow.Bytes()
	for i := uint32(0); i < this.stride; i++ {
		row[i] = src[i%this.bytes_per_pixel]
	}
	// Set the pointer to the strip and move y forward and ptr back for each strip
	ptr := this.dxrow.Ptr()
	rect := rpi.DXNewRect(0, 0, uint32(this.size.W), 1)
	for y := uint32(0); y < this.size.H; y++ {
		rpi.DXRectSet(rect, 0, int32(y), uint32(this.size.W), 1)
		rpi.DXResourceWriteData(this.handle, this.dxtype, this.stride, ptr, rect)

		// Offset pointer backwards - fudge!
		ptr -= uintptr(this.stride)
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION Image

func (this *bitmap) ColorModel() color.Model {
	return this.dxtype
}

func (this *bitmap) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{0, 0}, image.Point{int(this.size.W) - 1, int(this.size.H) - 1}}
}

func (this *bitmap) At(x, y int) color.Color {
	return gopi.ColorRed
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *bitmap) String() string {
	return fmt.Sprintf("<graphics.bitmap id=0x%08X type=%v size=%v stride=%v>", this.handle, this.flags.ConfigString(), this.size, this.stride)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func color_to_bytes(c gopi.Color, t rpi.DXImageType) []byte {
	// Returns color 0000 <= v <= FFFF
	r, g, b, a := c.RGBA()

	// Convert to []byte
	switch t {
	case rpi.DX_IMAGE_TYPE_RGB888:
		return []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8)}
	case rpi.DX_IMAGE_TYPE_RGB565:
		r := uint16(r>>(8+3)) << (5 + 6)
		g := uint16(g>>(8+2)) << 5
		b := uint16(b >> (8 + 3))
		v := r | g | b
		return []byte{byte(v), byte(v >> 8)}
	case rpi.DX_IMAGE_TYPE_RGBA32:
		return []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}
	default:
		return nil
	}
}
