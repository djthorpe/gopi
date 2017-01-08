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

	// Path methods
	CreatePath() (VGPath, error)
	DestroyPath(VGPath) error

	// Paint methods
	CreatePaint(color VGColor) (VGPaint, error)
	DestroyPaint(VGPaint) error

	// Clear surface to color
	Clear(color VGColor) error

	// Return point on screen
	GetPoint(flags EGLFrameAlignFlag) VGPoint
}

// Color with Alpha value
type VGColor struct {
	R, G, B, A float32
}

// Drawing Path
type VGPath interface {
	// Draw the path with both stroke and fill
	Draw(stroke, fill VGPaint) error

	// Stroke the path
	Stroke(stroke VGPaint) error

	// Fill the path
	Fill(fill VGPaint) error

	// Reset to empty path
	Clear() error

	// Append a line to the path
	Line(start, end VGPoint) error

	// Append a rectangle to the path
	Rect(origin, size VGPoint) error

	// Append an ellipse to the path
	Ellipse(origin, diameter VGPoint) error

	// Append a circle to the path
	Circle(origin VGPoint, diameter float32) error

}

// Point
type VGPoint struct {
	X, Y float32
}

// Paint Brush for Fill and Stroke
type VGPaint interface {
	// Set color
	SetColor(color VGColor) error

	// Set stroke line width
	SetLineWidth(width float32) error
}

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
