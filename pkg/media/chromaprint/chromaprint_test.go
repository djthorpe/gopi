// +build chromaprint

package chromaprint_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	chromaprint "github.com/djthorpe/gopi/v3/pkg/media/chromaprint"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

const (
	SAMPLE_FILE = "../../../etc/sample.mp4"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*chromaprint.Manager
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Chromaprint_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Manager == nil {
			t.Error("manager is nil")
		} else {
			t.Log(app.Manager)
		}
	})
}
