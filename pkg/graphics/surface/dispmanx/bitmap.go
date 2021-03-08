// +build dispmanx

package surface

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Bitmap struct {
	sync.RWMutex
	dx.Resource
	dx.PixFormat

	w, h, stride uint32
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewBitmap(f dx.PixFormat, w, h uint32) (*Bitmap, error) {
	this := new(Bitmap)

	// Check parameters
	if w == 0 || h == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewBitmap")
	}
	// Create resource
	if resource, err := dx.ResourceCreate(f, w, h); err != nil {
		return nil, err
	} else {
		this.Resource = resource
		this.PixFormat = f
		this.w, this.h = w, h
		this.stride = dx.ResourceStride(this.w)
	}

	// Return success
	return this, nil
}

func (this *Bitmap) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if err := dx.ResourceDelete(this.Resource); err != nil {
		return err
	}

	// Release resources
	this.Resource = 0
	this.PixFormat = 0
	this.w, this.h, this.stride = 0, 0, 0

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Bitmap) Format() gopi.SurfaceFormat {
	return surfaceFormat(this.PixFormat)
}

func (this *Bitmap) Size() gopi.Size {
	return gopi.Size{float32(this.w), float32(this.h)}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Bitmap) At(x, y int) color.Color {
	// TODO
	return color.Black
}

func (this *Bitmap) ColorModel() color.Model {
	// TODO
	return nil
}

func (this *Bitmap) Bounds() image.Rectangle {
	return image.Rectangle{
		image.Point{0, 0},
		image.Point{int(this.w) - 1, int(this.h) - 1},
	}
}

func (this *Bitmap) ClearToColor(color.Color) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// TODO
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Bitmap) String() string {
	str := "<bitmap"
	str += fmt.Sprint("fmt=", this.PixFormat)
	str += fmt.Sprintf("size={%d,%d} stride=%d", this.w, this.h, this.stride)
	return str + ">"
}
