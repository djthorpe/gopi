/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
)

func Test_Color_000(t *testing.T) {
	black := bitmap.Color{[]byte{0, 0, 0}, gopi.SURFACE_FLAG_RGB888}
	if r, g, b, a := black.RGBA(); r != 0 || g != 0 || b != 0 || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", black)
	}
	white := bitmap.Color{[]byte{0xFF, 0xFF, 0xFF}, gopi.SURFACE_FLAG_RGB888}
	if r, g, b, a := white.RGBA(); r != 0xFFFF || g != 0xFFFF || b != 0xFFFF || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", white)
	}
	red := bitmap.Color{[]byte{0xFF, 0x00, 0x00}, gopi.SURFACE_FLAG_RGB888}
	if r, g, b, a := red.RGBA(); r != 0xFFFF || g != 0 || b != 0 || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", red)
	}
	green := bitmap.Color{[]byte{0x00, 0xFF, 0x00}, gopi.SURFACE_FLAG_RGB888}
	if r, g, b, a := green.RGBA(); r != 0 || g != 0xFFFF || b != 0 || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", green)
	}
	blue := bitmap.Color{[]byte{0x00, 0x00, 0xFF}, gopi.SURFACE_FLAG_RGB888}
	if r, g, b, a := blue.RGBA(); r != 0 || g != 0 || b != 0xFFFF || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", blue)
	}
}

func Test_Color_001(t *testing.T) {
	transparent := bitmap.Color{[]byte{0, 0, 0, 0}, gopi.SURFACE_FLAG_RGBA32}
	if r, g, b, a := transparent.RGBA(); r != 0 || g != 0 || b != 0 || a != 0 {
		t.Error("Unexpected RGBA values for", transparent)
	}
	black := bitmap.Color{[]byte{0, 0, 0, 0xFF}, gopi.SURFACE_FLAG_RGBA32}
	if r, g, b, a := black.RGBA(); r != 0 || g != 0 || b != 0 || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", black)
	}
	white := bitmap.Color{[]byte{0xFF, 0xFF, 0xFF, 0xFF}, gopi.SURFACE_FLAG_RGBA32}
	if r, g, b, a := white.RGBA(); r != 0xFFFF || g != 0xFFFF || b != 0xFFFF || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", white)
	}
	red := bitmap.Color{[]byte{0xFF, 0x00, 0x00, 0xFF}, gopi.SURFACE_FLAG_RGBA32}
	if r, g, b, a := red.RGBA(); r != 0xFFFF || g != 0 || b != 0 || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", red)
	}
	green := bitmap.Color{[]byte{0x00, 0xFF, 0x00, 0xFF}, gopi.SURFACE_FLAG_RGBA32}
	if r, g, b, a := green.RGBA(); r != 0 || g != 0xFFFF || b != 0 || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", green)
	}
	blue := bitmap.Color{[]byte{0x00, 0x00, 0xFF, 0xFF}, gopi.SURFACE_FLAG_RGBA32}
	if r, g, b, a := blue.RGBA(); r != 0 || g != 0 || b != 0xFFFF || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", blue)
	}
}

func Test_Color_002(t *testing.T) {
	black := bitmap.Color{[]byte{0, 0}, gopi.SURFACE_FLAG_RGB565}
	if r, g, b, a := black.RGBA(); r != 0 || g != 0 || b != 0 || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", black)
	}
	white := bitmap.Color{[]byte{0xFF, 0xFF}, gopi.SURFACE_FLAG_RGB565}
	if r, g, b, a := white.RGBA(); r != 0xFFFF || g != 0xFFFF || b != 0xFFFF || a != 0xFFFF {
		t.Error("Unexpected RGBA values for", white)
	}
}
