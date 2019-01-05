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

// SurfaceManagerCallback is a function callback for
// performing surface operations
type SurfaceManagerCallback func(SurfaceManager) error

// SurfaceCallback is a function callback for
// performing drawing operations. If it's a bitmap
// the second argument will contain the bitmap
// resource
type SurfaceCallback func(Surface, Bitmap) error

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

// SurfaceManager allows you to open, close and move
// surfaces around an open display
type SurfaceManager interface {
	Driver
	SurfaceManagerSurfaceMethods
	//SurfaceManagerBitmapMethods

	// Return the display associated with the surface manager
	Display() Display

	// Return the name of the surface manager. It's basically the
	// GPU driver
	Name() string

	// Return capabilities for the GPU
	Types() []SurfaceType
}

type SurfaceManagerSurfaceMethods interface {
	// Perform surface operations (create, destroy, move, set) within
	// a 'Do' method to ensure atomic updates to the display
	Do(SurfaceManagerCallback) error

	// Create & destroy surfaces
	CreateSurface(api SurfaceType, flags SurfaceFlags, opacity float32, layer uint16, origin Point, size Size) (Surface, error)
	CreateSurfaceWithBitmap(bitmap Bitmap, flags SurfaceFlags, opacity float32, layer uint16, origin Point, size Size) (Surface, error)
	DestroySurface(Surface) error

	/*
		// Create background, surface and cursors
		CreateBackground(api SurfaceType, flags SurfaceFlags, opacity float32) (Surface, error)
		CreateCursor(api SurfaceType, flags SurfaceFlags, opacity float32, origin Point, cursor SurfaceCursor) (Surface, error)
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
	CreateBitmap(SurfaceType, Size) (Bitmap, error)
	DestroyBitmap(Bitmap) error
}

// Surface is manipulated by surface manager, and used by
// a GPU API (bitmap or vector drawing mostly)
type Surface interface {
	Type() SurfaceType
	Size() Size
	Origin() Point
	Opacity() float32
	Layer() uint16

	//Perform drawing operations on a surface - if it's a bitmap
	//then this reads the bitmap data from the GPU ready to draw on
	//Do(SurfaceCallback) error
}

// Bitmap defines a rectangular bitmap which can be used
// by the GPU
type Bitmap interface {
	Type() SurfaceType
	Size() Size

	// Bitmap operations - requires you to use Do() on
	// the surface in order to operate on the bitmap
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
	SURFACE_LAYER_MAX        uint16 = 0xFFFE
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
