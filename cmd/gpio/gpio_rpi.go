/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Watches for a GPIO pin rising and/or falling
package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/rpi"
	_ "github.com/djthorpe/gopi/sys/logger"
	"github.com/olekukonko/tablewriter"
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

func main_inner() int {
	// Create the configuration
	config := gopi.NewAppConfig("gpio")
	config.AppFlags.FlagUint("pin", 27, "Logical GPIO Pin to watch")

	// Create the application
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(runLoop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
