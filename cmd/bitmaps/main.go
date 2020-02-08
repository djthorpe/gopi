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

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	if len(args) > 0 {
		return gopi.ErrHelp
	}

	surfaces := app.Surfaces()
	fmt.Println(surfaces)

	// Wait for key press
	app.Log().Info("Press CTRL+C to close")
	app.WaitForSignal(context.Background(), os.Interrupt)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, nil, "surfaces"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Run and exit
		os.Exit(app.Run())
	}
}
