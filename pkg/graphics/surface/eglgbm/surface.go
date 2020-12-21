// +build egl,gbm

package surface

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	egl "github.com/djthorpe/gopi/v3/pkg/sys/egl"
	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Surface struct {
	sync.RWMutex

	ctx        *gbm.GBMSurface
	egl        egl.EGLSurface
	x, y, w, h uint32
	dirty      bool
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) NewSurface(api gopi.SurfaceFlags, display egl.EGLDisplay, format gopi.SurfaceFormat, width, height uint32) (*Surface, error) {
	surface := new(Surface)

	// Check parameters
	egl_api, supported := isSupportedApi(api)
	if supported == false {
		return nil, gopi.ErrBadParameter.WithPrefix(api)
	}
	if width == 0 || height == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("size")
	}
	gbm_format := gbmSurfaceFormat(format)
	if gbm_format == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("format")
	}

	// Create surface
	if ctx, err := this.gpu.SurfaceCreate(width, height, gbm_format, gbm.GBM_BO_USE_SCANOUT|gbm.GBM_BO_USE_RENDERING); err != nil {
		return nil, err
	} else {
		surface.ctx = ctx
		surface.w = width
		surface.h = height
		surface.dirty = true
	}

	// Create EGL surface
	if api != gopi.SURFACE_FLAG_BITMAP {
		if err := egl.EGLBindAPI(egl_api); err != nil {
			surface.ctx.Free()
			return nil, err
		}
		// TODO
		/*
		   eglGetConfigs(display, NULL, 0, &count);
		   configs = malloc(count * sizeof *configs);
		   eglChooseConfig (display, attributes, configs, count, &num_config);
		   config_index = match_config_to_visual(display,GBM_FORMAT_XRGB8888,configs,num_config);

		   context = eglCreateContext (display, configs[config_index], EGL_NO_CONTEXT, context_attribs);
		   egl_surface = eglCreateWindowSurface (display, configs[config_index], gbm_surface, NULL);
		   eglMakeCurrent (display, egl_surface, egl_surface, context);
		*/
	}

	return surface, nil
}

func (this *Surface) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	/* TODO
	if (previous_bo) {
		drmModeRmFB (device, previous_fb);
		gbm_surface_release_buffer (gbm_surface, previous_bo);
		}
	  eglDestroySurface (display, egl_surface);
	*/

	if this.ctx != nil {
		this.ctx.Free()
	}

	// Release resources
	this.ctx = nil
	this.w, this.h = 0, 0

	// Return success
	return nil
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

	if this.egl != nil {
		if err := egl.EGLSwapBuffers(display, this.egl); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

func (this *Surface) HasFreeBuffers() bool {
	if this.ctx == nil {
		return false
	} else {
		return this.ctx.HasFreeBuffers()
	}
}

func (this *Surface) RetainBuffer() *gbm.GBMBuffer {
	if this.ctx == nil {
		return nil
	} else {
		return this.ctx.RetainBuffer()
	}
}

func (this *Surface) ReleaseBuffer(buffer *gbm.GBMBuffer) {
	if this.ctx != nil && buffer != nil {
		this.ctx.ReleaseBuffer(buffer)
	}
}

func (this *Surface) Draw() error {
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
