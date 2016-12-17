/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows how to draw Hello, World on a DX bitmap surface
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

func MyRunLoop(app *app.App) error {
	app.Logger.Info("Device=%v", app.Device)
	app.Logger.Info("Display=%v", app.Display)
	app.Logger.Info("EGL=%v", app.EGL)
	app.Logger.Info("OpenVG=%v", app.OpenVG)
	app.Logger.Info("Fonts=%v", app.Fonts)

	// Create a background with opacity 0.75
	bg, err := app.EGL.CreateBackground("DX", 0.75)
	if err != nil {
		return app.Logger.Error("Error: %v", err)
	}
	defer app.EGL.DestroySurface(bg)

	// Load in font face
	face := GetFontFace(app)
	if face == nil {
		return app.Logger.Error("Error: Missing or invalid -font flag")
	}

	app.Logger.Info("FACE=%v",face)

	// Draw at {0,0}
	bitmap, err := bg.GetBitmap()
	if err != nil {
		return err
	}
	if err := bitmap.PaintText("Hello, world!",face,khronos.EGLPoint{ 0, 0 },128.0); err != nil {
		return err
	}

	// Wait until CTRL+C is pressed
	app.WaitUntilDone()

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func GetFontFace(app *app.App) khronos.VGFace {
	family, exists := app.FlagSet.GetString("font")
	if exists == false {
		return nil
	}
	faces := app.Fonts.GetFaces(family,khronos.VG_FONT_STYLE_REGULAR)
	if len(faces) == 0 {
		return nil
	}
	// return the first face
	return faces[0]
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_VGFONT|app.APP_EGL)

	// Font
	config.FlagSet.FlagString("font", "", "Font to use")

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
