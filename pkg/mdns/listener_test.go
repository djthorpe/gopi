package mdns_test

import (
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/event"
	mdns "github.com/djthorpe/gopi/v3/pkg/mdns"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

type App struct {
	gopi.Unit
	*mdns.Listener
}

func Test_Listener_001(t *testing.T) {
	tool.Test(t, []string{"-mdns.domain=test"}, new(App), func(app *App) {
		if app.Listener == nil {
			t.Error("Expected non-nil listener")
		}
		if domain := app.Listener.Domain(); domain != "test." {
			t.Errorf("Unexpected domain: %q", domain)
		}
	})
}

func Test_Listener_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if domain := app.Listener.Domain(); domain != "local." {
			t.Errorf("Unexpected domain: %q", domain)
		}
		t.Log(app.Listener)
		time.Sleep(5 * time.Second)
	})
}
