//+build mmal

package mmal

import (
	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_connection.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - CONNECTIONS

func MMALPortConnectionCreate(handle *MMAL_PortConnection, output_port, input_port MMAL_PortHandle, flags hw.MMALPortConnectionFlags) error {
	var cHandle (*C.MMAL_CONNECTION_T)
	if status := MMAL_Status(C.mmal_connection_create(&cHandle, output_port, input_port, C.uint(flags))); status == MMAL_SUCCESS {
		*handle = MMAL_PortConnection(cHandle)
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionAcquire(handle MMAL_PortConnection) error {
	C.mmal_connection_acquire(handle)
	return nil
}

func MMALPortConnectionRelease(handle MMAL_PortConnection) error {
	if status := MMAL_Status(C.mmal_connection_release(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionDestroy(handle MMAL_PortConnection) error {
	if status := MMAL_Status(C.mmal_connection_destroy(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionEnable(handle MMAL_PortConnection) error {
	if status := MMAL_Status(C.mmal_connection_enable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionDisable(handle MMAL_PortConnection) error {
	if status := MMAL_Status(C.mmal_connection_disable(handle)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionEventFormatChanged(handle MMAL_PortConnection, buffer MMAL_Buffer) error {
	if status := MMAL_Status(C.mmal_connection_event_format_changed(handle, buffer)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func MMALPortConnectionEnabled(handle MMAL_PortConnection) bool {
	return (handle.is_enabled != 0)
}
