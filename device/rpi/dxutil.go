/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"unsafe"
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host
    #cgo LDFLAGS:  -L/opt/vc/lib -lbcm_host
	#include "vc_dispmanx.h"
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	DXInputFormat    uint32
	DXTransform      int
	DXColorModel     int
	DXProtection     uint32
)

type DXPoint struct {
	X int32
	Y int32
}

type DXSize struct {
	Width  uint32
	Height uint32
}

type DXFrame struct {
	DXPoint
	DXSize
}

type (
	dxClampMode      int
)

type dxClamp struct {
	Mode    dxClampMode
	Flags   int
	Opacity uint32
	Mask    dxResourceHandle
}

type dxAlpha struct {
	Flags   uint32
	Opacity uint32
	Mask    dxResourceHandle
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	/* Success and failure conditions */
	DX_SUCCESS   = 0
	DX_INVALID   = -1
	DX_NO_HANDLE = 0
)

const (
	// DXTransform values
	DX_NO_ROTATE DXTransform = iota
	DX_ROTATE_90
	DX_ROTATE_180
	DX_ROTATE_270
)

const (
	// DXInputFormat values
	DX_INPUT_FORMAT_INVALID DXInputFormat = iota
	DX_INPUT_FORMAT_RGB888
	DX_INPUT_FORMAT_RGB565
)

const (
	// DXColorModel values
	// We only list defaults for supported color models on the Raspberry Pi
	DX_IMAGE_RGB565    DXColorModel = 1
	DX_IMAGE_YUV420    DXColorModel = 3
	DX_IMAGE_RGB888    DXColorModel = 5
	DX_IMAGE_4BPP      DXColorModel = 7  // 4bpp palettised image
	DX_IMAGE_RGBA32    DXColorModel = 15 /* RGB888 0xAABBGGRR */
	DX_IMAGE_YUV422    DXColorModel = 16 /* a line of Y (32-byte padded), a line of U (16-byte padded), and a line of V (16-byte padded) */
	DX_IMAGE_RGBA565   DXColorModel = 17 /* RGB565 with a transparent patch */
	DX_IMAGE_RGBA16    DXColorModel = 18 /* Compressed (4444) version of RGBA32 */
	DX_IMAGE_YUV_UV    DXColorModel = 19 /* VCIII codec format */
	DX_IMAGE_TF_RGBA32 DXColorModel = 20 /* VCIII T-format RGBA8888 */
	DX_IMAGE_TF_RGBX32 DXColorModel = 21 /* VCIII T-format RGBx8888 */
	DX_IMAGE_TF_RGBA16 DXColorModel = 23 /* VCIII T-format RGBA4444 */
	DX_IMAGE_TF_RGB565 DXColorModel = 25 /* VCIII T-format RGB565 */
)

const (
	/* Alpha flags */
	DX_FLAGS_ALPHA_FROM_SOURCE       uint32 = 0 /* Bottom 2 bits sets the alpha mode */
	DX_FLAGS_ALPHA_FIXED_ALL_PIXELS  uint32 = 1
	DX_FLAGS_ALPHA_FIXED_NON_ZERO    uint32 = 2
	DX_FLAGS_ALPHA_FIXED_EXCEED_0X07 uint32 = 3
	DX_FLAGS_ALPHA_PREMULT           uint32 = 1 << 16
	DX_FLAGS_ALPHA_MIX               uint32 = 1 << 17
)

const (
	/* Clamp values */
	DX_FLAGS_CLAMP_NONE               dxClampMode = 0
	DX_FLAGS_CLAMP_LUMA_TRANSPARENT   dxClampMode = 1
	DX_FLAGS_CLAMP_TRANSPARENT        dxClampMode = 2
	DX_FLAGS_CLAMP_CHROMA_TRANSPARENT dxClampMode = 2
	DX_FLAGS_CLAMP_REPLACE            dxClampMode = 3
)

const (
	/* Protection values */
	DX_PROTECTION_NONE DXProtection = 0
	DX_PROTECTION_HDCP DXProtection = 11
)

////////////////////////////////////////////////////////////////////////////////
// DXTransform

// Provide human-readable version of DXTransform value
func (t DXTransform) String() string {
	switch(t) {
	case DX_NO_ROTATE:
		return "DX_NO_ROTATE"
	case DX_ROTATE_90:
		return "DX_ROTATE_90"
	case DX_ROTATE_180:
		return "DX_ROTATE_180"
	case DX_ROTATE_270:
		return "DX_ROTATE_270"
	default:
		return "[Invalid DXTransform value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// DXInputFormat

func (f DXInputFormat) String() string {
	switch(f) {
	case DX_INPUT_FORMAT_RGB888:
		return "DX_INPUT_FORMAT_RGB888"
	case DX_INPUT_FORMAT_RGB565:
		return "DX_INPUT_FORMAT_RGB565"
	default:
		return "DX_INPUT_FORMAT_INVALID"
	}
}

////////////////////////////////////////////////////////////////////////////////
// DXFrame

func (this *DXFrame) Set(point DXPoint, size DXSize) {
	C.vc_dispmanx_rect_set((*C.VC_RECT_T)(unsafe.Pointer(this)), C.uint32_t(point.X), C.uint32_t(point.Y), C.uint32_t(size.Width), C.uint32_t(size.Height))
}

func (this *DXFrame) String() string {
	return fmt.Sprintf("<rpi.DXFrame>{%v,%v}",this.DXPoint,this.DXSize)
}

////////////////////////////////////////////////////////////////////////////////
// DXColorModel

func (m DXColorModel) String() string {
	switch(m) {
	case DX_IMAGE_RGB565:
		return "DX_IMAGE_RGB565"
	case DX_IMAGE_YUV420:
		return "DX_IMAGE_YUV420"
	case DX_IMAGE_RGB888:
		return "DX_IMAGE_RGB888"
	case DX_IMAGE_4BPP:
		return "DX_IMAGE_4BPP"
	case DX_IMAGE_RGBA32:
		return "DX_IMAGE_RGBA32"
	case DX_IMAGE_YUV422:
		return "DX_IMAGE_YUV422"
	case DX_IMAGE_RGBA565:
		return "DX_IMAGE_RGBA565"
	case DX_IMAGE_RGBA16:
		return "DX_IMAGE_RGBA16"
	case DX_IMAGE_YUV_UV:
		return "DX_IMAGE_YUV_UV"
	case DX_IMAGE_TF_RGBA32:
		return "DX_IMAGE_TF_RGBA32"
	case DX_IMAGE_TF_RGBX32:
		return "DX_IMAGE_TF_RGBX32"
	case DX_IMAGE_TF_RGBA16:
		return "DX_IMAGE_TF_RGBA16"
	case DX_IMAGE_TF_RGB565:
		return "DX_IMAGE_TF_RGB565"
	default:
		return "[Invalid DXColorModel value]"
	}
}



