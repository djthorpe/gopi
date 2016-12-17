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
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////

func DeviceInfo(app *app.App) error {
	for _, tuple := range app.Device.GetCapabilities() {
		fmt.Printf("%v: %v\n", tuple.GetKey(), tuple)
	}

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
	if err := myapp.Run(DeviceInfo); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
