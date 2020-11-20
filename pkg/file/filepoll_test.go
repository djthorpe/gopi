// +build linux

package file_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.FilePoll
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_FilePoll_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.FilePoll == nil {
			t.Error("nil FilePoll unit")
		} else {
			t.Log(app.FilePoll)
		}
	})
}
