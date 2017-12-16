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
// GPIO MODULE

func TestGPIO_001(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("gpio")
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	// Get GPIO
	if app.GPIO == nil {
		t.Fatal("Expecting app.GPIO object")
	}

	app.Logger.Info("app=%v", app)
	app.Logger.Info("gpio=%v", app.GPIO)
}

func TestGPIO_002(t *testing.T) {
	// Create a configuration with debug
	config := gopi.NewAppConfig("gpio")
	config.Debug = true

	// Create an application with a hardware module
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		t.Fatal(err)
	}
	defer app.Close()

	// Get GPIO
	if app.GPIO == nil {
		t.Fatal("Expecting app.GPIO object")
	}

	// Read all pins
	pins := app.GPIO.Pins()
	if len(pins) == 0 {
		t.Error("Expecting non-zero pins array")
	}

	// Expect number of pins less than value
	if app.GPIO.NumberOfPhysicalPins() == 0 {
		t.Error("Expecting non-zero return from NumberOfPhysicalPins")
	}
	if app.GPIO.NumberOfPhysicalPins() < uint(len(pins)) {
		t.Error("Expecting NumberOfPhysicalPins to be greater or equal to number of logical pins")
	}

	// Print out the pins
	for _, pin := range pins {
		physical := app.GPIO.PhysicalPinForPin(pin)
		app.Logger.Info("pin=%v physical=%v", pin, physical)
	}

	// Set mode for pin
	for _, pin := range pins {
		physical := app.GPIO.PhysicalPinForPin(pin)
		app.Logger.Info("pin=%v physical=%v", pin, physical)

		expected_mode := gopi.GPIO_ALT0
		app.GPIO.SetPinMode(pin, expected_mode)
		if mode := app.GPIO.GetPinMode(pin); mode != expected_mode {
			t.Error("For pin %v, Expecting mode=%v got mode=%v", pin, mode, expected_mode)
		}
	}

	// Set state low for pin
	for _, pin := range pins {
		physical := app.GPIO.PhysicalPinForPin(pin)
		app.Logger.Info("pin=%v physical=%v", pin, physical)

		expected_state := gopi.GPIO_LOW
		app.GPIO.WritePin(pin, expected_state)
		if state := app.GPIO.ReadPin(pin); state != expected_state {
			t.Error("For pin %v, Expecting state=%v got state=%v", pin, state, expected_state)
		}
	}

	// Set state high for pin
	for _, pin := range pins {
		physical := app.GPIO.PhysicalPinForPin(pin)
		app.Logger.Info("pin=%v physical=%v", pin, physical)

		expected_state := gopi.GPIO_HIGH
		app.GPIO.WritePin(pin, expected_state)
		if state := app.GPIO.ReadPin(pin); state != expected_state {
			t.Error("For pin %v, Expecting state=%v got state=%v", pin, state, expected_state)
		}
	}

}
