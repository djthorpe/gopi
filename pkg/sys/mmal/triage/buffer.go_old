//+build mmal

package mmal

import (
	"fmt"
	"reflect"
	"strings"
	"unsafe"

	// Frameworks
	hw "github.com/djthorpe/gopi-hw"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func MMALBufferCommand(handle MMAL_Buffer) hw.MMALEncodingType {
	return hw.MMALEncodingType(handle.cmd)
}

// Return complete allocated buffer
func MMALBufferBytes(handle MMAL_Buffer) []byte {
	var value []byte
	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
	sliceHeader.Cap = int(handle.alloc_size)
	sliceHeader.Len = int(handle.alloc_size)
	sliceHeader.Data = uintptr(unsafe.Pointer(handle.data))
	// Return data
	return value
}

// Return data from buffer
func MMALBufferData(handle MMAL_Buffer) []byte {
	var value []byte
	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&value)))
	sliceHeader.Cap = int(handle.alloc_size)
	sliceHeader.Len = int(handle.length)
	sliceHeader.Data = uintptr(unsafe.Pointer(handle.data))
	// Return data
	return value
}

func MMALBufferFlags(handle MMAL_Buffer) hw.MMALBufferFlag {
	return hw.MMALBufferFlag(handle.flags)
}

func MMALBufferSetFlags(handle MMAL_Buffer, value hw.MMALBufferFlag) {
	handle.flags = C.uint32_t(value)
}

func MMALBufferLength(handle MMAL_Buffer) uint32 {
	return uint32(handle.length)
}

func MMALBufferSetLength(handle MMAL_Buffer, length uint32) error {
	if length > uint32(handle.alloc_size) {
		return MMAL_EINVAL
	}
	handle.length = C.uint32_t(length)
	return nil
}

func MMALBufferOffset(handle MMAL_Buffer) uint32 {
	return uint32(handle.offset)
}

func MMALBufferString(handle MMAL_Buffer) string {
	if handle == nil {
		return fmt.Sprintf("<MMAL_Buffer>{ nil }")
	} else {
		parts := ""
		parts += fmt.Sprintf("alloc_size=%v ", handle.alloc_size)
		parts += fmt.Sprintf("length=%v ", handle.length)
		if handle.offset != 0 {
			parts += fmt.Sprintf("offset=%v ", handle.offset)
		}
		if handle.flags != 0 {
			parts += fmt.Sprintf("flags=%v ", hw.MMALBufferFlag(handle.flags))
		}
		if handle.cmd != 0 {
			parts += fmt.Sprintf("cmd=%v ", hw.MMALEncodingType(handle.cmd))
		}
		return fmt.Sprintf("<MMAL_Buffer>{ %v }", strings.TrimSpace(parts))
	}
}

func MMALBufferAcquire(handle MMAL_Buffer) error {
	C.mmal_buffer_header_acquire(handle)
	return nil
}

func MMALBufferRelease(handle MMAL_Buffer) error {
	C.mmal_buffer_header_release(handle)
	return nil
}

func MMALBufferReset(handle MMAL_Buffer) error {
	C.mmal_buffer_header_reset(handle)
	return nil
}
