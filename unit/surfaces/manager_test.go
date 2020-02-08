/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces_test

import (
	"image/png"
	"io/ioutil"
	"testing"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/display"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/platform"
	_ "github.com/djthorpe/gopi/v2/unit/surfaces"
)

func Test_Manager_000(t *testing.T) {
	t.Log("Test_Manager_000")
}

/////////////////////////////////////////////////////////////////////
// Create a bitmap and write to file

func Test_Manager_001(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Manager_001, nil, "surfaces"); err != nil {
		t.Error(err)
	} else {
		app.Run()
	}
}

func Main_Test_Manager_001(app gopi.App, t *testing.T) {
	surfaces := app.Surfaces()
	t.Log(surfaces)
	if bm, err := surfaces.CreateBitmap(0, gopi.Size{100, 100}); err != nil {
		t.Error(err)
	} else if f, err := ioutil.TempFile("", "*.png"); err != nil {
		t.Error(err)
	} else {
		bm.ClearToColor(gopi.ColorRed)
		if err := png.Encode(f, bm); err != nil {
			t.Error(err)
		} else {
			t.Log(f.Name())
			f.Close()
		}
	}
}

/////////////////////////////////////////////////////////////////////
// Create a surface with an existing bitmap

func Test_Manager_002(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Manager_002, nil, "surfaces"); err != nil {
		t.Error(err)
	} else {
		app.Run()
	}
}

func Main_Test_Manager_002(app gopi.App, t *testing.T) {
	var surface gopi.Surface
	app.Surfaces().Do(func() error {
		if bm, err := app.Surfaces().CreateBitmap(0, gopi.Size{100, 100}); err != nil {
			t.Error(err)
		} else if surface, err = app.Surfaces().CreateSurfaceWithBitmap(bm, 0, 1.0, 1, gopi.ZeroPoint, gopi.ZeroSize); err != nil {
			t.Error(err)
		} else {
			bm.ClearToColor(gopi.ColorPurple)
		}
		return nil
	})
	time.Sleep(time.Second)
	app.Surfaces().Do(func() error {
		return app.Surfaces().DestroySurface(surface)
	})
}

/////////////////////////////////////////////////////////////////////
// Create background with new bitmap and paint a circle

func Test_Manager_003(t *testing.T) {
	if app, err := app.NewTestTool(t, Main_Test_Manager_003, nil, "surfaces"); err != nil {
		t.Error(err)
	} else {
		app.Run()
	}
}

func Main_Test_Manager_003(app gopi.App, t *testing.T) {
	app.Surfaces().Do(func() error {
		if surface, err := app.Surfaces().CreateBackground(0, 1.0); err != nil {
			t.Error(err)
		} else {
			surface.Bitmap().ClearToColor(gopi.ColorBlue)
			surface.Bitmap().Line(gopi.ColorWhite, gopi.ZeroPoint, gopi.Point{100, 100})
		}
		return nil
	})
	time.Sleep(time.Second)
}
