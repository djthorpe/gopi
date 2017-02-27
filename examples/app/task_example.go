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
)

import (
	app "github.com/djthorpe/gopi/app"
)

////////////////////////////////////////////////////////////////////////////////

func Task(app *app.App,tick_duration time.Duration,name string,task_done chan bool) {
	// Tick every second
	ticker := time.Tick(tick_duration)

	// Get done channel
	finish := app.GetDoneChannel()

	// Loop until app is done
	outer_loop: for {
		select {
		case <- ticker:
			app.Logger.Info("Task %v: Tick",name)
		case <- finish:
			app.Logger.Info("Task %v: App Done Signal",name)
			break outer_loop
		}
	}

	// Cleanup task
	app.Logger.Info("Task %v: Cleanup",name)

	// Close
	task_done <- true
	app.Logger.Info("Task %v: Closed",name)
}

////////////////////////////////////////////////////////////////////////////////

func RunTasks(app *app.App) error {

	// Start task A
	task_done_a := make(chan bool)
	go Task(app,time.Second,"A",task_done_a)

	// Start task B
	task_done_b := make(chan bool)
	go Task(app,time.Millisecond * 2500,"B",task_done_b)

	// Wait until CTRL+C is pressed and all tasks have signalled completion
	app.WaitUntilDone(task_done_a,task_done_b)

	app.Logger.Info("Returning from RunTasks")

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Create the configuration, we want to use the DEVICE
	// subsystem
	config := app.Config(app.APP_DEVICE)

	// Create the application
	myapp, err := app.NewApp(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(RunTasks); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
