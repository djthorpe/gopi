// +build rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"unsafe"
	"fmt"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: bcm_host
#include <bcm_host.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	DXDisplayId uint16
	DXDisplayHandle C.DISPMANX_DISPLAY_HANDLE_T
	DXInputFormat   C.uint32_t
	DXTransform     C.int
	DXRect          (*C.VC_RECT_T)
	DXResource      C.DISPMANX_RESOURCE_HANDLE_T
	DXImageType     C.VC_IMAGE_TYPE_T
	DXElement       C.DISPMANX_ELEMENT_HANDLE_T
	DXUpdate        C.DISPMANX_UPDATE_HANDLE_T
	DXProtection    C.uint32_t
	DXAlphaFlags    C.uint32_t
	DXClampMode     C.int
	DXChangeFlags   C.int
)

type DXDisplayModeInfo struct {
	Size        DXSize
	Transform   DXTransform
	InputFormat DXInputFormat
	extra       C.uint32_t
}

type DXPoint struct {
	X int32
	Y int32
}

type DXSize struct {
	W uint32
	H uint32
}

type DXClamp struct {
	Mode    DXClampMode
	Flags   int
	Opacity uint32
	Mask    DXResource
}

type DXAlpha struct {
	Flags   DXAlphaFlags
	Opacity uint32
	Mask    DXResource
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// DX_DisplayId values
	DX_DISPLAYID_MAIN_LCD DXDisplayId = iota
	DX_DISPLAYID_AUX_LCD
	DX_DISPLAYID_HDMI
	DX_DISPLAYID_SDTV
	DX_DISPLAYID_FORCE_LCD
	DX_DISPLAYID_FORCE_TV
	DX_DISPLAYID_FORCE_OTHER
	DX_DISPLAYID_MAX = DX_DISPLAYID_FORCE_OTHER
	DX_DISPLAYID_MIN = DX_DISPLAYID_MAIN_LCD
)

const (
	/* Success and failure conditions */
	DX_SUCCESS   = 0
	DX_INVALID   = -1
	DX_NO_HANDLE = 0
)

const (
	// DX_Transform values
	DX_TRANSFORM_NONE DXTransform = iota
	DX_TRANSFORM_ROTATE_90
	DX_TRANSFORM_ROTATE_180
	DX_TRANSFORM_ROTATE_270
	DX_TRANSFORM_MAX = DX_TRANSFORM_ROTATE_270
)

const (
	// DX_ChangeFlags values
	DX_CHANGE_FLAG_LAYER     DXChangeFlags = (1 << 0)
	DX_CHANGE_FLAG_OPACITY   DXChangeFlags = (1 << 1)
	DX_CHANGE_FLAG_DEST_RECT DXChangeFlags = (1 << 2)
	DX_CHANGE_FLAG_SRC_RECT  DXChangeFlags = (1 << 3)
	DX_CHANGE_FLAG_MASK      DXChangeFlags = (1 << 4)
	DX_CHANGE_FLAG_TRANSFORM DXChangeFlags = (1 << 5)
	DX_CHANGE_FLAG_MIN                      = DX_CHANGE_FLAG_LAYER
	DX_CHANGE_FLAG_MAX                      = DX_CHANGE_FLAG_TRANSFORM
)

const (
	// DX_InputFormat values
	DX_INPUT_FORMAT_INVALID DXInputFormat = iota
	DX_INPUT_FORMAT_RGB888
	DX_INPUT_FORMAT_RGB565
	DX_INPUT_FORMAT_MAX = DX_INPUT_FORMAT_RGB565
)

