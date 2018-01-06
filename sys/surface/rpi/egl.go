// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
  #cgo CFLAGS:   -I/opt/vc/include -DUSE_VCHIQ_ARM
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL_static -lGLESv2_static -lkhrn_static -lvcos -lvchiq_arm -lbcm_host -lm
  #include <EGL/egl.h>
*/
import "C"
import "unsafe"

////////////////////////////////////////////////////////////////////////////////
// EGL Types

type (
	eglDisplay           C.EGLDisplay
	eglNativeDisplayType C.EGLNativeDisplayType
	eglBoolean           C.EGLBoolean
	eglInt               C.EGLint
)

////////////////////////////////////////////////////////////////////////////////
// EGL Errors

const (
	EGL_FALSE = eglBoolean(0)
	EGL_TRUE  = eglBoolean(1)
)

const (
	EGL_SUCCESS             eglInt = 0x3000
	EGL_NOT_INITIALIZED            = 0x3001
	EGL_BAD_ACCESS                 = 0x3002
	EGL_BAD_ALLOC                  = 0x3003
	EGL_BAD_ATTRIBUTE              = 0x3004
	EGL_BAD_CONFIG                 = 0x3005
	EGL_BAD_CONTEXT                = 0x3006
	EGL_BAD_CURRENT_SURFACE        = 0x3007
	EGL_BAD_DISPLAY                = 0x3008
	EGL_BAD_MATCH                  = 0x3009
	EGL_BAD_NATIVE_PIXMAP          = 0x300A
	EGL_BAD_NATIVE_WINDOW          = 0x300B
	EGL_BAD_PARAMETER              = 0x300C
	EGL_BAD_SURFACE                = 0x300D
	EGL_CONTEXT_LOST               = 0x300E /* EGL 1.1 - IMG_power_management */
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func to_eglDisplay(display uint) eglDisplay {
	return eglDisplay(uintptr(display))
}

////////////////////////////////////////////////////////////////////////////////
// EGL Methods

func eglGetError() eglInt {
	return eglInt(C.eglGetError())
}

func eglInitialize(display eglDisplay) (eglInt, eglInt, error) {
	var major, minor C.EGLint
	if C.eglInitialize(C.EGLDisplay(display), (*C.EGLint)(unsafe.Pointer(&major)), (*C.EGLint)(unsafe.Pointer(&minor))) != C.EGLBoolean(EGL_TRUE) {
		return 0, 0, eglGetError()
	} else {
		return eglInt(major), eglInt(minor), EGL_SUCCESS
	}
}

func eglTerminate(display eglDisplay) error {
	if C.eglTerminate(C.EGLDisplay(display)) != C.EGLBoolean(EGL_TRUE) {
		return eglGetError()
	} else {
		return EGL_SUCCESS
	}
}

func eglGetDisplay(display_id eglNativeDisplayType) eglDisplay {
	return eglDisplay(C.eglGetDisplay(C.EGLNativeDisplayType(display_id)))
}

////////////////////////////////////////////////////////////////////////////////
// Stringify

func (e eglInt) Error() string {
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
		return "Unknown EGL error"
	}
}
