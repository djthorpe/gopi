// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap

import (
	"image/color"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// CLEAR TO COLOR

func (this *bitmap) ClearToColor(c color.Color) {
	// Check some parameters
	if this.handle == rpi.DX_NO_HANDLE {
		return
	}
	// Create a row
	if err := this.Data.SetCapacity(uint(this.stride)); err != nil {
		this.Log.Warn(err)
		return
	}
	// Set the pixels on the row
	color := colorToBytes(this.mode, c)
	bytesPerPixel := len(color)
	if color == nil {
		this.Log.Warn("colorToBytes returned nil")
		return
	}
	bytes := this.Data.Bytes()
	for i := range bytes {
		bytes[i] = color[i%bytesPerPixel]
	}
	// Write out to each row in the bitmap
	for y := uint(0); y < uint(this.size.H); y++ {
		if err := this.Data.Write(this.handle, this.dxmode, y, 1, this.stride); err != nil {
			this.Log.Warn(err)
			return
		}
	}
	// Set modified as the whole bitmap area
	this.setDirty(this.bounds)
}

////////////////////////////////////////////////////////////////////////////////
// PAINT PIXEL

func (this *bitmap) PaintPixel(c color.Color, p gopi.Point) {
	// Check some parameters
	if this.handle == rpi.DX_NO_HANDLE {
		return
	}
	// Paint pixel
	this.paintPixel(c, rpi.DXPoint{int32(p.X), int32(p.Y)})
	// Add dirty area
	this.addDirty(rpi.DXNewRect(int32(p.X), int32(p.Y), 1, 1))
}

func (this *bitmap) paintPixel(c color.Color, p rpi.DXPoint) {
	// Check for point in rect
	if rpi.DXRectContainsPoint(this.bounds, p) == false {
		return
	}
	// Read a row of data
	// TODO: Indicate a row is cached in some circumstances
	if err := this.Data.SetCapacity(uint(this.stride)); err != nil {
		this.Log.Warn(err)
		return
	} else if err := this.Data.Read(this.handle, uint(p.Y), 1, this.stride); err != nil {
		this.Log.Warn(err)
		return
	}
	// Get the pixel
	color := colorToBytes(this.mode, c)
	bytesPerPixel := len(color)
	if color == nil {
		this.Log.Warn("colorToBytes returned nil")
		return
	}
	bytes := this.Data.Bytes()
	offset := int(p.X) * bytesPerPixel
	for i := range color {
		bytes[offset+i] = color[i]
	}
	// Write out the row
	if err := this.Data.Write(this.handle, this.dxmode, uint(p.Y), 1, this.stride); err != nil {
		this.Log.Warn(err)
		return
	}
}

////////////////////////////////////////////////////////////////////////////////
// PAINT LINE
// REF: https://rosettacode.org/wiki/Bitmap/Bresenham%27s_line_algorithm#Go

func (this *bitmap) PaintLine(c color.Color, from gopi.Point, to gopi.Point) {
	p0 := rpi.DXPoint{int32(from.X), int32(from.Y)}
	p1 := rpi.DXPoint{int32(to.X), int32(to.Y)}
	// Add dirty area before we start
	this.addDirty(rpi.DXRectFromPoints(p0, p1))
	dx := p1.X - p0.X
	if dx < 0 {
		dx = -dx
	}
	dy := p1.Y - p0.Y
	if dy < 0 {
		dy = -dy
	}
	var sx, sy int32
	if p0.X < p1.X {
		sx = 1
	} else {
		sx = -1
	}
	if p0.Y < p1.Y {
		sy = 1
	} else {
		sy = -1
	}
	err := dx - dy
	for {
		this.paintPixel(c, p0)
		if p0.Equals(p1) {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			p0.X += sx
		}
		if e2 < dx {
			err += dx
			p0.Y += sy
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// CIRCLE OUTLINE
// REF: https://rosettacode.org/wiki/Bitmap/Midpoint_circle_algorithm#Go

func (this *bitmap) PaintCircle(c color.Color, p gopi.Point, r uint32) {
	centre := rpi.DXPoint{int32(p.X), int32(p.Y)}
	// Add dirty area
	this.addDirty(rpi.DXNewRect(centre.X-int32(r), centre.Y-int32(r), r*2, r*2))
	// Deal with r==0
	if r == 0 {
		this.paintPixel(c, centre)
		return
	}
	r1 := int32(r)
	x1, y1 := int32(-r1), int32(0)
	err := int32(2) - int32(2)*int32(r)
	for {
		this.paintPixel(c, centre.Add(rpi.DXPoint{-x1, +y1}))
		this.paintPixel(c, centre.Add(rpi.DXPoint{-y1, -x1}))
		this.paintPixel(c, centre.Add(rpi.DXPoint{+x1, -y1}))
		this.paintPixel(c, centre.Add(rpi.DXPoint{+y1, +x1}))
		r1 = err
		if r1 > x1 {
			x1++
			err += x1*2 + 1
		}
		if r1 <= y1 {
			y1++
			err += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// RUNE

func (this *bitmap) PaintRune(c color.Color, p gopi.Point, r rune) {

}
