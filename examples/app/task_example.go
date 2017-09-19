/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// This example shows running two independent tasks until CTRL+C is
// pressed
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func taskA(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Info("In taskA (which blocks until done)")

	select {
	case <-done:
		break
	}

	app.Logger.Info("taskA done")

	// Return success
	return nil
}

func taskB(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Info("In taskB (which does something every second until done)")

	// Tick every second
	ticker := time.Tick(time.Second)

	outer_loop: for {
		select {
		case <-ticker:
			app.Logger.Info("taskB: Tick")
		case <-done:
			break outer_loop
		}
	}

	app.Logger.Info("taskB done")

	// Return success
	return nil
}

func taskMain(app *gopi.AppInstance, done chan struct{}) error {
	app.Logger.Info("In taskMain")

	// Wait for interrupt signal (INT or TERM)
	app.WaitForSignal()
	app.Logger.Info("taskMain done")

	// Signal other routines that we are DONE, and return
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the application
	app, err := gopi.NewAppInstance(gopi.NewAppConfig())
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application - one foreground and two background tasks
	if err := app.Run(taskMain, taskA, taskB); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
