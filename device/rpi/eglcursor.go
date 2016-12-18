/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// EGLCURSOR
//
// Implements a cursor sprite on the topmost layer which can be
// moved and changed
//
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Create a cursor
func (this *eglDriver) CreateCursor() (khronos.EGLSurface, error) {
	this.log.Debug2("<rpi.EGL>CreateCursor")

	// Place cursor in the middle of the screen
	screen := this.GetFrame()
	frame := khronos.EGLFrame{khronos.EGLPoint{0, 0}, khronos.EGLSize{32, 32}}.AlignTo(&screen, khronos.EGL_ALIGN_VCENTER|khronos.EGL_ALIGN_HCENTER)

	// Create the surface
	surface, err := this.createWindow("DX", frame.Size(), frame.Origin(), EGL_LAYER_CURSOR, 1.0, nil)
	if err != nil {
		return nil, err
	}

	// Temporarily clear to red
	bitmap, err := surface.GetBitmap()
	if err != nil {
		this.DestroySurface(surface)
		return nil, err
	}

	bitmap.ClearToColor(khronos.EGLRedColor)

	// Return the cursor surface
	return surface, nil
}
