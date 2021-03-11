package rotel_test

import (
	"context"
	"testing"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/event"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.RotelManager
	gopi.Publisher
	gopi.Logger
}

func (app *App) Run(ctx context.Context) error {
	ch := app.Publisher.Subscribe()
	defer app.Publisher.Unsubscribe(ch)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt := <-ch:
			app.Print(evt)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Manager_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.RotelManager == nil {
			t.Error("nil RotelManager unit")
		} else {
			t.Log(app.RotelManager)
		}
	})
}

func Test_Manager_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		t.Log(app.RotelManager)
		time.Sleep(5 * time.Second)
	})
}
