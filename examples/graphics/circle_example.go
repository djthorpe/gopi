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
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

import (
	gopi "github.com/djthorpe/gopi"
	app "github.com/djthorpe/gopi/app"
	rpi "github.com/djthorpe/gopi/device/rpi"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

func CreateBackground(app *app.App) (khronos.EGLSurface, error) {
	// Check for background flag, return nil if no background
	bg, exists := app.FlagSet.GetString("bg")
	if exists == false {
		return nil, nil
	}

	// Get color
	color, err := GetBackgroundColor(bg)
	if err != nil {
		return nil, err
	}

	// Create surface
	bgsurface, err := app.EGL.CreateBackground("DX", 1.0)
	if err != nil {
		return nil, err
	}

	// Clear background to color
	bgbitmap, err := bgsurface.GetBitmap()
	if err != nil {
		app.EGL.DestroySurface(bgsurface)
		return nil, err
	}
	if err := bgbitmap.ClearToColor(color); err != nil {
		app.EGL.DestroySurface(bgsurface)
		return nil, err
	}

	// Success
	app.Logger.Info("Background=%v", bgsurface)
	return bgsurface, nil
}

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

func DrawCircle(driver khronos.VGDriver, surface khronos.EGLSurface) error {

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	app.Logger.Info("Device=%v", app.Device)
	app.Logger.Info("Display=%v", app.Display)
	app.Logger.Info("EGL=%v", app.EGL)
	app.Logger.Info("OpenVG=%v", app.OpenVG)

	// Create background
	bgsurface, err := CreateBackground(app)
	if err != nil {
		return err
	}
	if bgsurface != nil {
		defer app.EGL.DestroySurface(bgsurface)
	}

	screen_rect := app.EGL.GetFrame()
	app.Logger.Info("Screen Rect = %v", screen_rect)

	// Create window
	surface, err := app.EGL.CreateSurface("OpenVG", screen_rect.Size(), screen_rect.Origin(), 2, 1.0)
	if err != nil {
		return err
	}
	defer app.EGL.DestroySurface(surface)
	app.Logger.Info("Surface=%v", surface)

	// Open Fonts
	fonts, err := gopi.Open(rpi.VGFont{}, app.Logger)
	if err != nil {
		return err
	}
	defer fonts.Close()

	// Open Faces
	basepath, exists := app.FlagSet.GetString("fontdir")
	if exists {
		err = fonts.(khronos.VGFontDriver).OpenFacesAtPath(basepath, func(filename string, info os.FileInfo) bool {
			if strings.HasPrefix(info.Name(), ".") {
				return false
			}
			if info.IsDir() {
				// Recurse into folders
				return true
			}
			if path.Ext(filename) == ".ttf" || path.Ext(filename) == ".TTF" {

				return true
			}
			app.Logger.Warn("Ignoring file %v", filename)
			return false
		})
		if err != nil {
			return err
		}
	}

	app.Logger.Info("Fonts=%v", fonts)

	// Wait until done (which means CTRL+C)
	app.WaitUntilDone()

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL | app.APP_OPENVG)

	// Add on command-line flags
	config.FlagSet.FlagString("bg", "", "Background color. One of red, green, blue, black, white, grey")
	config.FlagSet.FlagString("fontdir", "", "Font directory for loading of fonts")

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
