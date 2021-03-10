// +build dispmanx,egl

package dispmanx

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
	egl "github.com/djthorpe/gopi/v3/pkg/sys/egl"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Surface struct {
	sync.RWMutex
	dx.Element

	x, y    int32
	w, h    uint32
	opacity uint8
	layer   uint16
	bitmap  *Bitmap
	context egl.EGLContext
	surface egl.EGLSurface
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSurface(ctx dx.Update, display dx.Display, bitmap *Bitmap, x, y int32, w, h uint32, layer uint16, opacity uint8) (*Surface, error) {
	this := new(Surface)

	// Check parameters
	if ctx == 0 || w == 0 || h == 0 || w > 0xFFFF || h > 0xFFFF {
		return nil, gopi.ErrBadParameter.WithPrefix("NewSurface")
	}

	// Set bounds for src and dest
	dest := dx.NewRect(x, y, w, h)
	src := dx.NewRect(0, 0, w<<16, h<<16)

	// Set src to bitmap size
	var resource dx.Resource
	if bitmap != nil {
		resource = bitmap.Resource
		src = dx.NewRect(0, 0, bitmap.w<<16, bitmap.h<<16)
	}

	// Create native surface
	if element, err := dx.ElementAdd(ctx, display, layer, dest, resource, src, 0, dx.NewAlphaFromSource(opacity), nil, dx.DISPMANX_NO_ROTATE); err != nil {
		dx.ResourceDelete(resource)
		return nil, err
	} else {
		this.Element = element
	}

	// Set surface parameters
	this.bitmap = bitmap
	this.context = context
	this.x, this.y = x, y
	this.w, this.h = w, h
	this.opacity = opacity
	this.layer = layer

	// Return success
	return this, nil
}

// TODO: Dispose of EGLSurface and EGLContext
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

	// Release any resources
	this.Element = 0
	this.bitmap = nil
	this.x, this.y, this.w, this.h = 0, 0, 0, 0
	this.layer, this.opacity = 0, 0

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
