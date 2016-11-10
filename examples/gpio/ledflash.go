/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This sample program shows how you can interact with the GPIO device on
// your Raspberry Pi. Firstly it enumerates all the pins (showing both
// physical pin number and logical name) with the state of those pins, then
// it can blink an LED on/off if the LED is connected to physical pin 40
// in series with a resistor to pin 39 or any 0V pin
package main

import (
	"fmt"
	"os"
)

import (
	app "github.com/djthorpe/gopi/app"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// Get pin states
	gpio := app.GPIO

	app.Logger.Debug("Device=%v", app.Device)
	app.Logger.Debug("GPIO=%v", gpio)

	for _, logical := range gpio.Pins() {
		if physical := gpio.PhysicalPinForPin(logical); physical != 0 {
			app.Logger.Info("%v [Pin %v] => %v %v", logical, physical, gpio.ReadPin(logical), gpio.GetPinMode(logical))
		}
	}

/*

	led_pin := gpio.PhysicalPin(40)
	gpio.SetPinMode(led_pin, hw.GPIO_OUTPUT)

	for {
		gpio.WritePin(led_pin, hw.GPIO_LOW)
		app.Logger.Info("%v => %v %v", led_pin, gpio.ReadPin(led_pin), gpio.GetPinMode(led_pin))
		time.Sleep(1.0 * time.Second)
		gpio.WritePin(led_pin, hw.GPIO_HIGH)
		app.Logger.Info("%v => %v %v", led_pin, gpio.ReadPin(led_pin), gpio.GetPinMode(led_pin))
		time.Sleep(1.0 * time.Second)
	}
*/

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Configuration
	config := app.Config(app.APP_GPIO)

	// Create the application
	myapp, err := app.NewApp(config)
	if err != nil {
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
