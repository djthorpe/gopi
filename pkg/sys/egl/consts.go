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
	EGL_ALPHA_SIZE                                 EGLConfigAttrib = C.EGL_ALPHA_SIZE
	EGL_BLUE_SIZE                                  EGLConfigAttrib = C.EGL_BLUE_SIZE
	EGL_BUFFER_SIZE                                EGLConfigAttrib = C.EGL_BUFFER_SIZE
	EGL_CONFIG_CAVEAT                              EGLConfigAttrib = C.EGL_CONFIG_CAVEAT
	EGL_CONFIG_ID                                  EGLConfigAttrib = C.EGL_CONFIG_ID
	EGL_CORE_NATIVE_ENGINE                         EGLConfigAttrib = C.EGL_CORE_NATIVE_ENGINE
	EGL_DEPTH_SIZE                                 EGLConfigAttrib = C.EGL_DEPTH_SIZE
	EGL_DONT_CARE                                  EGLConfigAttrib = C.EGL_DONT_CARE
	EGL_DRAW                                       EGLConfigAttrib = C.EGL_DRAW
	EGL_EXTENSIONS                                 EGLConfigAttrib = C.EGL_EXTENSIONS
	EGL_GREEN_SIZE                                 EGLConfigAttrib = C.EGL_GREEN_SIZE
	EGL_HEIGHT                                     EGLConfigAttrib = C.EGL_HEIGHT
	EGL_LARGEST_PBUFFER                            EGLConfigAttrib = C.EGL_LARGEST_PBUFFER
	EGL_LEVEL                                      EGLConfigAttrib = C.EGL_LEVEL
	EGL_MAX_PBUFFER_HEIGHT                         EGLConfigAttrib = C.EGL_MAX_PBUFFER_HEIGHT
	EGL_MAX_PBUFFER_PIXELS                         EGLConfigAttrib = C.EGL_MAX_PBUFFER_PIXELS
	EGL_MAX_PBUFFER_WIDTH                          EGLConfigAttrib = C.EGL_MAX_PBUFFER_WIDTH
	EGL_NATIVE_RENDERABLE                          EGLConfigAttrib = C.EGL_NATIVE_RENDERABLE
	EGL_NATIVE_VISUAL_ID                           EGLConfigAttrib = C.EGL_NATIVE_VISUAL_ID
	EGL_NATIVE_VISUAL_TYPE                         EGLConfigAttrib = C.EGL_NATIVE_VISUAL_TYPE
	EGL_NONE                                       EGLConfigAttrib = C.EGL_NONE
	EGL_NON_CONFORMANT_CONFIG                      EGLConfigAttrib = C.EGL_NON_CONFORMANT_CONFIG
	EGL_PBUFFER_BIT                                EGLConfigAttrib = C.EGL_PBUFFER_BIT
	EGL_PIXMAP_BIT                                 EGLConfigAttrib = C.EGL_PIXMAP_BIT
	EGL_READ                                       EGLConfigAttrib = C.EGL_READ
	EGL_RED_SIZE                                   EGLConfigAttrib = C.EGL_RED_SIZE
	EGL_SAMPLES                                    EGLConfigAttrib = C.EGL_SAMPLES
	EGL_SAMPLE_BUFFERS                             EGLConfigAttrib = C.EGL_SAMPLE_BUFFERS
	EGL_SLOW_CONFIG                                EGLConfigAttrib = C.EGL_SLOW_CONFIG
	EGL_STENCIL_SIZE                               EGLConfigAttrib = C.EGL_STENCIL_SIZE
	EGL_SURFACE_TYPE                               EGLConfigAttrib = C.EGL_SURFACE_TYPE
	EGL_TRANSPARENT_BLUE_VALUE                     EGLConfigAttrib = C.EGL_TRANSPARENT_BLUE_VALUE
	EGL_TRANSPARENT_GREEN_VALUE                    EGLConfigAttrib = C.EGL_TRANSPARENT_GREEN_VALUE
	EGL_TRANSPARENT_RED_VALUE                      EGLConfigAttrib = C.EGL_TRANSPARENT_RED_VALUE
	EGL_TRANSPARENT_RGB                            EGLConfigAttrib = C.EGL_TRANSPARENT_RGB
	EGL_TRANSPARENT_TYPE                           EGLConfigAttrib = C.EGL_TRANSPARENT_TYPE
	EGL_VENDOR                                     EGLConfigAttrib = C.EGL_VENDOR
	EGL_VERSION                                    EGLConfigAttrib = C.EGL_VERSION
	EGL_WIDTH                                      EGLConfigAttrib = C.EGL_WIDTH
	EGL_WINDOW_BIT                                 EGLConfigAttrib = C.EGL_WINDOW_BIT
	EGL_BACK_BUFFER                                EGLConfigAttrib = C.EGL_BACK_BUFFER
	EGL_BIND_TO_TEXTURE_RGB                        EGLConfigAttrib = C.EGL_BIND_TO_TEXTURE_RGB
	EGL_BIND_TO_TEXTURE_RGBA                       EGLConfigAttrib = C.EGL_BIND_TO_TEXTURE_RGBA
	EGL_MIN_SWAP_INTERVAL                          EGLConfigAttrib = C.EGL_MIN_SWAP_INTERVAL
	EGL_MAX_SWAP_INTERVAL                          EGLConfigAttrib = C.EGL_MAX_SWAP_INTERVAL
	EGL_MIPMAP_TEXTURE                             EGLConfigAttrib = C.EGL_MIPMAP_TEXTURE
	EGL_MIPMAP_LEVEL                               EGLConfigAttrib = C.EGL_MIPMAP_LEVEL
	EGL_NO_TEXTURE                                 EGLConfigAttrib = C.EGL_NO_TEXTURE
	EGL_TEXTURE_2D                                 EGLConfigAttrib = C.EGL_TEXTURE_2D
	EGL_TEXTURE_FORMAT                             EGLConfigAttrib = C.EGL_TEXTURE_FORMAT
	EGL_TEXTURE_RGB                                EGLConfigAttrib = C.EGL_TEXTURE_RGB
	EGL_TEXTURE_RGBA                               EGLConfigAttrib = C.EGL_TEXTURE_RGBA
	EGL_TEXTURE_TARGET                             EGLConfigAttrib = C.EGL_TEXTURE_TARGET
	EGL_ALPHA_FORMAT                               EGLConfigAttrib = C.EGL_ALPHA_FORMAT
	EGL_ALPHA_FORMAT_NONPRE                        EGLConfigAttrib = C.EGL_ALPHA_FORMAT_NONPRE
	EGL_ALPHA_FORMAT_PRE                           EGLConfigAttrib = C.EGL_ALPHA_FORMAT_PRE
	EGL_ALPHA_MASK_SIZE                            EGLConfigAttrib = C.EGL_ALPHA_MASK_SIZE
	EGL_BUFFER_PRESERVED                           EGLConfigAttrib = C.EGL_BUFFER_PRESERVED
	EGL_BUFFER_DESTROYED                           EGLConfigAttrib = C.EGL_BUFFER_DESTROYED
	EGL_CLIENT_APIS                                EGLConfigAttrib = C.EGL_CLIENT_APIS
	EGL_COLORSPACE                                 EGLConfigAttrib = C.EGL_COLORSPACE
	EGL_COLORSPACE_sRGB                            EGLConfigAttrib = C.EGL_COLORSPACE_sRGB
	EGL_COLORSPACE_LINEAR                          EGLConfigAttrib = C.EGL_COLORSPACE_LINEAR
	EGL_COLOR_BUFFER_TYPE                          EGLConfigAttrib = C.EGL_COLOR_BUFFER_TYPE
	EGL_CONTEXT_CLIENT_TYPE                        EGLConfigAttrib = C.EGL_CONTEXT_CLIENT_TYPE
	EGL_DISPLAY_SCALING                            EGLConfigAttrib = C.EGL_DISPLAY_SCALING
	EGL_HORIZONTAL_RESOLUTION                      EGLConfigAttrib = C.EGL_HORIZONTAL_RESOLUTION
	EGL_LUMINANCE_BUFFER                           EGLConfigAttrib = C.EGL_LUMINANCE_BUFFER
	EGL_LUMINANCE_SIZE                             EGLConfigAttrib = C.EGL_LUMINANCE_SIZE
	EGL_OPENGL_ES_BIT                              EGLConfigAttrib = C.EGL_OPENGL_ES_BIT
	EGL_OPENVG_BIT                                 EGLConfigAttrib = C.EGL_OPENVG_BIT
	EGL_OPENGL_ES_API                              EGLConfigAttrib = C.EGL_OPENGL_ES_API
	EGL_OPENVG_API                                 EGLConfigAttrib = C.EGL_OPENVG_API
	EGL_OPENVG_IMAGE                               EGLConfigAttrib = C.EGL_OPENVG_IMAGE
	EGL_PIXEL_ASPECT_RATIO                         EGLConfigAttrib = C.EGL_PIXEL_ASPECT_RATIO
	EGL_RENDERABLE_TYPE                            EGLConfigAttrib = C.EGL_RENDERABLE_TYPE
	EGL_RENDER_BUFFER                              EGLConfigAttrib = C.EGL_RENDER_BUFFER
	EGL_RGB_BUFFER                                 EGLConfigAttrib = C.EGL_RGB_BUFFER
	EGL_SINGLE_BUFFER                              EGLConfigAttrib = C.EGL_SINGLE_BUFFER
	EGL_SWAP_BEHAVIOR                              EGLConfigAttrib = C.EGL_SWAP_BEHAVIOR
	EGL_UNKNOWN                                    EGLConfigAttrib = C.EGL_UNKNOWN
	EGL_VERTICAL_RESOLUTION                        EGLConfigAttrib = C.EGL_VERTICAL_RESOLUTION
	EGL_CONFORMANT                                 EGLConfigAttrib = C.EGL_CONFORMANT
	EGL_CONTEXT_CLIENT_VERSION                     EGLConfigAttrib = C.EGL_CONTEXT_CLIENT_VERSION
	EGL_MATCH_NATIVE_PIXMAP                        EGLConfigAttrib = C.EGL_MATCH_NATIVE_PIXMAP
	EGL_OPENGL_ES2_BIT                             EGLConfigAttrib = C.EGL_OPENGL_ES2_BIT
	EGL_VG_ALPHA_FORMAT                            EGLConfigAttrib = C.EGL_VG_ALPHA_FORMAT
	EGL_VG_ALPHA_FORMAT_NONPRE                     EGLConfigAttrib = C.EGL_VG_ALPHA_FORMAT_NONPRE
	EGL_VG_ALPHA_FORMAT_PRE                        EGLConfigAttrib = C.EGL_VG_ALPHA_FORMAT_PRE
	EGL_VG_ALPHA_FORMAT_PRE_BIT                    EGLConfigAttrib = C.EGL_VG_ALPHA_FORMAT_PRE_BIT
	EGL_VG_COLORSPACE                              EGLConfigAttrib = C.EGL_VG_COLORSPACE
	EGL_VG_COLORSPACE_sRGB                         EGLConfigAttrib = C.EGL_VG_COLORSPACE_sRGB
	EGL_VG_COLORSPACE_LINEAR                       EGLConfigAttrib = C.EGL_VG_COLORSPACE_LINEAR
	EGL_VG_COLORSPACE_LINEAR_BIT                   EGLConfigAttrib = C.EGL_VG_COLORSPACE_LINEAR_BIT
	EGL_MULTISAMPLE_RESOLVE_BOX_BIT                EGLConfigAttrib = C.EGL_MULTISAMPLE_RESOLVE_BOX_BIT
	EGL_MULTISAMPLE_RESOLVE                        EGLConfigAttrib = C.EGL_MULTISAMPLE_RESOLVE
	EGL_MULTISAMPLE_RESOLVE_DEFAULT                EGLConfigAttrib = C.EGL_MULTISAMPLE_RESOLVE_DEFAULT
	EGL_MULTISAMPLE_RESOLVE_BOX                    EGLConfigAttrib = C.EGL_MULTISAMPLE_RESOLVE_BOX
	EGL_OPENGL_API                                 EGLConfigAttrib = C.EGL_OPENGL_API
	EGL_OPENGL_BIT                                 EGLConfigAttrib = C.EGL_OPENGL_BIT
	EGL_SWAP_BEHAVIOR_PRESERVED_BIT                EGLConfigAttrib = C.EGL_SWAP_BEHAVIOR_PRESERVED_BIT
	EGL_CONTEXT_MAJOR_VERSION                      EGLConfigAttrib = C.EGL_CONTEXT_MAJOR_VERSION
	EGL_CONTEXT_MINOR_VERSION                      EGLConfigAttrib = C.EGL_CONTEXT_MINOR_VERSION
	EGL_CONTEXT_OPENGL_PROFILE_MASK                EGLConfigAttrib = C.EGL_CONTEXT_OPENGL_PROFILE_MASK
	EGL_CONTEXT_OPENGL_RESET_NOTIFICATION_STRATEGY EGLConfigAttrib = C.EGL_CONTEXT_OPENGL_RESET_NOTIFICATION_STRATEGY
	EGL_NO_RESET_NOTIFICATION                      EGLConfigAttrib = C.EGL_NO_RESET_NOTIFICATION
	EGL_LOSE_CONTEXT_ON_RESET                      EGLConfigAttrib = C.EGL_LOSE_CONTEXT_ON_RESET
	EGL_CONTEXT_OPENGL_CORE_PROFILE_BIT            EGLConfigAttrib = C.EGL_CONTEXT_OPENGL_CORE_PROFILE_BIT
	EGL_CONTEXT_OPENGL_COMPATIBILITY_PROFILE_BIT   EGLConfigAttrib = C.EGL_CONTEXT_OPENGL_COMPATIBILITY_PROFILE_BIT
	EGL_CONTEXT_OPENGL_DEBUG                       EGLConfigAttrib = C.EGL_CONTEXT_OPENGL_DEBUG
	EGL_CONTEXT_OPENGL_FORWARD_COMPATIBLE          EGLConfigAttrib = C.EGL_CONTEXT_OPENGL_FORWARD_COMPATIBLE
	EGL_CONTEXT_OPENGL_ROBUST_ACCESS               EGLConfigAttrib = C.EGL_CONTEXT_OPENGL_ROBUST_ACCESS
	EGL_OPENGL_ES3_BIT                             EGLConfigAttrib = C.EGL_OPENGL_ES3_BIT
	EGL_CL_EVENT_HANDLE                            EGLConfigAttrib = C.EGL_CL_EVENT_HANDLE
	EGL_SYNC_CL_EVENT                              EGLConfigAttrib = C.EGL_SYNC_CL_EVENT
	EGL_SYNC_CL_EVENT_COMPLETE                     EGLConfigAttrib = C.EGL_SYNC_CL_EVENT_COMPLETE
	EGL_SYNC_PRIOR_COMMANDS_COMPLETE               EGLConfigAttrib = C.EGL_SYNC_PRIOR_COMMANDS_COMPLETE
	EGL_SYNC_TYPE                                  EGLConfigAttrib = C.EGL_SYNC_TYPE
	EGL_SYNC_STATUS                                EGLConfigAttrib = C.EGL_SYNC_STATUS
	EGL_SYNC_CONDITION                             EGLConfigAttrib = C.EGL_SYNC_CONDITION
	EGL_SIGNALED                                   EGLConfigAttrib = C.EGL_SIGNALED
	EGL_UNSIGNALED                                 EGLConfigAttrib = C.EGL_UNSIGNALED
	EGL_SYNC_FLUSH_COMMANDS_BIT                    EGLConfigAttrib = C.EGL_SYNC_FLUSH_COMMANDS_BIT
	EGL_TIMEOUT_EXPIRED                            EGLConfigAttrib = C.EGL_TIMEOUT_EXPIRED
	EGL_CONDITION_SATISFIED                        EGLConfigAttrib = C.EGL_CONDITION_SATISFIED
	EGL_SYNC_FENCE                                 EGLConfigAttrib = C.EGL_SYNC_FENCE
	EGL_GL_COLORSPACE                              EGLConfigAttrib = C.EGL_GL_COLORSPACE
	EGL_GL_COLORSPACE_SRGB                         EGLConfigAttrib = C.EGL_GL_COLORSPACE_SRGB
	EGL_GL_COLORSPACE_LINEAR                       EGLConfigAttrib = C.EGL_GL_COLORSPACE_LINEAR
	EGL_GL_RENDERBUFFER                            EGLConfigAttrib = C.EGL_GL_RENDERBUFFER
	EGL_GL_TEXTURE_2D                              EGLConfigAttrib = C.EGL_GL_TEXTURE_2D
	EGL_GL_TEXTURE_LEVEL                           EGLConfigAttrib = C.EGL_GL_TEXTURE_LEVEL
	EGL_GL_TEXTURE_3D                              EGLConfigAttrib = C.EGL_GL_TEXTURE_3D
	EGL_GL_TEXTURE_ZOFFSET                         EGLConfigAttrib = C.EGL_GL_TEXTURE_ZOFFSET
	EGL_GL_TEXTURE_CUBE_MAP_POSITIVE_X             EGLConfigAttrib = C.EGL_GL_TEXTURE_CUBE_MAP_POSITIVE_X
	EGL_GL_TEXTURE_CUBE_MAP_NEGATIVE_X             EGLConfigAttrib = C.EGL_GL_TEXTURE_CUBE_MAP_NEGATIVE_X
	EGL_GL_TEXTURE_CUBE_MAP_POSITIVE_Y             EGLConfigAttrib = C.EGL_GL_TEXTURE_CUBE_MAP_POSITIVE_Y
	EGL_GL_TEXTURE_CUBE_MAP_NEGATIVE_Y             EGLConfigAttrib = C.EGL_GL_TEXTURE_CUBE_MAP_NEGATIVE_Y
	EGL_GL_TEXTURE_CUBE_MAP_POSITIVE_Z             EGLConfigAttrib = C.EGL_GL_TEXTURE_CUBE_MAP_POSITIVE_Z
	EGL_GL_TEXTURE_CUBE_MAP_NEGATIVE_Z             EGLConfigAttrib = C.EGL_GL_TEXTURE_CUBE_MAP_NEGATIVE_Z
	EGL_IMAGE_PRESERVED                            EGLConfigAttrib = C.EGL_IMAGE_PRESERVED

	EGL_CONFIG_ATTRIB_MIN EGLConfigAttrib = 0x0000
	EGL_CONFIG_ATTRIB_MAX EGLConfigAttrib = 0xFFFF
	//	EGL_FOREVER                                    EGLConfigAttrib = C.EGL_FOREVER
	// 	EGL_NO_SYNC                                    EGLConfigAttrib = C.EGL_NO_SYNC
)

const (
	// EGLRenderableFlag
	EGL_RENDERABLE_FLAG_OPENGL_ES  EGLRenderableFlag = 0x0001
	EGL_RENDERABLE_FLAG_OPENVG     EGLRenderableFlag = 0x0002
	EGL_RENDERABLE_FLAG_OPENGL_ES2 EGLRenderableFlag = 0x0004
	EGL_RENDERABLE_FLAG_OPENGL     EGLRenderableFlag = 0x0008
	EGL_RENDERABLE_FLAG_OPENGL_ES3 EGLRenderableFlag = C.EGL_OPENGL_ES3_BIT
	EGL_RENDERABLE_FLAG_MIN                          = EGL_RENDERABLE_FLAG_OPENGL_ES
	EGL_RENDERABLE_FLAG_MAX                          = EGL_RENDERABLE_FLAG_OPENGL_ES3
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
