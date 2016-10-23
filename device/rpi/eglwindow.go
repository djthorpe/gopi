/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"errors"
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
	config  eglConfig
	context eglContext
	surface eglSurface
	element *DXElement
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EGL_LAYER_BG        uint16 = 0x0000
	EGL_LAYER_MIN       uint16 = 0x0001
	EGL_LAYER_MAX       uint16 = 0xFFFE
	EGL_WINDOW_SIZE_MAX uint32 = 0xFFFF
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS: Windows

// Create a background
func (this *eglDriver) CreateBackground(api string) (khronos.EGLWindow, error) {
	this.log.Debug2("<rpi.EGL>CreateBackground api=%v", api)

	frame := this.GetFrame()

	return this.createWindow(api, khronos.EGLSize{frame.Width, frame.Height}, khronos.EGLPoint{frame.X, frame.Y}, EGL_LAYER_BG)
}

// Create a window
func (this *eglDriver) CreateWindow(api string, size khronos.EGLSize, origin khronos.EGLPoint, layer uint16) (khronos.EGLWindow, error) {
	this.log.Debug2("<rpi.EGL>CreateWindow api=%v size=%v origin=%v layer=%v", api, size, origin, layer)

	// Check layer is not background or topmost (which will be for the pointer)
	if layer < EGL_LAYER_MIN || layer > EGL_LAYER_MAX {
		return nil, errors.New("Invalid layer parameter")
	}

	window, err := this.createWindow(api, size, origin, layer)
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

// Flush window contents to screen
func (this *eglDriver) FlushWindow(window khronos.EGLWindow) error {
	return this.swapWindowBuffer(window.(*eglWindow))
}

// Human-readble string for window
func (window *eglWindow) String() string {
	return fmt.Sprintf("<rpi.EGLWindow>{ config=%v context=%v surface=%v element=%v }", window.config, window.context, window.surface, window.element)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS: Windows

func (this *eglDriver) createWindow(api string, size khronos.EGLSize, origin khronos.EGLPoint, layer uint16) (*eglWindow, error) {
	var err error

	this.log.Debug2("<rpi.EGL>createWindow api=%v size=%v origin=%v", api, size, origin)

	// CREATE WINDOW
	window := new(eglWindow)

	// CHECK SIZE PARAMETERS
	if uint32(size.Width) > EGL_WINDOW_SIZE_MAX || uint32(size.Height) > EGL_WINDOW_SIZE_MAX {
		this.closeWindow(window)
		return nil, EGLErrorInvalidParameter
	}

	// CREATE CONTEXT
	window.config, window.context, err = this.createContext(api)
	if err != nil {
		this.closeWindow(window)
		return nil, err
	}

	// CREATE SCREEN ELEMENT
	update, err := this.dx.UpdateBegin()
	if err != nil {
		this.closeWindow(window)
		return nil, err
	}

	source_frame := &DXFrame{}
	window_frame := &DXFrame{}
	source_frame.Set(DXPoint{0, 0}, DXSize{uint32(size.Width) << 16, uint32(size.Height) << 16})
	window_frame.Set(DXPoint{int32(origin.X), int32(origin.Y)}, DXSize{uint32(size.Width), uint32(size.Height)})
	window.element, err = this.dx.AddElement(update, layer, window_frame, source_frame)
	if err != nil {
		this.closeWindow(window)
		return nil, err
	}
	if err := this.dx.UpdateSubmit(update); err != nil {
		this.closeWindow(window)
		return nil, err
	}

	// CREATE SURFACE
	nativewindow := &eglNativeWindow{window.element.GetHandle(), int(window_frame.Width), int(window_frame.Height)}
	window.surface, err = this.createSurface(window.config, nativewindow)
	if err != nil {
		this.closeWindow(window)
		return nil, err
	}

	// Attach context to surface
	if err := this.attachContextToSurface(window.context, window.surface); err != nil {
		this.destroyContext(window.context)
		return nil, err
	}

	// Success
	return window, nil
}

func (this *eglDriver) closeWindow(window *eglWindow) error {

	// remove element
	if window.element != nil {
		update, err := this.dx.UpdateBegin()
		if err != nil {
			return err
		}
		if err := this.dx.RemoveElement(update, window.element); err != nil {
			return err
		}
		if err := this.dx.UpdateSubmit(update); err != nil {
			return err
		}
	}
	window.element = nil

	// remove surface
	if window.surface != EGL_NO_SURFACE {
		if err := this.destroySurface(window.surface); err != nil {
			return err
		}
	}
	window.surface = EGL_NO_SURFACE

	// remove context
	if window.context != EGL_NO_CONTEXT {
		if err := this.destroyContext(window.context); err != nil {
			return err
		}
	}
	window.context = EGL_NO_CONTEXT

	// return success
	return nil
}

func (this *eglDriver) swapWindowBuffer(window *eglWindow) error {
	return this.swapBuffer(window.surface)
}


