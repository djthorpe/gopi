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
	"time"
)

import (
	app "github.com/djthorpe/gopi/app"
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

func Draw(diameter float32,surface khronos.EGLSurface, vg khronos.VGDriver) error {
	vg.Begin(surface)

	// Clear to red
	vg.Clear(khronos.VGColorRed)

	// Paints
	fill, err := vg.CreatePaint(khronos.VGColorWhite)
	if err != nil {
		return err
	}
	defer vg.DestroyPaint(fill)
	stroke, err := vg.CreatePaint(khronos.VGColorMidGrey)
	if err != nil {
		return err
	}
	defer vg.DestroyPaint(stroke)
	stroke.SetLineWidth(10.0)

	// Paths
	path, err := vg.CreatePath()
	if err != nil {
		return err
	}
	defer vg.DestroyPath(path)
	path.Circle(vg.GetPoint(khronos.EGL_ALIGN_CENTER), diameter)

	// Draw
	path.Draw(stroke,fill)

	// Flush graphics
	vg.Flush()

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	app.Logger.Info("Device=%v", app.Device)
	app.Logger.Info("Display=%v", app.Display)
	app.Logger.Info("EGL=%v", app.EGL)
	app.Logger.Info("OpenVG=%v", app.OpenVG)

	opacity, _ := app.FlagSet.GetFloat64("opacity")
	surface, err := app.EGL.CreateBackground("OpenVG", float32(opacity))
	if err != nil {
		return err
	}
	defer app.EGL.DestroySurface(surface)

	var diam float32 = 1
	var inc float32 = 1
	var max float32 = 400

	for app.GetDone() == false {
		if err := Draw(diam, surface, app.OpenVG); err != nil {
			app.Logger.Error("%v",err)
			app.Done()
		}
		diam = diam + inc
		if diam <= 1.0 || diam >= max  {
			inc = -inc
		}
	}

	// Wait until done (which means CTRL+C)
	app.WaitUntilDone()

	// Wait for a while
	time.Sleep(1 * time.Second)

	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL | app.APP_OPENVG)
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
