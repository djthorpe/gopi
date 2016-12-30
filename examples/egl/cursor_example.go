/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This example draws a cursor and allows it to be moved on the
// screen by mouse or touchscreen
package main

import (
	"fmt"
	"os"
	"time"
)

import (
	app "github.com/djthorpe/gopi/app"
	khronos "github.com/djthorpe/gopi/khronos"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func ProcessEvents(event *hw.InputEvent,app *app.App,cursor khronos.EGLSurface) {
	app.EGL.MoveSurfaceOriginTo(cursor,event.Position)
}

func StartWatching(app *app.App,cursor khronos.EGLSurface) (chan bool,chan bool) {
	// Watch for events and check for completed every 100 milliseconds
	finished_channel := make(chan bool)
	finished_watch := make(chan bool)
	go func() {
		for {
			select {
			case _ = <-finished_channel:
				finished_watch <- true
				return
			default:
				app.Input.Watch(time.Millisecond * 100,func (event *hw.InputEvent, device hw.InputDevice) {
					ProcessEvents(event,app,cursor)
				})
			}
		}
	}()

	// Return the channels use for completing
	return finished_channel, finished_watch
}

func MyRunLoop(app *app.App) error {
	egl := app.EGL.(khronos.EGLDriver)

	// Open mouse
	devices, err := app.Input.OpenDevicesByName("", hw.INPUT_TYPE_MOUSE, hw.INPUT_BUS_ANY)
	if err != nil {
		return err
	}
	if len(devices) == 0 {
		return app.Logger.Error("No mouse found")
	}

	// Create a cursor
	cursor, err := egl.CreateCursor()
	if err != nil {
		return err
	}
	defer egl.DestroySurface(cursor)

	// Set cursor position on all devices
	for _, device := range devices {
		device.SetPosition(cursor.GetOrigin())
	}

	// Start watching for mouse events
	finished_channel, finished_watch := StartWatching(app,cursor)

	app.WaitUntilDone()

	// Shutdown goroutine
	finished_channel <- true
	_ = <-finished_watch

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the config
	config := app.Config(app.APP_EGL | app.APP_INPUT)

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
