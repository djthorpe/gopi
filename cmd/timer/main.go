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
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	fmt.Println("app=", app)
	fmt.Println("args=", args)
	return gopi.ErrNotImplemented
}

func main() {
	if app, err := app.NewCommandLineTool(Main, "timer"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		os.Exit(app.Run())
	}
}
