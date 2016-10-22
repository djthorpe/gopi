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

type eglNativeWindow struct {
	element ElementHandle
	width   int
	height  int
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Create a window
func (this *EGLState) CreateWindow(api string,size *khronos.EGLSize,origin *khronos.EGLPoint) (khronos.EGLWindow,error) {
	window := new(eglWindow)

	return window,errors.New("CreateWindow: NOT IMPLEMENTED")
}

// Create Background
func (this *EGLState) CreateBackground(api string) (khronos.EGLWindow,error) {
	window := new(eglWindow)

	return window,errors.New("CreateBackground: NOT IMPLEMENTED")
}

// Close a window
func (this *EGLState) CloseWindow(window khronos.EGLWindow) error {
	return errors.New("CloseWindow: NOT IMPLEMENTED")
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

/*
func (this *EGLState) createWindow()
	var err error

	// CREATE WINDOW
	window := new(eglWindow)

	// CREATE CONTEXT
	window.config, window.context, err = this.createContext(api)
	if err != nil {
		return nil,err
	}

	// CREATE SCREEN ELEMENT
	update, err := this.vc.UpdateBegin()
	if err != nil {
		return nil,err
	}
	source_frame := &Rectangle{}
	source_frame.Set(Point{  0, 0 }, Size{ frame.Size.Width << 16, frame.Size.Height << 16})
	window.element, err = this.vc.AddElement(update, 0, frame, nil, source_frame)
	if err != nil {
		this.destroyContext(window.context)
		return nil,err
	}
	if err := this.vc.UpdateSubmit(update); err != nil {
		this.destroyContext(window.context)
		return nil,err
	}

	// CREATE SURFACE
	nativewindow := &eglNativeWindow{ window.element.GetHandle(), int(frame.Size.Width), int(frame.Size.Height)}
	window.surface, err = this.createSurface(window.config, nativewindow)
	if err != nil {
		this.destroyContext(window.context)
		return nil,err
	}

	// Attach context to surface
	if err := this.attachContextToSurface(window.context, window.surface); err != nil {
		this.destroyContext(window.context)
		return nil,err
	}

	// Success
	return window,nil
*/

