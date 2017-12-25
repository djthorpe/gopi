package gopi_test

import (
	"testing"

	gopi "github.com/djthorpe/gopi"
	input "github.com/djthorpe/gopi/sys/input/mock"
	logger "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////
// CREATE INPUT

func TestCreateConfig_000(t *testing.T) {
	if logger, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Fatal("Unable to create logger driver")
	} else {
		defer logger.Close()
		if input, err := gopi.Open(input.Input{}, logger.(gopi.Logger)); err != nil {
			t.Fatal("Unable to create input driver:", err)
		} else {
			defer input.Close()

			if input2, ok := input.(gopi.Input); ok == false {
				t.Fatal("Unable to cast input driver to gopi.Input")
			} else {
				t.Log("Input Driver=", input2)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// CREATE DEVICE

func TestCreateDevice_000(t *testing.T) {
	if logger, err := gopi.Open(logger.Config{}, nil); err != nil {
		t.Fatal("Unable to create logger driver")
	} else {
		defer logger.Close()
		if device, err := gopi.Open(input.Device{
			Name: "test", Type: gopi.INPUT_TYPE_KEYBOARD, Bus: gopi.INPUT_BUS_NONE,
		}, logger.(gopi.Logger)); err != nil {
			t.Fatal("Unable to create input device:", err)
		} else {
			defer device.Close()

			if device2, ok := device.(gopi.InputDevice); ok == false {
				t.Fatal("Unable to cast input driver to gopi.InputDevice")
			} else {
				t.Log("Input Device=", device2)
			}
		}
	}
}

func TestCreateInput_001(t *testing.T) {
	logger, err := gopi.Open(logger.Config{}, nil)
	if err != nil {
		t.Fatal("Unable to create logger driver")
	}
	defer logger.Close()
	/*driver := CreateDriver(t, logger.(gopi.Logger))
	if driver == nil {
		t.Fatal("Unable to create input driver")
	}
	defer driver.Close()*/
	keyboard := CreateDevice(t, logger.(gopi.Logger), "keyboard", gopi.INPUT_TYPE_KEYBOARD, gopi.INPUT_BUS_USB)
	if keyboard == nil {
		t.Fatal("Unable to create keyboard")
	}
	defer keyboard.Close()

	/* check keyboard parameters */
	if keyboard.Bus() != gopi.INPUT_BUS_USB {
		t.Error("Unexpected bus")
	}
	if keyboard.Type() != gopi.INPUT_TYPE_KEYBOARD {
		t.Error("Unexpected type")
	}
	if keyboard.Name() != "keyboard" {
		t.Error("Unexpected name")
	}
}

func TestCreateInput_002(t *testing.T) {
	logger, err := gopi.Open(logger.Config{}, nil)
	if err != nil {
		t.Fatal("Unable to create logger driver")
	}
	defer logger.Close()
	/*driver := CreateDriver(t, logger.(gopi.Logger))
	if driver == nil {
		t.Fatal("Unable to create input driver")
	}
	defer driver.Close()*/
	mouse := CreateDevice(t, logger.(gopi.Logger), "mouse", gopi.INPUT_TYPE_MOUSE, gopi.INPUT_BUS_NONE)
	if mouse == nil {
		t.Fatal("Unable to create mouse")
	}
	defer mouse.Close()

	/* check mouse parameters */
	if mouse.Bus() != gopi.INPUT_BUS_NONE {
		t.Error("Unexpected bus")
	}
	if mouse.Type() != gopi.INPUT_TYPE_MOUSE {
		t.Error("Unexpected type")
	}
	if mouse.Name() != "mouse" {
		t.Error("Unexpected name")
	}
	mouse.SetPosition(gopi.Point{-1, -1})
	t.Log("Position=", mouse.Position())
	if mouse.Position().X != -1 && mouse.Position().Y != -1 {
		t.Error("Unexpected position,", mouse.Position())
	}
	mouse.SetPosition(gopi.Point{-2, -2})
	if mouse.Position().X != -2 && mouse.Position().Y != -2 {
		t.Error("Unexpected position,", mouse.Position())
	}
}

func TestCreateInput_003(t *testing.T) {
	logger, err := gopi.Open(logger.Config{}, nil)
	if err != nil {
		t.Fatal("Unable to create logger driver")
	}
	defer logger.Close()

	// Create driver
	driver := CreateDriver(t, logger.(gopi.Logger))
	if driver == nil {
		t.Fatal("Unable to create input driver")
	}
	defer driver.Close()

	// Add mouse
	mouse := CreateDevice(t, logger.(gopi.Logger), "mouse", gopi.INPUT_TYPE_MOUSE, gopi.INPUT_BUS_NONE)
	if mouse == nil {
		t.Fatal("Unable to create mouse")
	}
	if err := driver.AddDevice(mouse); err != nil {
		t.Fatal("Unable to add mouse:", err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func CreateDriver(t *testing.T, logger gopi.Logger) gopi.Input {
	if driver, err := gopi.Open(input.Input{}, logger); err != nil {
		t.Fatal("Unable to create input driver:", err)
		return nil
	} else {
		return driver.(gopi.Input)
	}
}

func CreateDevice(t *testing.T, logger gopi.Logger, name string, device_type gopi.InputDeviceType, device_bus gopi.InputDeviceBus) gopi.InputDevice {
	if device, err := gopi.Open(input.Device{
		Name: name, Type: device_type, Bus: device_bus,
	}, logger); err != nil {
		t.Fatal("Unable to create input device:", err)
		return nil
	} else {
		return device.(gopi.InputDevice)
	}
}
