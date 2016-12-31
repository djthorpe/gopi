/*
    GOPI Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing information, please see LICENSE.md
	For Documentation, see http://djthorpe.github.io/gopi/
*/

// This example outputs a table of detected input devices, their types
// and other information about them.
package main

import (
	"fmt"
	"os"
	"errors"
	"strings"
	"time"
)

import (
	app "github.com/djthorpe/gopi/app"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func ParseDeviceBus(value string) hw.InputDeviceBus {
	switch(value) {
	case "usb":
		return hw.INPUT_BUS_USB
	case "bluetooth":
		return hw.INPUT_BUS_BLUETOOTH
	case "any":
		return hw.INPUT_BUS_ANY
	default:
		return hw.INPUT_BUS_NONE
	}

}

func ParseDeviceType(value string) hw.InputDeviceType {
	switch(value) {
	case "mouse":
		return hw.INPUT_TYPE_MOUSE
	case "keyboard":
		return hw.INPUT_TYPE_KEYBOARD
	case "joystick":
		return hw.INPUT_TYPE_JOYSTICK
	case "touchscreen":
		return hw.INPUT_TYPE_TOUCHSCREEN
	case "any":
		return hw.INPUT_TYPE_ANY
	default:
		return hw.INPUT_TYPE_NONE
	}
}

func ParseFlags(flags *app.Flags) (string,hw.InputDeviceType,hw.InputDeviceBus,error) {
	// Bus
	bus_string, _ := flags.GetString("bus")
	bus_value := ParseDeviceBus(strings.ToLower(strings.TrimSpace(bus_string)))
	if bus_value == hw.INPUT_BUS_NONE {
		return "",hw.INPUT_TYPE_NONE,hw.INPUT_BUS_NONE,errors.New("Invalid -bus flag")
	}

	// Type
	type_string, _ := flags.GetString("type")
	type_value := ParseDeviceType(strings.ToLower(strings.TrimSpace(type_string)))
	if type_value == hw.INPUT_TYPE_NONE {
		return "",hw.INPUT_TYPE_NONE,hw.INPUT_BUS_NONE,errors.New("Invalid -type flag")
	}

	// Name
	name_string, _ := flags.GetString("name")

	// Return success
	return name_string, type_value, bus_value, nil
}

////////////////////////////////////////////////////////////////////////////////



func Watch(app *app.App) {

	format := "%-30s %-25s %-25s\n"
	i := 0

	for app.GetDone() == false {

		err := app.Input.Watch(time.Millisecond * 200,func (event *hw.InputEvent, device hw.InputDevice) {
			// Print table header every 40 invocations
			if (i % 40) == 0 {
				fmt.Println("")
				fmt.Printf(format, "Device", "Type", "Value")
				fmt.Printf(format, "------------------------------", "-------------------------", "-------------------------")
			}
			i += 1

			// Print out the event
			switch(event.EventType) {
			case hw.INPUT_EVENT_KEYPRESS, hw.INPUT_EVENT_KEYRELEASE, hw.INPUT_EVENT_KEYREPEAT:
				fmt.Printf(format,device.GetName(),event.EventType,event.Keycode)
			case hw.INPUT_EVENT_ABSPOSITION:
				fmt.Printf(format,device.GetName(),event.EventType,event.Position)
			case hw.INPUT_EVENT_RELPOSITION:
				fmt.Printf(format,device.GetName(),event.EventType,event.Relative)
			default:
				fmt.Printf(format,device.GetName(),event.EventType,"")
			}
		})
		if err != nil {
			// Report any errors
			app.Logger.Error("Error: %v",err)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {

	device_name, device_type, device_bus, err := ParseFlags(app.FlagSet)
	if err != nil {
		return err
	}

	// Opens devices
	app.Logger.Info("input=%v", app.Input)
	devices, err := app.Input.OpenDevicesByName(device_name, device_type, device_bus)
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return errors.New("No devices found")
	}

	format := "%-30s %-25s %-25s\n"
	fmt.Printf(format, "Name", "Type", "Bus")
	fmt.Printf(format, "------------------------------", "-------------------------", "-------------------------")

	for _, device := range devices {
		fmt.Printf(format, device.GetName(), device.GetType(), device.GetBus())
	}

	// Watch in background
	if watch, _ := app.FlagSet.GetBool("watch"); watch {
		 go Watch(app)
	}

	app.WaitUntilDone()

	// TODO: Wait for shutdown of the Watch goroutine

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_INPUT)

	// Flags
	config.FlagSet.FlagString("type", "any", "Input type (any, mouse, keyboard, joystick, touchscreen)")
	config.FlagSet.FlagString("bus", "any", "Input bus (any, usb, bluetooth)")
	config.FlagSet.FlagString("name", "", "Name of input device")
	config.FlagSet.FlagBool("watch", false, "Watch for events from devices until CTRL+C is pressed")

	// Create the application
	myapp, err := app.NewApp(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(MyRunLoop); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
