/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func runLoop(app *gopi.AppInstance, done chan struct{}) error {

	if app.GPIO == nil {
		return app.Logger.Error("Missing GPIO module instance")
	}

	// watch pin
	app.GPIO.SetPinMode(gopi.GPIOPin(27), gopi.GPIO_INPUT)
	app.GPIO.Watch(gopi.GPIOPin(27), gopi.GPIO_EDGE_FALLING)

	// wait until done
	app.WaitForSignal()

	done <- gopi.DONE
	return nil
}

func eventLoop(app *gopi.AppInstance, done chan struct{}) error {
	if app.GPIO == nil {
		return app.Logger.Error("Missing GPIO module instance")
	}

	app.Logger.Debug("eventLoop waiting for incoming events")
	edge := app.GPIO.Subscribe()

FOR_LOOP:
	for {
		select {
		case evt := <-edge:
			fmt.Println("EVENT: ", evt)
		case <-done:
			app.GPIO.Unsubscribe(edge)
			break FOR_LOOP
		}
	}
	app.Logger.Info("END eventLoop")
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the configuration
	config := gopi.NewAppConfig("gpio")
	// Create the application
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(runLoop, eventLoop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
