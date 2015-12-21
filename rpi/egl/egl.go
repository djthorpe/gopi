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
    "github.com/djthorpe/gopi/rpi/dispmanx"
)

type (
	Display uintptr
	Config  uintptr
	Context uintptr
	Surface uintptr
	Enum    uint32
)

////////////////////////////////////////////////////////////////////////////////

func Initialize(disp Display, major, minor *int32) error {
	if C.eglInitialize(C.EGLDisplay(unsafe.Pointer(disp)), (*C.EGLint)(major), (*C.EGLint)(minor)) == EGL_TRUE {
		return nil
	}
	return toError(GetLastError())
}

func Terminate(disp Display) error {
	if C.eglTerminate(C.EGLDisplay(unsafe.Pointer(disp))) == EGL_TRUE {
		return nil
	}
	return toError(GetLastError())
}

func GetDisplay() Display {
	return Display(C.eglGetDisplay(C.EGLNativeDisplayType(unsafe.Pointer(nil))))
}

////////////////////////////////////////////////////////////////////////////////

func getConfigs(disp Display, configs *Config, configSize int32, numConfig *int32) error {
	if C.eglGetConfigs(C.EGLDisplay(unsafe.Pointer(disp)), (*C.EGLConfig)(unsafe.Pointer(configs)), C.EGLint(configSize), (*C.EGLint)(unsafe.Pointer(numConfig))) == EGL_TRUE {
		return nil
	}
	return toError(GetLastError())
}

func chooseConfig(disp Display, attribList []int32, configs *Config, configSize int32, numConfig *int32) error {
	if C.eglChooseConfig(C.EGLDisplay(unsafe.Pointer(disp)), (*C.EGLint)(&attribList[0]), (*C.EGLConfig)(unsafe.Pointer(configs)), C.EGLint(configSize), (*C.EGLint)(numConfig)) == EGL_TRUE {
		return nil
	}
	return toError(GetLastError())
}

func GetNumConfigs(disp Display) (int32, error) {
	var numConfigs int32
	if err := getConfigs(disp, nil, 0, &numConfigs); err != nil {
		return -1, err
	}
	return numConfigs, nil
}

func GetConfigs(disp Display) ([]Config, error) {
	numConfigs, err := GetNumConfigs(disp)
	if err != nil {
		return nil, err
	}

	var configs = make([]Config, numConfigs)
	if err := getConfigs(disp, &configs[0], numConfigs, &numConfigs); err != nil {
		return nil, err
	}
	return configs, nil
}

func GetConfigAttribute(disp Display, config Config, attribute int32) (int32, error) {
	var value int32
	if C.eglGetConfigAttrib(C.EGLDisplay(unsafe.Pointer(disp)), C.EGLConfig(config), C.EGLint(attribute), (*C.EGLint)(unsafe.Pointer(&value))) == EGL_TRUE {
		return value, nil
	}
	return -1, toError(GetLastError())
}

func ChooseConfig(disp Display, attribList []int32) ([]Config, error) {
	numConfigs, err := GetNumConfigs(disp)
	if err != nil {
		return nil, err
	}
	var configs = make([]Config, numConfigs)
	if err := chooseConfig(disp, attribList, &configs[0], numConfigs, &numConfigs); err != nil {
		return nil, err
	}
	// reduce size of slice, return it
	configs = append([]Config(nil), configs[:numConfigs]...)
	return configs, nil
}

////////////////////////////////////////////////////////////////////////////////

func CreateContext(disp Display, config Config, shareContext Context, attribList *int32) (Context, error) {
	context := Context(C.eglCreateContext(C.EGLDisplay(unsafe.Pointer(disp)), C.EGLConfig(unsafe.Pointer(config)), C.EGLContext(unsafe.Pointer(shareContext)), (*C.EGLint)(attribList)))
	if context != EGL_NO_CONTEXT {
		return context, nil
	}
	return context, toError(GetLastError())
}

func DestroyContext(disp Display, ctx Context) error {
	if C.eglDestroyContext(C.EGLDisplay(unsafe.Pointer(disp)), C.EGLContext(unsafe.Pointer(ctx))) == EGL_TRUE {
		return nil
	}
	return toError(GetLastError())
}

////////////////////////////////////////////////////////////////////////////////

func CreateWindowSurface(disp Display, config Config, window dispmanx.Window, attribList *int32) (Surface, error) {
	surface := Surface(C.eglCreateWindowSurface(C.EGLDisplay(unsafe.Pointer(disp)), C.EGLConfig(unsafe.Pointer(config)), C.EGLNativeWindowType(unsafe.Pointer(&window)), (*C.EGLint)(attribList)))
	if surface != EGL_NO_SURFACE {
		return surface, nil
	}
	return surface, toError(GetLastError())
}

func DestroySurface(disp Display, surface Surface) error {
	if C.eglDestroySurface(C.EGLDisplay(unsafe.Pointer(disp)), C.EGLSurface(unsafe.Pointer(surface))) == EGL_TRUE {
		return nil
	}
	return toError(GetLastError())
}

func MakeCurrent(disp Display, draw Surface, read Surface, ctx Context) error {
	if C.eglMakeCurrent(C.EGLDisplay(unsafe.Pointer(disp)), C.EGLSurface(unsafe.Pointer(draw)), C.EGLSurface(unsafe.Pointer(read)), C.EGLContext(unsafe.Pointer(ctx))) == EGL_TRUE {
		return nil
	}
	return toError(GetLastError())
}

////////////////////////////////////////////////////////////////////////////////

func BindAPI(api Enum) error {
	if C.eglBindAPI(C.EGLenum(api)) == EGL_TRUE {
		return nil
	}
	return toError(GetLastError())
}

func QueryAPI() Enum {
	return Enum(C.eglQueryAPI())
}

////////////////////////////////////////////////////////////////////////////////

func GetLastError() int32 {
	return int32(C.eglGetError())
}
