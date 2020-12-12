package platform_test

import (
	"context"
	"sync"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/platform"
)

type App struct {
	gopi.Unit
	gopi.Logger
	gopi.Platform
}

func (this *App) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func Test_Power_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Platform == nil {
			t.Error("No Platform object")
		}
	})
}

func Test_Power_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if err := app.Platform.SetPowerState(); err != nil {
			t.Error(err)
		}
	})
}
