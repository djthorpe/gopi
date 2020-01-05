/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

// The canonical hello world example demonstrates printing hello world and then exiting.
// Here we use the 'generic' set of modules which provide generic system services
package main

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func helloWorld(app *gopi.AppInstance, done chan<- struct{}) error {
	// If -name argument is used then use that, else output generic message
	if name, exists := app.AppFlags.GetString("name"); exists {
		fmt.Println("Hello,", name)
	} else {
		fmt.Println("Hello, World (use -name flag to specify your name)")
	}

	// If wait flag is set, then wait until CTRL+C is pressed to continue
	if wait, _ := app.AppFlags.GetBool("wait"); wait {
		fmt.Println("Press CTRL+C to exit")
		app.WaitForSignal()
	}

	// Signal that main thread is done
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig()
	config.AppFlags.FlagString("name", "", "Your name")
	config.AppFlags.FlagBool("wait", false, "Wait for CTRL+C interrupt to end")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, helloWorld))
}
