// +build dispmanx

package dispmanx

import (
	"image/color"

	gopi "github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Model struct {
	fmt  dx.PixFormat
	fn   func(color.Color) color.Color
	size uint8 // bits per pixel
}

type RGBA32 struct {
	r, g, b, a uint8
}

type RGBX32 struct {
	r, g, b uint8
}

type RGB24 struct {
	r, g, b uint8
}

type RGB16 struct {
	r, g, b uint8 // 565
}

type Bit bool

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	colormodel = map[dx.PixFormat]*Model{
		dx.VC_IMAGE_RGBA32: newModel(dx.VC_IMAGE_RGBA32, rgba32Convert, 32),
		dx.VC_IMAGE_RGBX32: newModel(dx.VC_IMAGE_RGBX32, rgbx32Convert, 32),
		dx.VC_IMAGE_RGB888: newModel(dx.VC_IMAGE_RGB888, rgb888Convert, 24),
		dx.VC_IMAGE_RGB565: newModel(dx.VC_IMAGE_RGB565, rgb565Convert, 16),
		dx.VC_IMAGE_1BPP:   newModel(dx.VC_IMAGE_1BPP, bitConvert, 1),
	}
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func ColorModel(f dx.PixFormat) *Model {
	if model, exists := colormodel[f]; exists {
		return model
	} else {
		return nil
	}
}

func newModel(fmt dx.PixFormat, fn func(color.Color) color.Color, size uint8) *Model {
	this := new(Model)
	this.fmt = fmt
	this.fn = fn
	this.size = size
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Model) Format() dx.PixFormat {
	return this.fmt
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - MODEL

func (this *Model) Convert(src color.Color) color.Color {
	return this.fn(src)
}

func (this *Model) BytesPerLine(width uint32) uint32 {
	// Get number of bits per line, aligned on 8-bit boundaries and divide by 8
	bits := width << 3 * uint32(this.size)
	return dx.AlignUp(bits, 8) >> 3
}

// XOffset returns the byte and bit offset on a line
func (this *Model) XOffset(x uint32) (uint32, uint8) {
	// Calculate bit offset, use int64 to ensure no overrun
	bitoffset := int64(x) << 3 * int64(this.size)
	return uint32(bitoffset >> 3), uint8(bitoffset & 0x08)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - COLOR

func (c RGBA32) RGBA() (r, g, b, a uint32) {
	return uint32(c.r) * 0x0101, uint32(c.g) * 0x0101, uint32(c.b) * 0x0101, uint32(c.a) * 0x0101
}

func (c RGBX32) RGBA() (r, g, b, a uint32) {
	return uint32(c.r) * 0x0101, uint32(c.g) * 0x0101, uint32(c.b) * 0x0101, 0xFFFF
}

func (c RGB24) RGBA() (r, g, b, a uint32) {
	return uint32(c.r) * 0x0101, uint32(c.g) * 0x0101, uint32(c.b) * 0x0101, 0xFFFF
}

func (c RGB16) RGBA() (r, g, b, a uint32) {
	return uint32(c.r)*0x0842 + 1, uint32(c.g) * 0x0410, uint32(c.b)*0x0842 + 1, 0xFFFF
}

func (c Bit) RGBA() (r, g, b, a uint32) {
	if c {
		return 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF
	} else {
		return 0x0000, 0x0000, 0x0000, 0xFFFF
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - CONVERT COLOR TO NATIVE

func rgba32Convert(src color.Color) color.Color {
	if _, ok := src.(RGBA32); ok {
		return src
	}
	r, g, b, a := src.RGBA()
	return RGBA32{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}

func rgbx32Convert(src color.Color) color.Color {
	if _, ok := src.(RGBX32); ok {
		return src
	}
	r, g, b, _ := src.RGBA()
	return RGBX32{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

func rgb888Convert(src color.Color) color.Color {
	if _, ok := src.(RGB24); ok {
		return src
	}
	r, g, b, _ := src.RGBA()
	return RGB24{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

func rgb565Convert(src color.Color) color.Color {
	if _, ok := src.(RGB16); ok {
		return src
	}
	r, g, b, _ := src.RGBA()
	// TODO
	return RGB16{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

func bitConvert(src color.Color) color.Color {
	if _, ok := src.(Bit); ok {
		return src
	}
	r, g, b, _ := src.RGBA()
	y := (19595*r + 38470*g + 7471*b + 1<<15) >> 15 // 0x0000 -> 0xFFFF
	return Bit(y>>15 != 0)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS - CONVERT FORMATS TO NATIVE

// Return dispmanx PixFormat for SurfaceFormat or zero
func pixFormat(format gopi.SurfaceFormat) dx.PixFormat {
	switch format {
	case gopi.SURFACE_FMT_RGBA32:
		return dx.VC_IMAGE_RGBA32
	case gopi.SURFACE_FMT_XRGB32:
		return dx.VC_IMAGE_RGBX32
	case gopi.SURFACE_FMT_RGB888:
		return dx.VC_IMAGE_RGB888
	case gopi.SURFACE_FMT_RGB565:
		return dx.VC_IMAGE_RGB565
	case gopi.SURFACE_FMT_1BPP:
		return dx.VC_IMAGE_1BPP
	default:
		return 0
	}
}

// Return SurfaceFormat for dispmanx PixFormat or zero
func surfaceFormat(format dx.PixFormat) gopi.SurfaceFormat {
	switch format {
	case dx.VC_IMAGE_RGBA32:
		return gopi.SURFACE_FMT_RGBA32
	case dx.VC_IMAGE_RGBX32:
		return gopi.SURFACE_FMT_XRGB32
	case dx.VC_IMAGE_RGB888:
		return gopi.SURFACE_FMT_RGB888
	case dx.VC_IMAGE_RGB565:
		return gopi.SURFACE_FMT_RGB565
	case dx.VC_IMAGE_1BPP:
		return gopi.SURFACE_FMT_1BPP
	default:
		return 0
	}
}
