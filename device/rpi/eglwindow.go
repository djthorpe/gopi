/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
)

import (
	khronos "../../khronos" /* import "github.com/djthorpe/gopi/khronos" */
)

////////////////////////////////////////////////////////////////////////////////

/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2
  #include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Native window structre
type eglNativeWindow struct {
	element dxElementHandle
	width   int
	height  int
}

// Internal window structure
type eglWindow struct {
	config eglConfig
	context eglContext
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EGL_LAYER_BG  uint16 = 0x0000
	EGL_LAYER_MIN uint16 = 0x0001
	EGL_LAYER_MAX uint16 = 0xFFFE
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS: Windows

// Create a background
func (this *eglDriver) CreateBackground(api string) (khronos.EGLWindow,error) {
	this.log.Debug2("<rpi.EGL>CreateBackground api=%v",api)

	frame := this.GetFrame()

	return this.createWindow(api,khronos.EGLSize{ frame.Width, frame.Height },khronos.EGLPoint{ frame.X, frame.Y },EGL_LAYER_BG)
}

// Create a window
func (this *eglDriver) CreateWindow(api string,size khronos.EGLSize,origin khronos.EGLPoint,layer uint16) (khronos.EGLWindow,error) {
	this.log.Debug2("<rpi.EGL>CreateWindow api=%v size=%v origin=%v layer=%v",api,size,origin,layer)

	// Check layer is not background or topmost (which will be for the pointer)
	if layer < EGL_LAYER_MIN || layer > EGL_LAYER_MAX {
		return nil, errors.New("Invalid layer parameter")
	}

	window, err := this.createWindow(api,size,origin,layer)
	if err != nil {
		return nil, err
	}

	// success
	return khronos.EGLWindow(window), nil
}

// Close a window
func (this *eglDriver) CloseWindow(window khronos.EGLWindow) error {
	this.log.Debug2("<rpi.EGL>CloseWindow")
	return this.closeWindow(window.(*eglWindow))
}

// Human-readble string for window
func (window eglWindow) String() string {
	return fmt.Sprintf("<rpi.EGLWindow>{ config=%08X context=%08X }",window.config,window.context)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS: Windows

func (this *eglDriver) createWindow(api string,size khronos.EGLSize,origin khronos.EGLPoint,layer uint16) (*eglWindow,error) {
	var err error

	this.log.Debug2("<rpi.EGL>createWindow api=%v size=%v origin=%v",api,size,origin)

	// CREATE WINDOW
	window := new(eglWindow)

	// CREATE CONTEXT
	window.config, window.context, err = this.createContext(api)
	if err != nil {
		return nil,err
	}

	// CREATE SCREEN ELEMENT
	update, err := this.dx.UpdateBegin()
	if err != nil {
		return nil,err
	}
	// CREATE ELEMENT
	source_frame := &Rectangle{}
	source_frame.Set(Point{  0, 0 }, Size{ frame.Size.Width << 16, frame.Size.Height << 16})
	window.element, err = this.vc.AddElement(update, 0, frame, nil, source_frame)
	if err != nil {
		this.destroyContext(window.context)
		return nil,err
	}

	if err := this.dx.UpdateSubmit(update); err != nil {
		this.destroyContext(window.context)
		return nil,err
	}

	// Success
	return window,nil
}

func (this *eglDriver) closeWindow(window *eglWindow) error {
	if window.context == EGL_NO_CONTEXT {
		if err := this.destroyContext(window.context); err != nil {
			return err
		}
	}
	window.context = EGL_NO_CONTEXT
	return nil
}

