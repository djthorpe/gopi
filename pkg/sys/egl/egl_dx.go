// +build egl,dispmanx

package egl

// Frameworks

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
  #cgo pkg-config: brcmegl
  #include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func EGLGetDisplay(display uint) EGLDisplay {
	return EGLDisplay(C.eglGetDisplay(C.EGLNativeDisplayType(uintptr(display))))
}
