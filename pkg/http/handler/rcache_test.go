package handler_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/http/handler"
	"github.com/djthorpe/gopi/v3/pkg/http/renderer"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type RApp struct {
	gopi.Unit
	*handler.RenderCache
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	CONTENT = "../../../etc/http"
)

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_rcache_001(t *testing.T) {
	tool.Test(t, nil, new(RApp), func(app *RApp) {
		if app.RenderCache == nil {
			t.Error("nil RenderCache unit")
		} else {
			t.Log(app.RenderCache)
		}
	})
}

func Test_rcache_002(t *testing.T) {
	tool.Test(t, nil, new(RApp), func(app *RApp) {
		if err := app.RenderCache.Register(renderer.NewTextRenderer(CONTENT, "page.tmpl")); err != nil {
			t.Error(err)
		} else {
			t.Log(app.RenderCache)
		}
	})
}

func Test_rcache_003(t *testing.T) {
	tool.Test(t, nil, new(RApp), func(app *RApp) {
		app.RenderCache.Register(renderer.NewTextRenderer(CONTENT, "page.tmpl"))
		for i := 0; i < 2000; i++ {
			url := fmt.Sprint("http://localhost/page.tmpl?", i)
			if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
				t.Error(err)
			} else if r := app.RenderCache.Get(req); r == nil {
				t.Error("Expected renderer for", url)
			} else {
				t.Log(url, "=>", r)
			}
		}
	})
}

func Test_rcache_004(t *testing.T) {
	tool.Test(t, nil, new(RApp), func(app *RApp) {
		app.RenderCache.Register(renderer.NewTextRenderer(CONTENT, "page.tmpl"))
		for i := 0; i < 2000; i++ {
			url := fmt.Sprint("http://localhost/page.tmpl?", i)
			if req, err := http.NewRequest(http.MethodGet, url, nil); err != nil {
				t.Error(err)
			} else if r := app.RenderCache.Get(req); r == nil {
				t.Error("Expected renderer for", url)
			} else if ctx, err := app.RenderCache.Render(r, req); err != nil {
				t.Error(r)
			} else if ctx.Content == nil {
				t.Error("Expected content to be returned")
			} else if ctx.Modified.IsZero() {
				t.Error("Expected modified field to be returned")
			} else {
				t.Log("Req ", url)
				t.Log("  Template ", ctx.Template)
				t.Log("  Type ", ctx.Type)
				t.Log("  Modified ", ctx.Modified)
			}
		}
	})
}
