/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example takes a snapshot of the screen and writes to a file as a PNG
// image
package main

import (
	"fmt"
	"os"
)

import (
	app "github.com/djthorpe/gopi/app"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(application *app.App) error {
	egl := application.EGL.(khronos.EGLDriver)

	// Obtain filename from command line
	args := application.FlagSet.Args()
	if len(args) != 1 {
		return app.ErrHelp
	}

	// Open file for writing
	handle, err := os.Create(args[0])
	if err != nil {
		return err
	}
	defer handle.Close()

	// Do snapshot
	bitmap, err := egl.SnapshotImage()
	if err != nil {
		return err
	}
	defer egl.DestroyImage(bitmap)

	// Save file as PNG
	if err := egl.WriteImagePNG(handle,bitmap); err != nil {
		return err
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL)

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
