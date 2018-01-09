/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"image/color"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// SurfaceType of surface (which API it's bound to)
type SurfaceType uint

// SurfaceFlags are flags associated with surface
// usually during operations
type SurfaceFlags uint32

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// SurfaceManager allows you to open, close and move
// surfaces around an open display
type SurfaceManager interface {
	Driver

	// Return the display associated with the surface manager
	Display() Display

	// Return the name of the surface manager. It's basically the
	// GPU driver
	Name() string

	// Return capabilities for the GPU
	Types() []SurfaceType

	// Create & destroy surfaces
	CreateSurface(api SurfaceType, flags SurfaceFlags, opacity float32, layer uint16, origin Point, size Size) (Surface, error)
	DestroySurface(Surface) error

	// Create and destroy bitmaps
	CreateBitmap(Size) (Bitmap, error)
	DestroyBitmap(Bitmap) error

	/*
		// Create background, surface and cursors
		CreateBackground(api SurfaceType, flags SurfaceFlags, opacity float32) (Surface, error)
		CreateCursor(api SurfaceType, flags SurfaceFlags, opacity float32, origin Point, cursor SurfaceCursor) (Surface, error)

		// Change surface properties (size, position, etc)
		MoveOriginBy(Surface, SurfaceFlags, Point)
		SetOrigin(Surface, SurfaceFlags, Point)
		SetSize(Surface, SurfaceFlags, Size)
		SetOpacity(Surface, SurfaceFlags, float32)
		SetLayer(Surface, uint)

		// Surface operations to start and end drawing or other
		// surface operations
		SetCurrentContext(Surface)
		FlushSurface(Surface)
	*/
}

// Surface is manipulated by surface manager, and used by
// a GPU API (bitmap or vector drawing mostly)
type Surface interface {
	Driver

	Type() SurfaceType
	Size() Size
	Origin() Point
	Opacity() float32
	Layer() uint16
}

// Bitmap defines a rectangular bitmap which can be used
// by the GPU
type Bitmap interface {
	Driver

	Type() SurfaceType
	Size() Size

	// Bitmap operations
	ClearToColorRGBA(color color.RGBA) error
}

/*
type SurfaceCursor interface {
	API()
	Hotspot()
	Size()
}
*/

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// SurfaceType
	SURFACE_TYPE_NONE SurfaceType = iota
	SURFACE_TYPE_OPENGL
	SURFACE_TYPE_OPENGL_ES
	SURFACE_TYPE_OPENGL_ES2
	SURFACE_TYPE_OPENVG
	SURFACE_TYPE_RGBA32
)

const (
	// SurfaceType
	SURFACE_FLAG_NONE              SurfaceFlags = (1 << iota)
	SURFACE_FLAG_ALPHA_FROM_SOURCE              = (1 << iota)
)

const (
	// SurfaceLayer
	SURFACE_LAYER_BACKGROUND uint16 = 0x0000
	SURFACE_LAYER_DEFAULT    uint16 = 0x0001
	SURFACE_LAYER_CURSOR     uint16 = 0xFFFF
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
	default:
		return "[Invalid SurfaceType value]"
	}
}

func (f SurfaceFlags) String() string {
	if f == SURFACE_FLAG_NONE {
		return "SURFACE_FLAG_NONE"
	}
	flags := ""
	// Add flags here
	return strings.Trim(flags, "|")
}
