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
	app "../app"   /* import "github.com/djthorpe/gopi/app" */
	util "../util" /* import "github.com/djthorpe/gopi/util" */
)

////////////////////////////////////////////////////////////////////////////////

func HelloWorld(app *app.App) error {

	// Get Serial Number of the Raspberry Pi
	serial_number, err := app.Device.GetSerialNumber()
	if err != nil {
		return err
	}

	// Output information to logging
	app.Logger.Info("Hello, %v", serial_number)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Create the application
	myapp, err := app.NewApp(app.AppConfig{
		Features: app.APP_DEVICE,
		LogLevel: util.LOG_INFO,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(HelloWorld); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
