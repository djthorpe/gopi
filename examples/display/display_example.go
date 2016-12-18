/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example opens a display and returns information about the display.
// That's it!
package main

import (
	"fmt"
	"os"
)

import (
	gopi "github.com/djthorpe/gopi"
	app "github.com/djthorpe/gopi/app"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	display := app.Display.(gopi.DisplayDriver)
	fmt.Println("DISPLAY:", display.GetDisplay())
	w, h := display.GetDisplaySize()
	fmt.Println("   SIZE:", khronos.EGLSize{uint(w), uint(h)})
	fmt.Println("    PPI:", display.GetPixelsPerInch())

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_DISPLAY)

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
