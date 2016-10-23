/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package khronos /* import "github.com/djthorpe/gopi/khronos" */

import (
	gopi ".." /* import "github.com/djthorpe/gopi" */
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract driver interface
type VGDriver interface {
	// Inherit general driver interface
	gopi.Driver

	// Start drawing
	Begin(window EGLWindow) error

	// Flush
	Flush() error

	// Clear window to color
	Clear(color VGColor)

	// Draw a line from one point to another
	Line(a VGPoint,b VGPoint)
}

// Color with Alpha value
type VGColor struct {
	R,G,B,A float32
}

// Point
type VGPoint struct {
	X,Y float32
}


