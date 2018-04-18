/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Example command for discovery of RPC microservices using mDNS
package main

import (
	// Frameworks
	"errors"
	"fmt"
	"os"

	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func eventLoop(app *gopi.AppInstance, done <-chan struct{}) error {
	if app.LIRC == nil {
		return errors.New("Missing LIRC module")
	}

	messages := app.LIRC.Subscribe()

	fmt.Printf("%20s %12s\n", "Type", "Value")
	fmt.Printf("%20s %12s\n", "--------------------", "------------")

FOR_LOOP:
	for {
		select {
		case evt := <-messages:
			if event, ok := evt.(gopi.LIRCEvent); ok {
				fmt.Printf("%20s %10sms\n", event.Type(), fmt.Sprint(event.Value()))
			} else {
				fmt.Println(evt)
			}
		case <-done:
			break FOR_LOOP
		}
	}

	// End of routine
	app.LIRC.Unsubscribe(messages)
	return nil
}

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if app.LIRC == nil {
		return errors.New("Missing LIRC module")
	}

	// Set receive mode to be MODE2
	// Ref: https://linuxtv.org/downloads/v4l-dvb-apis/uapi/rc/lirc-dev-intro.html#lirc-modes
	if err := app.LIRC.SetRcvMode(gopi.LIRC_MODE_MODE2); err != nil {
		return err
	}

	// Wait for interrupt
	app.WaitForSignal()

	// Finish gracefully
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the lirc instance
	config := gopi.NewAppConfig("lirc")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop, eventLoop))
}
