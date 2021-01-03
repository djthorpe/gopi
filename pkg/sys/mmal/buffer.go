//+build mmal

package mmal

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"strings"
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
	MMALBuffer     C.MMAL_BUFFER_HEADER_T
	MMALBufferFlag C.uint32_t
)

////////////////////////////////////////////////////////////////////////////////
// CONSTS

const (
	MMAL_BUFFER_HEADER_FLAG_EOS                 MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_EOS
	MMAL_BUFFER_HEADER_FLAG_FRAME_START         MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_FRAME_START
	MMAL_BUFFER_HEADER_FLAG_FRAME_END           MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_FRAME_END
	MMAL_BUFFER_HEADER_FLAG_FRAME               MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_FRAME
	MMAL_BUFFER_HEADER_FLAG_KEYFRAME            MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_KEYFRAME
	MMAL_BUFFER_HEADER_FLAG_DISCONTINUITY       MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_DISCONTINUITY
	MMAL_BUFFER_HEADER_FLAG_CONFIG              MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_CONFIG
	MMAL_BUFFER_HEADER_FLAG_ENCRYPTED           MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_ENCRYPTED
	MMAL_BUFFER_HEADER_FLAG_CODECSIDEINFO       MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_CODECSIDEINFO
	MMAL_BUFFER_HEADER_FLAGS_SNAPSHOT           MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAGS_SNAPSHOT
	MMAL_BUFFER_HEADER_FLAG_CORRUPTED           MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_CORRUPTED
	MMAL_BUFFER_HEADER_FLAG_TRANSMISSION_FAILED MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_TRANSMISSION_FAILED
	MMAL_BUFFER_HEADER_FLAG_DECODEONLY          MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_DECODEONLY
	MMAL_BUFFER_HEADER_FLAG_NAL_END             MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_NAL_END
	MMAL_BUFFER_HEADER_FLAG_USER0               MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_USER0
	MMAL_BUFFER_HEADER_FLAG_USER1               MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_USER1
	MMAL_BUFFER_HEADER_FLAG_USER2               MMALBufferFlag = C.MMAL_BUFFER_HEADER_FLAG_USER2
	MMAL_BUFFER_HEADER_FLAG_USER3               MMALBufferFlag = (1 << 31)
	MMAL_BUFFER_HEADER_FLAG_MIN                                = MMAL_BUFFER_HEADER_FLAG_EOS
	MMAL_BUFFER_HEADER_FLAG_MAX                                = MMAL_BUFFER_HEADER_FLAG_USER3
	MMAL_BUFFER_HEADER_FLAG_NONE                MMALBufferFlag = 0
)

const (
	MMAL_TIME_UNKNOWN int64 = C.MMAL_TIME_UNKNOWN
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

func (this *MMALBuffer) Size() uint32 {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	return uint32(ctx.alloc_size)
}

// Get timestamps
func (this *MMALBuffer) PtsDts() (int64, int64) {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	return int64(ctx.pts), int64(ctx.dts)
}

// Set timestamps to "unknown"
func (this *MMALBuffer) ClearPtsDts() {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	ctx.pts, ctx.dts = C.MMAL_TIME_UNKNOWN, C.MMAL_TIME_UNKNOWN
}

// Flags returns any flags for buffer
func (this *MMALBuffer) Flags() MMALBufferFlag {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	return MMALBufferFlag(ctx.flags)
}

// HasFlag returns true if all flags presented are set
func (this *MMALBuffer) HasFlags(f MMALBufferFlag) bool {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	return MMALBufferFlag(ctx.flags)&f == f
}

// SetFlags for buffer
func (this *MMALBuffer) SetFlags(f MMALBufferFlag) {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	ctx.flags = C.uint32_t(f)
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

// AsData returns data as []byte array (where Event==NONE)
func (this *MMALBuffer) AsData() []byte {
	var result []byte

	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(ctx.alloc_size)
	sliceHeader.Len = int(ctx.length)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.data))

	// Return data
	return result
}

// AsError returns data as an error (where Event==ERROR)
func (this *MMALBuffer) AsError() error {
	var result Error
	buf := bytes.NewReader(this.AsData())
	if err := binary.Read(buf, binary.LittleEndian, &result); err != nil {
		return err
	} else {
		return result
	}
}

