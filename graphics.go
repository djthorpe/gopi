/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"io"
	"os"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Color including opacity
type Color struct {
	R, G, B, A float32
}

type (
	// SurfaceFlags are flags associated with surface
	SurfaceFlags uint16
)

// SurfaceManagerCallback is a function callback for
// performing surface operations
type SurfaceManagerCallback func(SurfaceManager) error

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// SurfaceManager allows you to open, close and move
// surfaces around an open display
type SurfaceManager interface {
	Driver
	SurfaceManagerSurfaceMethods
	SurfaceManagerBitmapMethods

	// Return the display associated with the surface manager
	Display() Display

	// Return the name of the surface manager. It's basically the
	// GPU driver
	Name() string

	// Return capabilities for the GPU
	Types() []SurfaceFlags
}

type SurfaceManagerSurfaceMethods interface {
	// Perform all surface operations (create, destroy, move, set, paint) within the 'Do' method
	// to ensure atomic updates to the display. When Do returns, the display is updated and any error
	// from the callback is returned
	Do(SurfaceManagerCallback) error

	// Create & destroy surfaces
	CreateSurface(flags SurfaceFlags, opacity float32, layer uint16, origin Point, size Size) (Surface, error)
	CreateSurfaceWithBitmap(bitmap Bitmap, flags SurfaceFlags, opacity float32, layer uint16, origin Point, size Size) (Surface, error)
	CreateBackground(flags SurfaceFlags, opacity float32) (Surface, error)
	CreateCursor(cursor Sprite, flags SurfaceFlags, origin Point) (Surface, error)
	DestroySurface(Surface) error

	// Change surface properties (size, position, etc)
	SetOrigin(Surface, Point) error
	MoveOriginBy(Surface, Point) error
	SetSize(Surface, Size) error
	SetLayer(Surface, uint16) error
	SetOpacity(Surface, float32) error
	SetBitmap(Bitmap) error
}

type SurfaceManagerBitmapMethods interface {
	// Create and destroy bitmaps
	CreateBitmap(SurfaceFlags, Size) (Bitmap, error)
	CreateSnapshot(SurfaceFlags) (Bitmap, error)
	DestroyBitmap(Bitmap) error
}

// Surface is manipulated by surface manager, and used by
// a GPU API (bitmap or vector drawing mostly)
type Surface interface {
	Type() SurfaceFlags
	Size() Size
	Origin() Point
	Opacity() float32
	Layer() uint16
}

// Bitmap defines a rectangular bitmap which can be used by the GPU
type Bitmap interface {
	Type() SurfaceFlags
	Size() Size

	// Bitmap operations
	ClearToColor(Color) error
	FillRectToColor(Color, Point, Size) error
	PaintPixel(Color, Point) error
	//PaintImage(image.Image, Point, Size) error
	//PaintText(Color, string, FontFace, FontSize, Point) error
}

// SpriteManager loads sprites from io.Reader buffers
type SpriteManager interface {
	Driver

	// Open one or more sprites from a stream and return them
	OpenSprites(io.Reader) ([]Sprite, error)

	// Open sprites from path, checking to see if individual files should
	// be opened through a callback function
	OpenSpritesAtPath(path string, callback func(manager SpriteManager, path string, info os.FileInfo) bool) error

	// Return loaded sprites, or a specific sprite
	Sprites(name string) []Sprite
}

// Sprite implemnts a bitmap with a unique name and hotspot location (for cursors)
type Sprite interface {
	Bitmap

	Name() string
	Hotspot() Point
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

// Standard Colors
var (
	ColorRed       = Color{1.0, 0.0, 0.0, 1.0}
	ColorGreen     = Color{0.0, 1.0, 0.0, 1.0}
	ColorBlue      = Color{0.0, 0.0, 1.0, 1.0}
	ColorWhite     = Color{1.0, 1.0, 1.0, 1.0}
	ColorBlack     = Color{0.0, 0.0, 0.0, 1.0}
	ColorPurple    = Color{1.0, 0.0, 1.0, 1.0}
	ColorCyan      = Color{0.0, 1.0, 1.0, 1.0}
	ColorYellow    = Color{1.0, 1.0, 0.0, 1.0}
	ColorDarkGrey  = Color{0.25, 0.25, 0.25, 1.0}
	ColorLightGrey = Color{0.75, 0.75, 0.75, 1.0}
	ColorMidGrey   = Color{0.5, 0.5, 0.5, 1.0}
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATIONS

func (c Color) RGBA() (r, g, b, a uint32) {
	return uint32(c.R*float32(0xFFFF)) & uint32(0xFFFF), uint32(c.G*float32(0xFFFF)) & uint32(0xFFFF), uint32(c.B*float32(0xFFFF)) & uint32(0xFFFF), uint32(c.A*float32(0xFFFF)) & uint32(0xFFFF)
}

// Type() returns the type of the surface
func (f SurfaceFlags) Type() SurfaceFlags {
	return f & SURFACE_FLAG_TYPEMASK
}

// Config() returns the configuration of the surface
func (f SurfaceFlags) Config() SurfaceFlags {
	return f & SURFACE_FLAG_CONFIGMASK
}

// Mod() returns surface modifiers
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
	parts += "|" + f.TypeString()
	parts += "|" + f.ConfigString()
	parts += "|" + f.ModString()
	return strings.Trim(parts, "|")
}

func (c Color) String() string {
	return fmt.Sprintf("Color{ %.1f,%.1f,%.1f,%.1f }", c.R, c.G, c.B, c.A)
}
