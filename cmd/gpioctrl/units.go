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
	"github.com/djthorpe/gopi/v2/app"

	// Units
	_ "github.com/djthorpe/gopi/v2/unit/files"
	_ "github.com/djthorpe/gopi/v2/unit/gpio"
	_ "github.com/djthorpe/gopi/v2/unit/gpio/gpiorpi"
	_ "github.com/djthorpe/gopi/v2/unit/gpio/gpiosysfs"
	_ "github.com/djthorpe/gopi/v2/unit/logger"
	_ "github.com/djthorpe/gopi/v2/unit/platform"
)

func main() {
	if app, err := app.NewCommandLineTool(Main, nil, "gpio"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		app.Flags().FlagString("low", "", "Comma-separated list of pins to set to low")
		app.Flags().FlagString("high", "", "Comma-separated list of pins to set to high")
		app.Flags().FlagString("in", "", "Comma-separated list of pins to set to input")
		app.Flags().FlagString("edge", "", "Comma-separated list of pins to monitor for changes")
		os.Exit(app.Run())
	}
}
