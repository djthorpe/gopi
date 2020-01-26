// +build freetype

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package freetype

import (
	"fmt"
	"image"
	"image/color"

	// Frameworks
	ft "github.com/djthorpe/gopi/v2/sys/freetype"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

type bitmap struct {
	handle  ft.FT_Bitmap
	advance struct{ x, y uint }
	row     struct {
		data []uint32
		y    uint
	}
}

type pixel struct {
	value uint32
	mode  ft.FT_PixelMode
}

type colormodel struct {
	mode ft.FT_PixelMode
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewBitmap(handle ft.FT_Bitmap, x, y uint) (image.Image, error) {
	this := new(bitmap)
	this.handle = handle
	this.advance.x = x
	this.advance.y = y
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *bitmap) At(x, y int) color.Color {
	if this.row.data == nil || this.row.y != uint(y) {
		this.row.data, this.row.y = ft.FT_BitmapPixelsForRow(this.handle, uint(y)), uint(y)
	}
	if this.row.data == nil {
		return color.Transparent
	}
	if x < 0 || x >= len(this.row.data) {
		return color.Transparent
	} else {
		return &pixel{this.row.data[x], ft.FT_BitmapPixelMode(this.handle)}
	}
	return color.Transparent
}

func (this *bitmap) Bounds() image.Rectangle {
	w, h := ft.FT_BitmapSize(this.handle)
	return image.Rect(0, 0, int(w)-1, int(h)-1)
}

func (this *bitmap) ColorModel() color.Model {
	return &colormodel{ft.FT_BitmapPixelMode(this.handle)}
}

// Advance to the next location
func (this *bitmap) Advance() (uint, uint) {
	return this.advance.x, this.advance.y
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *bitmap) String() string {
	return "<gopi.fonts.bitmap" +
		" " + this.handle.String() +
		" advance=" + fmt.Sprintf("{%d,%d}", this.advance.x, this.advance.y) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// COLOR

func (this *pixel) RGBA() (uint32, uint32, uint32, uint32) {
	switch this.mode {
	case ft.FT_PIXEL_MODE_MONO:
		// 1 bit per pixel
		if this.value == 0 {
			// Black
			return 0, 0, 0, 0xFFFF
		} else {
			// White
			return 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF
		}
	case ft.FT_PIXEL_MODE_LCD, ft.FT_PIXEL_MODE_LCD_V, ft.FT_PIXEL_MODE_GRAY:
		// 8 bits per pixel
		value16 := (this.value * 0x0101) & 0xFFFF
		return value16, value16, value16, 0xFFFF
	case ft.FT_PIXEL_MODE_GRAY2: // 2 bits per pixel
		value16 := (this.value * 0x5555) & 0xFFFF
		return value16, value16, value16, 0xFFFF
	case ft.FT_PIXEL_MODE_GRAY4: // 4 bits per pixel
		value16 := (this.value * 0x1111) & 0xFFFF
		return value16, value16, value16, 0xFFFF
	case ft.FT_PIXEL_MODE_BGRA: // 32 bits per pixel
		b := ((this.value >> 24) * 257) & 0xFFFF
		g := ((this.value >> 16) * 257) & 0xFFFF
		r := ((this.value >> 8) * 257) & 0xFFFF
		a := ((this.value >> 0) * 257) & 0xFFFF
		return r, g, b, a
	default:
		// Black
		return 0, 0, 0, 0xFFFF
	}
}

////////////////////////////////////////////////////////////////////////////////
// COLOR MODEL

func (this *colormodel) Convert(color.Color) color.Color {
	return color.White
}
