// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
  #cgo CFLAGS:   -I/opt/vc/include -DUSE_VCHIQ_ARM
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL_static -lGLESv2_static -lkhrn_static -lvcos -lvchiq_arm -lbcm_host -lm
  #include <EGL/egl.h>
*/
import "C"
import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	eglDisplay           C.EGLDisplay
	eglSurface           C.EGLSurface
	eglNativeDisplayType C.EGLNativeDisplayType
	eglBoolean           C.EGLBoolean
	eglInt               C.EGLint
	eglError             C.EGLint
	eglConfig            uintptr
	eglConfigAttrib      C.EGLint
	eglAPI               C.EGLenum
	eglRenderableType    C.EGLint
	eglSurfaceType       C.EGLint
)

// Native window structure
type eglNativeWindowType struct {
	// TODO	element dxElement
	width  int
	height int
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EGL_FALSE = eglBoolean(0)
	EGL_TRUE  = eglBoolean(1)
)

const (
	EGL_NO_DISPLAY = uintptr(0)
	EGL_NO_CONFIG  = eglConfig(0)
)

const (
	// Errors
	EGL_SUCCESS             eglError = 0x3000
	EGL_NOT_INITIALIZED              = 0x3001
	EGL_BAD_ACCESS                   = 0x3002
	EGL_BAD_ALLOC                    = 0x3003
	EGL_BAD_ATTRIBUTE                = 0x3004
	EGL_BAD_CONFIG                   = 0x3005
	EGL_BAD_CONTEXT                  = 0x3006
	EGL_BAD_CURRENT_SURFACE          = 0x3007
	EGL_BAD_DISPLAY                  = 0x3008
	EGL_BAD_MATCH                    = 0x3009
	EGL_BAD_NATIVE_PIXMAP            = 0x300A
	EGL_BAD_NATIVE_WINDOW            = 0x300B
	EGL_BAD_PARAMETER                = 0x300C
	EGL_BAD_SURFACE                  = 0x300D
	EGL_CONTEXT_LOST                 = 0x300E /* EGL 1.1 - IMG_power_management */
)

const (
	// QueryString targets
	EGL_VENDOR      eglInt = 0x3053
	EGL_VERSION            = 0x3054
	EGL_EXTENSIONS         = 0x3055
	EGL_CLIENT_APIS        = 0x308D
)

const (
	/* Config attributes */
	EGL_BUFFER_SIZE             eglConfigAttrib = 0x3020
	EGL_ALPHA_SIZE                              = 0x3021
	EGL_BLUE_SIZE                               = 0x3022
	EGL_GREEN_SIZE                              = 0x3023
	EGL_RED_SIZE                                = 0x3024
	EGL_DEPTH_SIZE                              = 0x3025
	EGL_STENCIL_SIZE                            = 0x3026
	EGL_CONFIG_CAVEAT                           = 0x3027
	EGL_CONFIG_ID                               = 0x3028
	EGL_LEVEL                                   = 0x3029
	EGL_MAX_PBUFFER_HEIGHT                      = 0x302A
	EGL_MAX_PBUFFER_PIXELS                      = 0x302B
	EGL_MAX_PBUFFER_WIDTH                       = 0x302C
	EGL_NATIVE_RENDERABLE                       = 0x302D
	EGL_NATIVE_VISUAL_ID                        = 0x302E
	EGL_NATIVE_VISUAL_TYPE                      = 0x302F
	EGL_SAMPLES                                 = 0x3031
	EGL_SAMPLE_BUFFERS                          = 0x3032
	EGL_SURFACE_TYPE                            = 0x3033
	EGL_TRANSPARENT_TYPE                        = 0x3034
	EGL_TRANSPARENT_BLUE_VALUE                  = 0x3035
	EGL_TRANSPARENT_GREEN_VALUE                 = 0x3036
	EGL_TRANSPARENT_RED_VALUE                   = 0x3037
	EGL_NONE                                    = 0x3038 /* Attrib list terminator */
	EGL_BIND_TO_TEXTURE_RGB                     = 0x3039
	EGL_BIND_TO_TEXTURE_RGBA                    = 0x303A
	EGL_MIN_SWAP_INTERVAL                       = 0x303B
	EGL_MAX_SWAP_INTERVAL                       = 0x303C
	EGL_LUMINANCE_SIZE                          = 0x303D
	EGL_ALPHA_MASK_SIZE                         = 0x303E
	EGL_COLOR_BUFFER_TYPE                       = 0x303F
	EGL_RENDERABLE_TYPE                         = 0x3040
	EGL_MATCH_NATIVE_PIXMAP                     = 0x3041 /* Pseudo-attribute (not queryable) */
	EGL_CONFORMANT                              = 0x3042

	/* Minimum and maximum attribute values */
	EGL_ATTRIB_FIRST = EGL_BUFFER_SIZE
	EGL_ATTRIB_MAX   = EGL_CONFORMANT
)

