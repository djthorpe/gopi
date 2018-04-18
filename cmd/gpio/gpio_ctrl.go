/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Outputs GPIO status
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/olekukonko/tablewriter"

	// Modules
	"./gpio_sys"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

var edge <-chan gopi.Event

////////////////////////////////////////////////////////////////////////////////

func pins(pins string) ([]gopi.GPIOPin, error) {
	logical_pins := make([]gopi.GPIOPin, 0)
	for _, value := range strings.Split(pins, ",") {
		if pin, err := strconv.ParseUint(value, 10, 64); err != nil {
			return nil, err
		} else {
			logical_pins = append(logical_pins, gopi.GPIOPin(pin))
		}
	}
	return logical_pins, nil
}

func inPins(app *gopi.AppInstance) ([]gopi.GPIOPin, error) {
	if flag_value, exists := app.AppFlags.GetString("in"); exists {
		return pins(flag_value)
	} else {
		return nil, nil
	}
}

func lowPins(app *gopi.AppInstance) ([]gopi.GPIOPin, error) {
	if flag_value, exists := app.AppFlags.GetString("low"); exists {
		return pins(flag_value)
	} else {
		return nil, nil
	}
}

func highPins(app *gopi.AppInstance) ([]gopi.GPIOPin, error) {
	if flag_value, exists := app.AppFlags.GetString("high"); exists {
		return pins(flag_value)
	} else {
		return nil, nil
	}
}

func edgePins(app *gopi.AppInstance) ([]gopi.GPIOPin, error) {
	if flag_value, exists := app.AppFlags.GetString("edge"); exists {
		return pins(flag_value)
	} else {
		return nil, nil
	}
}

////////////////////////////////////////////////////////////////////////////////

func eventLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	app.Logger.Debug("Started eventLoop")
FOR_LOOP:
	for {
		select {
		case evt := <-edge:
			if evt != nil {
				fmt.Println("EVENT: ", evt)
			}
		case <-done:
			break FOR_LOOP
		}
	}
	app.Logger.Debug("Ended eventLoop")
	return nil
}

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	watching := false

	// Get logical pins
	if in_pins, err := inPins(app); err != nil {
		return err
	} else if low_pins, err := lowPins(app); err != nil {
		return err
	} else if high_pins, err := highPins(app); err != nil {
		return err
	} else if edge_pins, err := edgePins(app); err != nil {
		return err
	} else {
		for _, logical_pin := range in_pins {
			app.GPIO.SetPinMode(logical_pin, gopi.GPIO_INPUT)
		}
		for _, logical_pin := range low_pins {
			app.GPIO.SetPinMode(logical_pin, gopi.GPIO_OUTPUT)
			app.GPIO.WritePin(logical_pin, gopi.GPIO_LOW)
		}
		for _, logical_pin := range high_pins {
			app.GPIO.SetPinMode(logical_pin, gopi.GPIO_OUTPUT)
			app.GPIO.WritePin(logical_pin, gopi.GPIO_HIGH)
		}
		for _, logical_pin := range edge_pins {
			watching = true
			app.GPIO.SetPinMode(logical_pin, gopi.GPIO_INPUT)
			app.GPIO.Watch(logical_pin, gopi.GPIO_EDGE_BOTH)
		}
	}

	// Output current state of pins
	if app.GPIO.NumberOfPhysicalPins() > 0 {
		table := tablewriter.NewWriter(os.Stdout)

		table.SetHeader([]string{"Physical", "Logical", "Direction", "Value"})

		// Physical pins start at index 1
		for pin := uint(1); pin <= app.GPIO.NumberOfPhysicalPins(); pin++ {
			var l, d, v string
			if logical := app.GPIO.PhysicalPin(pin); logical != gopi.GPIO_PIN_NONE {
				l = fmt.Sprint(logical)
				d = fmt.Sprint(app.GPIO.GetPinMode(logical))
				v = fmt.Sprint(app.GPIO.ReadPin(logical))
			}
			table.Append([]string{
				fmt.Sprintf("%v", pin), l, d, v,
			})
		}

		table.Render()
	}

	// If we are watching for changes to pins then wait for CTRL+C to stop watching
	if watching {
		edge = app.GPIO.Subscribe()
		if edge != nil {
			fmt.Printf("Watching for input pin changes, press CTRL+C to abort\n")
			app.WaitForSignal()
			app.GPIO.Unsubscribe(edge)
		} else {
			fmt.Printf("Edge detection not supported\n")
		}
	}

	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the gpio instance
	config := gopi.NewAppConfig(gpio_sys.MODULE_NAMES)
	config.AppFlags.FlagString("low", "", "Comma-separated list of pins to set to low")
	config.AppFlags.FlagString("high", "", "Comma-separated list of pins to set to high")
	config.AppFlags.FlagString("in", "", "Comma-separated list of pins to set to input")
	config.AppFlags.FlagString("edge", "", "Comma-separated list of pins to monitor for changes")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop, eventLoop))
}
