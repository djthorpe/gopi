// +build dispmanx

package rgba32dx

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	bitmap "github.com/djthorpe/gopi/v3/pkg/graphics/bitmap"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
	multierror "github.com/hashicorp/go-multierror"
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
	sync.Mutex
	dx.Resource
	Buffer

	model bitmap.ColorModel
	w, h  uint32
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Factory) New(model bitmap.ColorModel, w, h uint32) (gopi.Bitmap, error) {
	handle := new(RGBA32)

	// Check parameters
	if model.Format() != gopi.SURFACE_FMT_RGBA32 {
		return nil, gopi.ErrBadParameter.WithPrefix("RGBA32DX")
	} else if w == 0 || h == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("RGBA32DX")
	} else {
		handle.model = model
		handle.w = w
		handle.h = h
	}

	// Create resource
	if resource, err := dx.ResourceCreate(dx.VC_IMAGE_RGBA32, w, h); err != nil {
		return nil, err
	} else {
		handle.Resource = resource
	}

	// Create buffer equal to stride (4 bytes per pixel)
	if err := handle.Buffer.Init(dx.ResourceStride(handle.w * 4)); err != nil {
		dx.ResourceDelete(handle.Resource)
		return nil, err
	}

	// Return success
	return handle, nil
}

func (this *Factory) Dispose(bitmap gopi.Bitmap) error {
	// Delete resource
	var result error
	handle := bitmap.(*RGBA32)
	if err := dx.ResourceDelete(handle.Resource); err != nil {
		result = multierror.Append(result, err)
	} else {
		handle.Resource = 0
		handle.w, handle.h = 0, 0
	}

	// Dispose of buffer
	handle.Buffer.Dispose()

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *RGBA32) String() string {
	str := "<bitmap.rgba32dx"
	if this.Resource != 0 {
		str += fmt.Sprintf(" handle=0x%08X", this.Resource)
		str += fmt.Sprintf(" model=%v", this.ColorModel())
		str += fmt.Sprintf(" size=%v", this.Size())
		str += fmt.Sprintf(" stride=%v bytes", this.stride)
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
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Convert pixel to uint32, write color
	pixel := this.model.Convert(c).(bitmap.RGBA32)
	this.Buffer.Fill(uint32(pixel))

	// Write in all rows with the same data
	for y := uint32(0); y < this.h; y++ {
		this.Buffer.WriteRow(this.Resource, y)
	}
}

func (this *RGBA32) At(x, y int) color.Color {
	if x < 0 || y < 0 || uint32(x) >= this.w || uint32(y) >= this.h {
		return bitmap.RGBA32(0x808080FF)
	} else if err := this.Buffer.ReadRow(this.Resource, uint32(y)); err != nil {
		return bitmap.RGBA32(0x808080FF)
	} else {
		return bitmap.RGBA32(this.Buffer.GetAt(uint32(x)))
	}
}

func (this *RGBA32) SetAt(c color.Color, x, y int) error {
	if x < 0 || y < 0 || uint32(x) >= this.w || uint32(y) >= this.h {
		return gopi.ErrBadParameter
	} else if err := this.Buffer.ReadRow(this.Resource, uint32(y)); err != nil {
		return err
	}
	pixel := this.model.Convert(c).(bitmap.RGBA32)
	this.Buffer.SetAt(uint32(x), uint32(pixel))
	return this.Buffer.WriteRow(this.Resource, uint32(y))
}

func (this *RGBA32) ColorModel() color.Model {
	return this.model
}

func (this *RGBA32) Bounds() image.Rectangle {
	return image.Rectangle{image.Point{0, 0}, image.Point{int(this.w) - 1, int(this.h) - 1}}
}
