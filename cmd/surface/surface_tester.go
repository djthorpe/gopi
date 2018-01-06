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
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/hw/rpi"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/surface/rpi"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	return nil
}

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("surface")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
