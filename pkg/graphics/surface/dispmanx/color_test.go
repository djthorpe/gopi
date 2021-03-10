package dispmanx_test

import (
	"fmt"
	"image/color"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
	surface "github.com/djthorpe/gopi/v3/pkg/graphics/surface/dispmanx"
)

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Color_001(t *testing.T) {
	tests := []struct {
		fmt     gopi.SurfaceFormat
		pitch1  uint32
		pitch10 uint32
	}{
		{gopi.SURFACE_FMT_RGBA32, 4, 40},
		{gopi.SURFACE_FMT_XRGB32, 4, 40},
		{gopi.SURFACE_FMT_RGB888, 3, 30},
		{gopi.SURFACE_FMT_RGB565, 2, 20},
		{gopi.SURFACE_FMT_1BPP, 1, 2},
	}
	for _, test := range tests {
		if model := surface.ColorModel(test.fmt); model == nil {
			t.Error("Inavlid Color Model:", test.fmt)
		} else if model.Format() != test.fmt {
			t.Error(test.fmt, "Unexpected return from Format()")
		} else if pitch1 := model.Pitch(1); pitch1 != test.pitch1 {
			t.Error(test.fmt, "Unexpected return from Pitch(1):", pitch1)
		} else if pitch10 := model.Pitch(10); pitch10 != test.pitch10 {
			t.Error(test.fmt, "Unexpected return from Pitch(10):", pitch10)
		} else {
			t.Log(model)
		}
	}
}

func Test_Color_002(t *testing.T) {
	// Test black and white conversion, which should always be
	// lossless
	tests := []struct {
		fmt gopi.SurfaceFormat
	}{
		{gopi.SURFACE_FMT_RGBA32},
		{gopi.SURFACE_FMT_XRGB32},
		{gopi.SURFACE_FMT_RGB888},
		{gopi.SURFACE_FMT_RGB565},
		{gopi.SURFACE_FMT_1BPP},
	}
	for _, test := range tests {
		model := surface.ColorModel(test.fmt)
		black := model.Convert(color.Black)
		white := model.Convert(color.White)
		fmt.Println(test.fmt, black, white)
	}
}
