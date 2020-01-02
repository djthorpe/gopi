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
	tablewriter "github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	if len(args) > 0 {
		return gopi.ErrHelp
	}
	platform := app.Platform()
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Append([]string{
		"PLATFORM", fmt.Sprint(platform.Type()),
	})
	table.Append([]string{
		"PRODUCT", platform.Product(),
	})
	table.Append([]string{
		"SERIAL NUMBER", fmt.Sprint(platform.SerialNumber()),
	})
	table.Append([]string{
		"UPTIME", fmt.Sprint(platform.Uptime().Truncate(time.Hour).Hours()) + " hrs",
	})
	l1, l5, l15 := platform.LoadAverages()
	table.Append([]string{
		"LOAD AVERAGES", fmt.Sprintf("%.2f %.2f %.2f", l1, l5, l15),
	})
	table.Append([]string{
		"NUMBER OF DISPLAYS", fmt.Sprint(platform.NumberOfDisplays()),
	})
	table.Render()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BOOTSTRAP

func main() {
	if app, err := app.NewCommandLineTool(Main, "platform"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		// Run and exit
		os.Exit(app.Run())
	}
}