const (
	EGL_OPENGL_ES_BIT  eglRenderableType = 0x0001 /* EGL_RENDERABLE_TYPE mask bits */
	EGL_OPENVG_BIT                       = 0x0002 /* EGL_RENDERABLE_TYPE mask bits */
	EGL_OPENGL_ES2_BIT                   = 0x0004 /* EGL_RENDERABLE_TYPE mask bits */
	EGL_OPENGL_BIT                       = 0x0008 /* EGL_RENDERABLE_TYPE mask bits */
)

const (
	EGL_PBUFFER_BIT                 eglSurfaceType = 0x0001 /* EGL_SURFACE_TYPE mask bits */
	EGL_PIXMAP_BIT                                 = 0x0002 /* EGL_SURFACE_TYPE mask bits */
	EGL_WINDOW_BIT                                 = 0x0004 /* EGL_SURFACE_TYPE mask bits */
	EGL_VG_COLORSPACE_LINEAR_BIT                   = 0x0020 /* EGL_SURFACE_TYPE mask bits */
	EGL_VG_ALPHA_FORMAT_PRE_BIT                    = 0x0040 /* EGL_SURFACE_TYPE mask bits */
	EGL_MULTISAMPLE_RESOLVE_BOX_BIT                = 0x0200 /* EGL_SURFACE_TYPE mask bits */
	EGL_SWAP_BEHAVIOR_PRESERVED_BIT                = 0x0400 /* EGL_SURFACE_TYPE mask bits */
)

const (
	EGL_OPENGL_ES_API eglAPI = 0x30A0
	EGL_OPENVG_API           = 0x30A1
	EGL_OPENGL_API           = 0x30A2
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	EGL_NULL_POINTER = unsafe.Pointer(uintptr(0))
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func to_eglNativeDisplayType(display uint) eglNativeDisplayType {
	return eglNativeDisplayType(uintptr(display))
}

////////////////////////////////////////////////////////////////////////////////
// EGL PRIVATE METHODS

func eglGetError() eglError {
	return eglError(C.eglGetError())
}

func eglInitialize(display eglDisplay) (eglInt, eglInt, eglError) {
	var major, minor C.EGLint
	if C.eglInitialize(C.EGLDisplay(display), (*C.EGLint)(unsafe.Pointer(&major)), (*C.EGLint)(unsafe.Pointer(&minor))) != C.EGLBoolean(EGL_TRUE) {
		return 0, 0, eglGetError()
	} else {
		return eglInt(major), eglInt(minor), EGL_SUCCESS
	}
}

func eglTerminate(display eglDisplay) eglError {
	if C.eglTerminate(C.EGLDisplay(display)) != C.EGLBoolean(EGL_TRUE) {
		return eglGetError()
	} else {
		return EGL_SUCCESS
	}
}

func eglGetDisplay(display_id eglNativeDisplayType) (eglDisplay, eglError) {
	if display := eglDisplay(C.eglGetDisplay(C.EGLNativeDisplayType(display_id))); display == eglDisplay(EGL_NO_DISPLAY) {
		return display, eglGetError()
	} else {
		return display, EGL_SUCCESS
	}
}

func eglQueryString(display eglDisplay, value eglInt) string {
	return C.GoString(C.eglQueryString(C.EGLDisplay(display), C.EGLint(value)))
}

func eglGetConfigs(display eglDisplay) ([]eglConfig, eglError) {
	var num_config C.EGLint
	if C.eglGetConfigs(C.EGLDisplay(display), (*C.EGLConfig)(EGL_NULL_POINTER), C.EGLint(0), &num_config) != C.EGLBoolean(EGL_TRUE) {
		return nil, eglGetError()
	}
	if num_config == C.EGLint(0) {
		return nil, EGL_BAD_CONFIG
	}
	// configs is a slice so we need to pass the slice pointer
	configs := make([]eglConfig, num_config)
	if C.eglGetConfigs(C.EGLDisplay(display), (*C.EGLConfig)(unsafe.Pointer(&configs[0])), num_config, &num_config) != C.EGLBoolean(EGL_TRUE) {
		return nil, eglGetError()
	} else {
		return configs, EGL_SUCCESS
	}
}

func eglGetConfigAttribs(display eglDisplay, config eglConfig) (map[eglConfigAttrib]eglInt, eglError) {
	attribs := make(map[eglConfigAttrib]eglInt, 0)
	for k := EGL_ATTRIB_FIRST; k <= EGL_ATTRIB_MAX; k++ {
		if v, err := eglGetConfigAttrib(display, config, k); err == EGL_BAD_ATTRIBUTE {
			continue
		} else if err != EGL_SUCCESS {
			return nil, err
		} else {
			attribs[k] = v
		}
	}
	return attribs, EGL_SUCCESS
}

func eglGetConfigAttrib(display eglDisplay, config eglConfig, attrib eglConfigAttrib) (eglInt, eglError) {
	var value C.EGLint
	if C.eglGetConfigAttrib(C.EGLDisplay(display), C.EGLConfig(config), C.EGLint(attrib), &value) != C.EGLBoolean(EGL_TRUE) {
		return eglInt(0), eglGetError()
	} else {
		return eglInt(value), EGL_SUCCESS
	}
}

func eglChooseConfig(display eglDisplay, attributes map[eglConfigAttrib]eglInt) ([]eglConfig, eglError) {
	var num_config C.EGLint

	// Make list of attributes as eglInt values
	attribute_list := make([]C.EGLint, len(attributes)*2+1)
	i := 0
	for k, v := range attributes {
		attribute_list[i] = C.EGLint(k)
		attribute_list[i+1] = C.EGLint(v)
		i = i + 2
	}
	attribute_list[i] = C.EGLint(EGL_NONE)

	// Get number of configurations this matches
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(EGL_NULL_POINTER), C.EGLint(0), &num_config) != C.EGLBoolean(EGL_TRUE) {
		return nil, eglGetError()
	}
	// Return EGL_BAD_ATTRIBUTE if the attribute set doesn't match
	if num_config == 0 {
		return nil, EGL_BAD_ATTRIBUTE
	}
	// Allocate an array
	configs := make([]eglConfig, num_config)
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(EGL_NULL_POINTER), C.EGLint(0), &num_config) != C.EGLBoolean(EGL_TRUE) {
		return nil, eglGetError()
	}
	// Return the configurations
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(unsafe.Pointer(&configs[0])), num_config, &num_config) != C.EGLBoolean(EGL_TRUE) {
		return nil, eglGetError()
	} else {
		return configs, EGL_SUCCESS
	}
}

