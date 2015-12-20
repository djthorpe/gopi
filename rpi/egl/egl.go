/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package egl

/*
	#cgo CFLAGS:   -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
	#cgo LDFLAGS:  -L/opt/vc/lib -lGLESv2
	#include "EGL/egl.h"
*/
import "C"

type (
	Display		uintptr
)

func Initialize(disp Display, major, minor *int32) bool {
	success := C.eglInitialize(C.EGLDisplay(unsafe.Pointer(disp)),(*C.EGLint)(major),(*C.EGLint)(minor))
	return bool(success)
}

func Terminate(disp Display) bool {
	success := C.eglTerminate(C.EGLDisplay(unsafe.Pointer(disp))))
	return bool(success)
}

func GetDisplay() Display {
	return Display(C.eglGetDisplay(C.EGLNativeDisplayType(unsafe.Pointer(nil))))
}

