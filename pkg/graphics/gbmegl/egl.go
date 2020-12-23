// +build egl

package gbmegl

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	gopi "github.com/djthorpe/gopi/v3"
	egl "github.com/djthorpe/gopi/v3/pkg/sys/egl"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type EGL struct {
	sync.RWMutex

	gbm      *GBM
	display  egl.EGLDisplay
	ext      map[string]bool
	surfaces []*Surface
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	apiMap = map[string]egl.EGLAPI{
		"OpenGL_ES": egl.EGL_API_OPENGL_ES,
		"OpenGL":    egl.EGL_API_OPENGL,
		"OpenVG":    egl.EGL_API_OPENVG,
	}
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEGL(gbm *GBM) (*EGL, error) {
	this := new(EGL)

	if display := egl.EGLGetDisplay(gbm.Device()); display == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewEGL")
	} else {
		this.display = display
		this.gbm = gbm
	}

	if _, _, err := egl.EGLInitialize(this.display); err != nil {
		return nil, err
	}

	this.ext = make(map[string]bool)
	for _, ext := range this.ClientExt() {
		this.ext[ext] = true
	}
	for _, ext := range this.Ext() {
		this.ext[ext] = true
	}

	// Return success
	return this, nil
}

func (this *EGL) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Destroy surfaces
	for _, surface := range this.surfaces {
		if surface == nil {
			continue
		}
		if err := surface.Dispose(this.display); err != nil {
			result = multierror.Append(result)
		}
	}

	// Destroy display
	if this.display != 0 {
		if err := egl.EGLTerminate(this.display); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.surfaces = nil
	this.display = 0
	this.gbm = nil
	this.ext = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *EGL) Vendor() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return egl.EGLQueryString(this.display, egl.EGL_QUERY_VENDOR)
}

func (this *EGL) Version() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return egl.EGLQueryString(this.display, egl.EGL_QUERY_VERSION)
}

func (this *EGL) API() []string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return strings.Fields(egl.EGLQueryString(this.display, egl.EGL_QUERY_CLIENT_APIS))
}

func (this *EGL) ClientExt() []string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return strings.Fields(egl.EGLQueryString(0, egl.EGL_QUERY_CLIENT_APIS))
}

func (this *EGL) Ext() []string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return strings.Fields(egl.EGLQueryString(this.display, egl.EGL_QUERY_CLIENT_APIS))
}

func (this *EGL) HasExt(name string) bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	ext, _ := this.ext[name]
	return ext
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *EGL) BindAPI(api string) error {
	if api, exists := apiMap[api]; exists == false {
		return gopi.ErrBadParameter.WithPrefix(api)
	} else if err := egl.EGLBindAPI(api); err != nil {
		return fmt.Errorf("EGLBindAPI: %w", err)
	}

	// Return success
	return nil
}

func (this *EGL) BoundAPI() (string, error) {
	api, err := egl.EGLQueryAPI()
	if err != nil {
		return "", err
	}
	for k, v := range apiMap {
		if v == api {
			return k, nil
		}
	}

	// API not found
	return "", gopi.ErrNotFound.WithPrefix("BoundAPI")
}

func (this *EGL) RenderableBit(api string, version uint) (egl.EGLRenderableFlag, map[egl.EGLConfigAttrib]int) {
	api_, exists := apiMap[api]
	if exists == false {
		return 0, nil
	}
	switch api_ {
	case egl.EGL_API_OPENGL_ES:
		if version == 1 {
			return egl.EGL_RENDERABLE_FLAG_OPENGL_ES, map[egl.EGLConfigAttrib]int{
				egl.EGL_CONTEXT_CLIENT_VERSION: 1,
			}
		}
		if version == 2 {
			return egl.EGL_RENDERABLE_FLAG_OPENGL_ES2, map[egl.EGLConfigAttrib]int{
				egl.EGL_CONTEXT_CLIENT_VERSION: 2,
			}
		}
		if version == 3 {
			return egl.EGL_RENDERABLE_FLAG_OPENGL_ES3, map[egl.EGLConfigAttrib]int{
				egl.EGL_CONTEXT_CLIENT_VERSION: 3,
			}
		}
	case egl.EGL_API_OPENGL:
		return egl.EGL_RENDERABLE_FLAG_OPENGL, nil
	case egl.EGL_API_OPENVG:
		return egl.EGL_RENDERABLE_FLAG_OPENVG, nil
	}

	// Not found
	return 0, nil
}

func (this *EGL) CreateContextForSurface(api string, version uint, r, g, b, a uint) (egl.EGLConfig, egl.EGLContext, error) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	flags := egl.EGL_SURFACETYPE_FLAG_WINDOW
	if renderable, version := this.RenderableBit(api, version); renderable == 0 {
		return 0, nil, gopi.ErrBadParameter.WithPrefix(api, " ", version)
	} else if config, err := egl.EGLChooseConfig(this.display, r, g, b, a, flags, renderable); err != nil {
		return 0, nil, err
	} else if context, err := egl.EGLCreateContext(this.display, config, nil, version); err != nil {
		return 0, nil, err
	} else {
		return config, context, nil
	}
}

func (this *EGL) CreateSurface(api string, version uint, w, h uint32, format Format) (*Surface, error) {
	if this.display == 0 || this.gbm == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("CreateSurface")
	} else if r, g, b, a := this.gbm.BitsForFormat(format); r == 0 || g == 0 || b == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateSurface: ", format)
	} else if config, context, err := this.CreateContextForSurface(api, version, r, g, b, a); err != nil {
		return nil, err
	} else if gbm_surface, err := this.gbm.NewSurface(w, h, format); err != nil {
		egl.EGLDestroyContext(this.display, context)
		return nil, err
	} else if egl_surface, err := egl.EGLCreateSurface(this.display, config, egl.EGLNativeWindow(unsafe.Pointer(gbm_surface))); err != nil {
		egl.EGLDestroyContext(this.display, context)
		gbm_surface.Free()
		return nil, err
	} else if surface := NewSurface(gbm_surface, egl_surface, context); surface == nil {
		egl.EGLDestroySurface(this.display, egl_surface)
		egl.EGLDestroyContext(this.display, context)
		gbm_surface.Free()
		return nil, gopi.ErrInternalAppError.WithPrefix("CreateSurface")
	} else {
		this.RWMutex.Lock()
		defer this.RWMutex.Unlock()

		this.surfaces = append(this.surfaces, surface)
		return surface, nil
	}
}

func (this *EGL) DestroySurface(surface *Surface) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error
	for i, surface_ := range this.surfaces {
		if surface_ != surface {
			continue
		}
		if err := surface.Dispose(this.display); err != nil {
			result = multierror.Append(result)
		}
		this.surfaces[i] = nil
	}

	// Return errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRIMGIFY

func (this *EGL) String() string {
	str := "<egl"
	if vendor := this.Vendor(); vendor != "" {
		str += " vendor=" + strconv.Quote(vendor)
	}
	if version := this.Version(); version != "" {
		str += " version=" + strconv.Quote(version)
	}
	if api := this.API(); len(api) > 0 {
		str += " apis=" + fmt.Sprint(api)
	}
	return str + ">"
}