func eglCreateWindowSurface(display eglDisplay, config eglConfig, native eglNativeWindowType) (eglSurface, eglError) {
	return nil, EGL_BAD_SURFACE
}

func eglCreatePbufferSurface(display eglDisplay, config eglConfig, native eglNativeWindowType) (eglSurface, eglError) {
	return nil, EGL_BAD_SURFACE
}

func eglCreatePixmapSurface(display eglDisplay, config eglConfig, native eglNativeWindowType) (eglSurface, eglError) {
	return nil, EGL_BAD_SURFACE
}

func eglDestroySurface(display eglDisplay, surface eglSurface) eglError {
	if C.eglDestroySurface(C.EGLDisplay(display), C.EGLSurface(surface)) != C.EGLBoolean(EGL_TRUE) {
		return eglGetError()
	} else {
		return EGL_SUCCESS
	}
}

/*
EGLAPI EGLSurface EGLAPIENTRY eglCreateWindowSurface(EGLDisplay dpy, EGLConfig config,
	EGLNativeWindowType win,
	const EGLint *attrib_list);
EGLAPI EGLSurface EGLAPIENTRY eglCreatePbufferSurface(EGLDisplay dpy, EGLConfig config,
	 const EGLint *attrib_list);
EGLAPI EGLSurface EGLAPIENTRY eglCreatePixmapSurface(EGLDisplay dpy, EGLConfig config,
	EGLNativePixmapType pixmap,
	const EGLint *attrib_list);
EGLAPI EGLBoolean EGLAPIENTRY eglDestroySurface(EGLDisplay dpy, EGLSurface surface);
EGLAPI EGLBoolean EGLAPIENTRY eglQuerySurface(EGLDisplay dpy, EGLSurface surface,
EGLint attribute, EGLint *value);

*/

func eglQueryAPI() eglAPI {
	return eglAPI(C.eglQueryAPI())
}

func eglBindAPI(api eglAPI) eglError {
	if C.eglBindAPI(C.EGLenum(api)) != C.EGLBoolean(EGL_TRUE) {
		return eglGetError()
	} else {
		return EGL_SUCCESS
	}
}

////////////////////////////////////////////////////////////////////////////////
// Stringify

