// +build egl

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package egl

import (
	"unsafe"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: brcmegl
#include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EGLDisplay         C.EGLDisplay
	EGLConfig          C.EGLConfig
	EGLConfigAttrib    C.EGLint
	EGLContext         C.EGLContext
	EGLSurface         C.EGLSurface
	EGLError           C.EGLint
	EGLQuery           C.EGLint
	EGLRenderableFlag  C.EGLint
	EGLSurfaceTypeFlag C.EGLint
	EGLAPI             C.EGLint
	EGLNativeWindow    uintptr
)


////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EGL_FALSE = C.EGLBoolean(0)
	EGL_TRUE  = C.EGLBoolean(1)
)

const (
	// EGLError
	EGL_SUCCESS             EGLError = 0x3000
	EGL_NOT_INITIALIZED     EGLError = 0x3001
	EGL_BAD_ACCESS          EGLError = 0x3002
	EGL_BAD_ALLOC           EGLError = 0x3003
	EGL_BAD_ATTRIBUTE       EGLError = 0x3004
	EGL_BAD_CONFIG          EGLError = 0x3005
	EGL_BAD_CONTEXT         EGLError = 0x3006
	EGL_BAD_CURRENT_SURFACE EGLError = 0x3007
	EGL_BAD_DISPLAY         EGLError = 0x3008
	EGL_BAD_MATCH           EGLError = 0x3009
	EGL_BAD_NATIVE_PIXMAP   EGLError = 0x300A
	EGL_BAD_NATIVE_WINDOW   EGLError = 0x300B
	EGL_BAD_PARAMETER       EGLError = 0x300C
	EGL_BAD_SURFACE         EGLError = 0x300D
	EGL_CONTEXT_LOST        EGLError = 0x300E // EGL 1.1 - IMG_power_management
)

const (
	// EGLQuery
	EGL_QUERY_VENDOR      EGLQuery = 0x3053
	EGL_QUERY_VERSION     EGLQuery = 0x3054
	EGL_QUERY_EXTENSIONS  EGLQuery = 0x3055
	EGL_QUERY_CLIENT_APIS EGLQuery = 0x308D
)

const (
	// EGLConfigAttrib
	EGL_BUFFER_SIZE             EGLConfigAttrib = 0x3020
	EGL_ALPHA_SIZE              EGLConfigAttrib = 0x3021
	EGL_BLUE_SIZE               EGLConfigAttrib = 0x3022
	EGL_GREEN_SIZE              EGLConfigAttrib = 0x3023
	EGL_RED_SIZE                EGLConfigAttrib = 0x3024
	EGL_DEPTH_SIZE              EGLConfigAttrib = 0x3025
	EGL_STENCIL_SIZE            EGLConfigAttrib = 0x3026
	EGL_CONFIG_CAVEAT           EGLConfigAttrib = 0x3027
	EGL_CONFIG_ID               EGLConfigAttrib = 0x3028
	EGL_LEVEL                   EGLConfigAttrib = 0x3029
	EGL_MAX_PBUFFER_HEIGHT      EGLConfigAttrib = 0x302A
	EGL_MAX_PBUFFER_PIXELS      EGLConfigAttrib = 0x302B
	EGL_MAX_PBUFFER_WIDTH       EGLConfigAttrib = 0x302C
	EGL_NATIVE_RENDERABLE       EGLConfigAttrib = 0x302D
	EGL_NATIVE_VISUAL_ID        EGLConfigAttrib = 0x302E
	EGL_NATIVE_VISUAL_TYPE      EGLConfigAttrib = 0x302F
	EGL_SAMPLES                 EGLConfigAttrib = 0x3031
	EGL_SAMPLE_BUFFERS          EGLConfigAttrib = 0x3032
	EGL_SURFACE_TYPE            EGLConfigAttrib = 0x3033
	EGL_TRANSPARENT_TYPE        EGLConfigAttrib = 0x3034
	EGL_TRANSPARENT_BLUE_VALUE  EGLConfigAttrib = 0x3035
	EGL_TRANSPARENT_GREEN_VALUE EGLConfigAttrib = 0x3036
	EGL_TRANSPARENT_RED_VALUE   EGLConfigAttrib = 0x3037
	EGL_NONE                    EGLConfigAttrib = 0x3038 // Attrib list terminator
	EGL_BIND_TO_TEXTURE_RGB     EGLConfigAttrib = 0x3039
	EGL_BIND_TO_TEXTURE_RGBA    EGLConfigAttrib = 0x303A
	EGL_MIN_SWAP_INTERVAL       EGLConfigAttrib = 0x303B
	EGL_MAX_SWAP_INTERVAL       EGLConfigAttrib = 0x303C
	EGL_LUMINANCE_SIZE          EGLConfigAttrib = 0x303D
	EGL_ALPHA_MASK_SIZE         EGLConfigAttrib = 0x303E
	EGL_COLOR_BUFFER_TYPE       EGLConfigAttrib = 0x303F
	EGL_RENDERABLE_TYPE         EGLConfigAttrib = 0x3040
	EGL_MATCH_NATIVE_PIXMAP     EGLConfigAttrib = 0x3041 // Pseudo-attribute (not queryable)
	EGL_CONFORMANT              EGLConfigAttrib = 0x3042
	EGL_COMFIG_ATTRIB_MIN                        = EGL_BUFFER_SIZE
	EGL_COMFIG_ATTRIB_MAX                        = EGL_CONFORMANT
)

