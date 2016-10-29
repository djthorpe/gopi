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
	origin  *khronos.EGLPoint
	size    *khronos.EGLSize
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
func (this *eglDriver) CreateBackground(api string, opacity float32) (khronos.EGLSurface, error) {
	this.log.Debug2("<rpi.EGL>CreateBackground api=%v opacity=%v", api, opacity)
	frame := this.GetFrame()
	return this.createWindow(api, khronos.EGLSize{frame.Width, frame.Height}, khronos.EGLPoint{frame.X, frame.Y}, EGL_LAYER_BG, opacity)
}

// Create a surface
func (this *eglDriver) CreateSurface(api string, size khronos.EGLSize, origin khronos.EGLPoint, layer uint16, opacity float32) (khronos.EGLSurface, error) {
	this.log.Debug2("<rpi.EGL>CreateSurface api=%v size=%v origin=%v layer=%v opacity=%v", api, size, origin, layer, opacity)

	// Check layer is not background or topmost (which will be for the pointer)
	if layer < EGL_LAYER_MIN || layer > EGL_LAYER_MAX {
		return nil, errors.New("Invalid layer parameter")
	}

	// Check opacity
	if opacity < 0.0 || opacity > 1.0 {
		return nil, errors.New("Invalid opacity parameter")
	}

	surface, err := this.createWindow(api, size, origin, layer, opacity)
	if err != nil {
		return nil, err
	}

	// success
	return khronos.EGLSurface(surface), nil
}

// Destroy a surface
func (this *eglDriver) CloseSurface(surface khronos.EGLSurface) error {
	this.log.Debug2("<rpi.EGL>CloseSurface")
	return this.closeWindow(surface.(*eglWindow))
}

// Flush surface contents to screen
func (this *eglDriver) FlushSurface(surface khronos.EGLSurface) error {
	return this.swapWindowBuffer(surface.(*eglWindow))
}

// Move surface origin relative to current origin
func (this *eglDriver) MoveSurfaceOriginBy(surface khronos.EGLSurface, rel khronos.EGLPoint) error {
	return this.setWindowOrigin(surface.(*eglWindow), surface.GetOrigin().Add(rel))
}

// Set current context
func (this *eglDriver) SetCurrentContext(surface khronos.EGLSurface) error {
	return this.setCurrentContext(surface.(*eglWindow))
}

// Human-readble string for window
func (surface *eglWindow) String() string {
	return fmt.Sprintf("<rpi.EGLSurface>{ config=%v context=%v surface=%v element=%v }", surface.config, surface.context, surface.surface, surface.element)
}

// Return window origin on screen compared to NW corner of screen
func (surface *eglWindow) GetOrigin() khronos.EGLPoint {
	return *surface.origin
}

// Return window size
func (surface *eglWindow) GetSize() khronos.EGLSize {
	return *surface.size
}

// Is surface the background surface?
func (surface *eglWindow) IsBackgroundSurface() bool {
	return surface.GetLayer() == EGL_LAYER_BG
}

// Return layer the surface is on
func (surface *eglWindow) GetLayer() uint16 {
	return surface.element.GetLayer()
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS: Windows

func (this *eglDriver) createWindow(api string, size khronos.EGLSize, origin khronos.EGLPoint, layer uint16, opacity float32) (*eglWindow, error) {
	var err error

	// CREATE WINDOW
	window := new(eglWindow)
	window.origin = &origin
	window.size = &size

	// CHECK SIZE PARAMETERS
	if uint32(size.Width) > EGL_WINDOW_SIZE_MAX || uint32(size.Height) > EGL_WINDOW_SIZE_MAX {
		this.closeWindow(window)
		return nil, this.log.Error("Invalid width or height parameters: %v", size)
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
	window.element, err = this.dx.AddElement(update, layer, uint32(opacity*255.0), window_frame, source_frame)
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
	if err := this.setCurrentContext(window); err != nil {
		this.closeWindow(window)
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

func (this *eglDriver) setWindowOrigin(window *eglWindow, new_origin khronos.EGLPoint) error {
	update, err := this.dx.UpdateBegin()
	if err != nil {
		return err
	}
	size := window.GetSize()
	frame := DXFrame{DXPoint{int32(new_origin.X), int32(new_origin.Y)}, DXSize{uint32(size.Width), uint32(size.Height)}}
	if err := this.dx.SetElementDestination(update, window.element, frame); err != nil {
		return err
	}
	window.origin = &new_origin
	if err := this.dx.UpdateSubmit(update); err != nil {
		return err
	}
	return nil
}

func (this *eglDriver) setCurrentContext(window *eglWindow) error {
	return this.makeCurrent(window.surface, window.context)
}
