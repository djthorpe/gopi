// +build egl,gbm

package surface

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	display "github.com/djthorpe/gopi/v3/pkg/graphics/display/drm"
	egl "github.com/djthorpe/gopi/v3/pkg/sys/egl"
	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.DisplayManager
	sync.RWMutex

	fd         uintptr
	minx, miny uint32
	maxx, maxy uint32
	gpu        *gbm.GBMDevice
	egl        egl.EGLDisplay
	name       string
	surfaces   []*Surface
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	if this.DisplayManager == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing DisplayManager")
	} else if drm, ok := this.DisplayManager.(display.DisplayManager); ok == false {
		return gopi.ErrInternalAppError.WithPrefix("Invalid DisplayManager")
	} else {
		this.fd = drm.Fd()
		this.minx, this.maxx = drm.Width()
		this.miny, this.maxy = drm.Height()
	}

	if gpu, err := gbm.GBMCreateDevice(this.fd); err != nil {
		return err
	} else {
		this.gpu = gpu
	}

	if major, minor, err := egl.EGLInitialize(egl.EGLGetDisplay(this.gpu)); err != nil {
		return err
	} else {
		this.egl = egl.EGLGetDisplay(this.gpu)
		this.name = fmt.Sprintf("EGL %v.%v", major, minor)
	}

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Dispose of surfaces
	for _, surface := range this.surfaces {
		if surface != nil {
			if err := surface.Dispose(); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	if this.egl != 0 {
		if err := egl.EGLTerminate(this.egl); err != nil {
			result = multierror.Append(result, err)
		}
	}
	if this.gpu != nil {
		this.gpu.Free()
	}

	// Release resources
	this.fd = 0
	this.gpu = nil
	this.egl = 0
	this.surfaces = nil

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Manager) Name() string {
	return this.name
}

func (this *Manager) EGLVendor() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.egl == 0 {
		return ""
	} else {
		return egl.EGLQueryString(this.egl, egl.EGL_QUERY_VENDOR)
	}
}

func (this *Manager) EGLVersion() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.egl == 0 {
		return ""
	} else {
		return egl.EGLQueryString(this.egl, egl.EGL_QUERY_VERSION)
	}
}

func (this *Manager) EGLDeviceExtensions() []string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return strings.Fields(strings.TrimSpace(egl.EGLQueryString(0, egl.EGL_QUERY_EXTENSIONS)))
}

func (this *Manager) EGLDisplayExtensions() []string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.egl == 0 {
		return nil
	} else {
		return strings.Fields(strings.TrimSpace(egl.EGLQueryString(this.egl, egl.EGL_QUERY_EXTENSIONS)))
	}
}

func (this *Manager) EGLClientApis() []egl.EGLAPI {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var result []egl.EGLAPI
	if this.egl == 0 {
		return nil
	} else {
		apis := strings.Fields(strings.TrimSpace(egl.EGLQueryString(this.egl, egl.EGL_QUERY_CLIENT_APIS)))
		for _, api := range apis {
			if surface_type, exists := egl.EGLSurfaceTypeMap[api]; exists == false {
				continue
			} else if api, exists := egl.EGLAPIMap[surface_type]; exists == false {
				continue
			} else {
				result = append(result, api)
			}
		}
	}

	// Return success
	return result
}

func (this *Manager) EGLHasExtension(name string) bool {
	for _, ext := range this.EGLDeviceExtensions() {
		if strings.ToLower(ext) == strings.ToLower(name) {
			return true
		}
	}
	for _, ext := range this.EGLDisplayExtensions() {
		if strings.ToLower(ext) == strings.ToLower(name) {
			return true
		}
	}
	return false
}

func (this *Manager) GBMFormats() []gopi.SurfaceFormat {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.gpu == nil {
		return nil
	}

	results := []gopi.SurfaceFormat{}
	for fmt := gopi.SURFACE_FMT_NONE + 1; fmt <= gopi.SURFACE_FMT_MAX; fmt++ {
		if gbm_fmt := gbmSurfaceFormat(fmt); gbm_fmt == 0 {
			continue
		} else if this.gpu.IsFormatSupported(gbm_fmt, 0) {
			results = append(results, fmt)
		}
	}
	return results
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) CreateBackground(display gopi.Display, flags gopi.SurfaceFlags) (gopi.Surface, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if display == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateBackground")
	}
	width, height := display.Size()
	if surface, err := this.NewSurface(flags, this.egl, gopi.SURFACE_FMT_XRGB32, width, height); err != nil {
		return nil, err
	} else {
		this.surfaces = append(this.surfaces, surface)
		return surface, nil
	}
}

func (this *Manager) DisposeSurface(surface gopi.Surface) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error
	for i, surface_ := range this.surfaces {
		if surface == surface_ {
			if err := surface_.Dispose(); err != nil {
				result = multierror.Append(err)
			}
			this.surfaces[i] = nil
		}
	}
	return result
}

