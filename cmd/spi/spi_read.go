/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Reads bytes from the SPI interface
package main

import (
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan struct{}) error {

	if app.SPI == nil {
		return app.Logger.Error("Missing SPI module instance")
	}

	// Finished
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the spi instance
	config := gopi.NewAppConfig("spi")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
