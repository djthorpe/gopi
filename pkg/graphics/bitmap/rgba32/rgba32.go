// +build !dispmanx

package rgba32

import (
	"fmt"
	"image"
	"image/color"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	bitmap "github.com/djthorpe/gopi/v3/pkg/graphics/bitmap"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	bitmap.RegisterFactory(new(Factory), gopi.SURFACE_FMT_RGBA32)
}

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Factory struct{}

type RGBA32 struct {
	w, h   uint32
	stride uint32
	buf    []Pixel
}

type Model struct{}

type Pixel uint32

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Factory) New(fmt gopi.SurfaceFormat, w, h uint32) (gopi.Bitmap, error) {
	handle := new(RGBA32)
	if fmt != gopi.SURFACE_FMT_RGBA32 {
		return nil, gopi.ErrBadParameter.WithPrefix("RGBA32")
	} else if w == 0 || h == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("RGBA32")
	} else {
		handle.w = w
		handle.h = h
	}

	// The stride is on 16-byte boundaries
	handle.stride = bitmap.AlignUp(handle.w<<2, 16)
	handle.buf = make([]Pixel, handle.h*handle.stride)

	// Return success
	return handle, nil
}

func (this *Factory) Dispose(bitmap gopi.Bitmap) error {
	handle := bitmap.(*RGBA32)
	handle.w, handle.h = 0, 0
	handle.buf = nil
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *RGBA32) String() string {
	str := "<bitmap.rgba32"
	if this.buf != nil {
		str += fmt.Sprintf(" format=%q", this.Format())
		str += fmt.Sprintf(" size=%v", this.Size())
		str += fmt.Sprintf(" stride=%v", this.stride)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (RGBA32) Format() gopi.SurfaceFormat {
	return gopi.SURFACE_FMT_RGBA32
}

func (this *RGBA32) Size() gopi.Size {
	return gopi.Size{float32(this.w), float32(this.h)}
}

func (this *RGBA32) ClearToColor(c color.Color) {
	pixel := this.ColorModel().Convert(c).(Pixel)
	for i := range this.buf {
		this.buf[i] = pixel
	}
}

func (this *RGBA32) At(x, y int) color.Color {
	if x < 0 || y < 0 || uint32(x) >= this.w || uint32(y) >= this.h || this.buf == nil {
		return Pixel(0x808080FF)
	} else {
		i := uint32(x) + uint32(y)*(this.stride>>2)
		return Pixel(this.buf[i])
	}
}

func (this *RGBA32) SetAt(c color.Color, x, y int) error {
	if x < 0 || y < 0 || uint32(x) >= this.w || uint32(y) >= this.h || this.buf == nil {
		return gopi.ErrBadParameter
	}
	pixel := this.ColorModel().Convert(c).(Pixel)
	i := uint32(x) + uint32(y)*(this.stride>>2)
	this.buf[i] = pixel
	return nil
}

func (this *RGBA32) ColorModel() color.Model {
	return Model{}
}

func (this *RGBA32) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{0, 0}, image.Point{int(this.w) - 1, int(this.h) - 1}}

}

////////////////////////////////////////////////////////////////////////////////
// COLOR MODEL

func (Model) Convert(c color.Color) color.Color {
	if c, ok := c.(Pixel); ok {
		return c
	}
	r, g, b, a := c.RGBA()
	return Pixel(r<<16&0xFF000000) | Pixel(g<<8&0x00FF0000) | Pixel(b<<0&0x0000FF00) | Pixel(a>>8&0x000000FF)
}

func (p Pixel) RGBA() (uint32, uint32, uint32, uint32) {
	r := uint32(byte(p>>24)) * 0x0101
	g := uint32(byte(p>>16)) * 0x0101
	b := uint32(byte(p>>8)) * 0x0101
	a := uint32(byte(p)) * 0x0101
	return r, g, b, a
}
