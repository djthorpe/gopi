// +build ffmpeg

package ffmpeg

import (
	"fmt"
	"syscall"
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
	AVCodecContext C.struct_AVCodecContext
)

////////////////////////////////////////////////////////////////////////////////
// AVCodecContext

// NewAVCodecContext allocates an AVCodecContext and set its fields to
// default values
func NewAVCodecContext(codec *AVCodec) *AVCodecContext {
	return (*AVCodecContext)(C.avcodec_alloc_context3((*C.AVCodec)(codec)))
}

// Free AVCodecContext
func (this *AVCodecContext) Free() {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	C.avcodec_free_context(&ctx)
}

// Open will initialize the AVCodecContext to use the given AVCodec
func (this *AVCodecContext) Open(codec *AVCodec, options *AVDictionary) error {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	if err := AVError(C.avcodec_open2(ctx, (*C.AVCodec)(codec), (**C.struct_AVDictionary)(unsafe.Pointer(options)))); err != 0 {
		return err
	} else {
		return nil
	}
}

// Close a given AVCodecContext and free all the data associated with it, but
// not the AVCodecContext itself
func (this *AVCodecContext) Close() error {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	if err := AVError(C.avcodec_close(ctx)); err != 0 {
		return err
	} else {
		return nil
	}
}

// DecodePacket does the packet decode
func (this *AVCodecContext) DecodePacket(packet *AVPacket) error {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	if err := AVError(C.avcodec_send_packet(ctx, (*C.AVPacket)(packet))); err != 0 {
		return err
	} else {
		return nil
	}
}

// DecodeFrame does the frame decoding
func (this *AVCodecContext) DecodeFrame(frame *AVFrame) error {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	if err := AVError(C.avcodec_receive_frame(ctx, (*C.AVFrame)(frame))); err != 0 {
		if err.IsErrno(syscall.EAGAIN) {
			return syscall.EAGAIN
		} else if err.IsErrno(syscall.EINVAL) {
			return syscall.EINVAL
		} else {
			return err
		}
	} else {
		return nil
	}
}

func (this *AVCodecContext) Type() AVMediaType {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	return AVMediaType(ctx.codec_type)
}

func (this *AVCodecContext) Codec() *AVCodec {
	ctx := (*C.AVCodecContext)(this)
	return (*AVCodec)(ctx.codec)
}

func (this *AVCodecContext) Frame() int {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	return int(ctx.frame_number)
}

func (this *AVCodecContext) PixelFormat() AVPixelFormat {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	if pix_fmt := AVPixelFormat(ctx.pix_fmt); pix_fmt <= AV_PIX_FMT_NONE {
		return AV_PIX_FMT_NONE
	} else {
		return pix_fmt
	}
}

func (this *AVCodecContext) SampleFormat() AVSampleFormat {
	ctx := (*C.AVCodecContext)(unsafe.Pointer(this))
	if sample_format := AVSampleFormat(ctx.sample_fmt); sample_format <= AV_SAMPLE_FMT_NONE {
		return AV_SAMPLE_FMT_NONE
	} else {
		return sample_format
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *AVCodecContext) String() string {
	str := "<AVCodecContext"
	media_type := this.Type()
	if media_type != AVMEDIA_TYPE_UNKNOWN {
		str += " type=" + fmt.Sprint(media_type)
	}
	if media_type == AVMEDIA_TYPE_VIDEO {
		if pix_fmt := this.PixelFormat(); pix_fmt != AV_PIX_FMT_NONE {
			str += " pix_fmt=" + fmt.Sprint(pix_fmt)
		}
	}
	if media_type == AVMEDIA_TYPE_AUDIO {
		if sample_fmt := this.SampleFormat(); sample_fmt != AV_SAMPLE_FMT_NONE {
			str += " sample_fmt=" + fmt.Sprint(sample_fmt)
		}
	}
	if frame_number := this.Frame(); frame_number >= 0 {
		str += " frame_number=" + fmt.Sprint(frame_number)
	}
	if codec := this.Codec(); codec != nil {
		str += " codec=" + fmt.Sprint(codec)
	}

	return str + ">"
}
