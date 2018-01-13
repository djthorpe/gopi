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

	size := gopi.Size{100, 100}

	if graphics := app.Graphics; graphics == nil {
		return fmt.Errorf("Missing Graphics Manager")
	} else if surface, err := graphics.CreateSurface(gopi.SURFACE_TYPE_RGBA32, gopi.SURFACE_FLAG_NONE, 1.0, gopi.SURFACE_LAYER_DEFAULT, gopi.ZeroPoint, size); err != nil {
		return err
	} else {
		defer graphics.DestroySurface(surface)
		fmt.Println(surface)
	}

	return nil
}

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("graphics")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
