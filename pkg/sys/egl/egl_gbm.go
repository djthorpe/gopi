// +build egl,gbm

package egl

import (
	"unsafe"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
	gbm "github.com/djthorpe/gopi/v3/pkg/sys/gbm"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
  #cgo pkg-config: egl
  #include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func EGLGetDisplay(device *gbm.GBMDevice) EGLDisplay {
	disp := C.EGLNativeDisplayType(unsafe.Pointer(device))
	return EGLDisplay(C.eglGetDisplay(disp))
}
