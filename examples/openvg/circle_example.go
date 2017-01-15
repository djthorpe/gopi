/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// This example draws a circle on the screen and changes size on each
// iteration
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

type State struct {
	center       khronos.VGPoint
	diameter     float32
	maximum      float32
	increment    float32
	stroke, fill khronos.VGPaint
}

////////////////////////////////////////////////////////////////////////////////

func Increment(state *State) {
	state.diameter += state.increment
	if state.diameter <= 1.0 || state.diameter >= state.maximum {
		state.increment = -state.increment
	}
}

func Draw(vg khronos.VGDriver, state *State) error {
	// Paths
	path, err := vg.CreatePath()
	if err != nil {
		return err
	}
	defer vg.DestroyPath(path)
	path.Circle(state.center, state.diameter)

	// Draw
	return path.Draw(state.stroke, state.fill)

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func MyRunLoop(app *app.App) error {
	// Create a surface on which to draw
	opacity, _ := app.FlagSet.GetFloat64("opacity")
	surface, err := app.EGL.CreateBackground("OpenVG", float32(opacity))
	if err != nil {
		return err
	}
	defer app.EGL.DestroySurface(surface)

	// Set up state
	state := &State{
		center:    khronos.AlignPoint(surface, khronos.EGL_ALIGN_CENTER),
		diameter:  1.0,
		maximum:   400.0,
		increment: 1.0,
	}

	// Create stroke and fill paint brushes
	if state.fill, err = app.OpenVG.CreatePaint(khronos.VGColorWhite); err != nil {
		return err
	}
	defer app.OpenVG.DestroyPaint(state.fill)
	state.stroke, err = app.OpenVG.CreatePaint(khronos.VGColorRed)
	if err != nil {
		return err
	}
	defer app.OpenVG.DestroyPaint(state.stroke)
	state.stroke.SetStrokeWidth(10)
	state.stroke.SetStrokeDash(2, 2)

	// Loop and redraw
	go func() {
		for app.GetDone() == false {
			err = app.OpenVG.Do(surface, func() error {
				app.OpenVG.Clear(surface, khronos.VGColorBlack)
				return Draw(app.OpenVG, state)
			})
			if err != nil {
				app.Logger.Error("%v", err)
				break
			}
			Increment(state)
		}
		app.Logger.Info("Main Loop Ending")
	}()

	// Wait until done
	app.WaitUntilDone()

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
