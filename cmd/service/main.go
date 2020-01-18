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
	app "github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi-rpc/v2/unit/g"
	_ "github.com/djthorpe/gopi/v2/unit/bus"
	_ "github.com/djthorpe/gopi/v2/unit/loggerrpc"
)

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewServer(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Run and exit
		os.Exit(app.Run())
	}
}
