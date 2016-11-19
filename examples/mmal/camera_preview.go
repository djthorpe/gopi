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
	gopi "github.com/djthorpe/gopi"
	app "github.com/djthorpe/gopi/app"
	mmal "github.com/djthorpe/gopi/device/rpi/mmal"
)

////////////////////////////////////////////////////////////////////////////////

func RunLoop(app *app.App) error {

	// create MMAL component
	camera, err := gopi.Open(mmal.MMAL{mmal.MMAL_COMPONENT_DEFAULT_VIDEO_DECODER}, app.Logger)
	if err != nil {
		return err
	}
	defer camera.Close()

	if err := camera.(*mmal.Component).Enable(); err != nil {
		return err
	}

	app.Logger.Info("Camera=%v", camera)

	// Wait until CTRL+C Pressed
	app.WaitUntilDone()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_DEVICE)

	// Create the application
	myapp, err := app.NewApp(config)
	if err != nil {
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
