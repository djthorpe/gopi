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

	fd   uintptr
	gpu  *gbm.GBMDevice
	egl  egl.EGLDisplay
	name string
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

	// Check for necessary extensions
	/*
		for _, ext := range []string{"EGL_KHR_create_context", "EGL_KHR_surfaceless_context"} {
			if this.EGLHasExtension(ext) == false {
				return gopi.ErrNotImplemented.WithPrefix("EGL Extension: ", ext)
			}
		}
	*/

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

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

func (this *Manager) EGLExtensions() []string {
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
	for _, ext := range this.EGLExtensions() {
		if strings.ToLower(ext) == strings.ToLower(name) {
			return true
		}
	}
	return false
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
	if extensions := this.EGLExtensions(); len(extensions) > 0 {
		str += " egl.extensions=" + fmt.Sprint(extensions)
	}
	if apis := this.EGLClientApis(); len(apis) > 0 {
		str += " egl.client_apis=" + fmt.Sprint(apis)
	}
	return str + ">"
}
