// +build egl,dispmanx

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

type Manager struct {
	sync.RWMutex
	gopi.Unit
	gopi.Logger
	gopi.Platform

	display *uint
	handle  dx.Display
	egl     egl.EGLDisplay
	info    dx.DisplayInfo
	bitmap  map[*Bitmap]bool
	surface map[*Surface]uint16
	l       sync.RWMutex // Guards bitmap and surface
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) Define(cfg gopi.Config) error {
	this.display = cfg.FlagUint("display", 0, "Graphics Display Number")
	return nil
}

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger, this.Platform)

	// Open display
	if handle, err := dx.DisplayOpen(uint32(*this.display)); err != nil {
		return err
	} else {
		this.handle = handle
	}

	// Get info on display
	if info, err := dx.DisplayGetInfo(this.handle); err != nil {
		return err
	} else {
		this.info = info
	}

	// Create EGL
	if egl := egl.EGLGetDisplay(uint(this.info.Num())); egl == 0 {
		return gopi.ErrInternalAppError
	} else {
		this.egl = egl
	}
	if _, _, err := egl.EGLInitialize(this.egl); err != nil {
		return err
	}

	// Create bitmapsand surfaces
	this.bitmap = make(map[*Bitmap]bool)
	this.surface = make(map[*Surface]uint16)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	this.l.Lock()
	defer this.RWMutex.Unlock()
	defer this.l.Unlock()

	// Remove surfaces
	var result error
	if update, err := dx.UpdateStart(0); err != nil {
		result = multierror.Append(result, err)
	} else {
		for surface := range this.surface {
			if err := surface.Dispose(update); err != nil {
				result = multierror.Append(result, err)
			}
		}
		// Submit
		if err := dx.UpdateSubmitSync(update); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Remove bitmaps
	for bitmap := range this.bitmap {
		if err := bitmap.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Terminate EGL
	if err := egl.EGLTerminate(this.egl); err != nil {
		result = multierror.Append(result, err)
	}

	// Close display
	if err := dx.DisplayClose(this.handle); err != nil {
		result = multierror.Append(result, err)
	}

	// Release resources
	this.surface = nil
	this.bitmap = nil
	this.egl = 0
	this.handle = 0

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Manager) Size() gopi.Size {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.handle == 0 {
		return gopi.ZeroSize
	} else {
		return gopi.Size{float32(this.info.Width()), float32(this.info.Height())}
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Manager) CreateBackground(gopi.GraphicsContext, gopi.SurfaceFlags) (gopi.Surface, error) {
	// TODO
	return nil, gopi.ErrNotImplemented
}

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

func (this *Manager) CreateBitmap(format gopi.SurfaceFormat, size gopi.Size) (gopi.Bitmap, error) {
	bitmap, err := NewBitmap(format, uint32(size.W), uint32(size.H))
	if err != nil {
		return nil, err
	}

	if err := this.addBitmap(bitmap); err != nil {
		bitmap.Dispose()
		return nil, err
	}

	// Retain bitmap
	bitmap.Retain()

	// Return success
	return bitmap, nil
}

func (this *Manager) DisposeBitmap(bitmap gopi.Bitmap) error {
	if bitmap_, ok := bitmap.(*Bitmap); ok == false {
		return gopi.ErrBadParameter.WithPrefix("DisposeBitmap")
	} else {
		return this.releaseBitmap(bitmap_)
	}
}

////////////////////////////////////////////////////////////////////////////////
// DO

func (this *Manager) Do(cb gopi.SurfaceManagerCallback) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	ctx, err := NewContext(this.handle, 0)
	if err != nil {
		return err
	}

	var result error
	if cb != nil {
		if err := cb(ctx); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Swap buffers by disposing of context
	if err := ctx.Dispose(); err != nil {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<surfacemanager"
	if size := this.Size(); size != gopi.ZeroSize {
		str += fmt.Sprint(" size=", size)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) addBitmap(bitmap *Bitmap) error {
	this.l.Lock()
	defer this.l.Unlock()

	if _, exists := this.bitmap[bitmap]; exists {
		return gopi.ErrDuplicateEntry
	} else {
		this.bitmap[bitmap] = true
	}

	// Return success
	return nil
}

func (this *Manager) addSurface(surface *Surface) error {
	this.l.Lock()
	defer this.l.Unlock()

	if _, exists := this.surface[surface]; exists {
		return gopi.ErrDuplicateEntry
	} else {
		this.surface[surface] = surface.Layer()
	}

	// Return success
	return nil
}

func (this *Manager) releaseBitmap(bitmap *Bitmap) error {
	var result error

	if bitmap.Release() {
		if err := bitmap.Dispose(); err != nil {
			result = multierror.Append(result, err)
		}
		if err := this.delBitmap(bitmap); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

func (this *Manager) delBitmap(bitmap *Bitmap) error {
	this.l.Lock()
	defer this.l.Unlock()

	if _, exists := this.bitmap[bitmap]; exists == false {
		return gopi.ErrNotFound
	} else {
		delete(this.bitmap, bitmap)
	}

	// Return success
	return nil
}

func (this *Manager) delSurface(surface *Surface) error {
	this.l.Lock()
	defer this.l.Unlock()

	if _, exists := this.surface[surface]; exists == false {
		return gopi.ErrNotFound
	} else {
		delete(this.surface, surface)
	}

	// Return success
	return nil
}
