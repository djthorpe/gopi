/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows how to create a DXWindow on the screen
package main

import (
	"flag"
	"fmt"
	"os"
	"bufio"
)

import (
	app "../app"   /* import "github.com/djthorpe/gopi/app" */
	util "../util" /* import "github.com/djthorpe/gopi/util" */
	khronos "../khronos" /* import "github.com/djthorpe/gopi/khronos" */
)

////////////////////////////////////////////////////////////////////////////////

var (
	flagDisplay = flag.Uint("display", 0, "Display number")
	flagVerbose = flag.Bool("verbose", false, "Output verbose logging messages")
	flagLogFile = flag.String("log", "", "Logging file. If empty, logs to stderr")
)

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	app.Logger.Debug("Device=%v", app.Device)
	app.Logger.Debug("Display=%v", app.Display)
	app.Logger.Debug("EGL=%v", app.EGL)
	app.Logger.Debug("OpenVG=%v", app.OpenVG)

	// Create a background
	bg, err := app.EGL.CreateBackground("OpenVG")
	if err != nil {
		return app.Logger.Error("Error: %v", err)
	}
	defer app.EGL.CloseWindow(bg)

	// Clear to white
	gfx := app.OpenVG
	gfx.Begin(bg)
	gfx.Clear(khronos.VGColor{ 1.0, 0.0, 0.0, 1.0 })
	gfx.Flush()

	// wait for a key press
	app.Logger.Info("Press a key to continue")
	bufio.NewReader(os.Stdin).ReadString('\n')

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Parse flags
	flag.Parse()

	// Determine level of logging
	var level util.LogLevel
	if *flagVerbose {
		level = util.LOG_ANY
	} else {
		level = util.LOG_INFO
	}

	// Create the application
	myapp, err := app.NewApp(app.AppConfig{
		Features:  app.APP_OPENVG,
		Display:   uint16(*flagDisplay),
		LogFile:   *flagLogFile,
		LogAppend: false,
		LogLevel:  level,
	})
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
