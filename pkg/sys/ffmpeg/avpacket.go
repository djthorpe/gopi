// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"reflect"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libavcodec
#include <libavcodec/avcodec.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	AVPacket C.struct_AVPacket
)

////////////////////////////////////////////////////////////////////////////////
// AVPACKET

// NewAVPacket allocates an AVPacket and set its fields to default values
func NewAVPacket() *AVPacket {
	return (*AVPacket)(C.av_packet_alloc())
}

// Free AVPacket, if the packet is reference counted, it will be unreferenced first
func (this *AVPacket) Free() {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	C.av_packet_free(&ctx)
}

// Release AVPacket, wiping packet data
func (this *AVPacket) Release() {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	C.av_packet_unref(ctx)
}

// Init optional fields of a packet with default values
func (this *AVPacket) Init() {
	C.av_init_packet((*C.AVPacket)(this))
}

func (this *AVPacket) Size() int {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	return int(ctx.size)
}

// Returns bytes for a packet
func (this *AVPacket) Bytes() []byte {
	var bytes []byte

	ctx := (*C.AVPacket)(unsafe.Pointer(this))

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&bytes)))
	sliceHeader.Cap = int(ctx.size)
	sliceHeader.Len = int(ctx.size)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.data))

	// Return slice
	return bytes
}

func (this *AVPacket) Stream() int {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	return int(ctx.stream_index)
}

func (this *AVPacket) Flags() AVPacketFlag {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	return AVPacketFlag(ctx.flags)
}

func (this *AVPacket) Pos() int64 {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	return int64(ctx.pos)
}

func (this *AVPacket) Duration() int64 {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	return int64(ctx.duration)
}

func (this *AVPacket) Pts() int64 {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	return int64(ctx.pts)
}

func (this *AVPacket) Dts() int64 {
	ctx := (*C.AVPacket)(unsafe.Pointer(this))
	return int64(ctx.dts)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *AVPacket) String() string {
	str := "<AVPacket"
	str += " size=" + fmt.Sprint(this.Size())
	if stream := this.Stream(); stream >= 0 {
		str += " stream=" + fmt.Sprint(stream)
	}
	if flags := this.Flags(); flags != 0 {
		str += " flags=" + fmt.Sprint(flags)
	}
	if pos := this.Pos(); pos >= 0 {
		str += " pos=" + fmt.Sprint(pos)
	}
	return str + ">"
}
