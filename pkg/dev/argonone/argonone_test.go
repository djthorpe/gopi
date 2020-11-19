// +build linux

package argonone_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.ArgonOne
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_ArgonOne_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.ArgonOne == nil {
			t.Error("nil ArgonOne unit")
		} else {
			t.Log(app.ArgonOne)
		}
	})
}

func Test_ArgonOne_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		for fan := uint8(0); fan <= uint8(100); fan++ {
			if err := app.ArgonOne.SetFan(fan); err != nil {
				t.Error(err)
			} else {
				t.Log("Set fan=", fan, "%")
			}
		}
		// Switch fan off
		if err := app.ArgonOne.SetFan(0); err != nil {
			t.Error(err)
		}
	})
}

func Test_ArgonOne_003(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if err := app.ArgonOne.SetPower(gopi.ARGONONE_POWER_ALWAYSON); err != nil {
			t.Error(err)
		} else if err := app.ArgonOne.SetPower(gopi.ARGONONE_POWER_DEFAULT); err != nil {
			t.Error(err)
		}
	})
}

func Test_ArgonOne_004(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if err := app.ArgonOne.SetPower(gopi.ARGONONE_POWER_ALWAYSON); err != nil {
			t.Error(err)
		} else if err := app.ArgonOne.SetPower(gopi.ARGONONE_POWER_DEFAULT); err != nil {
			t.Error(err)
		}
	})
}

func Test_ArgonOne_005(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		t.Log(app.ArgonOne)
	})
}