func (e eglError) Error() string {
	switch e {
	case EGL_SUCCESS:
		return "EGL_SUCCESS"
	case EGL_NOT_INITIALIZED:
		return "EGL_NOT_INITIALIZED"
	case EGL_BAD_ACCESS:
		return "EGL_BAD_ACCESS"
	case EGL_BAD_ALLOC:
		return "EGL_BAD_ALLOC"
	case EGL_BAD_ATTRIBUTE:
		return "EGL_BAD_ATTRIBUTE"
	case EGL_BAD_CONFIG:
		return "EGL_BAD_CONFIG"
	case EGL_BAD_CONTEXT:
		return "EGL_BAD_CONTEXT"
	case EGL_BAD_CURRENT_SURFACE:
		return "EGL_BAD_CURRENT_SURFACE"
	case EGL_BAD_DISPLAY:
		return "EGL_BAD_DISPLAY"
	case EGL_BAD_MATCH:
		return "EGL_BAD_MATCH"
	case EGL_BAD_NATIVE_PIXMAP:
		return "EGL_BAD_NATIVE_PIXMAP"
	case EGL_BAD_NATIVE_WINDOW:
		return "EGL_BAD_NATIVE_WINDOW"
	case EGL_BAD_PARAMETER:
		return "EGL_BAD_PARAMETER"
	case EGL_BAD_SURFACE:
		return "EGL_BAD_SURFACE"
	case EGL_CONTEXT_LOST:
		return "EGL_CONTEXT_LOST"
	default:
		return "Unknown EGL error"
	}
}

func (a eglConfigAttrib) String() string {
	switch a {
	case EGL_BUFFER_SIZE:
		return "EGL_BUFFER_SIZE"
	case EGL_ALPHA_SIZE:
		return "EGL_ALPHA_SIZE"
	case EGL_BLUE_SIZE:
		return "EGL_BLUE_SIZE"
	case EGL_GREEN_SIZE:
		return "EGL_GREEN_SIZE"
	case EGL_RED_SIZE:
		return "EGL_RED_SIZE"
	case EGL_DEPTH_SIZE:
		return "EGL_DEPTH_SIZE"
	case EGL_STENCIL_SIZE:
		return "EGL_STENCIL_SIZE"
	case EGL_CONFIG_CAVEAT:
		return "EGL_CONFIG_CAVEAT"
	case EGL_CONFIG_ID:
		return "EGL_CONFIG_ID"
	case EGL_LEVEL:
		return "EGL_LEVEL"
	case EGL_MAX_PBUFFER_HEIGHT:
		return "EGL_MAX_PBUFFER_HEIGHT"
	case EGL_MAX_PBUFFER_PIXELS:
		return "EGL_MAX_PBUFFER_PIXELS"
	case EGL_MAX_PBUFFER_WIDTH:
		return "EGL_MAX_PBUFFER_WIDTH"
	case EGL_NATIVE_RENDERABLE:
		return "EGL_NATIVE_RENDERABLE"
	case EGL_NATIVE_VISUAL_ID:
		return "EGL_NATIVE_VISUAL_ID"
	case EGL_NATIVE_VISUAL_TYPE:
		return "EGL_NATIVE_VISUAL_TYPE"
	case EGL_SAMPLES:
		return "EGL_SAMPLES"
	case EGL_SAMPLE_BUFFERS:
		return "EGL_SAMPLE_BUFFERS"
	case EGL_SURFACE_TYPE:
		return "EGL_SURFACE_TYPE"
	case EGL_TRANSPARENT_TYPE:
		return "EGL_TRANSPARENT_TYPE"
	case EGL_TRANSPARENT_BLUE_VALUE:
		return "EGL_TRANSPARENT_BLUE_VALUE"
	case EGL_TRANSPARENT_GREEN_VALUE:
		return "EGL_TRANSPARENT_GREEN_VALUE"
	case EGL_TRANSPARENT_RED_VALUE:
		return "EGL_TRANSPARENT_RED_VALUE"
	case EGL_NONE:
		return "EGL_NONE"
	case EGL_BIND_TO_TEXTURE_RGB:
		return "EGL_BIND_TO_TEXTURE_RGB"
	case EGL_BIND_TO_TEXTURE_RGBA:
		return "EGL_BIND_TO_TEXTURE_RGB"
	case EGL_MIN_SWAP_INTERVAL:
		return "EGL_MIN_SWAP_INTERVAL"
	case EGL_MAX_SWAP_INTERVAL:
		return "EGL_MAX_SWAP_INTERVAL"
	case EGL_LUMINANCE_SIZE:
		return "EGL_LUMINANCE_SIZE"
	case EGL_ALPHA_MASK_SIZE:
		return "EGL_ALPHA_MASK_SIZE"
	case EGL_COLOR_BUFFER_TYPE:
		return "EGL_COLOR_BUFFER_TYPE"
	case EGL_RENDERABLE_TYPE:
		return "EGL_RENDERABLE_TYPE"
	case EGL_MATCH_NATIVE_PIXMAP:
		return "EGL_MATCH_NATIVE_PIXMAP"
	case EGL_CONFORMANT:
		return "EGL_CONFORMANT"
	default:
		return "[?? Invalid eglConfigAttrib value]"
	}
}
