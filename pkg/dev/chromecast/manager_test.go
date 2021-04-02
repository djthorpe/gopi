package chromecast_test

import (
	"context"
	"testing"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	cast "github.com/djthorpe/gopi/v3/pkg/dev/chromecast"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	// Dependencies
	_ "github.com/djthorpe/gopi/v3/pkg/event"
	_ "github.com/djthorpe/gopi/v3/pkg/mdns"
)

type CastApp struct {
	gopi.Unit
	gopi.Logger
	*cast.Manager
}

func Test_Cast_000(t *testing.T) {
	tool.Test(t, nil, new(CastApp), func(app *CastApp) {
		if app.Manager == nil {
			t.Fatal("app.Manager == nil")
		}
	})
}

func Test_Cast_001(t *testing.T) {
	tool.Test(t, nil, new(CastApp), func(app *CastApp) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if casts, err := app.Devices(ctx); err != nil {
			t.Error(err)
		} else {
			t.Log(casts)
		}
	})
}
