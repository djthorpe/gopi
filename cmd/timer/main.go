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
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	app.Log().Debug("timer=", app.Timer())

	// Set up handlers
	app.Bus().NewHandler("gopi.TimerEvent", func(evt gopi.Event) {
		app.Log().Debug("event=", evt)
	})

	// Schedule a ticker which fires every second
	app.Timer().NewTicker(time.Second)

	// Wait for interrupt signal
	if err := app.WaitForSignal(context.Background(), os.Interrupt); err != nil {
		app.Log().Error(err)
	}

	// Return success
	return nil
}

func main() {
	if app, err := app.NewCommandLineTool(Main, "timer"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		os.Exit(app.Run())
	}
}
