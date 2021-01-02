//+build mmal

package mmal

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
*/
import "C"
import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALBuffer C.MMAL_BUFFER_HEADER_T
)

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *MMALBuffer) Event() MMALBufferEvent {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	return MMALBufferEvent(ctx.cmd)
}

func (this *MMALBuffer) Offset() uint32 {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	return uint32(ctx.offset)
}

func (this *MMALBuffer) Length() uint32 {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	return uint32(ctx.length)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *MMALBuffer) Acquire() {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	C.mmal_buffer_header_acquire(ctx)
}

func (this *MMALBuffer) Release() {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	C.mmal_buffer_header_release(ctx)
}

func (this *MMALBuffer) Reset() {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	C.mmal_buffer_header_reset(ctx)
}

// Return complete allocated buffer
func (this *MMALBuffer) Bytes() []byte {
	var result []byte

	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.alloc_size)
	sliceHeader.Len = int(ctx.alloc_size)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.data))

	// Return data
	return result
}

// Fill buffer with data from file
func (this *MMALBuffer) Fill(r io.Reader) error {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	n, err := r.Read(this.Bytes())
	ctx.offset = C.uint32_t(0)
	ctx.length = C.uint32_t(n)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *MMALBuffer) String() string {
	str := "<mmal.buffer"
	if e := this.Event(); e != 0 {
		str += " event=" + fmt.Sprint(e)
	}
	if offset := this.Offset(); offset != 0 {
		str += " offset=" + fmt.Sprint(offset)
	}
	if length := this.Length(); length != 0 {
		str += " length=" + fmt.Sprint(length)
	}
	return str + ">"
}

/*

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

func MMALBufferSetLength(handle MMAL_Buffer, length uint32) error {
	if length > uint32(handle.alloc_size) {
		return MMAL_EINVAL
	}
	handle.length = C.uint32_t(length)
	return nil
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

*/
