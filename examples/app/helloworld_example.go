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
)

////////////////////////////////////////////////////////////////////////////////

func HelloWorld(app *app.App) error {

	// Get serial number of the device
	serial_number, err := app.Device.GetSerialNumber()
	if err != nil {
		return err
	}

	// Output message to stdout
	fmt.Fprintf(os.Stdout, "Hello %08X!!\n", serial_number)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Create the configuration, we want to use the DEVICE
	// subsystem
	config := app.Config(app.APP_DEVICE)

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
