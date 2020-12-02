package keycode_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/hw/lirc"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.LIRCKeycodeManager
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Manager_001(t *testing.T) {
	tool.Test(t, []string{"-lirc.db=../../../../etc/keycode"}, new(App), func(app *App) {
		if app.LIRCKeycodeManager == nil {
			t.Error("nil LIRCKeycodeManager unit")
		} else {
			t.Log(app.LIRCKeycodeManager)
		}
	})
}

func Test_Manager_002(t *testing.T) {
	tool.Test(t, []string{"-lirc.db=../../../../etc/keycode"}, new(App), func(app *App) {
		names := []string{"power", "up", "down", "1"}
		for _, name := range names {
			if keycodes := app.LIRCKeycodeManager.Keycode(name); len(keycodes) > 0 {
				t.Log(name, "=>", keycodes)
			} else {
				t.Error("No keycodes found for", name)
			}
		}
	})
}
