/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

// Outputs a table of metrics
package main

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/olekukonko/tablewriter"

	// Modules
	_ "github.com/djthorpe/gopi/sys/hw/darwin"
	_ "github.com/djthorpe/gopi/sys/hw/linux"
	_ "github.com/djthorpe/gopi/sys/logger"
)

////////////////////////////////////////////////////////////////////////////////

func Main(app *gopi.AppInstance, done chan<- struct{}) error {
	metrics := app.ModuleInstance("metrics").(gopi.Metrics)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Value"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Append([]string{"Host uptime", fmt.Sprint(metrics.UptimeHost())})
	table.Append([]string{"App uptime", fmt.Sprint(metrics.UptimeApp())})

	l1, l5, l15 := metrics.LoadAverage()
	table.Append([]string{"1-min Load Average", fmt.Sprintf("%.2f", l1)})
	table.Append([]string{"5-min Load Average", fmt.Sprintf("%.2f", l5)})
	table.Append([]string{"15-min Load Average", fmt.Sprintf("%.2f", l15)})

	table.Render()

	return nil
}

func main() {
	config := gopi.NewAppConfig("metrics")

	// Run the command line tool
	os.Exit(gopi.CommandLineTool(config, Main))
}
