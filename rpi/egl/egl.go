/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package egl

/*
	#cgo CFLAGS:   -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
	#cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2
	#include "EGL/egl.h"
*/
import "C"

import (
	"unsafe"
)

type (
	Display		uintptr
)

func Initialize(disp Display, major, minor *int32) (error) {
	if C.eglInitialize(C.EGLDisplay(unsafe.Pointer(disp)),(*C.EGLint)(major),(*C.EGLint)(minor)) != EGL_FALSE {
		return nil
	}
	return toError(GetError())
}

func Terminate(disp Display) (error) {
	if C.eglTerminate(C.EGLDisplay(unsafe.Pointer(disp))) != EGL_FALSE {
		return nil
	}
	return toError(GetError())
}

func GetDisplay() Display {
	return Display(C.eglGetDisplay(C.EGLNativeDisplayType(unsafe.Pointer(nil))))
}

func GetError() int32 {
	return int32(C.eglGetError())
}
