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

	face_index := 0
	size, _ := app.FlagSet.GetFloat64("size")
	text, text_exists := app.FlagSet.GetString("text")
	origin := khronos.EGLPoint{20, int(size)}
	bitmap, err := bg.GetBitmap()
	color := khronos.EGLWhiteColor
	if err != nil {
		return err
	}

	bitmap.ClearToColor(khronos.EGLRedColor)

	for {
		// Load in font face
		face := GetFontFace(app, face_index)
		if face == nil {
			break
		}
		// Draw
		if text_exists == false {
			text = fmt.Sprintf("%s %s", face.GetFamily(), face.GetStyle())
		}
		if err := bitmap.PaintText(text, face, color, origin, float32(size)); err != nil {
			return err
		}
		// Increment
		face_index += 1
		origin.Y += int(size)
	}

	// Wait until CTRL+C is pressed
	app.WaitUntilDone()

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func GetFontFace(app *app.App, index int) khronos.VGFace {
	family, exists := app.FlagSet.GetString("font")
	if exists == false {
		return nil
	}
	faces := app.Fonts.GetFaces(family, khronos.VG_FONT_STYLE_REGULAR)
	if len(faces) == 0 || index >= len(faces) {
		return nil
	}
	// return the face
	return faces[index]
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_VGFONT | app.APP_EGL)

	// Font
	config.FlagSet.FlagString("font", "", "Font to use")
	config.FlagSet.FlagFloat64("size", 48.0, "Font size, in points")
	config.FlagSet.FlagString("text", "Hello, world!", "Message to display")

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
