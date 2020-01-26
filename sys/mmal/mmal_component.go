// +build mmal

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mmal

import (
	"reflect"
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
// PUBLIC METHODS - COMPONENTS

func MMALComponentCreate(name string, handle *MMAL_ComponentHandle) error {
	var cHandle (*C.MMAL_COMPONENT_T)

	cName := C.CString(name)
	defer C.free(unsafe.Pointer(cName))

	if status := MMAL_Status(C.mmal_component_create(cName, &cHandle)); status == MMAL_SUCCESS {
		*handle = MMAL_ComponentHandle(cHandle)
		return nil
	} else {
		return status
	}
}

func MMALComponentDestroy(handle MMAL_ComponentHandle) error {
	if status := MMAL_Status(C.mmal_component_destroy(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALComponentAcquire(handle MMAL_ComponentHandle) error {
	C.mmal_component_acquire(handle)
	return nil
}

func MMALComponentRelease(handle MMAL_ComponentHandle) error {
	if status := MMAL_Status(C.mmal_component_release(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALComponentEnable(handle MMAL_ComponentHandle) error {
	if status := MMAL_Status(C.mmal_component_enable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALComponentDisable(handle MMAL_ComponentHandle) error {
	if status := MMAL_Status(C.mmal_component_disable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALComponentName(handle MMAL_ComponentHandle) string {
	return C.GoString(handle.name)
}

func MMALComponentId(handle MMAL_ComponentHandle) uint32 {
	return uint32(handle.id)
}

func MMALComponentIsEnabled(handle MMAL_ComponentHandle) bool {
	return (handle.is_enabled != 0)
}

func MMALComponentControlPort(handle MMAL_ComponentHandle) MMAL_PortHandle {
	return handle.control
}

func MMALComponentInputPortNum(handle MMAL_ComponentHandle) uint {
	return uint(handle.input_num)
}

func MMALComponentInputPortAtIndex(handle MMAL_ComponentHandle, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.input, uint(handle.input_num), index)
}

func MMALComponentOutputPortNum(handle MMAL_ComponentHandle) uint {
	return uint(handle.output_num)
}

func MMALComponentOutputPortAtIndex(handle MMAL_ComponentHandle, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.output, uint(handle.output_num), index)
}

func MMALComponentClockPortNum(handle MMAL_ComponentHandle) uint {
	return uint(handle.clock_num)
}

func MMALComponentClockPortAtIndex(handle MMAL_ComponentHandle, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.clock, uint(handle.clock_num), index)
}

func MMALComponentPortNum(handle MMAL_ComponentHandle) uint {
	return uint(handle.port_num)
}

func MMALComponentPortAtIndex(handle MMAL_ComponentHandle, index uint) MMAL_PortHandle {
	return mmal_component_port_at_index(handle.port, uint(handle.port_num), index)
}

func mmal_component_port_at_index(array **C.MMAL_PORT_T, num, index uint) MMAL_PortHandle {
	var handles []MMAL_PortHandle
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&handles)))
	sliceHeader.Cap = int(num)
	sliceHeader.Len = int(num)
	sliceHeader.Data = uintptr(unsafe.Pointer(array))
	return handles[index]
}
