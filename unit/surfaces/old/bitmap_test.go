// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces_test

import (
	"fmt"
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/display"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/platform"
	_ "github.com/djthorpe/gopi/v2/unit/surfaces"
)

func Test_Bitmap_000(t *testing.T) {
	t.Log("Test_Bitmap_000")
}

func Test_Bitmap_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Bitmap_001, []string{"-debug"}, "surfaces"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Bitmap_001(app gopi.App, t *testing.T) {
	surfaces := app.Surfaces()
	types := []gopi.SurfaceFlags{
		gopi.SURFACE_FLAG_RGB565, gopi.SURFACE_FLAG_RGB888, gopi.SURFACE_FLAG_RGBA32,
	}
	for _, imageType := range types {
		fmt.Println("DOING", imageType)
		if bitmap, err := surfaces.CreateBitmap(imageType, gopi.Size{100, 100}); err != nil {
			t.Error(err)
		} else {
			// Diagnol blue stripe
			bitmap.ClearToColor(gopi.ColorRed)
			for y := float32(0); y < bitmap.Size().H; y += 1.0 {
				bitmap.PaintPixel(gopi.ColorBlue, gopi.Point{y, y})
			}
		}
		fmt.Println("DONE", imageType)
	}
}
