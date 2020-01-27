/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Color struct {
	Data []byte
	Mode gopi.SurfaceFlags
}

type ColorModel struct {
	Mode gopi.SurfaceFlags
}

////////////////////////////////////////////////////////////////////////////////
// COLOR

func (c Color) RGBA() (r, g, b, a uint32) {
	switch c.Mode {
	case gopi.SURFACE_FLAG_RGB888:
		// Three bytes in R,G,B order
		return uint8to16(c.Data[0]), uint8to16(c.Data[1]), uint8to16(c.Data[2]), 0xFFFF
	case gopi.SURFACE_FLAG_RGBA32:
		// Four bytes - compress down to 0xFFFF
		return uint8to16(c.Data[0]), uint8to16(c.Data[1]), uint8to16(c.Data[2]), uint8to16(c.Data[3])
	case gopi.SURFACE_FLAG_RGB565:
		// Two bytes
		r := c.Data[0] >> 3 & 0x1F
		g1 := c.Data[0] & 0x07
		g2 := c.Data[1] & 0xE0
		g := (g1 << 3) | (g2 >> 5)
		b := c.Data[1] & 0x1F
		return uint5to16(r), uint6to16(g), uint5to16(b), 0xFFFF
	default:
		return 0, 0, 0, 0
	}
}

func (c Color) String() string {
	r, g, b, a := c.RGBA()
	return "<Color" +
		" bytes=" + strings.ToUpper(hex.EncodeToString(c.Data)) +
		" mode=" + c.Mode.ConfigString() +
		" rgba=" + fmt.Sprintf("{ %04X,%04X,%04X,%04X }", r, g, b, a) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// COLORMODEL

func (m ColorModel) Convert(c color.Color) color.Color {
	if bytes := colorToBytes(m.Mode, c); bytes != nil {
		return Color{bytes, m.Mode}
	} else {
		return color.Transparent
	}
}

////////////////////////////////////////////////////////////////////////////////
// NEW AND FREE

func uint8to16(value uint8) uint32 {
	return uint32(value) * 0x0101 & 0xFFFF
}

func uint5to16(value uint8) uint32 {
	if value == 0x1F {
		return 0xFFFF
	} else {
		return uint32(value&0x1F) * 0x0842 & 0xFFFF
	}
}

func uint6to16(value uint8) uint32 {
	if value == 0x3F {
		return 0xFFFF
	} else {
		return uint32(value&0x3F) * 0x0410 & 0xFFFF
	}
}

func colorToBytes(mode gopi.SurfaceFlags, c color.Color) []byte {
	// Returns color 0000 <= v <= FFFF
	r, g, b, a := c.RGBA()
	// Convert to []byte
	switch mode {
	case gopi.SURFACE_FLAG_RGB888:
		return []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8)}
	case gopi.SURFACE_FLAG_RGB565:
		r := uint16(r>>(8+3)) << (5 + 6)
		g := uint16(g>>(8+2)) << 5
		b := uint16(b >> (8 + 3))
		v := r | g | b
		return []byte{byte(v), byte(v >> 8)}
	case gopi.SURFACE_FLAG_RGBA32:
		return []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}
	default:
		return nil
	}
}
