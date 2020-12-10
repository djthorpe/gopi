// +build ffmpeg

package ffmpeg

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libswresample
#include <libswresample/swresample.h>

*/
import "C"
import "fmt"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	SwrContext C.SwrContext
)

////////////////////////////////////////////////////////////////////////////////
// VERSION

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewSwrContext() *SwrContext {
	return (*SwrContext)(C.swr_alloc())
}

func NewSwrContextEx(in_sample_fmt, out_sample_fmt AVSampleFormat, in_rate, out_rate int, in_ch_layout, out_ch_layout AVChannelLayout) *SwrContext {
	return (*SwrContext)(C.swr_alloc_set_opts(
		nil,
		C.int64_t(out_ch_layout),
		C.enum_AVSampleFormat(out_sample_fmt),
		C.int(out_rate),
		C.int64_t(in_ch_layout),
		C.enum_AVSampleFormat(in_sample_fmt),
		C.int(in_rate),
		0, nil))
}

func (this *SwrContext) Init() error {
	ctx := (*C.SwrContext)(this)
	if err := AVError(C.swr_init(ctx)); err != 0 {
		return err
	} else {
		return nil
	}
}

func (this *SwrContext) Close() {
	ctx := (*C.SwrContext)(this)
	C.swr_close(ctx)
}

func (this *SwrContext) Free() {
	ctx := (*C.SwrContext)(this)
	C.swr_free(&ctx)
}

func (this *SwrContext) IsInitialized() bool {
	ctx := (*C.SwrContext)(this)
	return C.swr_is_initialized(ctx) != 0
}

////////////////////////////////////////////////////////////////////////////////
// AVFrame

func (this *SwrContext) ConfigFrame(out, in *AVFrame) error {
	ctx := (*C.SwrContext)(this)
	if err := AVError(C.swr_config_frame(ctx, (*C.AVFrame)(out), (*C.AVFrame)(in))); err != 0 {
		return err
	} else if err := this.Init(); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *SwrContext) ConvertFrame(out, in *AVFrame) error {
	ctx := (*C.SwrContext)(this)
	if err := AVError(C.swr_convert_frame(ctx, (*C.AVFrame)(out), (*C.AVFrame)(in))); err != 0 {
		return err
	} else {
		return nil
	}
}

func (this *SwrContext) FlushFrame(out *AVFrame) error {
	ctx := (*C.SwrContext)(this)
	if err := AVError(C.swr_convert_frame(ctx, (*C.AVFrame)(out), nil)); err != 0 {
		return err
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *SwrContext) String() string {
	str := "<ffmpeg.swrcontext"
	str += " is_initialized=" + fmt.Sprint(this.IsInitialized())
	return str + ">"
}
