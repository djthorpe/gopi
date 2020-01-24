// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpiorpi_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/gpio/gpiorpi"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/platform"
)

func Test_GPIO_RPI_000(t *testing.T) {
	t.Log("Test_GPIO_RPI_000")
}

func Test_GPIO_RPI_001(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewTestTool(t, Main_Test_GPIO_RPI_001, flags, "gpio"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_GPIO_RPI_001(app gopi.App, t *testing.T) {
	gpio := app.GPIO()
	if gpio == nil {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("GPIO() failed"))
	} else if numPins := gpio.NumberOfPhysicalPins(); numPins == 0 {
		t.Error("Expected numPins > 0")
	} else {
		t.Log(gpio)
	}
}


func Test_GPIO_RPI_002(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewTestTool(t, Main_Test_GPIO_RPI_002, flags, "gpio"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_GPIO_RPI_002(app gopi.App, t *testing.T) {
	gpio := app.GPIO()
	for _,pin := range gpio.Pins() {
		physical := gpio.PhysicalPinForPin(pin)
		if physical == 0 {
			continue
		}
		logical := gpio.PhysicalPin(physical)
		if logical != pin {
			t.Error("Bad mapping between",pin,physical,logical)
		} else {
			t.Log(pin,"=>",physical)
		}
	}
}


func Test_GPIO_RPI_003(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewTestTool(t, Main_Test_GPIO_RPI_003, flags, "gpio"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_GPIO_RPI_003(app gopi.App, t *testing.T) {
	gpio := app.GPIO()
	for _,pin := range gpio.Pins() {
		t.Log(pin,"=>",gpio.GetPinMode(pin))
	}
}
