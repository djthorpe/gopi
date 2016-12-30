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
	"time"
)

import (
	app "github.com/djthorpe/gopi/app"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func HelloWorld(app *app.App) error {

	// open devices
	devices, err := app.Input.OpenDevicesByName("",hw.INPUT_TYPE_ANY,hw.INPUT_BUS_ANY)
	if err != nil {
		return err
	}

	app.Logger.Info("DEVICES = %v",devices)

	app.Input.Watch(time.Second * 10)

	app.WaitUntilDone()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Create the configuration, we want to use the DEVICE
	// subsystem
	config := app.Config(app.APP_INPUT)

	// Create the application
	myapp, err := app.NewApp(config)
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
