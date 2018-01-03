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

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi/sys/hw/rpi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func runLoop(app *gopi.AppInstance, done chan struct{}) error {
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

	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the gpio instance
	config := gopi.NewAppConfig("gpio")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, runLoop))
}
