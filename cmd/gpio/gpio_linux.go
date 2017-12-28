/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Watches for a GPIO pin rising and/or falling
package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func eventLoop(app *gopi.AppInstance, done chan struct{}) error {

	if app.GPIO == nil {
		return app.Logger.Error("Missing GPIO module instance")
	}

	// Look for edges
	edge := app.GPIO.Subscribe()

FOR_LOOP:
	for {
		select {
		case evt := <-edge:
			fmt.Println("EVENT: ", evt)
		case <-done:
			break FOR_LOOP
		}
	}

	// Unsubscribe from edges
	app.GPIO.Unsubscribe(edge)
	return nil
}

func runLoop(app *gopi.AppInstance, done chan struct{}) error {

	if app.GPIO == nil {
		return app.Logger.Error("Missing GPIO module instance")
	}

	// watch pin
	pin, _ := app.AppFlags.GetUint("pin")
	app.GPIO.SetPinMode(gopi.GPIOPin(pin), gopi.GPIO_INPUT)
	app.GPIO.Watch(gopi.GPIOPin(pin), gopi.GPIO_EDGE_BOTH) // when button pressed or released

	// wait until done
	app.WaitForSignal()

	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the configuration
	config := gopi.NewAppConfig("gpio")
	config.AppFlags.FlagUint("pin", 27, "Logical GPIO Pin to watch")

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
