// +build dispmanx

package dispmanx

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Bitmap struct {
	sync.RWMutex
	dx.Resource
	*Model

	w, h   uint32
	stride uint32
	count  uint32
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewBitmap(fmt gopi.SurfaceFormat, w, h uint32) (*Bitmap, error) {
	// Check parameters
	if w == 0 || h == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewBitmap")
	}
	// Get color model, create resource
	if model := ColorModel(fmt); model == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("NewBitmap")
	} else if resource, err := dx.ResourceCreate(model.PixFormat(), w, h); err != nil {
		return nil, err
	} else if bitmap, err := NewBitmapFromResource(resource, model, w, h); err != nil {
		dx.ResourceDelete(resource)
		return nil, err
	} else {
		return bitmap, nil
	}
}

func NewBitmapFromResource(handle dx.Resource, model *Model, w, h uint32) (*Bitmap, error) {
	this := new(Bitmap)
	this.Resource = handle
	this.Model = model
	this.w, this.h = w, h

	// Align up to 16-byte boundary
	this.stride = dx.ResourceStride(this.Model.Pitch(this.w))

	// Return success
	return this, nil
}

func (this *Bitmap) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Delete resource
	var result error
	if err := dx.ResourceDelete(this.Resource); err != nil {
		result = multierror.Append(result, err)
	}

	// Release resources
	this.Resource = 0
	this.Model = nil
	this.w, this.h, this.stride = 0, 0, 0

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// RETAIN AND RELEASE

func (this *Bitmap) Retain() {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.count += 1
}

func (this *Bitmap) Release() bool {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.count -= 1
	return this.count == 0
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Bitmap) Format() gopi.SurfaceFormat {
	return this.Model.Format()
}

func (this *Bitmap) Size() gopi.Size {
	return gopi.Size{float32(this.w), float32(this.h)}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Bitmap) At(x, y int) color.Color {
	if x < 0 || y < 0 || x >= int(this.w) || y >= int(this.h) {
		return color.Black
	}
	// TODO
	return color.Black
}

func (this *Bitmap) ColorModel() color.Model {
	return this.Model
}

func (this *Bitmap) Bounds() image.Rectangle {
	return image.Rectangle{
		image.Point{0, 0},
		image.Point{int(this.w) - 1, int(this.h) - 1},
	}
}

func (this *Bitmap) ClearToColor(c color.Color) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Make a row of data (bytes per line, number of lines)
	data := dx.NewData(this.stride)
	if data == nil {
		return
	}

	// TODO: Write color
	buf := data.Byte(0)
	for i := range buf {
		buf[i] = 0xFF
	}

	// Write in all rows with the same data
	for y := uint32(0); y < this.h; y++ {
		if err := this.Write(data, this.Resource, y, 1); err != nil {
			return
		}
	}

	// Dispose of data
	data.Dispose()
}

// Write data to GPU memory with the y-axis bounds as y and h
func (this *Bitmap) Write(src *dx.Data, dest dx.Resource, y, h uint32) error {
	stride := src.Cap()
	rect := dx.NewRect(0, int32(y), stride, h)
	return dx.ResourceWrite(dest, 0, stride, src.PtrMinusOffset(y*stride), rect)
}

// Read data from GPU to buffer
func (this *Bitmap) Read(src dx.Resource, dest *dx.Data, y, h uint32) error {
	stride := dest.Cap()
	rect := dx.NewRect(0, int32(y), stride, h)
	return dx.ResourceRead(src, rect, dest.PtrMinusOffset(y*stride), stride)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Bitmap) String() string {
	str := "<bitmap"
	str += fmt.Sprint(" fmt=", this.Model.Format())
	str += fmt.Sprintf(" size={%d,%d} stride=%d", this.w, this.h, this.stride)
	return str + ">"
}
