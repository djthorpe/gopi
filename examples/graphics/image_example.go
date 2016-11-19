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

	app.Logger.Info("Image=%v", image)

	// Check for background
	bg, exists := app.FlagSet.GetString("bg")
	if exists {
		// Get color
		color, err := GetBackgroundColor(bg)
		if err != nil {
			return err
		}

		// Create surface
		bgsurface, err := app.EGL.CreateBackground("DX", 1.0)
		if err != nil {
			return err
		}
		defer app.EGL.DestroySurface(bgsurface)

		// Clear background to color
		bgbitmap, err := bgsurface.GetBitmap()
		if err != nil {
			return err
		}
		bgbitmap.ClearToColor(color)

		app.Logger.Info("Background=%v", bgsurface)
	}

	// Create window with image - set opacity
	opacity, _ := app.FlagSet.GetFloat64("opacity")
	surface, err := app.EGL.CreateSurfaceWithBitmap(image, khronos.EGLPoint{0, 0}, 2, float32(opacity))
	if err != nil {
		return err
	}
	defer app.EGL.DestroySurface(surface)

	app.Logger.Info("Surface=%v", surface)

	// Wait until done (which means CTRL+C)
	app.WaitUntilDone()

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func GetBackgroundColor(color string) (khronos.EGLColorRGBA32, error) {
	switch {
	case color == "white":
		return khronos.EGLWhiteColor, nil
	case color == "red":
		return khronos.EGLRedColor, nil
	case color == "green":
		return khronos.EGLGreenColor, nil
	case color == "blue":
		return khronos.EGLBlueColor, nil
	case color == "black":
		return khronos.EGLBlackColor, nil
	case color == "grey":
		return khronos.EGLGreyColor, nil
	default:
		return khronos.EGLBlackColor, errors.New("Unknown color value")
	}
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL)

	// Add on command-line flags
	config.FlagSet.FlagString("image", "", "Image filename")
	config.FlagSet.FlagString("bg", "", "Background color. One of red, green, blue, black, white, grey")
	config.FlagSet.FlagFloat64("opacity", 1.0, "Image opacity, 0.0 -> 1.0")

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
