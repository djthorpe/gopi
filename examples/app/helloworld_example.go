/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// The canonical hello world example
package main

import (
	"fmt"
	"os"
	"os/user"
)

import (
	app "github.com/djthorpe/gopi/app"
)

////////////////////////////////////////////////////////////////////////////////

func HelloWorld(app *app.App) error {

	// Get name argument
	name, _ := app.FlagSet.GetString("name")

	// Output message to stdout
	fmt.Fprintf(os.Stdout, "Hello %v!!\n", name)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Create the configuration
	config := app.Config(app.APP_NONE)

	// Get Current user
	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}

	// Add -name argument
	config.FlagSet.FlagString("name", usr.Username, "Your name")

	// Create the application
	myapp, err := app.NewApp(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
	defer myapp.Close()

	// Run the application
	if err := myapp.Run(HelloWorld); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		return
	}
}
