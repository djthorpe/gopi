// +build chromaprint

package chromaprint

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libchromaprint
#include <chromaprint.h>
*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

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

func (this *Context) Feed(data []int16) error {
	ctx := (*C.ChromaprintContext)(this)
	ptr := (*C.int16_t)(unsafe.Pointer(&data[0]))
	sz := C.int(len(data))
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

func (this *Context) Channels() int {
	ctx := (*C.ChromaprintContext)(this)
	return int(C.chromaprint_get_num_channels(ctx))
}

func (this *Context) Rate() int {
	ctx := (*C.ChromaprintContext)(this)
	return int(C.chromaprint_get_sample_rate(ctx))
}

/* Function not exported
func (this *Context) Algorithm() AlgorithmType {
	ctx := (*C.ChromaprintContext)(this)
	return AlgorithmType(C.chromaprint_get_algorithm(ctx))
}
*/

func (this *Context) Duration() int {
	ctx := (*C.ChromaprintContext)(this)
	return int(C.chromaprint_get_item_duration(ctx))
}

func (this *Context) DurationMs() time.Duration {
	ctx := (*C.ChromaprintContext)(this)
	return time.Duration((C.chromaprint_get_item_duration_ms(ctx))) * time.Millisecond
}

func (this *Context) Delay() int {
	ctx := (*C.ChromaprintContext)(this)
	return int(C.chromaprint_get_delay(ctx))
}

func (this *Context) DelayMs() time.Duration {
	ctx := (*C.ChromaprintContext)(this)
	return time.Duration(C.chromaprint_get_delay_ms(ctx)) * time.Millisecond
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
// STRINGIFY

func (this *Context) String() string {
	str := "<chromaprint.context"
	/*
		if a := this.Algorithm(); a >= 0 {
			str += " algorithm =" + fmt.Sprint(a)
		}
	*/
	if r := this.Rate(); r > 0 {
		str += " sample_rate=" + fmt.Sprint(r)
	}
	if ch := this.Channels(); ch > 0 {
		str += " channels=" + fmt.Sprint(ch)
	}
	if d := this.DurationMs(); d > 0 {
		str += " duration=" + fmt.Sprint(d)
	}
	if d := this.DelayMs(); d > 0 {
		str += " delay=" + fmt.Sprint(d)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e Error) Error() string {
	switch e {
	case errStart:
		return "Chromaprint Start() error"
	case errFeed:
		return "Chromaprint Feed() error"
	case errFinish:
		return "Chromaprint Finish() error"
	case errFingerprint:
		return "Chromaprint Fingerprinting error"
	default:
		return "Unknown Error"
	}
}

func (a AlgorithmType) String() string {
	switch a {
	case ALGORITHM_TEST1:
		return "ALGORITHM_TEST1"
	case ALGORITHM_TEST2:
		return "ALGORITHM_TEST2"
	case ALGORITHM_TEST3:
		return "ALGORITHM_TEST3"
	case ALGORITHM_TEST4:
		return "ALGORITHM_TEST4"
	default:
		return "[?? Invalid AlgorithmType value]"
	}
}
