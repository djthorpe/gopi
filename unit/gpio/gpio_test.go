// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpio_test

import (
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Modules
	_ "github.com/djthorpe/gopi/v2/unit/platform"
	_ "github.com/djthorpe/gopi/v2/unit/gpio/gpiosysfs"
	_ "github.com/djthorpe/gopi/v2/unit/gpio/gpiorpi"
	_ "github.com/djthorpe/gopi/v2/unit/gpio"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/files"
)

func Test_GPIO_000(t *testing.T) {
	t.Log("Test_GPIO_000")
}

func Test_GPIO_001(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewTestTool(t, Main_Test_GPIO_001, flags, "gpio"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_GPIO_001(app gopi.App, t *testing.T) {
	gpio := app.GPIO()
	if gpio == nil {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("GPIO() failed"))
	} else {
		t.Log(gpio)
	}
}
