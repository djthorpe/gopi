// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpiosysfs_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/files"
	_ "github.com/djthorpe/gopi/v2/unit/gpio/gpiosysfs"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

func Test_GPIO_Sysfs_000(t *testing.T) {
	t.Log("Test_GPIO_Sysfs_000")
}

func Test_GPIO_Sysfs_001(t *testing.T) {
	flags := []string{"-debug", "-gpio.unexport=F"}
	if app, err := app.NewTestTool(t, Main_Test_GPIO_Sysfs_001, flags, "gpio"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_GPIO_Sysfs_001(app gopi.App, t *testing.T) {
	gpio := app.GPIO()
	if gpio == nil {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("GPIO() failed"))
	} else if mode := gpio.GetPinMode(gopi.GPIOPin(13)); mode == gopi.GPIO_NONE {
		t.Error("Unexpected mode for pin", mode)
	} else {
		for i := 0; i < 100; i++ {
			gpio.SetPinMode(gopi.GPIOPin(13), gopi.GPIO_OUTPUT)
			if mode := gpio.GetPinMode(gopi.GPIOPin(13)); mode != gopi.GPIO_OUTPUT {
				t.Error("Unexpected mode for pin", mode)
				break
			}
			gpio.SetPinMode(gopi.GPIOPin(13), gopi.GPIO_INPUT)
			if mode := gpio.GetPinMode(gopi.GPIOPin(13)); mode != gopi.GPIO_INPUT {
				t.Error("Unexpected mode for pin", mode)
				break
			}
		}
	}
}


func Test_GPIO_Sysfs_002(t *testing.T) {
	flags := []string{"-debug", "-gpio.unexport=F"}
	if app, err := app.NewTestTool(t, Main_Test_GPIO_Sysfs_002, flags, "gpio"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_GPIO_Sysfs_002(app gopi.App, t *testing.T) {
	gpio := app.GPIO()
	if gpio == nil {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("GPIO() failed"))
	} else {
		gpio.SetPinMode(gopi.GPIOPin(13),gopi.GPIO_OUTPUT)
		if mode := gpio.GetPinMode(gopi.GPIOPin(13)); mode != gopi.GPIO_OUTPUT {			
			t.Error("Unexpected mode for pin", mode)
		}
		for i := 0; i < 100; i++ {
			gpio.WritePin(gopi.GPIOPin(13), gopi.GPIO_LOW)
			if state := gpio.ReadPin(gopi.GPIOPin(13)); state != gopi.GPIO_LOW {
				t.Error("Unexpected state for pin", state)
				break
			}
			gpio.WritePin(gopi.GPIOPin(13), gopi.GPIO_HIGH)
			if state := gpio.ReadPin(gopi.GPIOPin(13)); state != gopi.GPIO_HIGH {
				t.Error("Unexpected state for pin", state)
				break
			}
		}
	}
}
