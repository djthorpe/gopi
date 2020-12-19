// +build egl

package egl

import (
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: egl
#include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EGLQuery           C.EGLint
	EGLConfigAttrib    C.EGLint
	EGLRenderableFlag  C.EGLint
	EGLSurfaceTypeFlag C.EGLint
	EGLAPI             C.EGLint
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EGL_FALSE = C.EGLBoolean(0)
	EGL_TRUE  = C.EGLBoolean(1)
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
	EGL_COMFIG_ATTRIB_MIN                       = EGL_BUFFER_SIZE
	EGL_COMFIG_ATTRIB_MAX                       = EGL_CONFORMANT
)

const (
	// EGLRenderableFlag
	EGL_RENDERABLE_FLAG_OPENGL_ES  EGLRenderableFlag = 0x0001
	EGL_RENDERABLE_FLAG_OPENVG     EGLRenderableFlag = 0x0002
	EGL_RENDERABLE_FLAG_OPENGL_ES2 EGLRenderableFlag = 0x0004
	EGL_RENDERABLE_FLAG_OPENGL     EGLRenderableFlag = 0x0008
	EGL_RENDERABLE_FLAG_MIN                          = EGL_RENDERABLE_FLAG_OPENGL_ES
	EGL_RENDERABLE_FLAG_MAX                          = EGL_RENDERABLE_FLAG_OPENGL
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
	EGL_SURFACETYPE_FLAG_MIN                                        = EGL_SURFACETYPE_FLAG_PBUFFER
	EGL_SURFACETYPE_FLAG_MAX                                        = EGL_SURFACETYPE_FLAG_SWAP_BEHAVIOR_PRESERVED
)

const (
	// EGLAPI
	EGL_API_NONE      EGLAPI = 0
	EGL_API_OPENGL_ES EGLAPI = 0x30A0
	EGL_API_OPENVG    EGLAPI = 0x30A1
	EGL_API_OPENGL    EGLAPI = 0x30A2
	EGL_API_MIN              = EGL_API_OPENGL_ES
	EGL_API_MAX              = EGL_API_OPENGL
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (a EGLConfigAttrib) String() string {
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
