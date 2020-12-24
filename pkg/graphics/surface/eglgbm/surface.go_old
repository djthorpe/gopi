// +build egl,gbm

package surface

import (
	"fmt"
	"sync"
	"unsafe"

	gopi "github.com/djthorpe/gopi/v3"
	egl "github.com/djthorpe/gopi/v3/pkg/sys/egl"
	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Surface struct {
	sync.RWMutex

	// Surface & context handles
	gbm         *gbm.GBMSurface
	egl_surface egl.EGLSurface
	egl_context egl.EGLContext

	// Surface properties
	x, y, w, h uint32
	dirty      bool
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) NewSurface(display egl.EGLDisplay, api gopi.SurfaceFlags, format gopi.SurfaceFormat, width, height uint32) (*Surface, error) {
	surface := new(Surface)

	// Check parameters
	egl_api, supported := isSupportedApi(api)
	if supported == false {
		return nil, gopi.ErrBadParameter.WithPrefix(api)
	}
	if width == 0 || height == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("size")
	}

	// Determine GBM format
	gbm_format := gbmBufferFormat(format)
	if gbm_format == gbm.GBM_BO_FORMAT_NONE {
		return nil, gopi.ErrBadParameter.WithPrefix("format")
	}

	// Create GBM surface
	if ctx, err := this.gbm.SurfaceCreate(width, height, gbm_format, gbm.GBM_BO_USE_SCANOUT|gbm.GBM_BO_USE_RENDERING); err != nil {
		return nil, err
	} else {
		surface.gbm = ctx
		surface.w = width
		surface.h = height
		surface.dirty = true
	}

	// Create EGL surface
	if api != gopi.SURFACE_FLAG_BITMAP {
		if err := egl.EGLBindAPI(egl_api); err != nil {
			surface.gbm.Free()
			return nil, err
		} else if attrs := eglAttributesForParams(format, api); attrs == nil {
			surface.gbm.Free()
			return nil, gopi.ErrBadParameter.WithPrefix(format)
		} else if config, err := eglChooseConfig(display, attrs); err != nil {
			surface.gbm.Free()
			return nil, err
		} else if egl_context, err := egl.EGLCreateContext(display, config, nil, nil); err != nil {
			surface.gbm.Free()
			return nil, gopi.ErrNotFound.WithPrefix(format)
		} else if egl_surface, err := egl.EGLCreateSurface(display, config, egl.EGLNativeWindow(unsafe.Pointer(surface.gbm))); err != nil {
			egl.EGLDestroyContext(display, egl_context)
			surface.gbm.Free()
			return nil, err
		} else {
			surface.egl_surface = egl_surface
			surface.egl_context = egl_context
		}
	}

	return surface, nil
}

func (this *Surface) Dispose(display egl.EGLDisplay) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	if this.egl_context != nil {
		if err := egl.EGLDestroyContext(display, this.egl_context); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if this.egl_surface != nil {
		if err := egl.EGLDestroySurface(display, this.egl_surface); err != nil {
			result = multierror.Append(result, err)
		}
	}

	/* TODO
	if (previous_bo) {
		drmModeRmFB (device, previous_fb);
		gbm_surface_release_buffer (gbm_surface, previous_bo);
		}
	*/

	if this.gbm != nil {
		this.gbm.Free()
	}

	// Release resources
	this.gbm = nil
	this.egl_surface = nil
	this.egl_context = nil
	this.w, this.h = 0, 0

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Surface) Size() gopi.Size {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return gopi.Size{float32(this.w), float32(this.h)}
}

func (this *Surface) Origin() gopi.Point {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return gopi.Point{float32(this.x), float32(this.y)}
}

func (this *Surface) Dirty() bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.dirty
}

func (this *Surface) SetDirty() {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.dirty = true
}

func (this *Surface) SetClean() {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	this.dirty = false
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Surface) EGLSwapBuffers(display egl.EGLDisplay) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.egl_surface == nil {
		return nil
	}

	if err := egl.EGLSwapBuffers(display, this.egl_surface); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *Surface) GBMSwapBuffers() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.gbm == nil {
		return nil
	}

	if buffer := this.gbm.RetainBuffer(); buffer == nil {
		return gopi.ErrOutOfOrder.WithPrefix("GBMSwapBuffers")
	} else {
		fmt.Println("Buffer=", buffer)
		this.gbm.ReleaseBuffer(buffer)
	}

	// Return success
	return nil
}

func (this *Surface) Draw() error {
	fmt.Println("TODO: Draw on surface")
	/*
		glClearColor (1.0f-progress, progress, 0.0, 1.0);
		glClear (GL_COLOR_BUFFER_BIT);
	*/
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Surface) String() string {
	str := "<surface.eglgbm"

	if size := this.Size(); size != gopi.ZeroSize {
		origin := this.Origin()
		str += " origin=" + fmt.Sprint(origin)
		str += " size=" + fmt.Sprint(size)
	}
	if this.gbm != nil {
		str += " gbm=" + fmt.Sprint(this.gbm)
	}
	if this.egl_surface != nil {
		str += " egl=" + fmt.Sprint(this.egl_surface)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func isSupportedApi(api gopi.SurfaceFlags) (egl.EGLAPI, bool) {
	eglapi, exists := egl.EGLAPIMap[api]
	switch api {
	case gopi.SURFACE_FLAG_BITMAP:
		return eglapi, true
	default:
		return eglapi, exists
	}
}
