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

import (
	khronos "github.com/djthorpe/gopi/khronos"
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
	model  DXColorModel    // color model, which should be RGBA32 (4 bytes per pixel)
	size   khronos.EGLSize // size of the bitmap
	stride uint32          // number of bytes per row rounded up to 16-byte boundaries
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DX_RESOURCE_NONE    dxResourceHandle = 0
	DX_RESOURCE_SUCCESS                  = DX_SUCCESS
)

////////////////////////////////////////////////////////////////////////////////

func (this *DXDisplay) CreateResource(model DXColorModel, size khronos.EGLSize) (*DXResource, error) {
	resource := new(DXResource)
	resource.size = size
	resource.model = model
	resource.stride = dxAlignUp(uint32(size.Width), uint32(16)) * 4
	this.log.Debug2("<rpi.DX>CreateResource model=%v size=%v stride=%v", model, size)
	resource.handle = dxResourceCreate(model, size.Width, size.Height)
	if resource.handle == DX_RESOURCE_NONE {
		return nil, this.log.Error("dxResourceCreate failed")
	}
	return resource, nil
}

func (this *DXDisplay) DestroyResource(resource *DXResource) error {
	this.log.Debug2("<rpi.DX>DestroyResource")
	if dxResourceDelete(resource.handle) != true {
		return this.log.Error("dxResourceDelete failed")
	}
	resource.handle = DX_RESOURCE_NONE
	return nil
}

func (this *DXResource) String() string {
	return fmt.Sprintf("<rpi.DXResource>{ handle=%v model=%v size=%v stride=%v }", this.handle, this.model, this.size, this.stride)
}

func (h dxResourceHandle) String() string {
	return fmt.Sprintf("<rpi.DXResourceHandle>{%08X}", uint32(h))
}

func (this *DXResource) GetSize() khronos.EGLSize {
	return this.size
}

func (this *DXResource) GetHandle() dxResourceHandle {
	return this.handle
}

func (this *DXResource) ClearToColor(color khronos.EGLColorRGBA32) error {
	data := make([]uint32,uint(this.stride / 4) * uint(this.size.Height))
	value := color.Uint32()
	fmt.Printf("v=%08X\n",value)
	for i := 0; i < len(data); i++ {
		data[i] = value
	}
	dst_frame := DXFrame{ DXPoint{ int32(0), int32(0) }, DXSize{ uint32(this.size.Width), uint32(this.size.Height) } }
	dxResourceWriteData(this.handle, this.model, this.stride, unsafe.Pointer(&data[0]), &dst_frame)
	return nil
}

func (this *DXResource) SetPixel(pt khronos.EGLPoint,color khronos.EGLColorRGBA32) error {
	// TODO
	// dst_frame := DXFrame{ DXPoint{int32(0), int32(0)}, DXSize{uint32(1), uint32(1)} }
	// dxResourceWriteData(this.handle, this.model, this.stride, unsafe.Pointer(&source[0]), &dst_frame)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods

func dxResourceCreate(model DXColorModel, w, h uint) dxResourceHandle {
	var dummy C.uint32_t
	return dxResourceHandle(C.vc_dispmanx_resource_create(C.VC_IMAGE_TYPE_T(model), C.uint32_t(w), C.uint32_t(h), (*C.uint32_t)(unsafe.Pointer(&dummy))))
}

func dxResourceDelete(handle dxResourceHandle) bool {
	return C.vc_dispmanx_resource_delete(C.DISPMANX_RESOURCE_HANDLE_T(handle)) == DX_RESOURCE_SUCCESS
}

func dxResourceWriteData(handle dxResourceHandle, model DXColorModel, src_pitch uint32, src_buffer unsafe.Pointer, dst_rect *DXFrame) bool {
	return C.vc_dispmanx_resource_write_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), C.VC_IMAGE_TYPE_T(model), C.int(src_pitch), src_buffer, (*C.VC_RECT_T)(unsafe.Pointer(dst_rect))) == DX_RESOURCE_SUCCESS
}

func dxResourceReadData(handle dxResourceHandle, src_rect *DXFrame, dst_buffer uintptr, dst_pitch int) bool {
	return C.vc_dispmanx_resource_read_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), (*C.VC_RECT_T)(unsafe.Pointer(src_rect)), unsafe.Pointer(dst_buffer), C.uint32_t(dst_pitch)) == DX_RESOURCE_SUCCESS
}

func dxAlignUp(value, alignment uint32) uint32 {
	return ((value - 1) & ^(alignment - 1)) + alignment
}
