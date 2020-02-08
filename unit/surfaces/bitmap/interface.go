/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap

import (
	"image/color"
	"image/draw"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

type Config struct {
	Size gopi.Size
	Mode gopi.SurfaceFlags
}

type Bitmap interface {
	// dispmanx properties
	DXSize() rpi.DXSize
	DXRect() rpi.DXRect
	DXHandle() rpi.DXResource

	// Return image type
	Type() gopi.SurfaceFlags

	// Return bounds
	Origin() gopi.Point
	Size() gopi.Size

	// Return the bitmap as bytes
	// with bytes per row
	Bytes() ([]byte, uint32)

	// ClearToColor clears the screen to a single color
	ClearToColor(color.Color)

	// Paint a pixel
	Pixel(color.Color, gopi.Point)

	// PaintCircle paints an outlined circle with origin and radius
	CircleOutline(color.Color, gopi.Point, float32)

	// PaintLine paints a line
	Line(color.Color, gopi.Point, gopi.Point)

	// Paint a rune with a particular font face
	Rune(color.Color, gopi.Point, rune, gopi.FontFace, gopi.FontSize)

	// Retain and release
	Retain()
	Release() bool

	// Return points in the rectangle
	Centre() gopi.Point
	NorthWest() gopi.Point
	SouthWest() gopi.Point
	NorthEast() gopi.Point
	SouthEast() gopi.Point

	// Implements image.Image and draw.Image
	draw.Image

	// Implements gopi.Unit
	gopi.Unit
}
