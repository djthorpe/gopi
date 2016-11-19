/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

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

	// Debugging output
	app.Logger.Debug("Device=%v", app.Device)
	app.Logger.Debug("I2C=%v", app.I2C)

	// Create 8 x 8 grid for detecting slaves
	fmt.Fprintln(os.Stdout, "     -0 -1 -2 -3 -4 -5 -6 -7 -8 -9 -A -B -C -D -E -F")
	for h := uint8(0x0); h <= uint8(0x7); h++ {
		fmt.Fprintf(os.Stdout, "%01X- | ", h)
		for l := uint8(0x0); l <= uint8(0xF); l++ {
			slave := (h << 4) + l
			detected, err := app.I2C.DetectSlave(slave)
			if err != nil {
				fmt.Fprint(os.Stdout, "??")
			} else if detected {
				fmt.Fprintf(os.Stdout, "%02X", slave)
			} else {
				fmt.Fprint(os.Stdout, "--")
			}
			fmt.Fprint(os.Stdout, " ")
		}
		fmt.Fprint(os.Stdout, "\n")
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_I2C)

	// Create the application
	myapp, err := app.NewApp(config)
	if err == app.ErrHelp {
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
