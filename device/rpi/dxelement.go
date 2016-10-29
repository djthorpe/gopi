/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
	"unsafe"
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
	dxElementHandle uint32
)

type DXElement struct {
	handle dxElementHandle
	frame  *DXFrame
	layer  uint16
}

const (
	DX_ELEMENT_NONE    dxElementHandle = DX_NO_HANDLE
	DX_ELEMENT_SUCCESS                 = DX_SUCCESS
)

const (
	DX_ELEMENT_CHANGE_LAYER         uint32 = (1 << 0)
	DX_ELEMENT_CHANGE_OPACITY       uint32 = (1 << 1)
	DX_ELEMENT_CHANGE_DEST_RECT     uint32 = (1 << 2)
	DX_ELEMENT_CHANGE_SRC_RECT      uint32 = (1 << 3)
	DX_ELEMENT_CHANGE_MASK_RESOURCE uint32 = (1 << 4)
	DX_ELEMENT_CHANGE_TRANSFORM     uint32 = (1 << 5)
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// TODO: Allow DXResource to be set
func (this *DXDisplay) AddElement(update dxUpdateHandle, layer uint16, opacity uint32, dst_rect *DXFrame, src_rect *DXFrame) (*DXElement, error) {

	// destination frame
	if dst_rect == nil {
		size := this.GetSize()
		dst_rect = &DXFrame{DXPoint{0, 0}, size}
	}

	// set alpha
	alpha := dxAlpha{DX_FLAGS_ALPHA_FIXED_ALL_PIXELS, opacity, 0}

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

func (this *DXDisplay) SetElementDestination(update dxUpdateHandle, element *DXElement, frame DXFrame) error {
	if dxElementChangeDestinationFrame(update, element.handle, &frame) != true {
		return this.log.Error("SetElementDestination failed")
	}

	// update frame
	element.frame = &frame

	// return success
	return nil
}

func (this *DXElement) GetHandle() dxElementHandle {
	return this.handle
}

func (this *DXElement) GetLayer() uint16 {
	return this.layer
}

func (this *DXElement) String() string {
	return fmt.Sprintf("<rpi.DXElement>{ handle=%v frame=%v layer=%v }", this.handle, this.frame, this.layer)
}

func (h dxElementHandle) String() string {
	return fmt.Sprintf("<rpi.dxElementHandle>{%08X}", uint32(h))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func dxElementAdd(update dxUpdateHandle, display dxDisplayHandle, layer uint16, dest_rect *DXFrame, src_resource dxResourceHandle, src_rect *DXFrame, protection DXProtection, alpha *dxAlpha, clamp *dxClamp, transform DXTransform) dxElementHandle {
	return dxElementHandle(C.vc_dispmanx_element_add(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_DISPLAY_HANDLE_T(display), C.int32_t(layer), (*C.VC_RECT_T)(unsafe.Pointer(dest_rect)), C.DISPMANX_RESOURCE_HANDLE_T(src_resource), (*C.VC_RECT_T)(unsafe.Pointer(src_rect)), C.DISPMANX_PROTECTION_T(protection), (*C.VC_DISPMANX_ALPHA_T)(unsafe.Pointer(alpha)), (*C.DISPMANX_CLAMP_T)(unsafe.Pointer(clamp)), C.DISPMANX_TRANSFORM_T(transform)))
}

func dxElementRemove(update dxUpdateHandle, element dxElementHandle) bool {
	return C.vc_dispmanx_element_remove(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element)) == DX_ELEMENT_SUCCESS
}

func dxElementChangeDestinationFrame(update dxUpdateHandle, element dxElementHandle, frame *DXFrame) bool {
	return C.vc_dispmanx_element_change_attributes(
		C.DISPMANX_UPDATE_HANDLE_T(update),
		C.DISPMANX_ELEMENT_HANDLE_T(element),
		C.uint32_t(DX_ELEMENT_CHANGE_DEST_RECT),
		C.int32_t(0),                          // layer
		C.uint8_t(0),                          // opacity
		(*C.VC_RECT_T)(unsafe.Pointer(frame)), // dest_rect
		(*C.VC_RECT_T)(unsafe.Pointer(nil)),   // src_rect
		C.DISPMANX_RESOURCE_HANDLE_T(0),       // mask
		C.DISPMANX_TRANSFORM_T(0),             // transform
	) == DX_ELEMENT_SUCCESS
}
