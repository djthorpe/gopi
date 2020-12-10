// +build ffmpeg

package ffmpeg

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libavcodec
#include <libavcodec/avcodec.h>
*/
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	AVFrame     C.struct_AVFrame
	AVBufferRef C.AVBufferRef
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewFrame() *AVFrame {
	return (*AVFrame)(C.av_frame_alloc())
}

func (this *AVFrame) Free() {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	C.av_frame_free(&ctx)
}

func (this *AVFrame) Release() {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	C.av_frame_unref(ctx)
}

func NewAudioFrame(f AVSampleFormat, rate int, layout AVChannelLayout) *AVFrame {
	frame := NewFrame()
	if frame == nil {
		return nil
	}
	ctx := (*C.AVFrame)(frame)
	ctx.format = C.int(f)
	ctx.sample_rate = C.int(rate)
	ctx.channel_layout = C.uint64_t(layout)
	ctx.channels = C.av_get_channel_layout_nb_channels(C.uint64_t(layout))
	return frame
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *AVFrame) PixelFormat() AVPixelFormat {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	if ctx.format == -1 {
		return AV_PIX_FMT_NONE
	} else if ctx.channels != 0 {
		return AV_PIX_FMT_NONE
	} else {
		return AVPixelFormat(ctx.format)
	}
}

func (this *AVFrame) SampleFormat() AVSampleFormat {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	if ctx.format == -1 {
		return AV_SAMPLE_FMT_NONE
	} else if ctx.channels == 0 {
		return AV_SAMPLE_FMT_NONE
	} else {
		return AVSampleFormat(ctx.format)
	}
}

func (this *AVFrame) SampleRate() int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.sample_rate)
}

func (this *AVFrame) ChannelLayout() AVChannelLayout {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return AVChannelLayout(ctx.channel_layout)
}

func (this *AVFrame) Channels() int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.channels)
}

func (this *AVFrame) NumSamples() int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.nb_samples)
}

func (this *AVFrame) IsPlanar() bool {
	f := C.enum_AVSampleFormat(this.SampleFormat())
	if ret := C.av_sample_fmt_is_planar(f); ret == 0 {
		return false
	} else {
		return true
	}
}

func (this *AVFrame) KeyFrame() bool {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.key_frame) != 0
}

func (this *AVFrame) PictType() AVPictureType {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return AVPictureType(ctx.pict_type)
}

func (this *AVFrame) PictSize() (int, int) {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.width), int(ctx.height)
}

func (this *AVFrame) Buffer(plane int) *AVBufferRef {
	ctx := (*C.AVFrame)(this)
	if buf := (C.av_frame_get_plane_buffer(ctx, C.int(plane))); buf == nil {
		return nil
	} else {
		return (*AVBufferRef)(buf)
	}
}

func (this *AVFrame) StrideForPlane(i int) int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.linesize[i])
}

func (this *AVFrame) GetAudioBuffer(num_samples int) error {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))

	ctx.nb_samples = C.int(num_samples)
	if err := AVError(C.av_frame_get_buffer(ctx, 0)); err != 0 {
		return err
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// AVBufferRef

func (this *AVBufferRef) Data() []byte {
	var bytes []byte

	ctx := (*C.AVBufferRef)(this)
	if ctx.data == nil {
		return nil
	}
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&bytes)))
	sliceHeader.Cap = int(ctx.size)
	sliceHeader.Len = int(ctx.size)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.data))
	return bytes
}

func (this *AVBufferRef) Size() int {
	ctx := (*C.AVBufferRef)(this)
	return int(ctx.size)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *AVBufferRef) String() string {
	str := "<AVBufferRef"
	str += " size=" + fmt.Sprint(this.Size())
	return str + ">"
}

func (this *AVFrame) String() string {
	str := "<AVFrame"
	if f := this.SampleFormat(); f != AV_SAMPLE_FMT_NONE {
		str += " sample_format=" + fmt.Sprint(f)
		if sample_rate := this.SampleRate(); sample_rate > 0 {
			str += " sample_rate=" + fmt.Sprint(sample_rate)
		}
		if layout := this.ChannelLayout(); layout > 0 {
			str += " channel_layout=" + fmt.Sprint(layout)
		}
		if c := this.Channels(); c > 0 {
			str += " channels=" + fmt.Sprint(c)
		}
		if n := this.NumSamples(); n > 0 {
			str += " nb_samples=" + fmt.Sprint(n)
		}
		if this.IsPlanar() {
			str += " is_planar=true"
		}
	} else if f := this.PixelFormat(); f != AV_PIX_FMT_NONE {
		str += " pixel_format=" + fmt.Sprint(f)
		if key_frame := this.KeyFrame(); key_frame {
			str += " key_frame=true"
		}
		if pict_type := this.PictType(); pict_type != AV_PICTURE_TYPE_NONE {
			str += " pict_type=" + fmt.Sprint(pict_type)
		}
		if w, h := this.PictSize(); w >= 0 && h >= 0 {
			str += " pict_size={" + fmt.Sprint(w, ",", h) + "}"
		}
	}
	return str + ">"
}
