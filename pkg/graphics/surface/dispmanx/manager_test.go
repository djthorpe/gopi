// +build dispmanx,rpi,egl

package dispmanx_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
	surface "github.com/djthorpe/gopi/v3/pkg/graphics/surface/dispmanx"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	// Units
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

func Test_Manager_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Manager == nil {
			t.Error("nil SurfaceManager unit")
		} else {
			t.Log("manager=", app.Manager)
		}
	})
}

/*
func Test_Manager_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		err := app.Manager.Do(func(ctx gopi.GraphicsContext) error {
			t.Log("ctx=", ctx)
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	})
}

func Test_Manager_003(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		err := app.Manager.Do(func(ctx gopi.GraphicsContext) error {
			if surface, err := app.Manager.CreateSurface(ctx, gopi.SURFACE_FLAG_BITMAP, 1.0, 0, gopi.Point{0, 0}, gopi.Size{100, 100}); err != nil {
				t.Error("CreateSurface error:", err)
				return err
			} else {
				t.Log("surface=", surface)
				surface.Bitmap().ClearToColor(color.RGBA{0xFF, 0x00, 0x00, 0xFF})
			}
			return nil
		})
		if err != nil {
			t.Error(err)
		}
	})
}
*/
