// +build egl

package egl

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: egl
#include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EGLError C.EGLint
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EGL_SUCCESS             EGLError = 0x3000
	EGL_NOT_INITIALIZED     EGLError = 0x3001
	EGL_BAD_ACCESS          EGLError = 0x3002
	EGL_BAD_ALLOC           EGLError = 0x3003
	EGL_BAD_ATTRIBUTE       EGLError = 0x3004
	EGL_BAD_CONFIG          EGLError = 0x3005
	EGL_BAD_CONTEXT         EGLError = 0x3006
	EGL_BAD_CURRENT_SURFACE EGLError = 0x3007
	EGL_BAD_DISPLAY         EGLError = 0x3008
	EGL_BAD_MATCH           EGLError = 0x3009
	EGL_BAD_NATIVE_PIXMAP   EGLError = 0x300A
	EGL_BAD_NATIVE_WINDOW   EGLError = 0x300B
	EGL_BAD_PARAMETER       EGLError = 0x300C
	EGL_BAD_SURFACE         EGLError = 0x300D
	EGL_CONTEXT_LOST        EGLError = 0x300E // EGL 1.1 - IMG_power_management
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func EGLGetError() error {
	if err := EGLError(C.eglGetError()); err == EGL_SUCCESS {
		return nil
	} else {
		return err
	}
}

func (e EGLError) Error() string {
	switch e {
	case EGL_SUCCESS:
		return "EGL_SUCCESS"
	case EGL_NOT_INITIALIZED:
		return "EGL_NOT_INITIALIZED"
	case EGL_BAD_ACCESS:
		return "EGL_BAD_ACCESS"
	case EGL_BAD_ALLOC:
		return "EGL_BAD_ALLOC"
	case EGL_BAD_ATTRIBUTE:
		return "EGL_BAD_ATTRIBUTE"
	case EGL_BAD_CONFIG:
		return "EGL_BAD_CONFIG"
	case EGL_BAD_CONTEXT:
		return "EGL_BAD_CONTEXT"
	case EGL_BAD_CURRENT_SURFACE:
		return "EGL_BAD_CURRENT_SURFACE"
	case EGL_BAD_DISPLAY:
		return "EGL_BAD_DISPLAY"
	case EGL_BAD_MATCH:
		return "EGL_BAD_MATCH"
	case EGL_BAD_NATIVE_PIXMAP:
		return "EGL_BAD_NATIVE_PIXMAP"
	case EGL_BAD_NATIVE_WINDOW:
		return "EGL_BAD_NATIVE_WINDOW"
	case EGL_BAD_PARAMETER:
		return "EGL_BAD_PARAMETER"
	case EGL_BAD_SURFACE:
		return "EGL_BAD_SURFACE"
	case EGL_CONTEXT_LOST:
		return "EGL_CONTEXT_LOST"
	default:
		return "[?? Unknown EGLError value]"
	}
}
