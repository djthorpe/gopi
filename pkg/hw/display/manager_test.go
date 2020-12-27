// +build rpi

package display_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/hw/display"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.DisplayManager
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Manager_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.DisplayManager == nil {
			t.Error("nil DisplayManager unit")
		} else {
			t.Log(app.DisplayManager)
		}
	})
}

func Test_Manager_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		for _, display := range app.DisplayManager.Displays() {
			t.Log(display)
		}
	})
}

func Test_Manager_003(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		//displays := app.DisplayManager.Displays()

		if disp, err := app.DisplayManager.Display(7); err != nil {
			t.Error("Display", err)
		} else {
			t.Log("Powering on", disp)
			if err := app.DisplayManager.PowerOn(disp); err != nil {
				t.Error("PowerOn", err)
			}
		}
		/*

			if len(displays) == 0 {
				t.Log("Skipping tests due to no displays")
			}
			for _, display := range displays {
				if err := app.DisplayManager.PowerOff(display); err != nil {
					t.Error(err)
				} else {
					t.Log("powered off", display)
				}
			}
			time.Sleep(time.Second)
			for _, display := range displays {
				if err := app.DisplayManager.PowerOn(display); err != nil {
					t.Error(err)
				} else {
					t.Log("powered on", display)
				}
			}*/
	})
}
