/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces_test

import (
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

func Test_Surfaces_000(t *testing.T) {
	t.Log("Test_Surfaces_000")
}

func Test_Surfaces_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Surfaces_001, []string{"-debug"}, "surfaces"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Surfaces_001(app gopi.App, t *testing.T) {
	surfaces := app.Surfaces()
	if display := surfaces.Display(); display == nil {
		t.Error("Unexpected Display() value")
	} else if name := surfaces.Name(); name == "" {
		t.Error("Unexpected Name() value")
	} else if types := surfaces.Types(); types == nil {
		t.Error("Unexpected Types() value")
	} else {
		t.Log("display=", display)
		t.Log("name=", name)
		t.Log("types=", types)
	}
}

func Test_Surfaces_002(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Surfaces_002, []string{"-debug"}, "surfaces"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_Surfaces_002(app gopi.App, t *testing.T) {
	surfaces := app.Surfaces()
	types := []gopi.SurfaceFlags{
		gopi.SURFACE_FLAG_RGB565, gopi.SURFACE_FLAG_RGB888, gopi.SURFACE_FLAG_RGBA32,
	}
	size := gopi.Size{100, 100}
	for _, image_type := range types {
		if bitmap, err := surfaces.CreateBitmap(image_type, size); err != nil {
			t.Error(err)
		} else if bs := bitmap.Size(); bs.W != size.W || bs.H != size.H {
			t.Error("Unexpected size")
		} else if bitmap.Type() != image_type {
			t.Error("Unexpected type")
		} else {
			t.Log(bitmap)
		}
	}

	for _, image_type := range types {
		if bitmap, err := surfaces.CreateSnapshot(image_type); err != nil {
			t.Error(err)
		} else if bitmap.Type() != image_type {
			t.Error("Unexpected type")
		} else {
			t.Log(bitmap)
		}
	}
}
