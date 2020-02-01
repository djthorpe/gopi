/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	tablewriter "github.com/olekukonko/tablewriter"
)

func StringToType(flags string) (gopi.InputDeviceType, error) {
	value := gopi.INPUT_TYPE_NONE
	if flags == "" {
		return gopi.INPUT_TYPE_ANY, nil
	}
	for _, field := range strings.Split(flags, ",") {
		switch strings.ToLower(field) {
		case "mouse":
			value |= gopi.INPUT_TYPE_MOUSE
		case "keyboard", "kb":
			value |= gopi.INPUT_TYPE_KEYBOARD
		case "touch", "touchscreen":
			value |= gopi.INPUT_TYPE_TOUCHSCREEN
		case "joystick":
			value |= gopi.INPUT_TYPE_JOYSTICK
		case "remote":
			value |= gopi.INPUT_TYPE_REMOTE
		default:
			return value, fmt.Errorf("Invalid value %v", strconv.Quote(field))
		}
	}
	return value, nil
}

func PrintDevices(devices []gopi.InputDevice) {
	table := tablewriter.NewWriter(os.Stdout)
	for _, device := range devices {
		table.Append([]string{
			device.Name(),
			fmt.Sprint(device.Type()),
		})
	}
	table.Render()
}

func Main(app gopi.App, args []string) error {
	if len(args) > 0 {
		return gopi.ErrHelp
	}

	// Open devices
	deviceName := app.Flags().GetString("input.name", gopi.FLAG_NS_DEFAULT)
	exclusive := app.Flags().GetBool("input.exclusive", gopi.FLAG_NS_DEFAULT)
	if deviceFlags, err := StringToType(app.Flags().GetString("input.type", gopi.FLAG_NS_DEFAULT)); err != nil {
		return err
	} else if devices, err := app.Input().OpenDevicesByNameType(deviceName, deviceFlags, exclusive); err != nil {
		return err
	} else if len(devices) == 0 {
		return gopi.ErrNotFound
	} else {
		PrintDevices(devices)
	}

	// If we are watching for events, then wait for CTRL+C
	if app.Flags().GetBool("watch", gopi.FLAG_NS_DEFAULT) {
		fmt.Println("Press CTRL+C to abort")
		app.WaitForSignal(context.Background(), os.Interrupt)
	}

	// Return success
	return nil
}