const (
	DX_IMAGE_TYPE_NONE DXImageType = iota
	DX_IMAGE_TYPE_RGB565
	DX_IMAGE_TYPE_1BPP
	DX_IMAGE_TYPE_YUV420
	DX_IMAGE_TYPE_48BPP
	DX_IMAGE_TYPE_RGB888
	DX_IMAGE_TYPE_8BPP
	DX_IMAGE_TYPE_4BPP          // 4bpp palettised image
	DX_IMAGE_TYPE_3D32          // A separated format of 16 colour/light shorts followed by 16 z values
	DX_IMAGE_TYPE_3D32B         // 16 colours followed by 16 z values
	DX_IMAGE_TYPE_3D32MAT       // A separated format of 16 material/colour/light shorts followed by 16 z values
	DX_IMAGE_TYPE_RGB2X9        // 32 bit format containing 18 bits of 6.6.6 RGB 9 bits per short
	DX_IMAGE_TYPE_RGB666        // 32-bit format holding 18 bits of 6.6.6 RGB
	DX_IMAGE_TYPE_PAL4_OBSOLETE // 4bpp palettised image with embedded palette
	DX_IMAGE_TYPE_PAL8_OBSOLETE // 8bpp palettised image with embedded palette
	DX_IMAGE_TYPE_RGBA32        // RGB888 with an alpha byte after each pixel */ /* xxx: isn't it BEFORE each pixel?
	DX_IMAGE_TYPE_YUV422        // a line of Y (32-byte padded) a line of U (16-byte padded) and a line of V (16-byte padded)
	DX_IMAGE_TYPE_RGBA565       // RGB565 with a transparent patch
	DX_IMAGE_TYPE_RGBA16        // Compressed (4444) version of RGBA32
	DX_IMAGE_TYPE_YUV_UV        // VCIII codec format
	DX_IMAGE_TYPE_TF_RGBA32     // VCIII T-format RGBA8888
	DX_IMAGE_TYPE_TF_RGBX32     // VCIII T-format RGBx8888
	DX_IMAGE_TYPE_TF_FLOAT      // VCIII T-format float
	DX_IMAGE_TYPE_TF_RGBA16     // VCIII T-format RGBA4444
	DX_IMAGE_TYPE_TF_RGBA5551   // VCIII T-format RGB5551
	DX_IMAGE_TYPE_TF_RGB565     // VCIII T-format RGB565
	DX_IMAGE_TYPE_TF_YA88       // VCIII T-format 8-bit luma and 8-bit alpha
	DX_IMAGE_TYPE_TF_BYTE       // VCIII T-format 8 bit generic sample
	DX_IMAGE_TYPE_TF_PAL8       // VCIII T-format 8-bit palette
	DX_IMAGE_TYPE_TF_PAL4       // VCIII T-format 4-bit palette
	DX_IMAGE_TYPE_TF_ETC1       // VCIII T-format Ericsson Texture Compressed
	DX_IMAGE_TYPE_BGR888        // RGB888 with R & B swapped
	DX_IMAGE_TYPE_BGR888_NP     // RGB888 with R & B swapped but with no pitch i.e. no padding after each row of pixels
	DX_IMAGE_TYPE_BAYER         // Bayer image extra defines which variant is being used
	DX_IMAGE_TYPE_CODEC         // General wrapper for codec images e.g. JPEG from camera
	DX_IMAGE_TYPE_YUV_UV32      // VCIII codec format
	DX_IMAGE_TYPE_TF_Y8         // VCIII T-format 8-bit luma
	DX_IMAGE_TYPE_TF_A8         // VCIII T-format 8-bit alpha
	DX_IMAGE_TYPE_TF_SHORT      // VCIII T-format 16-bit generic sample
	DX_IMAGE_TYPE_TF_1BPP       // VCIII T-format 1bpp black/white
	DX_IMAGE_TYPE_OPENGL
	DX_IMAGE_TYPE_YUV444I      // VCIII-B0 HVS YUV 4:4:4 interleaved samples
	DX_IMAGE_TYPE_YUV422PLANAR // Y U & V planes separately (DX_IMAGE_TYPE_YUV422 has them interleaved on a per line basis)
	DX_IMAGE_TYPE_ARGB8888     // 32bpp with 8bit alpha at MS byte with R G B (LS byte)
	DX_IMAGE_TYPE_XRGB8888     // 32bpp with 8bit unused at MS byte with R G B (LS byte)
	DX_IMAGE_TYPE_YUV422YUYV   // interleaved 8 bit samples of Y U Y V
	DX_IMAGE_TYPE_YUV422YVYU   // interleaved 8 bit samples of Y V Y U
	DX_IMAGE_TYPE_YUV422UYVY   // interleaved 8 bit samples of U Y V Y
	DX_IMAGE_TYPE_YUV422VYUY   // interleaved 8 bit samples of V Y U Y
	DX_IMAGE_TYPE_RGBX32       // 32bpp like RGBA32 but with unused alpha
	DX_IMAGE_TYPE_RGBX8888     // 32bpp corresponding to RGBA with unused alpha
	DX_IMAGE_TYPE_BGRX8888     // 32bpp corresponding to BGRA with unused alpha
	DX_IMAGE_TYPE_YUV420SP     // Y as a plane then UV byte interleaved in plane with with same pitch half height
	DX_IMAGE_TYPE_YUV444PLANAR // Y U & V planes separately 4:4:4
	DX_IMAGE_TYPE_TF_U8        // T-format 8-bit U - same as TF_Y8 buf from U plane
	DX_IMAGE_TYPE_TF_V8        // T-format 8-bit U - same as TF_Y8 buf from V plane
	DX_IMAGE_TYPE_YUV420_16    // YUV4:2:0 planar 16bit values
	DX_IMAGE_TYPE_YUV_UV_16    // YUV4:2:0 codec format 16bit values
	DX_IMAGE_TYPE_YUV420_S     // YUV4:2:0 with UV in side-by-side format
	DX_IMAGE_TYPE_MIN          = DX_IMAGE_TYPE_RGB565
	DX_IMAGE_TYPE_MAX          = DX_IMAGE_TYPE_YUV420_S
)