// AsFormat returns data as an stream format event (where Event==FORMAT CHANGED)
func (this *MMALBuffer) AsFormatChangeEvent() *MMALStreamFormatEvent {
	ctx := (*C.MMAL_BUFFER_HEADER_T)(this)
	return (*MMALStreamFormatEvent)(C.mmal_event_format_changed_get(ctx))
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
	if size := this.Size(); size != 0 {
		str += " size=" + fmt.Sprint(size)
	}
	if f := this.Flags(); f != 0 {
		str += " flags=" + fmt.Sprint(f)
	}
	if p, d := this.PtsDts(); p != MMAL_TIME_UNKNOWN && d != MMAL_TIME_UNKNOWN {
		str += fmt.Sprintf(" pts=%v dts=%v", p, d)
	}
	return str + ">"
}

func (f MMALBufferFlag) String() string {
	if f == 0 {
		return f.FlagString()
	}
	str := ""
	for v := MMAL_BUFFER_HEADER_FLAG_MIN; v <= MMAL_BUFFER_HEADER_FLAG_MAX; v <<= 1 {
		if f&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.Trim(str, "|")
}

func (f MMALBufferFlag) FlagString() string {
	switch f {
	case MMAL_BUFFER_HEADER_FLAG_NONE:
		return "MMAL_BUFFER_HEADER_FLAG_NONE"
	case MMAL_BUFFER_HEADER_FLAG_EOS:
		return "MMAL_BUFFER_HEADER_FLAG_EOS"
	case MMAL_BUFFER_HEADER_FLAG_FRAME_START:
		return "MMAL_BUFFER_HEADER_FLAG_FRAME_START"
	case MMAL_BUFFER_HEADER_FLAG_FRAME_END:
		return "MMAL_BUFFER_HEADER_FLAG_FRAME_END"
	case MMAL_BUFFER_HEADER_FLAG_KEYFRAME:
		return "MMAL_BUFFER_HEADER_FLAG_KEYFRAME"
	case MMAL_BUFFER_HEADER_FLAG_DISCONTINUITY:
		return "MMAL_BUFFER_HEADER_FLAG_DISCONTINUITY"
	case MMAL_BUFFER_HEADER_FLAG_CONFIG:
		return "MMAL_BUFFER_HEADER_FLAG_CONFIG"
	case MMAL_BUFFER_HEADER_FLAG_ENCRYPTED:
		return "MMAL_BUFFER_HEADER_FLAG_ENCRYPTED"
	case MMAL_BUFFER_HEADER_FLAG_CODECSIDEINFO:
		return "MMAL_BUFFER_HEADER_FLAG_CODECSIDEINFO"
	case MMAL_BUFFER_HEADER_FLAGS_SNAPSHOT:
		return "MMAL_BUFFER_HEADER_FLAGS_SNAPSHOT"
	case MMAL_BUFFER_HEADER_FLAG_CORRUPTED:
		return "MMAL_BUFFER_HEADER_FLAG_CORRUPTED"
	case MMAL_BUFFER_HEADER_FLAG_TRANSMISSION_FAILED:
		return "MMAL_BUFFER_HEADER_FLAG_TRANSMISSION_FAILED"
	case MMAL_BUFFER_HEADER_FLAG_DECODEONLY:
		return "MMAL_BUFFER_HEADER_FLAG_DECODEONLY"
	case MMAL_BUFFER_HEADER_FLAG_NAL_END:
		return "MMAL_BUFFER_HEADER_FLAG_NAL_END"
	case MMAL_BUFFER_HEADER_FLAG_USER0:
		return "MMAL_BUFFER_HEADER_FLAG_USER0"
	case MMAL_BUFFER_HEADER_FLAG_USER1:
		return "MMAL_BUFFER_HEADER_FLAG_USER1"
	case MMAL_BUFFER_HEADER_FLAG_USER2:
		return "MMAL_BUFFER_HEADER_FLAG_USER2"
	case MMAL_BUFFER_HEADER_FLAG_USER3:
		return "MMAL_BUFFER_HEADER_FLAG_USER3"
	default:
		return "[?? Invalid MMALBufferFlag value]"
	}
}

/*

func MMALBufferSetLength(handle MMAL_Buffer, length uint32) error {
	if length > uint32(handle.alloc_size) {
		return MMAL_EINVAL
	}
	handle.length = C.uint32_t(length)
	return nil
}

*/
