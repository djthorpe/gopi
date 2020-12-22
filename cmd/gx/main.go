package main

import (
	"fmt"
	"os"
	"time"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/sys/drm"
	"github.com/djthorpe/gopi/v3/pkg/sys/egl"
	"github.com/djthorpe/gopi/v3/pkg/sys/gbm"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type gx struct {
	dev        *os.File
	device     *gbm.GBMDevice
	display    egl.EGLDisplay
	context    egl.EGLContext
	surface    *gbm.GBMSurface
	eglsurface egl.EGLSurface
	buffer     *gbm.GBMBuffer
	connector  *drm.ModeConnector
	encoder    *drm.ModeEncoder
	mode       drm.ModeInfo
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC

func main() {
	if err := Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func Run() error {
	var result error

	this := new(gx)
	if err := this.Init(); err != nil {
		result = multierror.Append(result, err)
	}

	if result == nil {
		if drawable, err := NewEGL3Drawable(); err != nil {
			result = multierror.Append(result, err)
		} else if err := this.Loop(drawable); err != nil {
			result = multierror.Append(result, err)
		} else {
			drawable.Dispose()
		}
	}

	if err := this.Terminate(); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

func (this *gx) Terminate() error {
	var result error

	if this.buffer != nil {
		this.surface.ReleaseBuffer(this.buffer)
	}

	if this.eglsurface != nil {
		if err := egl.EGLDestroySurface(this.display, this.eglsurface); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if this.surface != nil {
		this.surface.Free()
	}

	if this.context != egl.EGLContext(nil) {
		if err := egl.EGLDestroyContext(this.display, this.context); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if this.display != 0 {
		if err := egl.EGLTerminate(this.display); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if this.connector != nil {
		this.connector.Free()
	}

	if this.encoder != nil {
		this.encoder.Free()
	}

	if this.device != nil {
		this.device.Free()
	}

	if this.dev != nil {
		if err := this.dev.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	this.eglsurface = nil
	this.surface = nil
	this.context = egl.EGLContext(nil)
	this.display = 0
	this.connector = nil
	this.encoder = nil
	this.device = nil
	this.dev = nil
	return nil
}

func (this *gx) ChooseDisplay() (*drm.ModeConnector, *drm.ModeEncoder, drm.ModeInfo, error) {
	res, err := drm.GetResources(this.dev.Fd())
	if err != nil {
		return nil, nil, drm.ModeInfo{}, err
	}
	for _, id := range res.Connectors() {
		connector, err := drm.GetConnector(this.dev.Fd(), id)
		if err != nil {
			continue
		}
		if connector.Status() != drm.ModeConnectionConnected {
			connector.Free()
			continue
		}
		modes := connector.Modes()
		if len(modes) == 0 {
			connector.Free()
			continue
		}
		encoder, err := drm.GetEncoder(this.dev.Fd(), connector.Encoder())
		if err != nil {
			connector.Free()
			continue
		}
		return connector, encoder, modes[0], nil
	}
	return nil, nil, drm.ModeInfo{}, gopi.ErrNotFound
}

func (this *gx) Init() error {
	if dev, err := drm.OpenDevice("card1"); err != nil {
		return err
	} else {
		this.dev = dev
	}

	if connector, encoder, mode, err := this.ChooseDisplay(); err != nil {
		return fmt.Errorf("ChooseDisplay: %w", err)
	} else {
		this.connector = connector
		this.encoder = encoder
		this.mode = mode
	}

	if device, err := gbm.GBMCreateDevice(this.dev.Fd()); err != nil {
		return fmt.Errorf("GBMCreateDevice: %w", err)
	} else {
		this.device = device
	}

	this.display = egl.EGLGetDisplay(this.device)
	if _, _, err := egl.EGLInitialize(this.display); err != nil {
		return fmt.Errorf("EGLInitialize: %w", err)
	}

	if err := egl.EGLBindAPI(egl.EGL_API_OPENGL_ES); err != nil {
		return fmt.Errorf("EGLBindAPI: %w", err)
	}

	config, err := egl.EGLChooseConfig(this.display, 8, 8, 8, 8, egl.EGL_SURFACETYPE_FLAG_WINDOW, egl.EGL_RENDERABLE_FLAG_OPENGL_ES3)
	if err != nil {
		return fmt.Errorf("EGLChooseConfig: %w", err)
	} else if attrs, err := egl.EGLGetConfigAttribs(this.display, config); err != nil {
		return fmt.Errorf("EGLGetConfigAttribs: %w", err)
	} else {
		fmt.Println(attrs)
	}

	version := map[egl.EGLConfigAttrib]int{
		egl.EGL_CONTEXT_CLIENT_VERSION: 3,
	}
	if context, err := egl.EGLCreateContext(this.display, config, nil, version); err != nil {
		return fmt.Errorf("EGLCreateContext: %w", err)
	} else {
		this.context = context
	}

	if surface, err := this.device.SurfaceCreate(1920, 1080, gbm.GBM_BO_FORMAT_XRGB8888, gbm.GBM_BO_USE_SCANOUT|gbm.GBM_BO_USE_RENDERING); err != nil {
		return fmt.Errorf("SurfaceCreate: %w", err)
	} else {
		this.surface = surface
	}

	if eglsurface, err := egl.EGLCreateSurface(this.display, config, egl.EGLNativeWindow(unsafe.Pointer(this.surface))); err != nil {
		return fmt.Errorf("EGLCreateSurface: %w", err)
	} else {
		this.eglsurface = eglsurface
	}

	if err := egl.EGLMakeCurrent(this.display, this.eglsurface, this.eglsurface, this.context); err != nil {
		return fmt.Errorf("EGLMakeCurrent: %w", err)
	}

	return nil
}

func (this *gx) Loop(drawable Drawable) error {
	if err := egl.EGLSwapBuffers(this.display, this.eglsurface); err != nil {
		return fmt.Errorf("EGLSwapBuffers: %w", err)
	}

	// Get free buffer
	buffer := this.surface.RetainBuffer()
	if buffer == nil {
		return fmt.Errorf("RetainBuffer: returned nil")
	} else {
		this.buffer = buffer
	}

	// New Frame Buffer
	fb, err := buffer.NewFrameBuffer()
	if err != nil {
		return fmt.Errorf("GetFrameBuffer: %w", err)
	}

	// Set CRTC mode
	if err := drm.SetCrtc(this.dev.Fd(), this.encoder.Crtc(), this.connector.Id(), fb, 0, 0, &this.mode); err != nil {
		return fmt.Errorf("SetCrtc: %w", err)
	}

	now := time.Now()
	for i := 0; i < 1000; i++ {
		drawable.Draw()
		if err := this.Flip(); err != nil {
			return fmt.Errorf("Iteration %v: %v", i, err)
		}
	}
	fmt.Println("fps=", 1000.0/time.Since(now).Seconds())

	// Return success
	return nil
}

func (this *gx) Flip() error {

	// EGL Swap Buffers
	if err := egl.EGLSwapBuffers(this.display, this.eglsurface); err != nil {
		return fmt.Errorf("EGLSwapBuffers: %w", err)
	}

	// Retain a free buffer
	buffer := this.surface.RetainBuffer()
	if buffer == nil {
		return fmt.Errorf("RetainBuffer: returned nil")
	}

	// Frame Buffer
	fb, err := buffer.NewFrameBuffer()
	if err != nil {
		return fmt.Errorf("GetFrameBuffer: %w", err)
	}

	// Flip and wait until completed
	if err := drm.PageFlip(this.dev.Fd(), this.encoder.Crtc(), fb); err != nil {
		return fmt.Errorf("PageFlip: %w", err)
	}

	if this.buffer != nil {
		this.surface.ReleaseBuffer(this.buffer)
	}
	this.buffer = buffer

	// Return success
	return nil
}
