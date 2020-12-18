// +build linux

package waveshare_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/dev/waveshare"
	"github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/hw/gpio"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/spi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	*waveshare.EPD
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_EPD_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.EPD == nil {
			t.Error("nil EPD unit")
		} else {
			t.Log(app.EPD)
		}
	})
}
