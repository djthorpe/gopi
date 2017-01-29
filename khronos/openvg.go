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

	// Atomic drawing operation
	Do(surface EGLSurface, callback func() error) error

	// Create a path
	CreatePath() (VGPath, error)

	// Destroy a path
	DestroyPath(VGPath) error

	// Create a paintbrush
	CreatePaint(color VGColor) (VGPaint, error)

	// Destroy a paintbrush
	DestroyPaint(VGPaint) error

	// Clear surface to color
	Clear(surface EGLSurface, color VGColor) error

	// Translate co-ordinate system
	Translate(offset VGPoint) error

	// Scale co-ordinate system
	Scale(x, y float32) error

	// Shear co-ordinate system
	Shear(x, y float32) error

	// Rotate co-ordinate system
	Rotate(r float32) error

	// Load Identity Matrix (reset co-ordinate system)
	LoadIdentity() error
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

	// Close Path
	Close() error

	// Move To
	MoveTo(VGPoint) error

	// Line To
	LineTo(...VGPoint) error

	// Quad To
	QuadTo(p1, p2 VGPoint) error

	// Cubic To
	CubicTo(p1, p2, p3 VGPoint) error

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

	// Set fill rule
	SetFillRule(style VGFillRule) error

	// Set stroke width
	SetStrokeWidth(width float32) error

	// Set miter limit
	SetMiterLimit(value float32) error

	// Set stroke path endpoint styles (for joins and cap). Use
	// NONE as a value for no change
	SetStrokeStyle(VGStrokeJoinStyle, VGStrokeCapStyle) error

	// Set stroke dash pattern, call with no arguments to reset
	SetStrokeDash(...float32) error
}

// Color with Alpha value
type VGColor struct {
	R, G, B, A float32
}

// Point
type VGPoint struct {
	X, Y float32
}

// Stroke styles and Fill rule
type VGStrokeCapStyle uint16
type VGStrokeJoinStyle uint16
type VGFillRule uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_STYLE_CAP_NONE   VGStrokeCapStyle = 0x0000
	VG_STYLE_CAP_BUTT   VGStrokeCapStyle = 0x1700 // Default
	VG_STYLE_CAP_ROUND  VGStrokeCapStyle = 0x1701
	VG_STYLE_CAP_SQUARE VGStrokeCapStyle = 0x1702
)

const (
	VG_STYLE_JOIN_NONE  VGStrokeJoinStyle = 0x0000
	VG_STYLE_JOIN_MITER VGStrokeJoinStyle = 0x1800 // Default
	VG_STYLE_JOIN_ROUND VGStrokeJoinStyle = 0x1801
	VG_STYLE_JOIN_BEVEL VGStrokeJoinStyle = 0x1802
)

const (
	VG_STYLE_FILL_NONE    VGFillRule = 0x0000
	VG_STYLE_FILL_NONZERO VGFillRule = 0x1900
	VG_STYLE_FILL_EVENODD VGFillRule = 0x1901 // Default
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

// Zero Point
var (
	VGZeroPoint      = VGPoint{ 0, 0 }
)

////////////////////////////////////////////////////////////////////////////////
// FUNCTIONS

// Return point aligned to surface
func AlignPoint(surface EGLSurface, flags EGLFrameAlignFlag) VGPoint {
	var pt VGPoint

	size := surface.GetSize()

	switch { /* X */
	case flags&EGL_ALIGN_HCENTER != 0:
		pt.X = float32(size.Width >> 1)
	case flags&EGL_ALIGN_LEFT != 0:
		pt.X = 0
	case flags&EGL_ALIGN_RIGHT != 0:
		pt.X = float32(size.Width - 1)
	}
	switch { /* Y */
	case flags&EGL_ALIGN_VCENTER != 0:
		pt.Y = float32(size.Height >> 1)
	case flags&EGL_ALIGN_TOP != 0:
		pt.Y = 0
	case flags&EGL_ALIGN_BOTTOM != 0:
		pt.Y = float32(size.Height - 1)
	}
	return pt
}
