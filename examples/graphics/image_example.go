/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows how to load an image
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

import (
	app "github.com/djthorpe/gopi/app"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	app.Logger.Info("Device=%v", app.Device)
	app.Logger.Info("Display=%v", app.Display)
	app.Logger.Info("EGL=%v", app.EGL)

	// Fetch image filename flag
	filename, _ := app.FlagSet.GetString("image")
	if filename == "" {
		return errors.New("Missing -image flag")
	}

	// Open the image
	reader, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer reader.Close()
	image, err := app.EGL.CreateImage(reader)
	if err != nil {
		return err
	}
	defer app.EGL.DestroyImage(image)

	image.ClearToColor(khronos.EGLColorRGBA32{ 0, 255, 0, 255 })

	app.Logger.Info("Image=%v", image)

	// Create window with image
	surface, err := app.EGL.CreateSurfaceWithBitmap(image, khronos.EGLPoint{50, 50}, 2, 0.9)
	if err != nil {
		return err
	}
	defer app.EGL.DestroySurface(surface)

	app.Logger.Info("Surface=%v", surface)



	// Wait until done
	app.WaitUntilDone()

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL)

	// Add on command-line flags
	config.FlagSet.FlagString("image", "", "Image filename")

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
