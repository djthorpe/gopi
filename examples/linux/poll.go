/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package main

import (
	"fmt"
	"os"
	"time"
)

import (
	app "github.com/djthorpe/gopi/app"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////

func ProcessEvents(event *hw.InputEvent, device hw.InputDevice) {
	fmt.Println(event)
}

func HelloWorld(app *app.App) error {

	// open devices
	devices, err := app.Input.OpenDevicesByName("", hw.INPUT_TYPE_ANY, hw.INPUT_BUS_ANY)
	if err != nil {
		return err
	}

	app.Logger.Info("DEVICES = %v", devices)

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
				app.Input.Watch(time.Millisecond * 100,ProcessEvents)
			}
		}
	}()

	app.WaitUntilDone()

	// Shutdown goroutine
	finished_channel <- true
	_ = <-finished_watch

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Create the configuration, we want to use the DEVICE
	// subsystem
	config := app.Config(app.APP_INPUT)

	// Create the application
	myapp, err := app.NewApp(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(HelloWorld); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
