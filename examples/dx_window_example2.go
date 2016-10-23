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
	app.Logger.Info("Device=%v", app.Device)
	app.Logger.Info("Display=%v", app.Display)
	app.Logger.Info("EGL=%v", app.EGL)
	app.Logger.Info("OpenVG=%v", app.OpenVG)

	// Create a background
	bg, err := app.EGL.CreateBackground("OpenVG")
	if err != nil {
		return app.Logger.Error("Error: %v", err)
	}
	defer app.EGL.CloseWindow(bg)

	// Create a window
	fg, err := app.EGL.CreateWindow("OpenVG",khronos.EGLSize{ 100, 100 },khronos.EGLPoint{ 100, 100 },1)
	if err != nil {
		return err
	}
	defer app.EGL.CloseWindow(fg)

	gfx := app.OpenVG

	// Clear background to white
	if err := gfx.Begin(bg); err != nil {
		return err
	}
	gfx.Clear(khronos.VGColor{ 1.0, 0.0, 0.0, 1.0 })
	gfx.Flush()

	// Clear foreground to green
	gfx.Begin(fg)
	gfx.Clear(khronos.VGColor{ 0.0, 1.0, 0.0, 1.0 })
	gfx.Flush()

	// Move window
	for i := 0; i < 100; i++ {
		app.EGL.MoveWindowOriginBy(fg,khronos.EGLPoint{ 1, 1 })
	}

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
