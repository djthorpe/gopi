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

	// Start drawing, lock if already drawing
	Begin(surface EGLSurface) error

	// Start drawing, return error if already drawing
	//BeginNoWait(surface EGLSurface) error

	// Flush drawing
	Flush() error

	// Create a path
	CreatePath() (VGPath, error)

	// Destroy a path
	DestroyPath(VGPath) error

	// Create a paintbrush
	CreatePaint(color VGColor) (VGPaint, error)

	// Destroy a paintbrush
	DestroyPaint(VGPaint) error

	// Clear surface to color
	Clear(color VGColor) error

	// Return point on screen
	GetPoint(flags EGLFrameAlignFlag) VGPoint
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

// Paintbrush for Fill and Stroke
type VGPaint interface {
	// Set color
	SetColor(color VGColor) error

	// Set stroke width
	SetStrokeWidth(width float32) error

	// Set stroke path endpoint styles (for joins and cap)
	SetStrokeStyle(VGStrokeJoinStyle,VGStrokeCapStyle) error
}

// Color with Alpha value
type VGColor struct {
	R, G, B, A float32
}

// Point
type VGPoint struct {
	X, Y float32
}

// Stroke styles
type VGStrokeCapStyle uint16
type VGStrokeJoinStyle uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_STYLE_CAP_NONE   VGStrokeCapStyle = 0x0000
	VG_STYLE_CAP_BUTT   VGStrokeCapStyle = 0x1700
	VG_STYLE_CAP_ROUND  VGStrokeCapStyle = 0x1701
	VG_STYLE_CAP_SQUARE VGStrokeCapStyle = 0x1702
)

const (
	VG_STYLE_JOIN_NONE  VGStrokeJoinStyle = 0x0000
	VG_STYLE_JOIN_MITER VGStrokeJoinStyle = 0x1800
	VG_STYLE_JOIN_ROUND VGStrokeJoinStyle = 0x1801
	VG_STYLE_JOIN_BEVEL VGStrokeJoinStyle = 0x1802
)

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
