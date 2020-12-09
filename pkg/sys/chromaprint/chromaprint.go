// +build chromaprint

package chromaprint

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libchromaprint
#include <chromaprint.h>
*/
import "C"
import "unsafe"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Error         int
	AlgorithmType C.int
	Context       C.ChromaprintContext
)

////////////////////////////////////////////////////////////////////////////////
// CONSTS

const (
	ALGORITHM_TEST1   AlgorithmType = C.CHROMAPRINT_ALGORITHM_TEST1
	ALGORITHM_TEST2   AlgorithmType = C.CHROMAPRINT_ALGORITHM_TEST2
	ALGORITHM_TEST3   AlgorithmType = C.CHROMAPRINT_ALGORITHM_TEST3
	ALGORITHM_TEST4   AlgorithmType = C.CHROMAPRINT_ALGORITHM_TEST4
	ALGORITHM_DEFAULT AlgorithmType = C.CHROMAPRINT_ALGORITHM_DEFAULT
)

const (
	errNone Error = iota
	errStart
	errFeed
	errFinish
	errFingerprint
)

////////////////////////////////////////////////////////////////////////////////
// METHODS

func Version() string {
	return C.GoString(C.chromaprint_get_version())
}

func NewChromaprint(algorithm AlgorithmType) *Context {
	ctx := C.chromaprint_new(C.int(algorithm))
	return (*Context)(ctx)
}

func (this *Context) Free() {
	ctx := (*C.ChromaprintContext)(this)
	C.chromaprint_free(ctx)
}

func (this *Context) Start(rate, channels int) error {
	ctx := (*C.ChromaprintContext)(this)
	if res := C.chromaprint_start(ctx, C.int(rate), C.int(channels)); res < 1 {
		return errStart
	}
	return nil
}

func (this *Context) Feed(data []byte) error {
	ctx := (*C.ChromaprintContext)(this)
	ptr := (*C.int16_t)(unsafe.Pointer(&data[0]))
	sz := C.int(len(data) >> 1)
	if res := C.chromaprint_feed(ctx, ptr, sz); res < 1 {
		return errFeed
	}
	return nil
}

func (this *Context) Finish() error {
	ctx := (*C.ChromaprintContext)(this)
	if res := C.chromaprint_finish(ctx); res < 1 {
		return errFinish
	}
	return nil
}

func (this *Context) GetFingerprint() (string, error) {
	ctx := (*C.ChromaprintContext)(this)
	var ptr (*C.char)
	if res := C.chromaprint_get_fingerprint(ctx, &ptr); res < 1 {
		return "", errFingerprint
	}

	defer C.chromaprint_dealloc(unsafe.Pointer(ptr))
	return C.GoString((*C.char)(ptr)), nil
}

////////////////////////////////////////////////////////////////////////////////
// ERROR

func (e Error) Error() string {
	switch e {
	default:
		return "Unknown Error"
	}
}
