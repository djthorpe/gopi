/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"os"
)

// CommandLineTool is the basic form of running a command-line
// application, you generally call this from the main() function
func CommandLineTool(config AppConfig, main_task MainTask, background_tasks ...BackgroundTask) int {
	// Set the usage function
	config.AppFlags.SetUsageFunc(func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", config.AppFlags.Name())
		config.AppFlags.PrintDefaults()
	})

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
