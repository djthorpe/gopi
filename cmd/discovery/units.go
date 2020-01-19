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
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	app "github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/mdns"
)

////////////////////////////////////////////////////////////////////////////////
// EVENTS AND UNITS

var (
	Events = []gopi.EventHandler{
		gopi.EventHandler{Name: "gopi.RPCEvent", Handler: EventHandler},
	}
	Units = []string{"gopi/mdns/discovery", "register"}
)

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, Units...); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Flags
		app.Flags().FlagDuration("timeout", time.Second, "Timeout for discovery")
		app.Flags().FlagBool("register", false, "Register service")
		app.Flags().FlagBool("watch", false, "Watch for discovery events")

		// Run and exit
		os.Exit(app.Run())
	}
}
