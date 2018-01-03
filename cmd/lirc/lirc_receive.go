/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
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

	// Codecs
	lirc_codec "./lirc_codec"
)

////////////////////////////////////////////////////////////////////////////////

func eventLoop(app *gopi.AppInstance, done chan struct{}) error {
	if app.LIRC == nil {
		return errors.New("Missing LIRC module")
	}

	sony_decoder := lirc_codec.NewSonyDecoder(app.Logger)
	philips_decoder := lirc_codec.NewPhilipsRC5Decoder(app.Logger)
	edge := app.LIRC.Subscribe()

FOR_LOOP:
	for {
		select {
		case evt := <-edge:
			fmt.Println("EVENT: ", evt)
			sony_decoder.Receive(evt.(gopi.LIRCEvent))
			philips_decoder.Receive(evt.(gopi.LIRCEvent))
		case <-done:
			break FOR_LOOP
		}
	}

	// Unsubscribe from edges
	app.LIRC.Unsubscribe(edge)
	return nil
}

func mainLoop(app *gopi.AppInstance, done chan struct{}) error {

	if app.LIRC == nil {
		return errors.New("Missing LIRC module")
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
