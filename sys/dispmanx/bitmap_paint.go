// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package dispmanx

import (
	"fmt"
	"image"
	"image/color"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	"github.com/djthorpe/gopi/v2/unit/fonts"
)

////////////////////////////////////////////////////////////////////////////////
// CLEAR TO COLOR

func (this *bitmap) ClearToColor(c color.Color) {
	// Create buffer for one row
	data, err := this.ReadRows(0, 1, false)
	if err != nil {
		return
	}
	// Get size of bitmap
	size := rpi.DXRectSize(this.bounds)
	// Clear the row to the color
	buf := data.Bytes()
	pixel := colorToBytes(this.imageType, c)
	for i := uint32(0); i < this.stride; i++ {
		buf[i] = pixel[i%this.bytesPerPixel]
	}
	// Write data out for each row
	for offset := uint32(0); offset < size.H; offset++ {
		if err := this.WriteRows(data, offset); err != nil {
			return
		}
	}
	// Set modified as the whole bitmap area
	this.setDirty(this.bounds)
}

////////////////////////////////////////////////////////////////////////////////
// PAINT PIXEL

func (this *bitmap) PaintPixel(c color.Color, p rpi.DXPoint) {
	// Paint pixel
	this.paintPixel(c, p)
	// Add dirty area
	this.addDirty(rpi.DXNewRect(p.X, p.Y, 1, 1))
}

func (this *bitmap) paintPixel(c color.Color, p rpi.DXPoint) {
	// Check for point in rect
	if rpi.DXRectContainsPoint(this.bounds, p) == false {
		return
	}
	// Create buffer for one row
	data, err := this.ReadRows(uint32(p.Y), 1, true)
	if err != nil {
		fmt.Println("ReadRows Error:", err)
		return
	}
	// Record pixel to row
	buf := data.Bytes()
	offset := this.bytesPerPixel * uint32(p.X)
	for _, pixel := range colorToBytes(this.imageType, c) {
		buf[offset] = pixel
		offset += 1
	}
	// Write data out for the row
	if err := this.WriteRows(data, uint32(p.Y)); err != nil {
		return
	}
}

////////////////////////////////////////////////////////////////////////////////
// PAINT CIRCLE (OUTLINE)

func (this *bitmap) PaintCircle(c color.Color, p rpi.DXPoint, r uint32) {
	// Add dirty area
	this.addDirty(rpi.DXNewRect(p.X-int32(r), p.Y-int32(r), r*2, r*2))
	// Deal with r==0
	if r == 0 {
		this.paintPixel(c, p)
		return
	}
	// https://rosettacode.org/wiki/Bitmap/Midpoint_circle_algorithm#Go
	r1 := int32(r)
	x1, y1 := int32(-r1), int32(0)
	err := int32(2) - int32(2)*int32(r)
	for {
		this.paintPixel(c, p.Add(rpi.DXPoint{-x1, +y1}))
		this.paintPixel(c, p.Add(rpi.DXPoint{-y1, -x1}))
		this.paintPixel(c, p.Add(rpi.DXPoint{+x1, -y1}))
		this.paintPixel(c, p.Add(rpi.DXPoint{+y1, +x1}))
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
// PAINT LINE

func (this *bitmap) PaintLine(c color.Color, p0 rpi.DXPoint, p1 rpi.DXPoint) {
	// Add dirty area before we start
	this.addDirty(rpi.DXRectFromPoints(p0, p1))
	// https://rosettacode.org/wiki/Bitmap/Bresenham%27s_line_algorithm#Go
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
// PAINT RUNE

func (this *bitmap) PaintRune(c color.Color, ch rune, face gopi.FontFace) {
	if face_, ok := face.(fonts.Face); ok == false || face_ == nil {
		// Invalid face
		return
	}
	if bitmap, err := face_.BitmapForRunePixels(ch, 128)(image.Image, error); err != nil {
		// Error
		return
	}
}
