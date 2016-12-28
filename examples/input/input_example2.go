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
	"errors"
)

import (
	app "github.com/djthorpe/gopi/app"
	linux "github.com/djthorpe/gopi/device/linux"
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {

	path := app.FlagSet.Args()
	if len(path) != 1 {
		return errors.New("Invalid number of arguments")
	}

	device, err := gopi.Open(linux.InputDevice{ Path: path[0] },app.Logger)
	if err != nil {
		return err
	}
	defer device.Close()

	device.(hw.InputDevice).Watch(func (event hw.InputEvent,device hw.InputDevice) {
		fmt.Println("EVENT")
		fmt.Println("  DEVICE",device)
		fmt.Println("   EVENT",event)
	})

	// Wait for termination
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
	if err := myapp.Run(MyRunLoop); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
