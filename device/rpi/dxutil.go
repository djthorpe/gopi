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

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

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
	return fmt.Sprintf("<rpi.DXFrame>{ <Origin>{%v,%v} <Size>{%v,%v} }",this.DXPoint,this.DXSize)
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



