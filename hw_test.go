package gopi_test

import (
	"testing"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/mock"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE HARDWARE MODULE
func TestHardware_000(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("hw")
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
		return
	}
	if app == nil {
		t.Fatal("Expecting app object")
		return
	}
	if app.Logger == nil {
		t.Fatal("Expecting app.Logger object")
		return
	}
	if app.Hardware == nil {
		t.Fatal("Expecting app.Hardware object")
		return
	}
	app.Logger.Info("hardware=%v", app.Hardware)
}

func TestHardware_001(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("hw")
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}

	// Get name and serial number, neither should be empty
	if app.Hardware.Name() == "" {
		t.Fatal("Expecting a name")
	}
	if app.Hardware.SerialNumber() == "" {
		t.Fatal("Expecting a serial number")
	}
}

////////////////////////////////////////////////////////////////////////////////
// CREATE DISPLAY MODULE

func TestDisplay_001(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("hw", "display")
	config.Debug = true
	//	config.AppFlags.SetUint("display", 0)

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}

	// Get name and serial number, neither should be empty
	if app.Display == nil {
		t.Fatal("Expecting app.Display object")
	}

	app.Logger.Info("app=%v", app)
	app.Logger.Info("display=%v", app.Display)
}

////////////////////////////////////////////////////////////////////////////////
// CREATE GPIO MODULE
