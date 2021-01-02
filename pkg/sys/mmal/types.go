//+build mmal

package mmal

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_util_rational.h>
*/
import "C"
import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	MMALRect     C.MMAL_RECT_T
	MMALRational C.MMAL_RATIONAL_T
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_BOOL_FALSE = C.MMAL_BOOL_T(0)
	MMAL_BOOL_TRUE  = C.MMAL_BOOL_T(1)
)

////////////////////////////////////////////////////////////////////////////////
// RATIONAL

func NewRational(value float32) MMALRational {
	q16 := C.int32_t(value * float32(1<<16))
	return MMALRational(C.mmal_rational_from_fixed_16_16(q16))
}

func (a MMALRational) Add(b MMALRational) MMALRational {
	return MMALRational(C.mmal_rational_add(C.MMAL_RATIONAL_T(a), C.MMAL_RATIONAL_T(b)))
}

func (a MMALRational) Subtract(b MMALRational) MMALRational {
	return MMALRational(C.mmal_rational_subtract(C.MMAL_RATIONAL_T(a), C.MMAL_RATIONAL_T(b)))
}

func (a MMALRational) Divide(b MMALRational) MMALRational {
	return MMALRational(C.mmal_rational_divide(C.MMAL_RATIONAL_T(a), C.MMAL_RATIONAL_T(b)))
}

func (a MMALRational) Multiply(b MMALRational) MMALRational {
	return MMALRational(C.mmal_rational_multiply(C.MMAL_RATIONAL_T(a), C.MMAL_RATIONAL_T(b)))
}

func (a MMALRational) Equals(b MMALRational) bool {
	return mmal_to_bool(C.mmal_rational_equal(C.MMAL_RATIONAL_T(a), C.MMAL_RATIONAL_T(b)))
}

func (a MMALRational) Float32() float32 {
	return float32(a.num) / float32(a.den)
}

func (a MMALRational) IsZero() bool {
	return a.num == 0
}

func (a MMALRational) String() string {
	return fmt.Sprintf("<%d/%d>", int32(a.num), int32(a.den))
}

////////////////////////////////////////////////////////////////////////////////
// RECT

func NewRect(x, y, w, h int32) MMALRect {
	rect := C.MMAL_RECT_T{
		C.int32_t(x),
		C.int32_t(y),
		C.int32_t(w),
		C.int32_t(h),
	}
	return MMALRect(rect)
}

func (r MMALRect) Origin() (int32, int32) {
	ctx := (C.MMAL_RECT_T)(r)
	return int32(ctx.x), int32(ctx.y)
}

func (r MMALRect) Size() (int32, int32) {
	ctx := (C.MMAL_RECT_T)(r)
	return int32(ctx.width), int32(ctx.height)
}

func (r MMALRect) IsZero() bool {
	ctx := (C.MMAL_RECT_T)(r)
	return ctx.x == 0 && ctx.y == 0 && ctx.width == 0 && ctx.height == 0
}

func (r MMALRect) String() string {
	x, y := r.Origin()
	w, h := r.Size()
	return fmt.Sprintf("<mmal.rect origin={ %d,%d } size={%d,%d}>", x, y, w, h)
}

////////////////////////////////////////////////////////////////////////////////
// BOOL

func mmal_from_bool(value bool) C.MMAL_BOOL_T {
	if value {
		return MMAL_BOOL_TRUE
	} else {
		return MMAL_BOOL_FALSE
	}
}

func mmal_to_bool(value C.MMAL_BOOL_T) bool {
	return value != MMAL_BOOL_FALSE
}
