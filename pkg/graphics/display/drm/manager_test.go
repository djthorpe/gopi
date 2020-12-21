// +build drm

package display_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	display "github.com/djthorpe/gopi/v3/pkg/graphics/display/drm"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*display.Manager
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Display_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Manager == nil {
			t.Error("nil DisplayManager unit")
		} else {
			t.Log("manager=", app.Manager)
		}
	})
}

func Test_Display_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		displays := app.Manager.Displays()
		t.Log("displays=", displays)
	})
}
