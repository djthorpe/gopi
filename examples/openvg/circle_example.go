/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows how to load an image using dispmanx (DX) bitmaps onto
// a surface, also setting a background color.
package main

import (
	"flag"
	"fmt"
	"os"
)

import (
	app "github.com/djthorpe/gopi/app"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

func Draw(app *app.App,vg khronos.VGDriver) error {
	surface, err := app.EGL.CreateBackground("OpenVG", 1.0)
	if err != nil {
		return err
	}
	defer app.EGL.DestroySurface(surface)

	vg.Begin(surface)
	vg.Clear(khronos.VGColorRed)
	vg.Flush()

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	app.Logger.Info("Device=%v", app.Device)
	app.Logger.Info("Display=%v", app.Display)
	app.Logger.Info("EGL=%v", app.EGL)
	app.Logger.Info("OpenVG=%v", app.OpenVG)

	// Draw circle on background
	err := Draw(app,app.OpenVG)
	if err != nil {
		return err
	}

	// Wait until done (which means CTRL+C)
	app.WaitUntilDone()

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL | app.APP_OPENVG)

	// Create the application
	myapp, err := app.NewApp(config)
	if err == flag.ErrHelp {
		return
	} else if err != nil {
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

