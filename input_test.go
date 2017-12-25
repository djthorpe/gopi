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
			t.Fatal("Unable to create input driver")
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
		if device, err := gopi.Open(input.Device{}, logger.(gopi.Logger)); err != nil {
			t.Fatal("Unable to create input device")
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
