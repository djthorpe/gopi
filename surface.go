/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

// SurfaceType of surface (which API it's bound to)
type SurfaceType uint

// SurfaceFlags are flags associated with surface
// usually during operations
type SurfaceFlags uint32

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

	// Return a list of extensions the GPU provides
	Extensions() []string

	// Create background, surface and cursors
	CreateBackground(api SurfaceType, flags SurfaceFlags, opacity float32) (Surface, error)
	CreateSurface(api SurfaceType, flags SurfaceFlags, opacity float32, layer uint, origin Point, size Size) (Surface, error)
	//CreateCursor(api SurfaceType, flags SurfaceFlags, opacity float32, origin Point, cursor SurfaceCursor) (Surface, error)
	DestroySurface(Surface) error

	// Change surface properties (size, position, etc)
	MoveSurfaceOriginBy(Surface, SurfaceFlags, Point)
	MoveSurfaceOrigin(SurfaceFlags, Point)
	SetSurfaceSize(Surface, SurfaceFlags, Size)
	SetSurfaceOpacity(Surface, SurfaceFlags, float32)
	SetSurfaceLayer(Surface)

	// Surface operations to start and end drawing or other
	// surface operations
	SetCurrentContext(Surface)
	FlushSurface(Surface)
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
