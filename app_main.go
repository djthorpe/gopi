/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
    Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"os"
)

var (
	start_rpc chan struct{}
)

////////////////////////////////////////////////////////////////////////////////
// COMMAND LINE TOOL STARTUP

// CommandLineTool is the basic form of running a command-line
// application, you generally call this from the main() function
func CommandLineTool(config AppConfig, main_task MainTask, background_tasks ...BackgroundTask) int {

	// Create the application
	app, err := NewAppInstance(config)
	if err != nil {
		if err != ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(main_task, background_tasks...); err == ErrHelp {
		config.AppFlags.PrintUsage()
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

// CommandLineTool2 runs multiple background tasks, waits for them all to send a
// start signl and then runs the main task
func CommandLineTool2(config AppConfig, main_task MainTask, background_tasks ...BackgroundTask2) int {

	// Create the application
	app, err := NewAppInstance(config)
	if err != nil {
		if err != ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Start main task
	if err := app.Run2(main_task, background_tasks...); err == ErrHelp {
		config.AppFlags.PrintUsage()
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	} else {
		return 0
	}
}
