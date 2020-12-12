// +build ffmpeg

package ffmpeg

import (
	"net/url"
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
	AVIOContext C.struct_AVIOContext
)

////////////////////////////////////////////////////////////////////////////////
// AVIO

func NewAVIOContext(url *url.URL, flags AVIOFlag) (*AVIOContext, error) {
	ctx := (*C.AVIOContext)(nil)

	url_ := C.CString(url.String())
	defer C.free(unsafe.Pointer(url_))

	if err := AVError(C.avio_open(&ctx, url_, C.int(flags))); err != 0 {
		return nil, err
	} else {
		return (*AVIOContext)(ctx), nil
	}
}

func (this *AVIOContext) Close() error {
	ctx := (*C.AVIOContext)(this)
	if err := AVError(C.avio_close(ctx)); err != 0 {
		return err
	} else {
		return nil
	}
}

func (this *AVIOContext) Free() {
	ctx := (*C.AVIOContext)(this)
	C.avio_context_free(&ctx)
}

func (this *AVIOContext) Flush() {
	ctx := (*C.AVIOContext)(this)
	C.avio_flush(ctx)
}

func (this *AVIOContext) Read(buf []byte) (int, error) {
	ctx := (*C.AVIOContext)(this)
	data := unsafe.Pointer(nil)
	if buf != nil {
		data = unsafe.Pointer(&buf[0])
	}
	size := len(buf)
	if ret := C.avio_read(ctx, (*C.uint8_t)(data), C.int(size)); ret >= 0 {
		return int(ret), nil
	} else {
		return -1, AVError(ret)
	}
}

func (this *AVIOContext) Write(buf []byte) {
	ctx := (*C.AVIOContext)(this)
	data := unsafe.Pointer(nil)
	if buf != nil {
		data = unsafe.Pointer(&buf[0])
	}
	size := len(buf)
	C.avio_write(ctx, (*C.uint8_t)(data), C.int(size))
}
