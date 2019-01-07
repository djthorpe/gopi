/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
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

// SurfaceType of surface (which API it's bound to)
type SurfaceType uint

// SurfaceFlags are flags associated with surface
// usually during operations
type SurfaceFlags uint32

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
	Types() []SurfaceType
}

type SurfaceManagerSurfaceMethods interface {
	// Perform all surface operations (create, destroy, move, set, paint) within the 'Do' method
	// to ensure atomic updates to the display. When Do returns, the display is updated and any error
	// from the callback is returned
	Do(SurfaceManagerCallback) error

	// Create & destroy surfaces
	CreateSurface(api SurfaceType, flags SurfaceFlags, opacity float32, layer uint16, origin Point, size Size) (Surface, error)
	CreateSurfaceWithBitmap(bitmap Bitmap, flags SurfaceFlags, opacity float32, layer uint16, origin Point, size Size) (Surface, error)
	DestroySurface(Surface) error

	/*
		// Create background, surface and cursors
		CreateBackground(api SurfaceType, flags SurfaceFlags, opacity float32) (Surface, error)
		CreateCursor(cursor Sprite, flags SurfaceFlags, origin Point) (Surface, error)
	*/

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
	CreateBitmap(SurfaceType, SurfaceFlags, Size) (Bitmap, error)
	DestroyBitmap(Bitmap) error

	// Snapshot the display to a bitmap
	SnapshotDisplay() (Bitmap, error)
}

// Surface is manipulated by surface manager, and used by
// a GPU API (bitmap or vector drawing mostly)
type Surface interface {
	Type() SurfaceType
	Size() Size
	Origin() Point
	Opacity() float32
	Layer() uint16
}

// Bitmap defines a rectangular bitmap which can be used by the GPU
type Bitmap interface {
	Type() SurfaceType
	Size() Size

	// Bitmap operations
	ClearToColor(color Color) error
	FillRectToColor(Point, Size, Color) error
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
	// SurfaceType
	SURFACE_TYPE_NONE SurfaceType = iota
	SURFACE_TYPE_OPENGL
	SURFACE_TYPE_OPENGL_ES
	SURFACE_TYPE_OPENGL_ES2
	SURFACE_TYPE_OPENVG // 2D Vector
	SURFACE_TYPE_RGBA32 // Bitmap, 4 bytes per pixel
	SURFACE_TYPE_RGB888 // Bitmap, 3 bytes per pixel
	SURFACE_TYPE_RGB565 // Bitmap, 2 bytes per pixel
	SURFACE_TYPE_MIN    = SURFACE_TYPE_OPENGL
	SURFACE_TYPE_MAX    = SURFACE_TYPE_RGB888
)

const (
	// SurfaceFlags
	SURFACE_FLAG_NONE              SurfaceFlags = (1 << iota)
	SURFACE_FLAG_ALPHA_FROM_SOURCE SurfaceFlags = (1 << iota)
	SURFACE_FLAG_MIN                            = SURFACE_FLAG_ALPHA_FROM_SOURCE
	SURFACE_FLAG_MAX                            = SURFACE_FLAG_ALPHA_FROM_SOURCE
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
// STRINGIFY

func (t SurfaceType) String() string {
	switch t {
	case SURFACE_TYPE_OPENGL:
		return "SURFACE_TYPE_OPENGL"
	case SURFACE_TYPE_OPENGL_ES:
		return "SURFACE_TYPE_OPENGL_ES"
	case SURFACE_TYPE_OPENGL_ES2:
		return "SURFACE_TYPE_OPENGL_ES2"
	case SURFACE_TYPE_OPENVG:
		return "SURFACE_TYPE_OPENVG"
	case SURFACE_TYPE_RGBA32:
		return "SURFACE_TYPE_RGBA32"
	case SURFACE_TYPE_RGB565:
		return "SURFACE_TYPE_RGB565"
	case SURFACE_TYPE_RGB888:
		return "SURFACE_TYPE_RGB888"
	default:
		return "[Invalid SurfaceType value]"
	}
}

func (f SurfaceFlags) String() string {
	parts := ""
	if f == SURFACE_FLAG_NONE {
		return "SURFACE_FLAG_NONE"
	}
	for flag := SURFACE_FLAG_MIN; flag <= SURFACE_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case SURFACE_FLAG_ALPHA_FROM_SOURCE:
			parts += "|" + "SURFACE_FLAG_ALPHA_FROM_SOURCE"
		default:
			parts += "|" + "[?? Invalid SurfaceFlags value]"
		}
	}
	return strings.Trim(parts, "|")
}
