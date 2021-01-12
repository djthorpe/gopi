package handler_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/http/handler"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*handler.TemplateCache
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	TEMPLATES = "../../../etc/http"
)

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_tcache_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.TemplateCache == nil {
			t.Error("nil TemplateCache unit")
		} else {
			t.Log(app.TemplateCache)
		}
	})
}

func Test_tcache_002(t *testing.T) {
	tool.Test(t, []string{"-http.templates", TEMPLATES}, new(App), func(app *App) {
		if tmpl, _, err := app.TemplateCache.Get("page.tmpl"); err != nil {
			t.Error(err)
		} else {
			t.Log(tmpl)
		}
		if tmpl, _, err := app.TemplateCache.Get("index.tmpl"); err != nil {
			t.Error(err)
		} else {
			t.Log(tmpl)
		}
	})
}
