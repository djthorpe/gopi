/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Detects slaves on the I2C bus
package main

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func runLoop(app *gopi.AppInstance, done chan struct{}) error {

	if app.I2C == nil {
		return app.Logger.Error("Missing I2C module instance")
	}

	for slave := uint8(0); slave < 0x80; slave++ {
		if detected, err := app.I2C.DetectSlave(slave); err != nil {
			return err
		} else {
			fmt.Printf("0x%02X -> %v\n", slave, detected)
		}
	}

	// Finished
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main_inner() int {
	// Create the configuration
	config := gopi.NewAppConfig("i2c")

	// Create the application
	app, err := gopi.NewAppInstance(config)
	if err != nil {
		if err != gopi.ErrHelp {
			fmt.Fprintln(os.Stderr, err)
			return -1
		}
		return 0
	}
	defer app.Close()

	// Run the application
	if err := app.Run(runLoop); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return -1
	}
	return 0
}

func main() {
	os.Exit(main_inner())
}
