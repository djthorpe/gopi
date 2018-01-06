/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import "strings"

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
	/*
		// Return capabilities for the GPU
		Types() []SurfaceType

		// Return a list of extensions the GPU provides
		Extensions() []string

		// Create background, surface and cursors
		CreateBackground(api SurfaceType, flags SurfaceFlags, opacity float32) (Surface, error)
		CreateSurface(api SurfaceType, flags SurfaceFlags, opacity float32, layer uint, origin Point, size Size) (Surface, error)
		//CreateCursor(api SurfaceType, flags SurfaceFlags, opacity float32, origin Point, cursor SurfaceCursor) (Surface, error)
		DestroySurface(Surface) error

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
	Type() SurfaceType
	Opacity() float32
	Layer() uint
	Origin() Point
	Size() Size
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
	SURFACE_FLAG_NONE SurfaceFlags = 0
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
