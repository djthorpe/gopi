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
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"
)

func Main(app gopi.App, args []string) error {
	// Wait for interrupt signal
	fmt.Println("Waiting for CTRL+C")
	if err := app.WaitForSignal(context.Background(), os.Interrupt); err != nil {
		app.Log().Error(err)
	}

	// Return success
	return nil
}

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, "lirc"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		os.Exit(app.Run())
	}
}
