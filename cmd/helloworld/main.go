/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"fmt"
	"os"
	"os/user"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {

	// Print out name
	fmt.Println("Hello, " + app.Flags().GetString("name", gopi.FLAG_NS_DEFAULT))
	fmt.Println("Press CTRL+C to exit")

	// Wait for CTRL+C
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else if user, err := user.Current(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Register -name flag
		app.Flags().FlagString("name", user.Name, "Name of user to print")

		// Run and exit
		os.Exit(app.Run())
	}
}
