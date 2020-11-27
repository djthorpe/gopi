// +build linux

package i2c_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	gopi.I2C
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_I2C_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.I2C == nil {
			t.Error("nil I2C unit")
		} else {
			t.Log(app.I2C)
		}
	})
}

func Test_I2C_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		for _, device := range app.I2C.Devices() {
			for s := uint8(0x00); s <= uint8(0x7F); s++ {
				if available, err := app.I2C.DetectSlave(device, s); err != nil {
					t.Error(err)
				} else {
					t.Logf("<bus %d slave 0x%02X> available=%v", device, s, available)
				}
			}
		}
		t.Log(app.I2C)
	})
}
