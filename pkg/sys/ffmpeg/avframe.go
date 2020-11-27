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

func (this *AVFrame) BytesForPlane(i int) []byte {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))

	// Make a fake slice
	var bytes []byte
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&bytes)))
	sliceHeader.Cap = int(ctx.linesize[i] * ctx.height)
	sliceHeader.Len = int(ctx.linesize[i] * ctx.height)
	sliceHeader.Data = uintptr(unsafe.Pointer(ctx.data[i]))

	// Return slice
	return bytes
}

func (this *AVFrame) StrideForPlane(i int) int {
	ctx := (*C.AVFrame)(unsafe.Pointer(this))
	return int(ctx.linesize[i])
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *AVFrame) String() string {
	str := "<AVFrame"
	if key_frame := this.KeyFrame(); key_frame {
		str += " key_frame=true"
	}
	if pict_type := this.PictType(); pict_type != AV_PICTURE_TYPE_NONE {
		str += " pict_type=" + fmt.Sprint(pict_type)
	}
	if w, h := this.PictSize(); w >= 0 && h >= 0 {
		str += " pict_size={" + fmt.Sprint(w, ",", h) + "}"
	}
	return str + ">"
}
