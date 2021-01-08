//+build mmal

package mmal

import (
	"fmt"
	"reflect"
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
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_COMPONENT_DEFAULT_VIDEO_DECODER   = "vc.ril.video_decode"
	MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER   = "vc.ril.video_encode"
	MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER  = "vc.ril.video_render"
	MMAL_COMPONENT_DEFAULT_IMAGE_DECODER   = "vc.ril.image_decode"
	MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER   = "vc.ril.image_encode"
	MMAL_COMPONENT_DEFAULT_CAMERA          = "vc.ril.camera"
	MMAL_COMPONENT_DEFAULT_VIDEO_CONVERTER = "vc.video_convert"
	MMAL_COMPONENT_DEFAULT_SPLITTER        = "vc.splitter"
	MMAL_COMPONENT_DEFAULT_SCHEDULER       = "vc.scheduler"
	MMAL_COMPONENT_DEFAULT_VIDEO_INJECTER  = "vc.video_inject"
	MMAL_COMPONENT_DEFAULT_VIDEO_SPLITTER  = "vc.ril.video_splitter"
	MMAL_COMPONENT_DEFAULT_AUDIO_RENDERER  = "vc.ril.audio_render"
	MMAL_COMPONENT_DEFAULT_MIRACAST        = "vc.miracast"
	MMAL_COMPONENT_DEFAULT_CLOCK           = "vc.clock"
	MMAL_COMPONENT_DEFAULT_CAMERA_INFO     = "vc.camera_info"
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
// METHODS

func (this *MMALComponent) ControlPort() *MMALPort {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return (*MMALPort)(ctx.control)
}

func (this *MMALComponent) InputPorts() []*MMALPort {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return this.ports(ctx.input, uint(ctx.input_num))
}

func (this *MMALComponent) OutputPorts() []*MMALPort {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return this.ports(ctx.output, uint(ctx.output_num))
}

func (this *MMALComponent) ClockPorts() []*MMALPort {
	ctx := (*C.MMAL_COMPONENT_T)(this)
	return this.ports(ctx.clock, uint(ctx.clock_num))
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
	if ctrl := this.ControlPort(); ctrl != nil {
		str += " control=" + fmt.Sprint(ctrl)
	}
	if input := this.InputPorts(); len(input) != 0 {
		str += " input=" + fmt.Sprint(input)
	}
	if output := this.OutputPorts(); len(output) != 0 {
		str += " output=" + fmt.Sprint(output)
	}
	if clock := this.ClockPorts(); len(clock) != 0 {
		str += " clock=" + fmt.Sprint(clock)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *MMALComponent) ports(array **C.MMAL_PORT_T, count uint) []*MMALPort {
	var result []*MMALPort
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(count)
	sliceHeader.Len = int(count)
	sliceHeader.Data = uintptr(unsafe.Pointer(array))
	return result
}
