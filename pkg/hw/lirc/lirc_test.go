// +build linux

package lirc_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.LIRC
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_LIRC_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.LIRC == nil {
			t.Error("nil LIRC unit")
		} else {
			t.Log(app.LIRC)
		}
	})
}
