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
	dxResourceHandle uint32
)

type DXResource struct {
	handle dxResourceHandle
	model  DXColorModel
	size   DXSize
	buffer uintptr
}


////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DX_RESOURCE_NONE dxResourceHandle = 0
)

////////////////////////////////////////////////////////////////////////////////

func (this *DXDisplay) CreateResource(model DXColorModel,size DXSize) (*DXResource,error) {
	resource := new(DXResource)
	resource.size = size
	resource.model = model
	resource.handle, resource.buffer = resourceCreate(model,size.Width,size.Height)
	if resource.handle == DX_RESOURCE_NONE {
		return nil,this.log.Error("dxResourceCreate failed")
	}
	return resource,nil
}

func (this *DXDisplay) CloseResource(resource *DXResource) error {
	if resourceDelete(resource.handle) != true {
		return this.log.Error("dxResourceDelete failed")
	}
	resource.handle = DX_RESOURCE_NONE
	return nil
}

func (this *DXResource) String() string {
	return fmt.Sprintf("<rpi.DXResource>{ handle=%v model=%v size=%v buffer=%08X }",this.handle,this.model,this.size,this.buffer)
}

////////////////////////////////////////////////////////////////////////////////
// Private methods

func resourceCreate(model DXColorModel, w, h uint32) (dxResourceHandle,uintptr) {
	var ptr C.uint32_t
	handle := dxResourceHandle(C.vc_dispmanx_resource_create(C.VC_IMAGE_TYPE_T(model), C.uint32_t(w), C.uint32_t(h), (*C.uint32_t)(unsafe.Pointer(&ptr))))
	return handle,uintptr(ptr)
}

func resourceDelete(handle dxResourceHandle) bool {
	return C.vc_dispmanx_resource_delete(C.DISPMANX_RESOURCE_HANDLE_T(handle)) == DX_SUCCESS
}

func resourceWriteData(handle dxResourceHandle, model DXColorModel, src_pitch int, src_buffer uintptr, dst_rect *DXFrame) bool {
	return C.vc_dispmanx_resource_write_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), C.VC_IMAGE_TYPE_T(model), C.int(src_pitch), unsafe.Pointer(src_buffer), (*C.VC_RECT_T)(unsafe.Pointer(dst_rect))) == DX_SUCCESS
}

func resourceReadData(handle dxResourceHandle, src_rect *DXFrame, dst_buffer uintptr, dst_pitch int) bool {
	return C.vc_dispmanx_resource_read_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), (*C.VC_RECT_T)(unsafe.Pointer(src_rect)), unsafe.Pointer(dst_buffer), C.uint32_t(dst_pitch)) == DX_SUCCESS
}

