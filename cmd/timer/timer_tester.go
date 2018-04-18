/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Runs either a one-shot or interval timer
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/logger"
	_ "github.com/djthorpe/gopi/sys/timer"
)

////////////////////////////////////////////////////////////////////////////////

func handleEvent(app *gopi.AppInstance, evt gopi.TimerEvent) {
	fmt.Println("EVENT: ", evt)
	if evt.UserInfo().(string) == "Timeout" {
		app.Timer.NewTimeout(4*time.Second, "Timeout")
		app.Timer.NewTimeout(10*time.Second, "Timeout2")
	}
}

func eventLoop(app *gopi.AppInstance, done <-chan struct{}) error {

	// Subscribe to timers
	edge := app.Timer.Subscribe()

FOR_LOOP:
	for {
		select {
		case evt := <-edge:
			if evt != nil {
				handleEvent(app, evt.(gopi.TimerEvent))
			}
		case <-done:
			break FOR_LOOP
		}
	}

	// Unsubscribe from events
	app.Timer.Unsubscribe(edge)
	return nil
}

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	app.Timer.NewInterval(1*time.Second, "Periodic Timer", false)
	app.Timer.NewTimeout(4*time.Second, "Timeout")

	// wait until done
	app.WaitForSignal()

	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the timer instance
	config := gopi.NewAppConfig("timer")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop, eventLoop))
}
