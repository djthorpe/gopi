// +build egl,gbm,drm

package surface_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	surface "github.com/djthorpe/gopi/v3/pkg/graphics/surface/eglgbm"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*surface.Manager
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Manager_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Manager == nil {
			t.Error("nil SurfaceManager unit")
		} else {
			t.Log("manager=", app.Manager)
		}
	})
}
