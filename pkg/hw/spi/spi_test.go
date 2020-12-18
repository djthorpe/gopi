// +build linux

package spi_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.SPI
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_SPI_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.SPI == nil {
			t.Error("nil SPI unit")
		} else {
			t.Log(app.SPI)
		}
	})
}

func Test_SPI_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		for _, dev := range app.SPI.Devices() {
			if err := app.SPI.SetMode(dev, gopi.SPI_MODE_0); err != nil {
				t.Error(err)
			} else {
				t.Log(dev)
			}
		}
		t.Log(app.SPI)
	})
}
