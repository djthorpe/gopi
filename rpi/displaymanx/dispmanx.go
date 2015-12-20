/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package egl

/*
	#cgo CFLAGS: -I/opt/vc/include/interface/vmcs_host
	#include "vc_dispmanx.h"
*/
import "C"

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

type (
	DispmanxDisplayHandle uint32
	DispmanxElementHandle uint32
	DispmanxUpdateHandle uint32
	DispmanxResourceHandle uint32
	DispmanxProtection uint32
)

type Rect struct {
	X,Y uint32
	Width,Height uint32
}

type Window struct {
	Element DispmanxElementHandle
	Width,Height uint32
}

type Alpha struct {
	Flags int32
	Opacity uint32
	Mask DispmanxResourceHandle
}

type Clamp struct {
  Mode int // DISPMANX_FLAGS_CLAMP_T mode;
  KeyMask int // DISPMANX_FLAGS_KEYMASK_T key_mask;
  KeyValue int // DISPMANX_CLAMP_KEYS_T key_value;
  ReplaceValue uint32 // uint32_t replace_value;
}

////////////////////////////////////////////////////////////////////////////////

func VCDispmanxDisplayOpen(device uint32) DispmanxDisplayHandle {
	return DispmanxDisplayHandle(C.vc_dispmanx_display_open(C.uint32_t(device)))
}

func VCDispmanxDisplaySetBackground(update DispmanxUpdateHandle,display DispmanxDisplayHandle,red uint8,green uint8,blue uint8) {
	C.vc_dispmanx_display_set_background(C.DISPMANX_UPDATE_HANDLE_T(update),C.DISPMANX_DISPLAY_HANDLE_T(display),C.uint8_t(red),C.uint8_t(green),C.uint8_t(blue))
}

func VCDispmanxUpdateStart(priority uint32) DispmanxUpdateHandle {
	return DispmanxUpdateHandle(C.vc_dispmanx_update_start(C.int32_t(priority)))
}

func VCDispmanxElementAdd(update DispmanxUpdateHandle,display DispmanxDisplayHandle,layer int32,dstRect *Rect,src DispmanxResourceHandle,srcRect *Rect,protection DispmanxProtection,alpha *Alpha,clamp *Clamp,transform int) DispmanxElementHandle {
	return DispmanxElementHandle(C.vc_dispmanx_element_add(C.DISPMANX_UPDATE_HANDLE_T(update),C.DISPMANX_DISPLAY_HANDLE_T(display),C.int32_t(layer),(*C.VC_RECT_T)(unsafe.Pointer(dstRect)),C.DISPMANX_RESOURCE_HANDLE_T(src),(*C.VC_RECT_T)(unsafe.Pointer(srcRect)),C.DISPMANX_PROTECTION_T(protection),(*C.VC_DISPMANX_ALPHA_T)(unsafe.Pointer(alpha)),(*C.DISPMANX_CLAMP_T)(unsafe.Pointer(clamp)),C.DISPMANX_TRANSFORM_T(transform)))
}

func VCDispmanxUpdateSubmitSync(update DispmanxUpdateHandle) int {
	return int(C.vc_dispmanx_update_submit_sync(C.DISPMANX_UPDATE_HANDLE_T(update)))
}

