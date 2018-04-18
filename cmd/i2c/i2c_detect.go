/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
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
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func mainLoop(app *gopi.AppInstance, done chan<- struct{}) error {

	if app.I2C == nil {
		return app.Logger.Error("Missing I2C module instance")
	}

	// Output detected I2C addresses
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"", "-0", "-1", "-2", "-3", "-4", "-5", "-6", "-7", "-8", "-9", "-A", "-B", "-C", "-D", "-E", "-F"})
	row := make([]string, 0)

	for slave := uint8(0); slave < 0x80; slave++ {
		if len(row) == 0 {
			row = append(row, fmt.Sprintf("0x%02X", slave&0xF0))
		}
		if detected, err := app.I2C.DetectSlave(slave); err != nil {
			return err
		} else if detected {
			row = append(row, fmt.Sprintf("%02X", slave))
		} else {
			row = append(row, "--")
		}
		if len(row) >= 17 {
			table.Append(row)
			row = make([]string, 0)
		}
	}
	table.Render()

	// Finished
	done <- gopi.DONE
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration, load the i2c instance
	config := gopi.NewAppConfig("i2c")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, mainLoop))
}
