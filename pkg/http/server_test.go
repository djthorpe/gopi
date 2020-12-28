package http_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.Server
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Server_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Server == nil {
			t.Error("nil Server unit")
		} else {
			t.Log(app.Server)
		}
	})
}

func Test_Server_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if err := app.Server.StartInBackground("tcp", ""); err != nil {
			t.Error(err)
		} else if err := app.Server.Stop(false); err != nil {
			t.Error(err)
		}
	})
}
