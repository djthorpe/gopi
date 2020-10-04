// +build rpi

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
	for _, pin := range gpio.Pins() {
		physical := gpio.PhysicalPinForPin(pin)
		if physical == 0 {
			continue
		}
		logical := gpio.PhysicalPin(physical)
		if logical != pin {
			t.Error("Bad mapping between", pin, physical, logical)
		} else {
			t.Log(pin, "=>", physical)
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
	for _, pin := range gpio.Pins() {
		t.Log(pin, "=>", gpio.GetPinMode(pin))
	}
}

func Test_GPIO_RPI_004(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewTestTool(t, Main_Test_GPIO_RPI_004, flags, "gpio"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_GPIO_RPI_004(app gopi.App, t *testing.T) {
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

func Test_GPIO_RPI_005(t *testing.T) {
	flags := []string{"-debug"}
	if app, err := app.NewTestTool(t, Main_Test_GPIO_RPI_005, flags, "gpio"); err != nil {
		t.Error(err)
	} else if returnCode := app.Run(); returnCode != 0 {
		t.Error("Unexpected return code", returnCode)
	}
}

func Main_Test_GPIO_RPI_005(app gopi.App, t *testing.T) {
	gpio := app.GPIO()
	if gpio == nil {
		t.Fatal(gopi.ErrInternalAppError.WithPrefix("GPIO() failed"))
	} else {
		gpio.SetPinMode(gopi.GPIOPin(13), gopi.GPIO_OUTPUT)
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
