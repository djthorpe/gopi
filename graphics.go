/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// SurfaceFlags are flags associated with surface
	SurfaceFlags uint16

	// SurfaceCallback
	SurfaceCallback func(SurfaceManager) error
)

// Color including opacity
type Color struct {
	R, G, B, A float32
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// SurfaceManager allows you to open, close and move surfaces around an open display
type SurfaceManager interface {
	// Return the display associated with the surface manager
	Display() Display

	// Return the name of the surface manager
	Name() string

	// Return capabilities for the GPU
	Types() []SurfaceFlags

	// Do CreateSurface, updates and drawing within this method
	Do(SurfaceCallback) error

	// Create and destroy surfaces
	CreateSurfaceWithBitmap(Bitmap, SurfaceFlags, float32, uint16, Point, Size) (Surface, error)
	CreateSurface(SurfaceFlags, float32, uint16, Point, Size) (Surface, error)
	CreateBackground(SurfaceFlags, float32) (Surface, error)
	DestroySurface(Surface) error

	// Create and destroy bitmaps
	CreateBitmap(SurfaceFlags, Size) (Bitmap, error)
	CreateSnapshot(SurfaceFlags) (Bitmap, error)
	DestroyBitmap(Bitmap) error

	// Implements gopi.Unit
	Unit
}

// Surface defines an on-screen rectanglar surface
type Surface interface {
	Type() SurfaceFlags
	Size() Size
	Origin() Point
	Opacity() float32
	Layer() uint16
	Bitmap() Bitmap
}

// Bitmap defines a rectangular bitmap which can be used by the GPU
type Bitmap interface {
	Type() SurfaceFlags
	Size() Size

	// Clear bitmap
	ClearToColor(Color)

	// Paint a single pixel
	PaintPixel(Color, Point)
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// SurfaceFlags - surface binding
	SURFACE_FLAG_NONE       SurfaceFlags = 0x0000
	SURFACE_FLAG_BITMAP     SurfaceFlags = 0x0001 // Bitmap
	SURFACE_FLAG_OPENGL     SurfaceFlags = 0x0002
	SURFACE_FLAG_OPENGL_ES  SurfaceFlags = 0x0003
	SURFACE_FLAG_OPENGL_ES2 SurfaceFlags = 0x0004
	SURFACE_FLAG_OPENVG     SurfaceFlags = 0x0005 // 2D Vector
	SURFACE_FLAG_TYPEMASK   SurfaceFlags = 0x000F
	// SurfaceFlags - surface configuration
	SURFACE_FLAG_RGBA32     SurfaceFlags = 0x0000 // 4 bytes per pixel
	SURFACE_FLAG_RGB888     SurfaceFlags = 0x0010 // 3 bytes per pixel
	SURFACE_FLAG_RGB565     SurfaceFlags = 0x0020 // 2 bytes per pixel
	SURFACE_FLAG_CONFIGMASK SurfaceFlags = 0x00F0
	// SurfaceFlags - modifiers
	SURFACE_FLAG_ALPHA_FROM_SOURCE SurfaceFlags = 0x0100
	SURFACE_FLAG_MODMASK           SurfaceFlags = 0x0F00
)

const (
	// SurfaceLayer
	SURFACE_LAYER_BACKGROUND uint16 = 0x0000
	SURFACE_LAYER_DEFAULT    uint16 = 0x0001
	SURFACE_LAYER_MAX        uint16 = 0xFFFE
	SURFACE_LAYER_CURSOR     uint16 = 0xFFFF
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

// Standard Colors
var (
	ColorRed         = Color{1.0, 0.0, 0.0, 1.0}
	ColorGreen       = Color{0.0, 1.0, 0.0, 1.0}
	ColorBlue        = Color{0.0, 0.0, 1.0, 1.0}
	ColorWhite       = Color{1.0, 1.0, 1.0, 1.0}
	ColorBlack       = Color{0.0, 0.0, 0.0, 1.0}
	ColorPurple      = Color{1.0, 0.0, 1.0, 1.0}
	ColorCyan        = Color{0.0, 1.0, 1.0, 1.0}
	ColorYellow      = Color{1.0, 1.0, 0.0, 1.0}
	ColorDarkGrey    = Color{0.25, 0.25, 0.25, 1.0}
	ColorLightGrey   = Color{0.75, 0.75, 0.75, 1.0}
	ColorMidGrey     = Color{0.5, 0.5, 0.5, 1.0}
	ColorTransparent = Color{}
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATIONS

// RGBA returns uint32 values for color
func (c Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.R*float32(0xFFFF)) & uint32(0xFFFF), uint32(c.G*float32(0xFFFF)) & uint32(0xFFFF), uint32(c.B*float32(0xFFFF)) & uint32(0xFFFF), uint32(c.A*float32(0xFFFF)) & uint32(0xFFFF)
}

// Type returns the type of the surface
func (f SurfaceFlags) Type() SurfaceFlags {
	return f & SURFACE_FLAG_TYPEMASK
}

// Config returns the configuration of the surface
func (f SurfaceFlags) Config() SurfaceFlags {
	return f & SURFACE_FLAG_CONFIGMASK
}

// Mod returns surface modifiers
func (f SurfaceFlags) Mod() SurfaceFlags {
	return f & SURFACE_FLAG_MODMASK
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f SurfaceFlags) TypeString() string {
	switch f.Type() {
	case SURFACE_FLAG_BITMAP:
		return "SURFACE_FLAG_BITMAP"
	case SURFACE_FLAG_OPENGL:
		return "SURFACE_FLAG_OPENGL"
	case SURFACE_FLAG_OPENGL_ES:
		return "SURFACE_FLAG_OPENGL_ES"
	case SURFACE_FLAG_OPENGL_ES2:
		return "SURFACE_FLAG_OPENGL_ES2"
	case SURFACE_FLAG_OPENVG:
		return "SURFACE_FLAG_OPENVG"
	default:
		return "[?? Invalid SurfaceFlags value]"
	}
}

func (f SurfaceFlags) ConfigString() string {
	switch f.Config() {
	case SURFACE_FLAG_RGBA32:
		return "SURFACE_FLAG_RGBA32"
	case SURFACE_FLAG_RGB888:
		return "SURFACE_FLAG_RGB888"
	case SURFACE_FLAG_RGB565:
		return "SURFACE_FLAG_RGB565"
	default:
		return "[?? Invalid SurfaceFlags value]"
	}
}

func (f SurfaceFlags) ModString() string {
	m := f.Mod()
	switch {
	case m == 0:
		return ""
	case m&SURFACE_FLAG_ALPHA_FROM_SOURCE == SURFACE_FLAG_ALPHA_FROM_SOURCE:
		return "SURFACE_FLAG_ALPHA_FROM_SOURCE"
	default:
		return "[?? Invalid SurfaceFlags value]"
	}
}

func (f SurfaceFlags) String() string {
	parts := ""
	if f.Type() != SURFACE_FLAG_NONE {
		parts += "|" + f.TypeString()
	}
	parts += "|" + f.ConfigString()
	if f.Mod() != SURFACE_FLAG_NONE {
		parts += "|" + f.ModString()
	}
	return strings.Trim(parts, "|")
}

func (c Color) String() string {
	return fmt.Sprintf("Color{ %.1f,%.1f,%.1f,%.1f }", c.R, c.G, c.B, c.A)
}
