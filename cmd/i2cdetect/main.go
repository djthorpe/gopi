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
	app "github.com/djthorpe/gopi/v2/app"
	tablewriter "github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	if len(args) > 0 {
		return gopi.ErrHelp
	}

	i2c := app.I2C()

	// Output detected I2C addresses
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"", "-0", "-1", "-2", "-3", "-4", "-5", "-6", "-7", "-8", "-9", "-A", "-B", "-C", "-D", "-E", "-F"})
	row := make([]string, 0)

	for slave := uint8(0); slave < 0x80; slave++ {
		if len(row) == 0 {
			row = append(row, fmt.Sprintf("0x%02X", slave&0xF0))
		}
		if detected, err := i2c.DetectSlave(slave); err != nil {
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

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, nil, "i2c"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Run and exit
		os.Exit(app.Run())
	}
}
