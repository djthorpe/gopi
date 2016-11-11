/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This sample program shows how you can get information from the GPIO
// device, and set pins input, output, etc. To use the software, there are
// some flags. To enumerate the pins on the device with their current status
// and names:
//
//   gpioctrl
//
// To get status on an individual pin:
//
//   gpioctrl -pin <pin>
//
// The pin can be queried by the physical pin number or by the name of the pin,
// for example GPIO23. To set a pin to OUTPUT and set output high and low:
//
//   gpioctrl -pin <pin> -high
//
//   gpioctrl -pin <pin> -low
//
// To set a pin to INPUT:
//
//   gpioctrl -pin <pin> -input
//
// You can also set a pin mode to be alternate function:
//
//  gpioctrl -pin <pin> -alt 0
//
// And so forth...
//
package main

import (
	"fmt"
	"os"
	"errors"
)

import (
	app "github.com/djthorpe/gopi/app"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// Check flags
	app.Logger.Info("flags=%v",app.FlagSet)
	err := CheckFlags(app.FlagSet)
	if err != nil {
		return err
	}

	// Get pin states
	gpio := app.GPIO

	app.Logger.Debug("Device=%v", app.Device)
	app.Logger.Debug("GPIO=%v", gpio)

	for _, logical := range gpio.Pins() {
		if physical := gpio.PhysicalPinForPin(logical); physical != 0 {
			app.Logger.Info("%v [Pin %v] => %v %v", logical, physical, gpio.ReadPin(logical), gpio.GetPinMode(logical))
		}
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func CheckFlags(flagset *app.Flags) error {
	// if no flags, then return OK
	if len(flagset.Flags()) == 0 {
		return nil
	}

	// Check for either: low, high, input or alt which are mutually
	// exclusive flags
	c := 0
	for _,flag := range([]string{ "input","alt","low","high" }) {
		if flagset.HasFlag(flag) {
			c++
		}
	}
	if c != 1 {
		return errors.New("One of -low, -high, -input, or -alt required")
	}

	// check for alt being between 0 and 5
	alt, exists := flagset.GetUint("alt")
	if exists {
		if alt > 5 {
			return errors.New("-alt is required to be between 0 and 5")
		}
	}

	// Check for pin
	pin, exists := flagset.GetString("pin")
	if exists != true {
		return errors.New("-pin flag required")
	}

	_, err := ParsePinValue(pin)
	if err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func ParsePinValue(name string) (hw.GPIOPin,error) {
	return hw.GPIO_PIN_NONE,errors.New("NOT IMPLEMENTED")
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Configuration
	config := app.Config(app.APP_GPIO)

	// Set the flags
	config.FlagSet.FlagString("pin","","Physical Pin Number or name")
	config.FlagSet.FlagBool("low",false,"Set pin to OUTPUT and set pin level LOW")
	config.FlagSet.FlagBool("high",false,"Set pin to OUTPUT and set pin level HIGH")
	config.FlagSet.FlagBool("input",false,"Set pin to INPUT")
	config.FlagSet.FlagUint("alt",0,"Set pin to an alternate function 0-5")

	// Create the application
	myapp, err := app.NewApp(config)
	if err == app.ErrHelp {
		// Help requested
		return
	} else if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(RunLoop); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