const (
	// EGLRenderableFlag
	EGL_RENDERABLE_FLAG_OPENGL_ES  EGLRenderableFlag = 0x0001
	EGL_RENDERABLE_FLAG_OPENVG     EGLRenderableFlag = 0x0002
	EGL_RENDERABLE_FLAG_OPENGL_ES2 EGLRenderableFlag = 0x0004
	EGL_RENDERABLE_FLAG_OPENGL     EGLRenderableFlag = 0x0008
	EGL_RENDERABLE_FLAG_MIN                           = EGL_RENDERABLE_FLAG_OPENGL_ES
	EGL_RENDERABLE_FLAG_MAX                           = EGL_RENDERABLE_FLAG_OPENGL
)

const (
	// EGLSurfaceTypeFlag
	EGL_SURFACETYPE_FLAG_PBUFFER                 EGLSurfaceTypeFlag = 0x0001 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_PIXMAP                  EGLSurfaceTypeFlag = 0x0002 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_WINDOW                  EGLSurfaceTypeFlag = 0x0004 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_VG_COLORSPACE_LINEAR    EGLSurfaceTypeFlag = 0x0020 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_VG_ALPHA_FORMAT_PRE     EGLSurfaceTypeFlag = 0x0040 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_MULTISAMPLE_RESOLVE_BOX EGLSurfaceTypeFlag = 0x0200 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED EGLSurfaceTypeFlag = 0x0400 /* EGL_SURFACE_TYPE mask bits */
	EGL_SURFACETYPE_FLAG_MIN                                         = EGL_SURFACETYPE_FLAG_PBUFFER
	EGL_SURFACETYPE_FLAG_MAX                                         = EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED
)

