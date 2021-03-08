// +build dispmanx

package surface

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Surface struct {
	sync.RWMutex
	dx.Element

	x, y   int32
	w, h   uint32
	layer  uint16
	bitmap *Bitmap
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSurface(ctx dx.Update, display dx.Display, x, y int32, w, h uint32) (*Surface, error) {
	this := new(Surface)

	// Check parameters
	if ctx == 0 || w == 0 || h == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewSurface")
	}

	// Create resource for surface
	r := dx.NewRect(x, y, w, h)
	layer := uint16(100)
	if resource, err := dx.ResourceCreate(dx.VC_IMAGE_RGBA32, w, h); err != nil {
		return nil, err
	} else if element, err := dx.ElementAdd(ctx, display, layer, r, resource, r, 0, dx.NewAlphaFromSource(), nil, dx.DISPMANX_NO_ROTATE); err != nil {
		dx.ResourceDelete(resource)
		return nil, err
	} else if bitmap, err := NewBitmapFromResource(resource, dx.VC_IMAGE_RGBA32, w, h); err != nil {
		dx.ElementRemove(ctx, element)
		return nil, err
	} else {
		bitmap.Retain()
		this.Element = element
		this.bitmap = bitmap
		this.x, this.y = x, y
		this.w, this.h = w, h
		this.layer = layer
	}

	// Return success
	return this, nil
}

func (this *Surface) Dispose(ctx dx.Update) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check state
	if this.Element == 0 || ctx == 0 {
		return gopi.ErrOutOfOrder
	}

	// Remove the element
	var result error
	if err := dx.ElementRemove(ctx, this.Element); err != nil {
		result = multierror.Append(result, err)
	}

	// Release bitmap
	if this.bitmap != nil {
		if this.bitmap.Release() {
			// Dispose of bitmap
			if err := this.bitmap.Dispose(); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Release any resources
	this.Element = 0
	this.bitmap = nil
	this.x, this.y, this.w, this.h = 0, 0, 0, 0
	this.layer = 0

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Surface) Origin() gopi.Point {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return gopi.Point{float32(this.x), float32(this.y)}

}

func (this *Surface) Size() gopi.Size {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return gopi.Size{float32(this.w), float32(this.h)}
}

func (this *Surface) Layer() uint16 {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.layer
}

func (this *Surface) Bitmap() gopi.Bitmap {
	return this.bitmap
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Surface) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<surface"
	str += fmt.Sprintf(" origin={%d,%d} size={%d,%d}", this.x, this.y, this.w, this.h)
	str += fmt.Sprint(" layer=", this.layer)
	if this.bitmap != nil {
		str += fmt.Sprint(" bitmap=", this.bitmap)
	}
	return str + ">"
}
