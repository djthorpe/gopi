/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// LEDFLASH
//
// This sample program shows how you can control an LED. A very simple program
// indeed! To run it, indicate which physical pin the LED is connected to. You
// should connect the LED as follows:
//
// 1. Wire <GPIO physical pin> -> <led anode>
// 2. Wire <led cathode> -> <resistor>
// 3. Wire <resistor> -> <GPIO ground pin>
//
// The resistor should be about 220 ohms, or other low ohm value. It is simply
// there to limit the current draw from the GPIO pin and blow out your hardware.
// Once wired, use the following:
//
//   ledflash -pin <pin> -low
//   ledflash -pin <pin> -high
//   ledflash -pin <pin> -flash
//
package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

import (
	app "github.com/djthorpe/gopi/app"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// Debugging output
	app.Logger.Debug("Device=%v", app.Device)
	app.Logger.Debug("GPIO=%v", app.GPIO)

	// Check flags
	app.Logger.Debug("flags=%v", app.FlagSet)
	err := CheckFlags(app.FlagSet)
	if err != nil {
		return err
	}

	// Get pin
	pin, err := ParsePinFlag(app.GPIO, app.FlagSet)
	if err != nil {
		return err
	}
	app.Logger.Debug("Pin=%v", pin)

	switch {
	case app.FlagSet.HasFlag("off"):
		app.GPIO.SetPinMode(pin, hw.GPIO_OUTPUT)
		app.GPIO.WritePin(pin, hw.GPIO_LOW)
		return nil
	case app.FlagSet.HasFlag("on"):
		app.GPIO.SetPinMode(pin, hw.GPIO_OUTPUT)
		app.GPIO.WritePin(pin, hw.GPIO_HIGH)
		return nil
	case app.FlagSet.HasFlag("flash"):
		app.GPIO.SetPinMode(pin, hw.GPIO_OUTPUT)
		app.GPIO.WritePin(pin, hw.GPIO_LOW)
		go func() {
			for {
				time.Sleep(100 * time.Millisecond)
				app.GPIO.WritePin(pin, hw.GPIO_HIGH)
				time.Sleep(50 * time.Millisecond)
				app.GPIO.WritePin(pin, hw.GPIO_LOW)
			}
		}()
	default:
		return errors.New("NOT IMPLEMENTED")
	}

	app.WaitUntilDone()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func CheckFlags(flagset *app.Flags) error {
	// if no pin flag, then return error
	if flagset.HasFlag("pin") == false {
		return errors.New("Missing -pin flag")
	}

	// Check for either: on, off or flash
	c := 0
	for _, flag := range []string{"flash", "off", "on"} {
		if flagset.HasFlag(flag) {
			c++
		}
	}
	if c != 1 {
		return errors.New("One of -on, -off or -flash required")
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func ParsePinFlag(gpio hw.GPIODriver, flagset *app.Flags) (hw.GPIOPin, error) {
	value, exists := flagset.GetString("pin")
	if exists == false {
		return hw.GPIO_PIN_NONE, nil
	}

	// Check for physical pin
	pin, err := strconv.ParseUint(value, 10, 32)
	if err == nil {
		logical := gpio.PhysicalPin(uint(pin))
		if logical == hw.GPIO_PIN_NONE {
			return logical, errors.New("Invalid pin")
		}
		return logical, nil
	}

	// Check for logical pin
	for _, pin := range gpio.Pins() {
		if value == pin.String() {
			return pin, nil
		}
	}

	return hw.GPIO_PIN_NONE, errors.New("Unknown pin")
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Configuration
	config := app.Config(app.APP_GPIO)

	// Set the flags
	config.FlagSet.FlagString("pin", "", "Physical Pin Number or name")
	config.FlagSet.FlagBool("on", false, "Switch LED on")
	config.FlagSet.FlagBool("off", false, "Switch LED off")
	config.FlagSet.FlagBool("flash", false, "Flash LED until CTRL+C pressed")

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
