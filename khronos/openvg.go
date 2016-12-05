/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// OPENVG
//
// This package defines abstract interface for drawing vector graphics
// on an EGL surface
//
package khronos /* import "github.com/djthorpe/gopi/khronos" */

import (
	"os"
)

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract driver interface
type VGDriver interface {
	// Inherit general driver interface
	gopi.Driver

	// Start drawing
	Begin(surface EGLSurface) error

	// Flush
	Flush() error

	// Clear window to color
	Clear(color VGColor) error

	// Draw a line from one point to another
	Line(a VGPoint, b VGPoint) error
}

// Color with Alpha value
type VGColor struct {
	R, G, B, A float32
}

// Point
type VGPoint struct {
	X, Y float32
}

// Drawing Path
type VGPath uint64

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

// Standard Colors
var (
	VGColorRed       = VGColor{1.0, 0.0, 0.0, 1.0}
	VGColorGreen     = VGColor{0.0, 1.0, 0.0, 1.0}
	VGColorBlue      = VGColor{0.0, 0.0, 1.0, 1.0}
	VGColorWhite     = VGColor{1.0, 1.0, 1.0, 1.0}
	VGColorBlack     = VGColor{0.0, 0.0, 0.0, 1.0}
	VGColorPurple    = VGColor{1.0, 0.0, 1.0, 1.0}
	VGColorCyan      = VGColor{0.0, 1.0, 1.0, 1.0}
	VGColorYellow    = VGColor{1.0, 1.0, 0.0, 1.0}
	VGColorDarkGrey  = VGColor{0.25, 0.25, 0.25, 1.0}
	VGColorLightGrey = VGColor{0.75, 0.75, 0.75, 1.0}
	VGColorMidGrey   = VGColor{0.5, 0.5, 0.5, 1.0}
)
