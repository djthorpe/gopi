/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// SPIUTIL
//
// This example demonstrates reading from an SPI device
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
	app.Logger.Debug("SPI=%v", app.SPI)

	return nil
}

func main() {
	// Create the config
	config := app.Config(app.APP_SPI)

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
