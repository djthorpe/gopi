/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package mmal /* import "github.com/djthorpe/gopi/device/rpi/mmal" */

import (
	"fmt"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/mmal
    #cgo LDFLAGS:  -L/opt/vc/lib -lmmal -lmmal_components -lmmal_core
	#include <mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Port struct {
	handle *C.MMAL_PORT_T
}

type PortType uint32

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_PORT_TYPE_UNKNOWN PortType = iota
	MMAL_PORT_TYPE_CONTROL
	MMAL_PORT_TYPE_INPUT
	MMAL_PORT_TYPE_OUTPUT
	MMAL_PORT_TYPE_CLOCK
	MMAL_PORT_TYPE_INVALID PortType = 0xFFFFFFFF
)

////////////////////////////////////////////////////////////////////////////////
// GET PROPERTIES

func (this *Port) Type() PortType {
	return PortType(this.handle._type)
}

func (this *Port) Index() uint16 {
	return uint16(this.handle.index)
}

func (this *Port) Name() string {
	return C.GoString(this.handle.name)
}

func (this *Port) IsEnabled() bool {
	if uint32(this.handle.is_enabled) == 0 {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t PortType) String() string {
	switch t {
	case MMAL_PORT_TYPE_UNKNOWN:
		return "MMAL_PORT_TYPE_UNKNOWN"
	case MMAL_PORT_TYPE_CONTROL:
		return "MMAL_PORT_TYPE_CONTROL"
	case MMAL_PORT_TYPE_INPUT:
		return "MMAL_PORT_TYPE_INPUT"
	case MMAL_PORT_TYPE_OUTPUT:
		return "MMAL_PORT_TYPE_OUTPUT"
	case MMAL_PORT_TYPE_CLOCK:
		return "MMAL_PORT_TYPE_CLOCK"
	case MMAL_PORT_TYPE_INVALID:
		return "MMAL_PORT_TYPE_INVALID"
	default:
		return "[?? Unknown PortType value]"

	}
}

func (this *Port) String() string {
	return fmt.Sprintf("<mmal.Port>{ name=%v type=%v index=%v enabled=%v }", this.Name(), this.Type(), this.Index(), this.IsEnabled())
}

////////////////////////////////////////////////////////////////////////////////
// PORT CONTROL

func (this *Port) Flush() error {
	ret := status(C.mmal_port_flush(unsafe.Pointer(this.handle)))
	if ret != MMAL_SUCCESS {
		return ret.Error()
	}
	return nil
}

func (this *Port) Disable() error {
	ret := status(C.mmal_port_disable(unsafe.Pointer(this.handle)))
	if ret != MMAL_SUCCESS {
		return ret.Error()
	}
	return nil
}

func (this *Port) FormatCommit() error {
	ret := status(C.mmal_port_format_commit(unsafe.Pointer(this.handle)))
	if ret != MMAL_SUCCESS {
		return ret.Error()
	}
	return nil
}
