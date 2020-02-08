/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"image/color"
	"image/draw"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// SurfaceFlags are flags associated with surface
	SurfaceFlags uint16

	// SurfaceCallback
	SurfaceCallback func() error
)

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
	ClearToColor(color.Color)

	// Set Pixel
	Pixel(color.Color, Point)

	// Draw line
	Line(color.Color, Point, Point)

	// Outline circle with centre and radius
	CircleOutline(color.Color, Point, float32)

	// Paint a rune
	Rune(color.Color, Point, rune, FontFace, FontSize)

	// Implements image.Image and draw.Image
	draw.Image
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
	ColorRed         = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
	ColorGreen       = color.RGBA{0x00, 0xFF, 0x00, 0xFF}
	ColorBlue        = color.RGBA{0x00, 0x00, 0xFF, 0xFF}
	ColorWhite       = color.White
	ColorBlack       = color.Black
	ColorPurple      = color.RGBA{0xFF, 0x00, 0xFF, 0xFF}
	ColorCyan        = color.RGBA{0x00, 0xFF, 0xFF, 0xFF}
	ColorYellow      = color.RGBA{0xFF, 0xFF, 0x00, 0xFF}
	ColorDarkGrey    = color.RGBA{0x40, 0x40, 0x40, 0xFF}
	ColorLightGrey   = color.RGBA{0xBF, 0xBF, 0xBF, 0xFF}
	ColorMidGrey     = color.RGBA{0x80, 0x80, 0x80, 0xFF}
	ColorTransparent = color.Transparent
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATIONS

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
