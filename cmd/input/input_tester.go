/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Runs either a one-shot or interval timer
package main

import (
	"errors"
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi/sys/input/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func ParseDeviceBus(value string) gopi.InputDeviceBus {
	switch value {
	case "usb":
		return gopi.INPUT_BUS_USB
	case "bluetooth":
		return gopi.INPUT_BUS_BLUETOOTH
	case "any":
		return gopi.INPUT_BUS_ANY
	default:
		return gopi.INPUT_BUS_NONE
	}

}

func ParseDeviceType(value string) gopi.InputDeviceType {
	switch value {
	case "mouse":
		return gopi.INPUT_TYPE_MOUSE
	case "keyboard":
		return gopi.INPUT_TYPE_KEYBOARD
	case "joystick":
		return gopi.INPUT_TYPE_JOYSTICK
	case "touchscreen":
		return gopi.INPUT_TYPE_TOUCHSCREEN
	case "any":
		return gopi.INPUT_TYPE_ANY
	default:
		return gopi.INPUT_TYPE_NONE
	}
}

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if app.Input == nil {
		return errors.New("Missing Input instance")
	}

	var device_type gopi.InputDeviceType
	var device_bus gopi.InputDeviceBus

	device_name, _ := app.AppFlags.GetString("input.name")
	if flag_type, exists := app.AppFlags.GetString("input.type"); exists {
		device_type = ParseDeviceType(flag_type)
	} else {
		device_type = gopi.INPUT_TYPE_ANY
	}
	if flag_bus, exists := app.AppFlags.GetString("input.bus"); exists {
		device_bus = ParseDeviceBus(flag_bus)
	} else {
		device_bus = gopi.INPUT_BUS_ANY
	}

	// Open ALL input devices
	if devices, err := app.Input.OpenDevicesByName(device_name, device_type, device_bus); err != nil {
		return err
	} else if len(devices) == 0 {
		fmt.Println("No input devices found")
	} else {
		// Output detected I2C addresses
		table := tablewriter.NewWriter(os.Stdout)

		table.SetHeader([]string{"Name", "Type", "Bus", "Position"})
		for _, device := range devices {
			table.Append([]string{
				fmt.Sprint(device.Name()),
				fmt.Sprint(device.Type()),
				fmt.Sprint(device.Bus()),
				fmt.Sprint(device.Position()),
			})
		}
		table.Render()
	}

	if watch, _ := app.AppFlags.GetBool("watch"); watch {
		fmt.Printf("Watching for input events, press CTRL+C to abort\n")
		app.WaitForSignal()
	}

	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	config := gopi.NewAppConfig("input")

	// Command-Line Flags
	config.AppFlags.FlagString("input.type", "any", "Input type (any, mouse, keyboard, joystick, touchscreen)")
	config.AppFlags.FlagString("input.bus", "any", "Input bus (any, usb, bluetooth)")
	config.AppFlags.FlagString("input.name", "", "Name of input device")
	config.AppFlags.FlagBool("watch", false, "Watch for events from devices until CTRL+C is pressed")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
