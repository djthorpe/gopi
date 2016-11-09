/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows how to load an image
package main

import (
	"flag"
	"fmt"
	"os"
	"errors"
)

import (
	app "github.com/djthorpe/gopi/app"
)

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	app.Logger.Info("Device=%v", app.Device)
	app.Logger.Info("Display=%v", app.Display)
	app.Logger.Info("EGL=%v", app.EGL)

	// Fetch image filename flag
	filename := app.FlagSet.Lookup("image").Value.(flag.Getter).Get().(string)
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

	app.Logger.Info("Image=%v", image)

	// Wait until done
	app.WaitUntilDone()

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL)

	// Add on command-line flags
	config.FlagSet.String("image", "", "Image filename")

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
