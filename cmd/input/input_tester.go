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

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/input/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if app.Input == nil {
		return errors.New("Missing Input instance")
	}

	// Open ALL input devices
	if devices, err := app.Input.OpenDevicesByName("", gopi.INPUT_TYPE_ANY, gopi.INPUT_BUS_ANY); err != nil {
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

	// wait until done
	app.WaitForSignal()

	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	config := gopi.NewAppConfig("input")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
