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
	app "github.com/djthorpe/gopi/app"
)

////////////////////////////////////////////////////////////////////////////////

func Task(app *app.App, task_name string, task_done chan bool) {
	// Tick every second
	ticker := time.Tick(time.Second)

	// Get done channel
	finish := app.GetDoneChannel()

	// Loop until app is done
outer_loop:
	for {
		select {
		case <-ticker:
			app.Logger.Info("Task %v: Tick", task_name)
		case <-finish:
			app.Logger.Info("Task %v: App Done Signal", task_name)
			break outer_loop
		}
	}

	// Cleanup task
	app.Logger.Info("Task %v: Cleanup", task_name)

	// Close
	task_done <- true
	app.Logger.Info("Task %v: Closed", task_name)
}

////////////////////////////////////////////////////////////////////////////////

func RunTasks(app *app.App) error {

	// Wait until CTRL+C is pressed and all tasks have signalled completion
	app.WaitUntilDone(app.RunTask(Task, "app.TaskA"), app.RunTask(Task, "app.TaskB"))

	// Return success
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

	// Run the application
	if err := app.Run(helloWorld); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
