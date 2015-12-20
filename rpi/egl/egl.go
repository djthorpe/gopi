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

func Initialize(disp Display, major, minor *int32) (bool, error) {
	result := C.eglInitialize(C.EGLDisplay(unsafe.Pointer(disp)),(*C.EGLint)(major),(*C.EGLint)(minor))
	if result != EGL_FALSE {
		return true,nil
	} else {
		return false,toError(GetError())
	}
}

func Terminate(disp Display) (bool, error) {
	success := C.eglTerminate(C.EGLDisplay(unsafe.Pointer(disp)))
	if success != 0 {
		return true,nil
	} else {
		return false,toError(GetError())
	}
}

func GetDisplay() Display {
	return Display(C.eglGetDisplay(C.EGLNativeDisplayType(unsafe.Pointer(nil))))
}

func GetError() int32 {
	return int32(C.eglGetError())
}
