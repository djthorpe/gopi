package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CHECK FLAGS

func TestFlags_000(t *testing.T) {
	for f := gopi.SURFACE_FLAG_BITMAP; f <= gopi.SURFACE_FLAG_OPENVG; f++ {
		t.Logf("%v => type %v config %v mod %v", f, f.TypeString(), f.ConfigString(), f.ModString())
	}
}

func TestFlags_001(t *testing.T) {
	for f := gopi.SURFACE_FLAG_BITMAP; f <= gopi.SURFACE_FLAG_OPENVG; f++ {
		g := f | gopi.SURFACE_FLAG_RGB565 | gopi.SURFACE_FLAG_ALPHA_FROM_SOURCE
		t.Logf("%v => type %v config %v mod %v", g, g.TypeString(), g.ConfigString(), g.ModString())
	}
}

////////////////////////////////////////////////////////////////////////////////
// CHECK COLORS

func TestColors_000(t *testing.T) {
	all_colors := []gopi.Color{
		gopi.ColorRed, gopi.ColorGreen, gopi.ColorBlue, gopi.ColorWhite, gopi.ColorBlack,
		gopi.ColorPurple, gopi.ColorCyan, gopi.ColorYellow, gopi.ColorDarkGrey,
		gopi.ColorLightGrey, gopi.ColorMidGrey, gopi.ColorTransparent,
	}
	for _, color := range all_colors {
		r, g, b, a := color.RGBA()
		t.Logf("%v => [ %04X %04X %04X %04X ]", color, r, g, b, a)
	}
}
