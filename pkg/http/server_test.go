package http_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/http/renderer"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.Server
	gopi.HttpTemplate
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	HTTP_ROOT = "../../etc/http"
)

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Server_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Server == nil {
			t.Error("nil Server unit")
		} else if app.HttpTemplate == nil {
			t.Error("nil HttpTemplate unit")
		} else {
			t.Log(app.Server)
		}
	})
}

func Test_Server_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if err := app.Server.StartInBackground("tcp", ":0"); err != nil {
			t.Error(err)
		}
		t.Log(app.Server)
		if err := app.Server.Stop(true); err != nil {
			t.Error(err)
		}
	})
}

func Test_Server_003(t *testing.T) {
	tool.Test(t, []string{"-http.templates", HTTP_ROOT}, new(App), func(app *App) {
		if err := app.Server.StartInBackground("tcp", ":0"); err != nil {
			t.Error(err)
		} else if err := app.HttpTemplate.RegisterRenderer(renderer.NewTextRenderer(HTTP_ROOT, "page.tmpl")); err != nil {
			t.Error(err)
		} else if err := app.HttpTemplate.RegisterRenderer(renderer.NewIndexRenderer(HTTP_ROOT, "index.tmpl")); err != nil {
			t.Error(err)
		} else if err := app.HttpTemplate.Serve("/"); err != nil {
			t.Error(err)
		}

		// Make requests for pages and index
		for i := 0; i < 2000; i++ {
			var url string
			switch i % 7 {
			case 0:
				url = fmt.Sprintf("http://%v/page.tmpl?%v", app.Server.Addr(), i)
			case 1:
				url = fmt.Sprintf("http://%v/index.tmpl?%v", app.Server.Addr(), i)
			case 2:
				url = fmt.Sprintf("http://%v/index.tmpl?static", app.Server.Addr())
			case 4:
				url = fmt.Sprintf("http://%v", app.Server.Addr())
			case 5:
				url = fmt.Sprintf("http://%v?%v", app.Server.Addr(), i)
			default:
				url = fmt.Sprintf("http://%v/?%v", app.Server.Addr(), i)
			}

			if resp, err := http.Get(url); err != nil {
				t.Error(err)
			} else {
				resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					t.Error("Unexpected response", resp.Status)
				} else {
					t.Log(url, "=>", resp)
				}
			}
		}
	})
}