func (this *Manager) SwapBuffers() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error
	for _, surface := range this.surfaces {
		if surface != nil && surface.Dirty() {
			if err := this.swapBuffersForSurface(surface); err != nil {
				result = multierror.Append(result)
			}
		}
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<surfacemanager.eglgbm"
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if vendor := this.EGLVendor(); vendor != "" {
		str += " egl.vendor=" + strconv.Quote(vendor)
	}
	if version := this.EGLVersion(); version != "" {
		str += " egl.version=" + strconv.Quote(version)
	}
	if extensions := this.EGLDeviceExtensions(); len(extensions) > 0 {
		str += " egl.device_extensions=" + fmt.Sprint(extensions)
	}
	if extensions := this.EGLDisplayExtensions(); len(extensions) > 0 {
		str += " egl.display_extensions=" + fmt.Sprint(extensions)
	}
	if apis := this.EGLClientApis(); len(apis) > 0 {
		str += " egl.client_apis=" + fmt.Sprint(apis)
	}
	if fmts := this.GBMFormats(); len(fmts) > 0 {
		str += " gbm.supported_formats=" + fmt.Sprint(fmts)
	}
	if this.maxx > 0 && this.maxy > 0 {
		str += fmt.Sprintf(" size={%v,%v,%v,%v}", this.minx, this.miny, this.maxx, this.maxy)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Manager) swapBuffersForSurface(surface *Surface) error {
	var result error

	// Draw here and set the dirty flag
	if err := surface.Draw(); err != nil {
		result = multierror.Append(result, err)
	}

	if surface.Dirty() == false {
		return nil
	}
	if err := surface.EGLSwapBuffers(this.egl); err != nil {
		result = multierror.Append(result, err)
	}
	if surface.HasFreeBuffers() == false {
		result = multierror.Append(result, fmt.Errorf("swapBuffers: No free buffers"))
	}
	if buffer := surface.RetainBuffer(); buffer == nil {
		result = multierror.Append(result, fmt.Errorf("swapBuffers: Failed to lock front buffer"))
	} else {
		handle := buffer.Handle()
		stride := buffer.Stride()
		bpp := buffer.BitsPerPixel()
		depth := 32 // TODO
		if fb, err := drm.AddFrameBuffer(this.fd, surface.w, surface.h, depth, bpp, stride, handle); err != nil {
			result = multierror.Append(result, err)
		}
		/* else if err := drm.SetCrtc(...); err != nil {
			result = multierror.Append(result,err)
		} else if(previous_bo) {
		  drmModeRmFB (device, previous_fb);
		  gbm_surface_release_buffer (gbm_surface, previous_bo);
		}
		previous_bo = bo;
		previous_fb = fb;*/
	}

	// Indicate surface has been swapped
	surface.SetClean()

	// Return any errors
	return result
}

func gbmSurfaceFormat(fmt gopi.SurfaceFormat) gbm.GBMFormat {
	switch fmt {
	case gopi.SURFACE_FMT_RGBA32:
		return gbm.GBM_FORMAT_RGBA8888
	case gopi.SURFACE_FMT_XRGB32:
		return gbm.GBM_FORMAT_XRGB8888
	case gopi.SURFACE_FMT_RGB888:
		return gbm.GBM_FORMAT_RGB888
	case gopi.SURFACE_FMT_RGB565:
		return gbm.GBM_FORMAT_RGB565
	default:
		return 0
	}
}

func eglApiFlags(flags gopi.SurfaceFlags) int {
	value := int(0)
	if flags&gopi.SURFACE_FLAG_OPENGL != 0 {
		value |= int(egl.EGL_OPENGL_BIT)
	}
	if flags&gopi.SURFACE_FLAG_OPENGL_ES != 0 {
		value |= int(egl.EGL_OPENGL_ES_BIT)
	}
	if flags&gopi.SURFACE_FLAG_OPENGL_ES2 != 0 {
		value |= int(egl.EGL_OPENGL_ES2_BIT)
	}
	if flags&gopi.SURFACE_FLAG_OPENVG != 0 {
		value |= int(egl.EGL_OPENVG_BIT)
	}
	return value
}

func eglAttributesForParams(fmt gopi.SurfaceFormat, flags gopi.SurfaceFlags) map[egl.EGLConfigAttrib]int {
	attribs := make(map[egl.EGLConfigAttrib]int, 10)
	switch fmt {
	case gopi.SURFACE_FMT_RGBA32:
		attribs[egl.EGL_RED_SIZE] = 8
		attribs[egl.EGL_GREEN_SIZE] = 8
		attribs[egl.EGL_BLUE_SIZE] = 8
		attribs[egl.EGL_ALPHA_SIZE] = 8
		attribs[egl.EGL_BUFFER_SIZE] = 32
	case gopi.SURFACE_FMT_XRGB32:
		attribs[egl.EGL_RED_SIZE] = 8
		attribs[egl.EGL_GREEN_SIZE] = 8
		attribs[egl.EGL_BLUE_SIZE] = 8
		attribs[egl.EGL_ALPHA_SIZE] = 0
		attribs[egl.EGL_BUFFER_SIZE] = 32
	case gopi.SURFACE_FMT_RGB888:
		attribs[egl.EGL_RED_SIZE] = 8
		attribs[egl.EGL_GREEN_SIZE] = 8
		attribs[egl.EGL_BLUE_SIZE] = 8
		attribs[egl.EGL_ALPHA_SIZE] = 0
		attribs[egl.EGL_BUFFER_SIZE] = 24
	case gopi.SURFACE_FMT_RGB565:
		attribs[egl.EGL_RED_SIZE] = 5
		attribs[egl.EGL_GREEN_SIZE] = 6
		attribs[egl.EGL_BLUE_SIZE] = 5
		attribs[egl.EGL_ALPHA_SIZE] = 0
		attribs[egl.EGL_BUFFER_SIZE] = 16
	default:
		return nil
	}

	attribs[egl.EGL_SURFACE_TYPE] = int(egl.EGL_WINDOW_BIT)
	attribs[egl.EGL_RENDERABLE_TYPE] = eglApiFlags(flags)

	return attribs
}
