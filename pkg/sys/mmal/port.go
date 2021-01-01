//+build mmal

package mmal

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util.h>
#include <interface/mmal/util/mmal_util_params.h>

// Callback Functions
void mmal_port_callback(MMAL_PORT_T* port, MMAL_BUFFER_HEADER_T* buffer);
*/
import "C"
import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALPort           C.MMAL_PORT_T
	MMALPortCallback   func(port *MMALPort, buffer *MMALBuffer)
	MMALPortType       uint
	MMALPortCapability uint
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_PORT_TYPE_UNKNOWN MMALPortType = C.MMAL_PORT_TYPE_UNKNOWN
	MMAL_PORT_TYPE_CONTROL MMALPortType = C.MMAL_PORT_TYPE_CONTROL
	MMAL_PORT_TYPE_INPUT   MMALPortType = C.MMAL_PORT_TYPE_INPUT
	MMAL_PORT_TYPE_OUTPUT  MMALPortType = C.MMAL_PORT_TYPE_OUTPUT
	MMAL_PORT_TYPE_CLOCK   MMALPortType = C.MMAL_PORT_TYPE_CLOCK
)

const (
	MMAL_PORT_CAPABILITY_PASSTHROUGH                  MMALPortCapability = C.MMAL_PORT_CAPABILITY_PASSTHROUGH
	MMAL_PORT_CAPABILITY_ALLOCATION                   MMALPortCapability = C.MMAL_PORT_CAPABILITY_ALLOCATION
	MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE MMALPortCapability = C.MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE
	MMAL_PORT_CAPABILITY_MIN                                             = MMAL_PORT_CAPABILITY_PASSTHROUGH
	MMAL_PORT_CAPABILITY_MAX                                             = MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE
)

////////////////////////////////////////////////////////////////////////////////
// CALLBACK REGISTRATION

var (
	portCallback = make(map[*C.MMAL_PORT_T]MMALPortCallback)
)

