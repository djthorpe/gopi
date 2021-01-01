//+build mmal

package mmal

import (
	"fmt"
	"strconv"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALComponent C.MMAL_COMPONENT_T
	MMALPort      C.MMAL_PORT_T
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_COMPONENT_DEFAULT_VIDEO_DECODER    = "vc.ril.video_decode"
	MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER    = "vc.ril.video_encode"
	MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER   = "vc.ril.video_render"
	MMAL_COMPONENT_DEFAULT_IMAGE_DECODER    = "vc.ril.image_decode"
	MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER    = "vc.ril.image_encode"
	MMAL_COMPONENT_DEFAULT_CAMERA           = "vc.ril.camera"
	MMAL_COMPONENT_DEFAULT_VIDEO_CONVERTER  = "vc.video_convert"
	MMAL_COMPONENT_DEFAULT_SPLITTER         = "vc.splitter"
	MMAL_COMPONENT_DEFAULT_SCHEDULER        = "vc.scheduler"
	MMAL_COMPONENT_DEFAULT_VIDEO_INJECTER   = "vc.video_inject"
	MMAL_COMPONENT_DEFAULT_VIDEO_SPLITTER   = "vc.ril.video_splitter"
	MMAL_COMPONENT_DEFAULT_AUDIO_DECODER    = "none"
	MMAL_COMPONENT_DEFAULT_AUDIO_RENDERER   = "vc.ril.audio_render"
	MMAL_COMPONENT_DEFAULT_MIRACAST         = "vc.miracast"
	MMAL_COMPONENT_DEFAULT_CLOCK            = "vc.clock"
	MMAL_COMPONENT_DEFAULT_CAMERA_INFO      = "vc.camera_info"
	MMAL_COMPONENT_DEFAULT_CONTAINER_READER = "container_reader"
	MMAL_COMPONENT_DEFAULT_CONTAINER_WRITER = "container_writer"
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func MMALComponentCreate(name string) (*MMALComponent, error) {
	var ctx (*C.MMAL_COMPONENT_T)

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	if status := Error(C.mmal_component_create(cName, &ctx)); status == MMAL_SUCCESS {
		return (*MMALComponent)(ctx), nil
	} else {
		return nil, status
	}
}

func (this *MMALComponent) Free() error {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	if status := Error(C.mmal_component_destroy(ctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *MMALComponent) Id() uint32 {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return uint32(ctx.id)
}

func (this *MMALComponent) Name() string {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return C.GoString(ctx.name)
}

func (this *MMALComponent) Enabled() bool {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return (ctx.is_enabled != 0)
}

func (this *MMALComponent) CountPorts() uint {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return uint(ctx.port_num)
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *MMALComponent) Acquire() {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	C.mmal_component_acquire(ctx)
}

func (this *MMALComponent) Release() error {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	if status := Error(C.mmal_component_release(ctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func (this *MMALComponent) Enable() error {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	if status := Error(C.mmal_component_enable(ctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func (this *MMALComponent) Disable() error {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	if status := Error(C.mmal_component_disable(ctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS - PORTS

func (this *MMALComponent) ControlPort() *MMALPort {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return (*MMALPort)(ctx.control)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *MMALComponent) String() string {
	str := "<mmal.component"
	if id := this.Id(); id > 0 {
		str += " id=" + fmt.Sprint(id)
	}
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if enabled := this.Enabled(); enabled {
		str += " enabled=true"
	}
	if ports := this.CountPorts(); ports > 0 {
		str += " count_ports=" + fmt.Sprint(ports)
	}
	return str + ">"
}

/*

func MMALComponentInputPortNum(handle MMALComponent) uint {
	return uint(handle.input_num)
}

func MMALComponentInputPortAtIndex(handle MMALComponent, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.input, uint(handle.input_num), index)
}

func MMALComponentOutputPortNum(handle MMALComponent) uint {
	return uint(handle.output_num)
}

func MMALComponentOutputPortAtIndex(handle MMALComponent, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.output, uint(handle.output_num), index)
}

func MMALComponentClockPortNum(handle MMAL_Component) uint {
	return uint(handle.clock_num)
}

func MMALComponentClockPortAtIndex(handle MMAL_Component, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.clock, uint(handle.clock_num), index)
}


func mmal_component_port_at_index(array **C.MMAL_PORT_T, num, index uint) MMAL_PortHandle {
	var handles []MMAL_PortHandle
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&handles)))
	sliceHeader.Cap = int(num)
	sliceHeader.Len = int(num)
	sliceHeader.Data = uintptr(unsafe.Pointer(array))
	return handles[index]
}

*/
