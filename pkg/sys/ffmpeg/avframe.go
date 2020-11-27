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
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	AVFrame C.struct_AVFrame
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

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *AVFrame) Samples() int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.nb_samples)
}

func (this *AVFrame) PixelFormat() AVPixelFormat {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	if pixel_format := AVPixelFormat(ctx.format); pixel_format <= AV_PIX_FMT_NONE {
		return AV_PIX_FMT_NONE
	} else {
		return pixel_format
	}
}

func (this *AVFrame) SampleFormat() AVSampleFormat {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	if sample_format := AVSampleFormat(ctx.format); sample_format <= AV_SAMPLE_FMT_NONE {
		return AV_SAMPLE_FMT_NONE
	} else {
		return sample_format
	}
}

func (this *AVFrame) SampleRate() int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.sample_rate)
}

func (this *AVFrame) Channels() int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.channels)
}

func (this *AVFrame) KeyFrame() bool {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.key_frame) != 0
}

func (this *AVFrame) DisplayPictureNumber() int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.display_picture_number)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *AVFrame) String() string {
	str := "<AVFrame"
	if channels := this.Channels(); channels == 0 {
		str += " type=video"
		if pixel_format := this.PixelFormat(); pixel_format != AV_PIX_FMT_NONE {
			str += " pixel_format=" + fmt.Sprint(pixel_format)
		}
		if key_frame := this.KeyFrame(); key_frame {
			str += " key_frame=true"
		}
		str += " seq=" + fmt.Sprint(this.DisplayPictureNumber())
	} else {
		str += " type=audio"
		if sample_format := this.SampleFormat(); sample_format != AV_SAMPLE_FMT_NONE {
			str += " sample_format=" + fmt.Sprint(sample_format)
		}
		str += " audio_channels=" + fmt.Sprint(channels)
		if num_samples := this.Samples(); num_samples > 0 {
			str += " num_samples=" + fmt.Sprint(num_samples)
		}
		if sample_rate := this.SampleRate(); sample_rate > 0 {
			str += " sample_rate=" + fmt.Sprint(sample_rate)
		}
	}
	return str + ">"
}
