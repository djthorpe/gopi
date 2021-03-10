package dispmanx

import (
	"fmt"
	"image/color"

	gopi "github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Model represents a way to convert colors and how a bitmap is represented in memory
type Model struct {
	fmt    gopi.SurfaceFormat
	pixfmt dx.PixFormat
	fn     func(color.Color) color.Color
	size   uint8 // bits per pixel
}

// Set color types
type RGBA32 color.RGBA
type RGBX32 color.RGBA
type RGB24 color.RGBA
type RGB16 color.RGBA
type Bit bool

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	colormodel = map[gopi.SurfaceFormat]*Model{
		gopi.SURFACE_FMT_RGBA32: newModel(gopi.SURFACE_FMT_RGBA32, dx.VC_IMAGE_RGBA32, rgba32Convert, 32),
		gopi.SURFACE_FMT_XRGB32: newModel(gopi.SURFACE_FMT_XRGB32, dx.VC_IMAGE_RGBX32, rgbx32Convert, 32),
		gopi.SURFACE_FMT_RGB888: newModel(gopi.SURFACE_FMT_RGB888, dx.VC_IMAGE_RGB888, rgb24Convert, 24),
		gopi.SURFACE_FMT_RGB565: newModel(gopi.SURFACE_FMT_RGB565, dx.VC_IMAGE_RGB565, rgb16Convert, 16),
		gopi.SURFACE_FMT_1BPP:   newModel(gopi.SURFACE_FMT_1BPP, dx.VC_IMAGE_1BPP, bitConvert, 1),
	}
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func ColorModel(fmt gopi.SurfaceFormat) *Model {
	if model, exists := colormodel[fmt]; exists {
		return model
	} else {
		return nil
	}
}

func newModel(fmt gopi.SurfaceFormat, pixfmt dx.PixFormat, fn func(color.Color) color.Color, bits uint8) *Model {
	return &Model{fmt, pixfmt, fn, bits}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

// Format returns abstract bitmap format
func (this *Model) Format() gopi.SurfaceFormat {
	return this.fmt
}

// PixFormat returns native pixel format
func (this *Model) PixFormat() dx.PixFormat {
	return this.pixfmt
}

// EGLConfig returns number of bits for each plane or nil if unsupported
func (this *Model) EGLConfig() []uint {
	switch this.fmt {
	case gopi.SURFACE_FMT_RGBA32:
		return []uint{8, 8, 8, 8}
	case gopi.SURFACE_FMT_XRGB32:
		return []uint{8, 8, 8, 0}
	case gopi.SURFACE_FMT_RGB888:
		return []uint{8, 8, 8, 0}
	case gopi.SURFACE_FMT_RGB565:
		return []uint{5, 6, 5, 0}
	default:
		return nil
	}
}

// Pitch returns the number of bytes per for "width" pixels
func (this *Model) Pitch(pixels uint32) uint32 {
	bits := pixels * uint32(this.size)
	return dx.AlignUp(bits, 8) >> 3
}

// XOffset returns the byte and bit offset for an x value
func (this *Model) XOffset(x uint32) (uint32, uint8) {
	// Calculate bit offset, use int64 to ensure no overrun
	bitoffset := int64(x) << 3 * int64(this.size)
	return uint32(bitoffset >> 3), uint8(bitoffset & 0x08)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - MODEL

func (this *Model) Convert(src color.Color) color.Color {
	return this.fn(src)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - color.Color implementation

func (c RGBA32) RGBA() (uint32, uint32, uint32, uint32) {
	return color.RGBA(c).RGBA()
}

func (c RGBX32) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b, _ := color.RGBA(c).RGBA()
	return r, g, b, 0xFFFF
}

func (c RGB24) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b, _ := color.RGBA(c).RGBA()
	return r, g, b, 0xFFFF
}

func (c RGB16) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b, _ := color.RGBA(c).RGBA()
	return r, g, b, 0xFFFF
}

func (c Bit) RGBA() (uint32, uint32, uint32, uint32) {
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
	return RGBX32{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xFF}
}

func rgb24Convert(src color.Color) color.Color {
	if _, ok := src.(RGB24); ok {
		return src
	}
	r, g, b, _ := src.RGBA()
	return RGB24{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xFF}
}

func rgb16Convert(src color.Color) color.Color {
	if _, ok := src.(RGB16); ok {
		return src
	}
	r, g, b, _ := src.RGBA()
	return RGB16{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), 0xFF}
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
// STRINGIFY

func (this *Model) String() string {
	str := "<colormodel"
	str += fmt.Sprint(" fmt=", this.fmt)
	str += fmt.Sprint(" pixfmt=", this.pixfmt)
	str += fmt.Sprint(" bits_per_pixel=", this.size)
	return str + ">"
}
