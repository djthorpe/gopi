/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Outputs a table of displays - works on RPi at the moment
package main

import (
	"errors"
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/olekukonko/tablewriter"

	// Modules
	rpi "github.com/djthorpe/gopi/sys/hw/rpi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {
	if app.Hardware == nil || app.Hardware.NumberOfDisplays() == 0 {
		return errors.New("No displays detected")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Display", "Name", "Width", "height", "Pixels per inch"})
	for n := uint(0); n < app.Hardware.NumberOfDisplays(); n++ {
		if module, err := gopi.Open(rpi.Display{Display: n}, app.Logger); err != nil {
			return err
		} else if display, ok := module.(gopi.Display); !ok {
			module.Close()
			return err
		} else {
			w, h := display.Size()
			ppi := fmt.Sprint(display.PixelsPerInch())
			if ppi == "0" {
				ppi = "-"
			}
			table.Append([]string{
				fmt.Sprint(n),
				fmt.Sprint(display.Name()),
				fmt.Sprint(w),
				fmt.Sprint(h),
				fmt.Sprint(ppi),
			})
			module.Close()
		}
	}
	table.Render()

	return nil
}

func main() {
	// Create the configuration, load the gpio instance
	config := gopi.NewAppConfig("hw")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
