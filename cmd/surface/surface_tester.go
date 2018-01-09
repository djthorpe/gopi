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
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/graphics/rpi"
	_ "github.com/djthorpe/gopi/sys/hw/rpi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if surface_manager := app.Surface; surface_manager == nil {
		return fmt.Errorf("Missing Surface Manager")
	} else if bitmap, err := surface_manager.CreateBitmap(gopi.Size{100, 100}); err != nil {
		return err
	} else {
		defer bitmap.Close()
		fmt.Println(bitmap)
	}

	return nil
}

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("surface")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
