// +build linux

package sysfs_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi/v3"
	sysfs "github.com/djthorpe/gopi/v3/pkg/hw/gpio/sysfs"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"

	_ "github.com/djthorpe/gopi/v3/pkg/file"
)

type App struct {
	gopi.Unit
	*sysfs.GPIO
}

func Test_Sysfs_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.GPIO == nil {
			t.Error("GPIO is nil")
		} else {
			t.Log(app.GPIO)
		}
	})
}

func Test_Sysfs_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		v := app.GPIO.ReadPin(gopi.GPIOPin(4))
		t.Log(v)
		t.Log(app.GPIO)
	})
}
