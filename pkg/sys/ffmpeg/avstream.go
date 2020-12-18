// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libavformat
#include <libavformat/avformat.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	AVStream C.struct_AVStream
)

////////////////////////////////////////////////////////////////////////////////
// AVStream

func NewStream(ctx *AVFormatContext, codec *AVCodec) *AVStream {
	return (*AVStream)(C.avformat_new_stream(
		(*C.AVFormatContext)(ctx),
		(*C.AVCodec)(codec),
	))
}

func (this *AVStream) Index() int {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return int(ctx.index)
}

func (this *AVStream) Id() int {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return int(ctx.id)
}

func (this *AVStream) Metadata() *AVDictionary {
	return &AVDictionary{ctx: this.metadata}
}

func (this *AVStream) CodecPar() *AVCodecParameters {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return (*AVCodecParameters)(ctx.codecpar)
}

func (this *AVStream) Disposition() AVDisposition {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return AVDisposition(ctx.disposition)
}

func (this *AVStream) AttachedPicture() *AVPacket {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	if AVDisposition(ctx.disposition)&AV_DISPOSITION_ATTACHED_PIC == 0 {
		return nil
	} else {
		return (*AVPacket)(&this.attached_pic)
	}
}

func (this *AVStream) Duration() int64 {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return int64(ctx.duration)
}

func (this *AVStream) NumFrames() int64 {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return int64(ctx.nb_frames)
}

func (this *AVStream) StartTime() int64 {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return int64(ctx.start_time)
}

func (this *AVStream) TimeBase() AVRational {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return AVRational(ctx.time_base)
}

func (this *AVStream) MeanFrameRate() AVRational {
	ctx := (*C.AVStream)(unsafe.Pointer(this))
	return AVRational(ctx.avg_frame_rate)
}

func (this *AVStream) String() string {
	str := "<AVStream"
	str += " index=" + fmt.Sprint(this.Index())
	str += " id=" + fmt.Sprint(this.Id())
	str += " metadata=" + fmt.Sprint(this.Metadata())
	str += " codecpar=" + fmt.Sprint(this.CodecPar())
	if d := this.Disposition(); d != AV_DISPOSITION_NONE {
		str += " disposition=" + fmt.Sprint(this.Disposition())
	}
	return str + ">"
}
