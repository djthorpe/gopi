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
	dxElementHandle  uint32
)

type DXElement struct {
	handle dxElementHandle
	frame *DXFrame
	layer uint16
}

const (
	DX_ELEMENT_NONE dxElementHandle = DX_NO_HANDLE
	DX_ELEMENT_SUCCESS = DX_SUCCESS
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// TODO: Allow DXResource and alphas to be set
func (this *DXDisplay) AddElement(update dxUpdateHandle, layer uint16, dst_rect *DXFrame, src_rect *DXFrame) (*DXElement, error) {

	// destination frame
	if dst_rect == nil {
		size := this.GetSize()
		dst_rect = &DXFrame{ DXPoint{ 0, 0}, size }
	}

	// set alpha to 255
	// TODO: Allow Alpha to be set
	alpha := dxAlpha{ DX_FLAGS_ALPHA_FIXED_ALL_PIXELS, 255, 0}

	// create element structure
	element := new(DXElement)

	// add element
	element.handle = dxElementAdd(update, this.handle, layer, dst_rect, DX_RESOURCE_NONE, src_rect, DX_PROTECTION_NONE, &alpha, nil, 0)
	if element.handle == DX_ELEMENT_NONE {
		return nil, this.log.Error("dxElementAdd failed")
	}

	// set other members of the element
	element.layer = layer
	element.frame = dst_rect

	// success
	return element, nil
}

func (this *DXDisplay) RemoveElement(update dxUpdateHandle, element *DXElement) error {
	if dxElementRemove(update, element.handle) != true {
		return this.log.Error("RemoveElement failed")
	}
	return nil
}

func (this *DXElement) GetHandle() dxElementHandle {
	return this.handle
}

func (this *DXElement) String() string {
	return fmt.Sprintf("<rpi.DXElement>{ handle=%v frame=%v layer=%v }",this.handle,this.frame,this.layer)
}

func (h dxElementHandle) String() string {
	return fmt.Sprintf("<rpi.dxElementHandle>{%08X}",uint32(h))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func dxElementAdd(update dxUpdateHandle, display dxDisplayHandle, layer uint16, dest_rect *DXFrame, src_resource dxResourceHandle, src_rect *DXFrame, protection DXProtection, alpha *dxAlpha, clamp *dxClamp, transform DXTransform) dxElementHandle {
	return dxElementHandle(C.vc_dispmanx_element_add(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_DISPLAY_HANDLE_T(display), C.int32_t(layer), (*C.VC_RECT_T)(unsafe.Pointer(dest_rect)), C.DISPMANX_RESOURCE_HANDLE_T(src_resource), (*C.VC_RECT_T)(unsafe.Pointer(src_rect)), C.DISPMANX_PROTECTION_T(protection), (*C.VC_DISPMANX_ALPHA_T)(unsafe.Pointer(alpha)), (*C.DISPMANX_CLAMP_T)(unsafe.Pointer(clamp)), C.DISPMANX_TRANSFORM_T(transform)))
}

func dxElementRemove(update dxUpdateHandle, element dxElementHandle) bool {
	return C.vc_dispmanx_element_remove(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element)) == DX_ELEMENT_SUCCESS
}

