// +build dispmanx,egl

package dispmanx

import (
	"fmt"
	"sync"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	bitmap "github.com/djthorpe/gopi/v3/pkg/graphics/bitmap"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Surfaces struct {
	gopi.Unit
	sync.RWMutex
	*bitmap.Bitmaps

	surface map[*Surface]uint16
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	DEFAULT_UPDATE_PRIORITY = 0
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Surfaces) New(gopi.Config) error {
	this.Require(this.Bitmaps)

	// Create surfaces
	this.surface = make(map[*Surface]uint16)

	// Return success
	return nil
}

func (this *Surfaces) Dispose(gopi.Config) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Free surfaces
	var result error
	if update, err := dx.UpdateStart(DEFAULT_UPDATE_PRIORITY); err != nil {
		result = multierror.Append(result, err)
	} else {
		for surface := range this.surface {
			if dispose(update, surface); err != nil {
				result = multierror.Append(result, err)
			}
		}
		// Submit
		if err := dx.UpdateSubmitSync(update); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Free bitmaps TODO

	// Free resources
	this.surface = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Surfaces) String() string {
	str := "<dispmanx.surfaces"

	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	for surface := range this.surface {
		str += fmt.Sprint(" ", surface)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Surfaces) NewSurface(ctx *Context, flags gopi.SurfaceFlags, opacity uint8, layer uint16, x, y, w, h uint32) (*Surface, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *Surfaces) DisposeSurface(ctx *Context, surface *Surface) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check arguments
	if ctx == nil || ctx.Valid() == false {
		return gopi.ErrOutOfOrder.WithPrefix("DisposeSurface")
	}
	if _, exists := this.surface[surface]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("DisposeSurface")
	}

	// Dispose surface
	var result error
	if err := dispose(ctx.Update, surface); err != nil {
		result = multierror.Append(result, err)
	}

	// Delete surface
	delete(this.surface, surface)

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func dispose(dx.Update, *Surface) error {

}

/*
func (this *Manager) CreateSurface(ctx gopi.GraphicsContext, flags gopi.SurfaceFlags, opacity float32, layer uint16, origin gopi.Point, size gopi.Size) (gopi.Surface, error) {
	ctx_, ok := ctx.(*Context)
	if ok == false || ctx_.Valid() == false {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateSurface")
	}

	// Convert width and height
	w, h := uint32(size.W), uint32(size.H)
	if w == 0 || h == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateSurface")
	}

	// Get color model for format
	colormodel := ColorModel(gopi.SURFACE_FMT_RGBA32)
	opacity8 := byte(opacity * 255.0)

	// Initialise bitmap
	var bitmap *Bitmap
	var context egl.EGLContext
	var err error
	switch flags & gopi.SURFACE_FLAG_MASK {
	case gopi.SURFACE_FLAG_BITMAP:
		bitmap, err = NewBitmap(colormodel.Format(), w, h)
		if err != nil {
			return nil, err
		}
	case gopi.SURFACE_FLAG_OPENVG:
		if bits := colormodel.EGLConfig(); len(bits) == 0 {
			return nil, gopi.ErrNotImplemented.WithPrefix("CreateSurface")
		} else if config, err := egl.EGLChooseConfig(this.egl, bits[0], bits[1], bits[2], bits[3], egl.EGL_SURFACETYPE_FLAG_WINDOW, egl.EGL_RENDERABLE_FLAG_OPENVG); err != nil {
			return nil, err
		} else if err := egl.EGLBindAPI(egl.EGL_API_OPENVG); err != nil {
			return nil, err
		} else if ctx, err := egl.EGLCreateContext(this.egl, config, nil, nil); err != nil {
			return nil, err
		} else {
			context = ctx
		}
	default:
		return nil, gopi.ErrNotImplemented.WithPrefix("CreateSurface")
	}

	// Create surface, retain bitmap
	if surface, err := NewSurface(ctx_.Update, ctx_.Display, bitmap, context, int32(origin.X), int32(origin.Y), w, h, layer, opacity8); err != nil {
		return nil, err
	} else if err := this.addSurface(surface); err != nil {
		surface.Dispose(ctx_.Update)
		return nil, err
	} else {
		if bitmap != nil {
			bitmap.Retain()
		}
		return surface, nil
	}
}

func (this *Manager) DisposeSurface(ctx gopi.GraphicsContext, surface gopi.Surface) error {
	surface_, ok := surface.(*Surface)
	if ok == false {
		return gopi.ErrBadParameter.WithPrefix("DisposeSurface")
	}
	ctx_, ok := ctx.(*Context)
	if ok == false || ctx_.Valid() == false {
		return gopi.ErrBadParameter.WithPrefix("DisposeSurface")
	}

	// Get bitmap
	bitmap := surface_.bitmap

	// Dispose surface
	var result error
	if err := surface_.Dispose(ctx_.Update); err != nil {
		result = multierror.Append(result, err)
	}

	// Dispose bitmap
	if bitmap != nil {
		if err := this.releaseBitmap(bitmap); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Remove surface
	if err := this.delSurface(surface_); err != nil {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
}
*/