const (
	/* Protection values */
	DX_PROTECTION_NONE DXProtection = 0
	DX_PROTECTION_HDCP DXProtection = 11
)

const (
	/* Clamp flags */
	DX_CLAMP_MODE_NONE DXClampMode = iota
	DX_CLAMP_MODE_LUMA_TRANSPARENT
	DX_CLAMP_MODE_CHROMA_TRANSPARENT
	DX_CLAMP_MODE_REPLACE
)

const (
	/* Alpha flags */
	DX_ALPHA_FLAG_FROM_SOURCE DXAlphaFlags = iota // Bottom 2 bits sets the alpha mode
	DX_ALPHA_FLAG_FIXED_ALL_PIXELS
	DX_ALPHA_FLAG_FIXED_NON_ZERO
	DX_ALPHA_FLAG_FIXED_EXCEED_0X07
	DX_ALPHA_FLAG_PREMULT               DXAlphaFlags = 1 << 16
	DX_ALPHA_FLAG_MIX                   DXAlphaFlags = 1 << 17
	DX_ALPHA_FLAG__DISCARD_LOWER_LAYERS DXAlphaFlags = 1 << 18
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: DISPLAYS

func DXNumberOfDisplays() uint16 {
	return uint16(DX_DISPLAYID_MAX-DX_DISPLAYID_MIN) + 1
}

func DXInit() {
	C.bcm_host_init()
}

func DXStop() {
	C.vc_dispmanx_stop()
}

func DXDisplayOpen(display DXDisplayId) (DXDisplayHandle, error) {
	if handle := DXDisplayHandle(C.vc_dispmanx_display_open(C.uint32_t(display))); handle != 0 {
		return handle, nil
	} else {
		return 0, gopi.ErrBadParameter
	}
}

func DXDisplayClose(display DXDisplayHandle) error {
	if C.vc_dispmanx_display_close(C.DISPMANX_DISPLAY_HANDLE_T(display)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrUnexpectedResponse
	}
}

func DXDisplayGetInfo(display DXDisplayHandle) (DXDisplayModeInfo, error) {
	info := DXDisplayModeInfo{}
	if C.vc_dispmanx_display_get_info(C.DISPMANX_DISPLAY_HANDLE_T(display), (*C.DISPMANX_MODEINFO_T)(unsafe.Pointer(&info))) == DX_SUCCESS {
		return info, nil
	} else {
		return info, gopi.ErrUnexpectedResponse
	}
}

func DXDisplaySnapshot(display DXDisplayHandle, resource DXResource, transform DXTransform) error {
	if C.vc_dispmanx_snapshot(C.DISPMANX_DISPLAY_HANDLE_T(display), C.DISPMANX_RESOURCE_HANDLE_T(resource), C.DISPMANX_TRANSFORM_T(transform)) == DX_SUCCESS {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (h DXDisplayHandle) String() string {
	return "<DXDisplayHandle " + fmt.Sprint(C.DISPMANX_DISPLAY_HANDLE_T(h)) + ">"
}

func (d DXDisplayId) String() string {
	switch d {
	case DX_DISPLAYID_MAIN_LCD:
		return "DX_DISPLAYID_MAIN_LCD"
	case DX_DISPLAYID_AUX_LCD:
		return "DX_DISPLAYID_AUX_LCD"
	case DX_DISPLAYID_HDMI:
		return "DX_DISPLAYID_HDMI"
	case DX_DISPLAYID_SDTV:
		return "DX_DISPLAYID_SDTV"
	case DX_DISPLAYID_FORCE_LCD:
		return "DX_DISPLAYID_FORCE_LCD"
	case DX_DISPLAYID_FORCE_TV:
		return "DX_DISPLAYID_FORCE_TV"
	case DX_DISPLAYID_FORCE_OTHER:
		return "DX_DISPLAYID_FORCE_OTHER"
	default:
		return "[?? Invalid DXDisplayId value]"
	}
}

func (this DXDisplayModeInfo) String() string {
	return fmt.Sprintf("DXDisplayModeInfo{ size=%v transform=%v input_format=%v }", this.Size, this.Transform, this.InputFormat)
}

func (size DXSize) String() string {
	return fmt.Sprintf("DXSize{%v,%v}", size.W, size.H)
}

func (t DXTransform) String() string {
	switch t {
	case DX_TRANSFORM_NONE:
		return "DX_TRANSFORM_NONE"
	case DX_TRANSFORM_ROTATE_90:
		return "DX_TRANSFORM_ROTATE_90"
	case DX_TRANSFORM_ROTATE_180:
		return "DX_TRANSFORM_ROTATE_180"
	case DX_TRANSFORM_ROTATE_270:
		return "DX_TRANSFORM_ROTATE_270"
	default:
		return "[?? Invalid DX_Transform value]"
	}
}

func (f DXChangeFlags) String() string {
	parts := ""
	for flag := DX_CHANGE_FLAG_MIN; flag <= DX_CHANGE_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case DX_CHANGE_FLAG_LAYER:
			parts += "|" + "DX_CHANGE_FLAG_LAYER"
		case DX_CHANGE_FLAG_OPACITY:
			parts += "|" + "DX_CHANGE_FLAG_OPACITY"
		case DX_CHANGE_FLAG_DEST_RECT:
			parts += "|" + "DX_CHANGE_FLAG_DEST_RECT"
		case DX_CHANGE_FLAG_SRC_RECT:
			parts += "|" + "DX_CHANGE_FLAG_SRC_RECT"
		case DX_CHANGE_FLAG_MASK:
			parts += "|" + "DX_CHANGE_FLAG_MASK"
		case DX_CHANGE_FLAG_TRANSFORM:
			parts += "|" + "DX_CHANGE_FLAG_TRANSFORM"
		default:
			parts += "|" + "[?? Invalid DX_ChangeFlags value]"
		}
	}
	return strings.Trim(parts, "|")
}

func (f DXInputFormat) String() string {
	switch f {
	case DX_INPUT_FORMAT_RGB888:
		return "DX_INPUT_FORMAT_RGB888"
	case DX_INPUT_FORMAT_RGB565:
		return "DX_INPUT_FORMAT_RGB565"
	default:
		return "DX_INPUT_FORMAT_INVALID"
	}
}

func (t DXImageType) String() string {
	switch t {
	case DX_IMAGE_TYPE_NONE:
		return "DX_IMAGE_TYPE_NONE"
	case DX_IMAGE_TYPE_RGB565:
		return "DX_IMAGE_TYPE_RGB565"
	case DX_IMAGE_TYPE_1BPP:
		return "DX_IMAGE_TYPE_1BPP"
	case DX_IMAGE_TYPE_YUV420:
		return "DX_IMAGE_TYPE_YUV420"
	case DX_IMAGE_TYPE_48BPP:
		return "DX_IMAGE_TYPE_48BPP"
	case DX_IMAGE_TYPE_RGB888:
		return "DX_IMAGE_TYPE_RGB888"
	case DX_IMAGE_TYPE_8BPP:
		return "DX_IMAGE_TYPE_8BPP"
	case DX_IMAGE_TYPE_4BPP:
		return "DX_IMAGE_TYPE_4BPP"
	case DX_IMAGE_TYPE_3D32:
		return "DX_IMAGE_TYPE_3D32"
	case DX_IMAGE_TYPE_3D32B:
		return "DX_IMAGE_TYPE_3D32B"
	case DX_IMAGE_TYPE_3D32MAT:
		return "DX_IMAGE_TYPE_3D32MAT"
	case DX_IMAGE_TYPE_RGB2X9:
		return "DX_IMAGE_TYPE_RGB2X9"
	case DX_IMAGE_TYPE_RGB666:
		return "DX_IMAGE_TYPE_RGB666"
	case DX_IMAGE_TYPE_PAL4_OBSOLETE:
		return "DX_IMAGE_TYPE_PAL4_OBSOLETE"
	case DX_IMAGE_TYPE_PAL8_OBSOLETE:
		return "DX_IMAGE_TYPE_PAL8_OBSOLETE"
	case DX_IMAGE_TYPE_RGBA32:
		return "DX_IMAGE_TYPE_RGBA32"
	case DX_IMAGE_TYPE_YUV422:
		return "DX_IMAGE_TYPE_YUV422"
	case DX_IMAGE_TYPE_RGBA565:
		return "DX_IMAGE_TYPE_RGBA565"
	case DX_IMAGE_TYPE_RGBA16:
		return "DX_IMAGE_TYPE_RGBA16"
	case DX_IMAGE_TYPE_YUV_UV:
		return "DX_IMAGE_TYPE_YUV_UV"
	case DX_IMAGE_TYPE_TF_RGBA32:
		return "DX_IMAGE_TYPE_TF_RGBA32"
	case DX_IMAGE_TYPE_TF_RGBX32:
		return "DX_IMAGE_TYPE_TF_RGBX32"
	case DX_IMAGE_TYPE_TF_FLOAT:
		return "DX_IMAGE_TYPE_TF_FLOAT"
	case DX_IMAGE_TYPE_TF_RGBA16:
		return "DX_IMAGE_TYPE_TF_RGBA16"
	case DX_IMAGE_TYPE_TF_RGBA5551:
		return "DX_IMAGE_TYPE_TF_RGBA5551"
	case DX_IMAGE_TYPE_TF_RGB565:
		return "DX_IMAGE_TYPE_TF_RGB565"
	case DX_IMAGE_TYPE_TF_YA88:
		return "DX_IMAGE_TYPE_TF_YA88"
	case DX_IMAGE_TYPE_TF_BYTE:
		return "DX_IMAGE_TYPE_TF_BYTE"
	case DX_IMAGE_TYPE_TF_PAL8:
		return "DX_IMAGE_TYPE_TF_PAL8"
	case DX_IMAGE_TYPE_TF_PAL4:
		return "DX_IMAGE_TYPE_TF_PAL4"
	case DX_IMAGE_TYPE_TF_ETC1:
		return "DX_IMAGE_TYPE_TF_ETC1"
	case DX_IMAGE_TYPE_BGR888:
		return "DX_IMAGE_TYPE_BGR888"
	case DX_IMAGE_TYPE_BGR888_NP:
		return "DX_IMAGE_TYPE_BGR888_NP"
	case DX_IMAGE_TYPE_BAYER:
		return "DX_IMAGE_TYPE_BAYER"
	case DX_IMAGE_TYPE_CODEC:
		return "DX_IMAGE_TYPE_CODEC"
	case DX_IMAGE_TYPE_YUV_UV32:
		return "DX_IMAGE_TYPE_YUV_UV32"
	case DX_IMAGE_TYPE_TF_Y8:
		return "DX_IMAGE_TYPE_TF_Y8"
	case DX_IMAGE_TYPE_TF_A8:
		return "DX_IMAGE_TYPE_TF_A8"
	case DX_IMAGE_TYPE_TF_SHORT:
		return "DX_IMAGE_TYPE_TF_SHORT"
	case DX_IMAGE_TYPE_TF_1BPP:
		return "DX_IMAGE_TYPE_TF_1BPP"
	case DX_IMAGE_TYPE_OPENGL:
		return "DX_IMAGE_TYPE_OPENGL"
	case DX_IMAGE_TYPE_YUV444I:
		return "DX_IMAGE_TYPE_YUV444I"
	case DX_IMAGE_TYPE_YUV422PLANAR:
		return "DX_IMAGE_TYPE_YUV422PLANAR"
	case DX_IMAGE_TYPE_ARGB8888:
		return "DX_IMAGE_TYPE_ARGB8888"
	case DX_IMAGE_TYPE_XRGB8888:
		return "DX_IMAGE_TYPE_XRGB8888"
	case DX_IMAGE_TYPE_YUV422YUYV:
		return "DX_IMAGE_TYPE_YUV422YUYV"
	case DX_IMAGE_TYPE_YUV422YVYU:
		return "DX_IMAGE_TYPE_YUV422YVYU"
	case DX_IMAGE_TYPE_YUV422UYVY:
		return "DX_IMAGE_TYPE_YUV422UYVY"
	case DX_IMAGE_TYPE_YUV422VYUY:
		return "DX_IMAGE_TYPE_YUV422VYUY"
	case DX_IMAGE_TYPE_RGBX32:
		return "DX_IMAGE_TYPE_RGBX32"
	case DX_IMAGE_TYPE_RGBX8888:
		return "DX_IMAGE_TYPE_RGBX8888"
	case DX_IMAGE_TYPE_BGRX8888:
		return "DX_IMAGE_TYPE_BGRX8888"
	case DX_IMAGE_TYPE_YUV420SP:
		return "DX_IMAGE_TYPE_YUV420SP"
	case DX_IMAGE_TYPE_YUV444PLANAR:
		return "DX_IMAGE_TYPE_YUV444PLANAR"
	case DX_IMAGE_TYPE_TF_U8:
		return "DX_IMAGE_TYPE_TF_U8"
	case DX_IMAGE_TYPE_TF_V8:
		return "DX_IMAGE_TYPE_TF_V8"
	case DX_IMAGE_TYPE_YUV420_16:
		return "DX_IMAGE_TYPE_YUV420_16"
	case DX_IMAGE_TYPE_YUV_UV_16:
		return "DX_IMAGE_TYPE_YUV_UV_16"
	case DX_IMAGE_TYPE_YUV420_S:
		return "DX_IMAGE_TYPE_YUV420_S"
	default:
		return "[?? Invalid DX_ImageType value]"
	}
}