func MMALPortRegisterCallback(port *C.MMAL_PORT_T, fn MMALPortCallback) {
	if fn != nil {
		portCallback[port] = fn
	} else {
		delete(portCallback, port)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *MMALPort) Name() string {
	ctx := (*C.MMAL_PORT_T)(this)
	return C.GoString(ctx.name)
}

func (this *MMALPort) Type() MMALPortType {
	ctx := (*C.MMAL_PORT_T)(this)
	return MMALPortType(ctx._type)
}

func (this *MMALPort) Index() uint {
	ctx := (*C.MMAL_PORT_T)(this)
	return uint(ctx.index_all)
}

func (this *MMALPort) Enabled() bool {
	ctx := (*C.MMAL_PORT_T)(this)
	return (ctx.is_enabled != 0)
}

func (this *MMALPort) Capabilities() MMALPortCapability {
	ctx := (*C.MMAL_PORT_T)(this)
	return MMALPortCapability(ctx.capabilities)
}

func (this *MMALPort) Component() *MMALComponent {
	ctx := (*C.MMAL_PORT_T)(this)
	return (*MMALComponent)(ctx.component)
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *MMALPort) Enable() error {
	ctx := (*C.MMAL_PORT_T)(this)
	if status := Error(C.mmal_port_enable(ctx, C.MMAL_PORT_BH_CB_T(C.mmal_port_callback))); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func (this *MMALPort) Disable() error {
	ctx := (*C.MMAL_PORT_T)(this)
	if status := Error(C.mmal_port_disable(ctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func (this *MMALPort) Flush() error {
	ctx := (*C.MMAL_PORT_T)(this)
	if status := Error(C.mmal_port_flush(ctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func (this *MMALPort) Connect(other *MMALPort) error {
	ctx := (*C.MMAL_PORT_T)(this)
	otherctx := (*C.MMAL_PORT_T)(other)
	if status := Error(C.mmal_port_connect(ctx, otherctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func (this *MMALPort) Disconnect() error {
	ctx := (*C.MMAL_PORT_T)(this)
	if status := Error(C.mmal_port_disconnect(ctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func (this *MMALPort) FormatCommit() error {
	ctx := (*C.MMAL_PORT_T)(this)
	if status := Error(C.mmal_port_format_commit(ctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

// BufferCount returns current, minimum & recommended number of buffers the port requires
// A value of zero for recommendation means no special recommendation
func (this *MMALPort) BufferCount() (uint32, uint32, uint32) {
	ctx := (*C.MMAL_PORT_T)(this)
	return uint32(ctx.buffer_num), uint32(ctx.buffer_num_min), uint32(ctx.buffer_num_recommended)
}

// BufferSize returns current, minimum & recommended size of buffers the port requires
// A value of zero means no special recommendation
func (this *MMALPort) BufferSize() (uint32, uint32, uint32) {
	ctx := (*C.MMAL_PORT_T)(this)
	return uint32(ctx.buffer_size), uint32(ctx.buffer_size_min), uint32(ctx.buffer_size_recommended)
}

func (this *MMALPort) BufferSet(count, size uint32) {
	ctx := (*C.MMAL_PORT_T)(this)
	ctx.buffer_num = C.uint32_t(count)
	ctx.buffer_size = C.uint32_t(size)
}

// BufferAlignment returns minimum alignment requirement for the buffers.
// A value of zero means no special alignment requirements.
func (this *MMALPort) BufferAlignment() uint32 {
	ctx := (*C.MMAL_PORT_T)(this)
	return uint32(ctx.buffer_alignment_min)
}

func (this *MMALPort) SendBuffer(buffer *MMALBuffer) error {
	ctx := (*C.MMAL_PORT_T)(this)
	bufferctx := (*C.MMAL_BUFFER_HEADER_T)(buffer)
	if status := Error(C.mmal_port_send_buffer(ctx, bufferctx)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

func (this *MMALPort) SetURI(value string) error {
	ctx := (*C.MMAL_PORT_T)(this)
	cValue := C.CString(value)
	defer C.free(unsafe.Pointer(cValue))
	if status := Error(C.mmal_util_port_set_uri(ctx, cValue)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

/*

func MMALPortFormat(handle MMAL_PortHandle) MMAL_StreamFormat {
	return handle.format
}

func MMALPortSetDisplayRegion(handle MMAL_PortHandle, value MMAL_DisplayRegion) error {
	if status := MMAL_Status(C.mmal_util_set_display_region(handle, value)); status == MMAL_SUCCESS {
		return nil
	} else {
		return status
	}
}

*/

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *MMALPort) String() string {
	str := "<mmal.port"
	str += " index=" + fmt.Sprint(this.Index())
	str += " enabled=" + fmt.Sprint(this.Enabled())
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if typ := this.Type(); typ != 0 {
		str += " type=" + fmt.Sprint(typ)
	}
	if cap := this.Capabilities(); cap != 0 {
		str += " capabilities=" + fmt.Sprint(cap)
	}
	return str + ">"
}

func (p MMALPortType) String() string {
	switch p {
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
	default:
		return "[?? Invalid MMALPortType value]"
	}
}

func (c MMALPortCapability) String() string {
	parts := ""
	for flag := MMAL_PORT_CAPABILITY_MIN; flag <= MMAL_PORT_CAPABILITY_MAX; flag <<= 1 {
		if c&flag == 0 {
			continue
		}
		switch flag {
		case MMAL_PORT_CAPABILITY_PASSTHROUGH:
			parts += "|" + "MMAL_PORT_CAPABILITY_PASSTHROUGH"
		case MMAL_PORT_CAPABILITY_ALLOCATION:
			parts += "|" + "MMAL_PORT_CAPABILITY_ALLOCATION"
		case MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE:
			parts += "|" + "MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE"
		default:
			parts += "|" + "[?? Invalid MMALPortCapability value]"
		}
	}
	return strings.Trim(parts, "|")
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

//export mmal_port_callback
func mmal_port_callback(port *C.MMAL_PORT_T, buffer *C.MMAL_BUFFER_HEADER_T) {
	if fn, exists := portCallback[port]; exists {
		fn((*MMALPort)(port), (*MMALBuffer)(buffer))
	} else {
		// TODO
		fmt.Printf("mmal_port_callback{port=%v buffer=%v}\n", port, buffer)
	}
}
