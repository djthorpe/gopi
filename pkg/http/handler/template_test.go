package handler_test

import (
	"fmt"
	"net/http"
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/http"
	handler "github.com/djthorpe/gopi/v3/pkg/http/handler"
	"github.com/djthorpe/gopi/v3/pkg/http/renderer"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type TApp struct {
	gopi.Unit
	*handler.Templates
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_templates_001(t *testing.T) {
	tool.Test(t, nil, new(TApp), func(app *TApp) {
		if app.TemplateCache == nil {
			t.Error("nil Templates unit")
		} else {
			t.Log(app.Templates)
		}
	})
}

func Test_templates_002(t *testing.T) {
	tool.Test(t, []string{"-http.templates", TEMPLATES}, new(TApp), func(app *TApp) {
		// Register a text renderer
		if err := app.Templates.RegisterRenderer(renderer.NewTextRenderer(CONTENT, "page.tmpl")); err != nil {
			t.Error(err)
		}
		// Serve content from /
		if err := app.Templates.Serve("/"); err != nil {
			t.Error(err)
		}
		// Start Server
		if err := app.Server.StartInBackground("tcp", ":0"); err != nil {
			t.Error(err)
		} else {
			t.Log(app.Server.Addr())
		}
		// Run through 2000 URLS
		for i := 0; i < 2000; i++ {
			url := fmt.Sprintf("http://%v/page.tmpl?%v", app.Server.Addr(), i)
			if resp, err := http.Get(url); err != nil {
				t.Error(err)
			} else if resp.StatusCode != http.StatusOK {
				t.Error("Unexpected response", resp.Status)
			} else {
				resp.Body.Close()
				t.Log(url, "=>", resp)
			}
		}
	})
}

func Test_templates_003(t *testing.T) {
	tool.Test(t, []string{"-http.templates", TEMPLATES}, new(TApp), func(app *TApp) {
		// Register a text renderer
		if err := app.Templates.RegisterRenderer(renderer.NewTextRenderer(CONTENT, "page.tmpl")); err != nil {
			t.Error(err)
		}
		// Serve content from /
		if err := app.Templates.Serve("/"); err != nil {
			t.Error(err)
		}
		// Start Server
		if err := app.Server.StartInBackground("tcp", ":0"); err != nil {
			t.Error(err)
		} else {
			t.Log(app.Server.Addr())
		}
		// Run through 100 URLS
		url := fmt.Sprintf("http://%v/page.tmpl?static", app.Server.Addr())
		for i := 0; i < 100; i++ {
			if resp, err := http.Get(url); err != nil {
				t.Error(err)
			} else if resp.StatusCode != http.StatusOK {
				t.Error("Unexpected response", resp.Status)
			} else {
				resp.Body.Close()
				t.Log(url, "=>", resp)
			}
		}
	})
}