const (
	// EGLAPI
	EGL_API_NONE      EGLAPI = 0
	EGL_API_OPENGL_ES EGLAPI = 0x30A0
	EGL_API_OPENVG    EGLAPI = 0x30A1
	EGL_API_OPENGL    EGLAPI = 0x30A2
	EGL_API_MIN               = EGL_API_OPENGL_ES
	EGL_API_MAX               = EGL_API_OPENGL
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func EGLGetError() error {
	if err := EGLError(C.eglGetError()); err == EGL_SUCCESS {
		return nil
	} else {
		return err
	}
}

func EGLInitialize(display EGLDisplay) (int, int, error) {
	var major, minor C.EGLint
	if C.eglInitialize(C.EGLDisplay(display), (*C.EGLint)(unsafe.Pointer(&major)), (*C.EGLint)(unsafe.Pointer(&minor))) != EGL_TRUE {		
		return 0, 0, EGLGetError()
	} else {
		return int(major), int(minor), nil
	}
}

func EGLTerminate(display EGLDisplay) error {
	if C.eglTerminate(C.EGLDisplay(display)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

func EGLGetDisplay(display uint) EGLDisplay {
	return EGLDisplay(C.eglGetDisplay(C.EGLNativeDisplayType(uintptr(display))))
}

func EGLQueryString(display EGLDisplay, value EGLQuery) string {
	return C.GoString(C.eglQueryString(C.EGLDisplay(display), C.EGLint(value)))
}


////////////////////////////////////////////////////////////////////////////////
// SURFACE CONFIGS

func EGLGetConfigs(display EGLDisplay) ([]EGLConfig, error) {
	var num_config C.EGLint
	if C.eglGetConfigs(C.EGLDisplay(display), (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	}
	if num_config == C.EGLint(0) {
		return nil, EGL_BAD_CONFIG
	}
	// configs is a slice so we need to pass the slice pointer
	configs := make([]EGLConfig, num_config)
	if C.eglGetConfigs(C.EGLDisplay(display), (*C.EGLConfig)(unsafe.Pointer(&configs[0])), num_config, &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	} else {
		return configs, nil
	}
}

func EGLGetConfigAttrib(display EGLDisplay, config EGLConfig, attrib EGLConfigAttrib) (int, error) {
	var value C.EGLint
	if C.eglGetConfigAttrib(C.EGLDisplay(display), C.EGLConfig(config), C.EGLint(attrib), &value) != EGL_TRUE {
		return 0, EGLGetError()
	} else {
		return int(value), nil
	}
}

func EGLGetConfigAttribs(display EGLDisplay, config EGLConfig) (map[EGLConfigAttrib]int, error) {
	attribs := make(map[EGLConfigAttrib]int, 0)
	for k := EGL_COMFIG_ATTRIB_MIN; k <= EGL_COMFIG_ATTRIB_MAX; k++ {
		if v, err := EGLGetConfigAttrib(display, config, k); err == EGL_BAD_ATTRIBUTE {
			continue
		} else if err != nil {
			return nil, err
		} else {
			attribs[k] = v
		}
	}
	return attribs, nil
}

func EGLChooseConfig_(display EGLDisplay, attributes map[EGLConfigAttrib]int) ([]EGLConfig, error) {
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
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	}
	// Return EGL_BAD_ATTRIBUTE if the attribute set doesn't match
	if num_config == 0 {
		return nil, EGL_BAD_ATTRIBUTE
	}
	// Allocate an array
	configs := make([]EGLConfig, num_config)
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(nil), C.EGLint(0), &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	}
	// Return the configurations
	if C.eglChooseConfig(C.EGLDisplay(display), &attribute_list[0], (*C.EGLConfig)(unsafe.Pointer(&configs[0])), num_config, &num_config) != EGL_TRUE {
		return nil, EGLGetError()
	} else {
		return configs, nil
	}
}

func EGLChooseConfig(display EGLDisplay, r_bits, g_bits, b_bits, a_bits uint, surface_type EGLSurfaceTypeFlag, renderable_type EGLRenderableFlag) (EGLConfig, error) {
	if configs, err := EGLChooseConfig_(display, map[EGLConfigAttrib]int{
		EGL_RED_SIZE:        int(r_bits),
		EGL_GREEN_SIZE:      int(g_bits),
		EGL_BLUE_SIZE:       int(b_bits),
		EGL_ALPHA_SIZE:      int(a_bits),
		EGL_SURFACE_TYPE:    int(surface_type),
		EGL_RENDERABLE_TYPE: int(renderable_type),
	}); err != nil {
		return nil, err
	} else if len(configs) == 0 {
		return nil, EGL_BAD_CONFIG
	} else {
		return configs[0], nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// API

func EGLBindAPI(api EGLAPI) error {
	if success := C.eglBindAPI(C.EGLenum(api)); success != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

func EGLQueryAPI() (EGLAPI, error) {
	if api := EGLAPI(C.eglQueryAPI()); api == 0 {
		return EGL_API_NONE, EGLGetError()
	} else {
		return api, nil
	}
}


////////////////////////////////////////////////////////////////////////////////
// CONTEXT

func EGLCreateContext(display EGLDisplay, config EGLConfig, share_context EGLContext) (EGLContext, error) {
	if context := EGLContext(C.eglCreateContext(C.EGLDisplay(display), C.EGLConfig(config), C.EGLContext(share_context), nil)); context == nil {
		return nil, EGLGetError()
	} else {
		return context, nil
	}
}

func EGL_DestroyContext(display EGLDisplay, context EGLContext) error {
	if C.eglDestroyContext(C.EGLDisplay(display), C.EGLContext(context)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

func EGLMakeCurrent(display EGLDisplay, draw, read EGLSurface, context EGLContext) error {
	if C.eglMakeCurrent(C.EGLDisplay(display), C.EGLSurface(draw), C.EGLSurface(read), C.EGLContext(context)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}


////////////////////////////////////////////////////////////////////////////////
// SURFACE

func EGLCreateSurface(display EGLDisplay, config EGLConfig, window EGLNativeWindow) (EGLSurface, error) {
	if surface := EGLSurface(C.eglCreateWindowSurface(C.EGLDisplay(display), C.EGLConfig(config), C.EGLNativeWindowType(window), nil)); surface == nil {
		return nil, EGLGetError()
	} else {
		return surface, nil
	}
}

func EGLDestroySurface(display EGLDisplay, surface EGLSurface) error {
	if C.eglDestroySurface(C.EGLDisplay(display), C.EGLSurface(surface)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

func EGLSwapBuffers(display EGLDisplay, surface EGLSurface) error {
	if C.eglSwapBuffers(C.EGLDisplay(display), C.EGLSurface(surface)) != EGL_TRUE {
		return EGLGetError()
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e EGLError) Error() string {
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
		return "[?? Unknown EGL_Error value]"
	}
}

func (a EGLConfigAttrib) Error() string {
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
		return "EGL_BIND_TO_TEXTURE_RGBA"
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
		return "[?? Invalid EGL_ConfigAttrib value]"
	}
}

func (a EGLAPI) String() string {
	switch a {
	case EGL_API_OPENGL_ES:
		return "EGL_API_OPENGL_ES"
	case EGL_API_OPENGL:
		return "EGL_API_OPENGL"
	case EGL_API_OPENVG:
		return "EGL_API_OPENVG"
	default:
		return "[?? Invalid EGL_API value]"
	}
}

func (f EGLRenderableFlag) String() string {
	parts := ""
	for flag := EGL_RENDERABLE_FLAG_MIN; flag <= EGL_RENDERABLE_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case EGL_RENDERABLE_FLAG_OPENGL_ES:
			parts += "|" + "EGL_RENDERABLE_FLAG_OPENGL_ES"
		case EGL_RENDERABLE_FLAG_OPENVG:
			parts += "|" + "EGL_RENDERABLE_FLAG_OPENVG"
		case EGL_RENDERABLE_FLAG_OPENGL_ES2:
			parts += "|" + "EGL_RENDERABLE_FLAG_OPENGL_ES2"
		case EGL_RENDERABLE_FLAG_OPENGL:
			parts += "|" + "EGL_RENDERABLE_FLAG_OPENGL"
		default:
			parts += "|" + "[?? Invalid EGL_RenderableFlag value]"
		}
	}
	return strings.Trim(parts, "|")
}

func (f EGLSurfaceTypeFlag) String() string {
	parts := ""
	for flag := EGL_SURFACETYPE_FLAG_MIN; flag <= EGL_SURFACETYPE_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case EGL_SURFACETYPE_FLAG_PBUFFER:
			parts += "|" + "EGL_SURFACETYPE_FLAG_PBUFFER"
		case EGL_SURFACETYPE_FLAG_PIXMAP:
			parts += "|" + "EGL_SURFACETYPE_FLAG_PIXMAP"
		case EGL_SURFACETYPE_FLAG_WINDOW:
			parts += "|" + "EGL_SURFACETYPE_FLAG_WINDOW"
		case EGL_SURFACETYPE_FLAG_VG_COLORSPACE_LINEAR:
			parts += "|" + "EGL_SURFACETYPE_FLAG_VG_COLORSPACE_LINEAR"
		case EGL_SURFACETYPE_FLAG_VG_ALPHA_FORMAT_PRE:
			parts += "|" + "EGL_SURFACETYPE_FLAG_VG_ALPHA_FORMAT_PRE"
		case EGL_SURFACETYPE_FLAG_MULTISAMPLE_RESOLVE_BOX:
			parts += "|" + "EGL_SURFACETYPE_FLAG_MULTISAMPLE_RESOLVE_BOX"
		case EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED:
			parts += "|" + "EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED"
		default:
			parts += "|" + "[?? Invalid EGL_SurfaceTypeFlag value]"
		}
	}
	return strings.Trim(parts, "|")
}
