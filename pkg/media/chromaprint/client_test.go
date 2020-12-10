// +build chromaprint

package chromaprint_test

import (
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	chromaprint "github.com/djthorpe/gopi/v3/pkg/media/chromaprint"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*chromaprint.Manager
	*chromaprint.Client
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Client_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Client == nil {
			t.Error("client is nil")
		} else {
			t.Log(app.Client)
		}
	})
}

func Test_Client_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if matches, err := app.Client.Lookup("AQAAT0mUaEkSRZEGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", 5*time.Second, chromaprint.META_TRACK); err != nil {
			t.Error(err)
		} else if len(matches) != 0 {
			t.Error("Unexpected matches")
		}
	})
}
