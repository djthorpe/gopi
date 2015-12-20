/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package egl

/*
	#cgo CFLAGS: -I/opt/vc/include
	#include "EGL/egl.h"
*/
import "C"

const (
	EGL_FALSE = 0
	EGL_TRUE  = 1
)

const (
	/* Config attributes */
	EGL_BUFFER_SIZE             = 0x3020
	EGL_ALPHA_SIZE              = 0x3021
	EGL_BLUE_SIZE               = 0x3022
	EGL_GREEN_SIZE              = 0x3023
	EGL_RED_SIZE                = 0x3024
	EGL_DEPTH_SIZE              = 0x3025
	EGL_STENCIL_SIZE            = 0x3026
	EGL_CONFIG_CAVEAT           = 0x3027
	EGL_CONFIG_ID               = 0x3028
	EGL_LEVEL                   = 0x3029
	EGL_MAX_PBUFFER_HEIGHT      = 0x302A
	EGL_MAX_PBUFFER_PIXELS      = 0x302B
	EGL_MAX_PBUFFER_WIDTH       = 0x302C
	EGL_NATIVE_RENDERABLE       = 0x302D
	EGL_NATIVE_VISUAL_ID        = 0x302E
	EGL_NATIVE_VISUAL_TYPE      = 0x302F
	EGL_SAMPLES                 = 0x3031
	EGL_SAMPLE_BUFFERS          = 0x3032
	EGL_SURFACE_TYPE            = 0x3033
	EGL_TRANSPARENT_TYPE        = 0x3034
	EGL_TRANSPARENT_BLUE_VALUE  = 0x3035
	EGL_TRANSPARENT_GREEN_VALUE = 0x3036
	EGL_TRANSPARENT_RED_VALUE   = 0x3037
	EGL_NONE                    = 0x3038 /* Attrib list terminator */
	EGL_BIND_TO_TEXTURE_RGB     = 0x3039
	EGL_BIND_TO_TEXTURE_RGBA    = 0x303A
	EGL_MIN_SWAP_INTERVAL       = 0x303B
	EGL_MAX_SWAP_INTERVAL       = 0x303C
	EGL_LUMINANCE_SIZE          = 0x303D
	EGL_ALPHA_MASK_SIZE         = 0x303E
	EGL_COLOR_BUFFER_TYPE       = 0x303F
	EGL_RENDERABLE_TYPE         = 0x3040
	EGL_MATCH_NATIVE_PIXMAP     = 0x3041 /* Pseudo-attribute (not queryable) */
	EGL_CONFORMANT              = 0x3042

	/* Reserved 0x3041-0x304F for additional config attributes */

	/* Config attribute values */
	EGL_SLOW_CONFIG           = 0x3050 /* EGL_CONFIG_CAVEAT value */
	EGL_NON_CONFORMANT_CONFIG = 0x3051 /* EGL_CONFIG_CAVEAT value */
	EGL_TRANSPARENT_RGB       = 0x3052 /* EGL_TRANSPARENT_TYPE value */
	EGL_RGB_BUFFER            = 0x308E /* EGL_COLOR_BUFFER_TYPE value */
	EGL_LUMINANCE_BUFFER      = 0x308F /* EGL_COLOR_BUFFER_TYPE value */

	/* More config attribute values, for EGL_TEXTURE_FORMAT */
	EGL_NO_TEXTURE   = 0x305C
	EGL_TEXTURE_RGB  = 0x305D
	EGL_TEXTURE_RGBA = 0x305E
	EGL_TEXTURE_2D   = 0x305F

	/* Config attribute mask bits */
	EGL_PBUFFER_BIT                 = 0x0001 /* EGL_SURFACE_TYPE mask bits */
	EGL_PIXMAP_BIT                  = 0x0002 /* EGL_SURFACE_TYPE mask bits */
	EGL_WINDOW_BIT                  = 0x0004 /* EGL_SURFACE_TYPE mask bits */
	EGL_VG_COLORSPACE_LINEAR_BIT    = 0x0020 /* EGL_SURFACE_TYPE mask bits */
	EGL_VG_ALPHA_FORMAT_PRE_BIT     = 0x0040 /* EGL_SURFACE_TYPE mask bits */
	EGL_MULTISAMPLE_RESOLVE_BOX_BIT = 0x0200 /* EGL_SURFACE_TYPE mask bits */
	EGL_SWAP_BEHAVIOR_PRESERVED_BIT = 0x0400 /* EGL_SURFACE_TYPE mask bits */

	EGL_OPENGL_ES_BIT  = 0x0001 /* EGL_RENDERABLE_TYPE mask bits */
	EGL_OPENVG_BIT     = 0x0002 /* EGL_RENDERABLE_TYPE mask bits */
	EGL_OPENGL_ES2_BIT = 0x0004 /* EGL_RENDERABLE_TYPE mask bits */
	EGL_OPENGL_BIT     = 0x0008 /* EGL_RENDERABLE_TYPE mask bits */
)
