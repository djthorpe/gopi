package bitmap_test

import (
	"image/color"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
	bitmap "github.com/djthorpe/gopi/v3/pkg/graphics/bitmap"
)

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Color_001(t *testing.T) {
	tests := []struct {
		c          color.Color
		r, g, b, a uint32
	}{
		{color.Black, 0x0000, 0x0000, 0x0000, 0xFFFF},
		{color.White, 0xFFFF, 0xFFFF, 0xFFFF, 0xFFFF},
		{color.RGBA{0xFF, 0x00, 0x00, 0xFF}, 0xFFFF, 0x0000, 0x0000, 0xFFFF}, // Red
		{color.RGBA{0x00, 0xFF, 0x00, 0xFF}, 0x0000, 0xFFFF, 0x0000, 0xFFFF}, // Green
		{color.RGBA{0x00, 0x00, 0xFF, 0xFF}, 0x0000, 0x0000, 0xFFFF, 0xFFFF}, // Blue
	}
	for fmt := gopi.SURFACE_FMT_NONE; fmt <= gopi.SURFACE_FMT_MAX; fmt++ {
		model := bitmap.GetColorModel(fmt)
		if model == nil {
			continue
		}
		if v := model.Format(); v != fmt {
			t.Error("Unexpected return from Format()", v)
		} else {
			t.Log(model)
		}
		for _, test := range tests {
			color := model.Convert(test.c)
			r, g, b, a := color.RGBA()
			if r != test.r || g != test.g || b != test.b || a != test.a {
				t.Error("Unexpected return from Convert()", test.c, " != ", color)
			}
		}
	}
}
