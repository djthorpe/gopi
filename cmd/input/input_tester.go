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
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/input/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

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
