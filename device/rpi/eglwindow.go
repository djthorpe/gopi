/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package rpi
/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2
  #include <EGL/egl.h>
  #include <VG/openvg.h>
*/
import "C"

import (
	"errors"
)

import (
	"../../khronos"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type eglWindow struct {
	config EGLConfig
	context EGLContext
	surface EGLSurface
	element *Element
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Create a window
func (this *EGLState) CreateWindow(api string,origin *khronos.EGLPoint,size *khronos.EGLSize) (khronos.EGLWindow,error) {
	window := new(eglWindow)

	return window,errors.New("CreateWindow: NOT IMPLEMENTED")
}

// Close a window
func (this *EGLState) CloseWindow(window *eglWindow) error {
	return errors.New("CloseWindow: NOT IMPLEMENTED")
}
