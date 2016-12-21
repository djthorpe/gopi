/*
    GOPI Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing information, please see LICENSE.md
	For Documentation, see http://djthorpe.github.io/gopi/
*/

// This example outputs a table of detected input devices, their types
// and other information about them.
package main

import (
	"fmt"
	"os"
)

import (
	app "github.com/djthorpe/gopi/app"
)

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {

	devices, err := app.Input.OpenDevicesByName("",nil)
	if err != nil {
		return err
	}

	for _, device := range devices {
		fmt.Println(device)
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_INPUT)

	// Create the application
	myapp, err := app.NewApp(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(MyRunLoop); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
