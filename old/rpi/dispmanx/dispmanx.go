/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package dispmanx

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host -I/opt/vc/include/interface/vcos/pthreads
    #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lbcm_host -lGLESv2
	#include "vc_dispmanx.h"
*/
import "C"

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

type (
	DisplayHandle  uint32
	ElementHandle  uint32
	UpdateHandle   uint32
	ResourceHandle uint32
	Protection     uint32
)

type Rect struct {
	X, Y          uint32
	Width, Height uint32
}

type Window struct {
	Element       ElementHandle
	Width, Height uint32
}

type Alpha struct {
	Flags   int32
	Opacity uint32
	Mask    ResourceHandle
}

type Clamp struct {
	Mode         int    // DISPMANX_FLAGS_CLAMP_T mode;
	KeyMask      int    // DISPMANX_FLAGS_KEYMASK_T key_mask;
	KeyValue     int    // DISPMANX_CLAMP_KEYS_T key_value;
	ReplaceValue uint32 // uint32_t replace_value;
}

type ModeInfo struct {
	Width, Height int32
	Transform     int
	InputFormat   int
	Device        uint32
}

////////////////////////////////////////////////////////////////////////////////

func DisplayOpen(device uint32) (DisplayHandle, error) {
	display := DisplayHandle(C.vc_dispmanx_display_open(C.uint32_t(device)))
	if display >= DisplayHandle(0) {
		return display, nil
	} else {
		return DisplayHandle(0), ErrorDisplay
	}
}

func DisplayClose(display DisplayHandle) error {
	if C.vc_dispmanx_display_close(C.DISPMANX_DISPLAY_HANDLE_T(display)) != DISPMANX_SUCCESS {
		return ErrorDisplay
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func DisplayGetInfo(display DisplayHandle,info *ModeInfo) error {
	if C.vc_dispmanx_display_get_info(C.DISPMANX_DISPLAY_HANDLE_T(display),(*C.DISPMANX_MODEINFO_T)(unsafe.Pointer(info))) != DISPMANX_SUCCESS {
		return ErrorGetInfo
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////


func DisplaySetBackground(update UpdateHandle, display DisplayHandle, red uint8, green uint8, blue uint8) {
	C.vc_dispmanx_display_set_background(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_DISPLAY_HANDLE_T(display), C.uint8_t(red), C.uint8_t(green), C.uint8_t(blue))
}

func UpdateStart(priority uint32) UpdateHandle {
	return UpdateHandle(C.vc_dispmanx_update_start(C.int32_t(priority)))
}

func ElementAdd(update UpdateHandle, display DisplayHandle, layer int32, dstRect *Rect, src ResourceHandle, srcRect *Rect, protection Protection, alpha *Alpha, clamp *Clamp, transform int) ElementHandle {
	return ElementHandle(C.vc_dispmanx_element_add(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_DISPLAY_HANDLE_T(display), C.int32_t(layer), (*C.VC_RECT_T)(unsafe.Pointer(dstRect)), C.DISPMANX_RESOURCE_HANDLE_T(src), (*C.VC_RECT_T)(unsafe.Pointer(srcRect)), C.DISPMANX_PROTECTION_T(protection), (*C.VC_DISPMANX_ALPHA_T)(unsafe.Pointer(alpha)), (*C.DISPMANX_CLAMP_T)(unsafe.Pointer(clamp)), C.DISPMANX_TRANSFORM_T(transform)))
}

func UpdateSubmitSync(update UpdateHandle) int {
	return int(C.vc_dispmanx_update_submit_sync(C.DISPMANX_UPDATE_HANDLE_T(update)))
}
