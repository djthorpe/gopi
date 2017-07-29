package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/mock"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE HARDWARE MODULE

func TestHardware_000(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_HARDWARE)
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Error(err)
	}
	if app == nil {
		t.Error("Expecting app object")
	}
	if app.Logger == nil {
		t.Error("Expecting app.Logger object")
	}
	if app.Hardware == nil {
		t.Error("Expecting app.Hardware object")
	}
	app.Logger.Info("hardware=%v", app.Hardware)
}

func TestHardware_001(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig(gopi.MODULE_TYPE_HARDWARE)
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Error(err)
	}

	// Get name and serial number, neither should be empty
	if app.Hardware.Name() == "" {
		t.Error("Expecting a name")
	}
	if app.Hardware.SerialNumber() == "" {
		t.Error("Expecting a serial number")
	}
}
