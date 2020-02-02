/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"fmt"
	"os"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/files"
	_ "github.com/djthorpe/gopi/v2/unit/input"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
)

var (
	Events = []gopi.EventHandler{
		gopi.EventHandler{Name: "gopi.InputEvent", Handler: EventHandler},
	}
)

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, "input"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Set flags
		app.Flags().FlagString("input.name", "", "Open specific input device")
		app.Flags().FlagString("input.type", "", "Device type (mouse, keyboard, joystick or touch)")
		app.Flags().FlagBool("input.exclusive", true, "Grab input device")
		app.Flags().FlagBool("watch", false, "Watch for input events")

		// Run and exit
		os.Exit(app.Run())
	}
}
