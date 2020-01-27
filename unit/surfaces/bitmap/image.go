// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap

import (
	"image"
	"image/color"

	// Frameworks
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION image.Image and draw.Image

// Return bounds
func (this *bitmap) Bounds() image.Rectangle {
	return image.Rectangle{
		image.Point{0, 0}, image.Point{int(this.size.W - 1), int(this.size.H - 1)},
	}
}

// Return colour model
func (this *bitmap) ColorModel() color.Model {
	return ColorModel{this.mode}
}

// Get a pixel to a color
func (this *bitmap) At(x, y int) color.Color {
	// Check some parameters
	if this.handle == rpi.DX_NO_HANDLE {
		return color.Transparent
	}
	if x < 0 || y < 0 || uint32(x) >= this.size.W || uint32(y) >= this.size.H {
		return color.Transparent
	}
	// Read a row of data
	// TODO: Indicate a row is cached in some circumstances
	if err := this.Data.SetCapacity(uint(this.stride)); err != nil {
		this.Log.Warn(err)
		return color.Transparent
	} else if err := this.Data.Read(this.handle, uint(y), 1, this.stride); err != nil {
		this.Log.Warn(err)
		return color.Transparent
	}
	// Return the pixel
	bytes := this.Data.Bytes()
	low := x * int(this.pixelSize)
	high := low + int(this.pixelSize)
	return Color{bytes[low:high], this.mode}
}

// Set a single pixel to a color
func (this *bitmap) Set(x, y int, c color.Color) {
	// Check some parameters
	if this.handle == rpi.DX_NO_HANDLE {
		return
	}	
	// Paint pixel
	this.paintPixel(c, rpi.DXPoint{ int32(x), int32(y) })
	// Add dirty area
	this.addDirty(rpi.DXNewRect(int32(x),int32(y), 1, 1))
}
