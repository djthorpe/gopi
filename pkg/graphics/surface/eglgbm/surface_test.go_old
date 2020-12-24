// +build egl,gbm,drm

package surface_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	surface "github.com/djthorpe/gopi/v3/pkg/graphics/surface/eglgbm"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/graphics/display/drm"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*surface.Manager
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Surface_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Manager == nil {
			t.Error("nil SurfaceManager unit")
		} else {
			t.Log(app.Manager)
		}
	})
}

func Test_Surface_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if surface, err := app.Manager.CreateSurface(100, 100); err != nil {
			t.Error(err)
		} else {
			t.Log("surface=", surface)
		}
	})
}
