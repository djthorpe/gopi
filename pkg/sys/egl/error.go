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
	EGL_SUCCESS             EGLError = C.EGL_SUCCESS
	EGL_NOT_INITIALIZED     EGLError = C.EGL_NOT_INITIALIZED
	EGL_BAD_ACCESS          EGLError = C.EGL_BAD_ACCESS
	EGL_BAD_ALLOC           EGLError = C.EGL_BAD_ALLOC
	EGL_BAD_ATTRIBUTE       EGLError = C.EGL_BAD_ATTRIBUTE
	EGL_BAD_CONFIG          EGLError = C.EGL_BAD_CONFIG
	EGL_BAD_CONTEXT         EGLError = C.EGL_BAD_CONTEXT
	EGL_BAD_CURRENT_SURFACE EGLError = C.EGL_BAD_CURRENT_SURFACE
	EGL_BAD_DISPLAY         EGLError = C.EGL_BAD_DISPLAY
	EGL_BAD_MATCH           EGLError = C.EGL_BAD_MATCH
	EGL_BAD_NATIVE_PIXMAP   EGLError = C.EGL_BAD_NATIVE_PIXMAP
	EGL_BAD_NATIVE_WINDOW   EGLError = C.EGL_BAD_NATIVE_WINDOW
	EGL_BAD_PARAMETER       EGLError = C.EGL_BAD_PARAMETER
	EGL_BAD_SURFACE         EGLError = C.EGL_BAD_SURFACE
	EGL_CONTEXT_LOST        EGLError = C.EGL_CONTEXT_LOST
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
