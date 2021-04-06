package bitmap

import (
	"fmt"
	"image/color"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

type ColorModel interface {
	color.Model

	Format() gopi.SurfaceFormat
	BitsPerPixel() uint8
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

type RGBA32 uint32

type model struct {
	fmt gopi.SurfaceFormat
	bpp uint8
	fn  func(color.Color) color.Color
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	// Register Color Models
	RegisterColorModel(gopi.SURFACE_FMT_RGBA32, NewModel(gopi.SURFACE_FMT_RGBA32, 32, toRGBA32))
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewModel(fmt gopi.SurfaceFormat, bpp uint8, fn func(color.Color) color.Color) ColorModel {
	m := new(model)
	m.fmt = fmt
	m.fn = fn
	m.bpp = bpp
	return m
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *model) Convert(c color.Color) color.Color {
	return this.fn(c)
}

func (this *model) Format() gopi.SurfaceFormat {
	return this.fmt
}

func (this *model) BitsPerPixel() uint8 {
	return this.bpp
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *model) String() string {
	str := "<bitmap.colormodel"
	str += fmt.Sprint(" fmt=", this.Format())
	str += fmt.Sprint(" bits_per_pixel=", this.BitsPerPixel())
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// RGBA32

func toRGBA32(c color.Color) color.Color {
	if c, ok := c.(RGBA32); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return RGBA32(r<<16&0xFF000000) | RGBA32(g<<8&0x00FF0000) | RGBA32(b<<0&0x0000FF00) | RGBA32(a>>8&0x000000FF)
}

func (p RGBA32) RGBA() (uint32, uint32, uint32, uint32) {
	r := uint32(byte(p>>24)) * 0x0101
	g := uint32(byte(p>>16)) * 0x0101
	b := uint32(byte(p>>8)) * 0x0101
	a := uint32(byte(p)) * 0x0101
	return r, g, b, a
}
