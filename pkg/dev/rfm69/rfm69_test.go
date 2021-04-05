package rfm69_test

import (
	"context"
	"testing"

	// Modules
	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/dev/rfm69"
	"github.com/djthorpe/gopi/v3/pkg/tool"

	// Dependencies
	_ "github.com/djthorpe/gopi/v3/pkg/hw/spi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*rfm69.RFM69
}

func (app *App) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_RFM69_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.RFM69 == nil {
			t.Error("nil RFM69 unit")
		} else {
			t.Log(app.RFM69)
		}
	})
}
