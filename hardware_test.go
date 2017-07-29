package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/mock"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE HARDWARE MODULE

func TestHardware_000(t *testing.T) {
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_HARDWARE)
	config.Debug = true

	// Create an application with a hardware module
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else {
		app.Logger.Info("hardware=%v", app.Hardware)
	}
}

func TestHardware_001(t *testing.T) {
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_HARDWARE)
	config.Debug = true

	// Create an application with a hardware module
	if app, err := gopi.NewAppInstance(config); err != nil {
		t.Error(err)
	} else {
		app.Logger.Info("hardware=%v", app.Hardware)
	}
}
