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
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

func RenderTable(app gopi.App) {
	gpio := app.GPIO()
	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader([]string{"Physical", "Logical", "Direction", "Value"})

	// Physical pins start at index 1
	for pin := uint(1); pin <= gpio.NumberOfPhysicalPins(); pin++ {
		var l, d, v string
		if logical := gpio.PhysicalPin(pin); logical != gopi.GPIO_PIN_NONE {
			l = fmt.Sprint(logical)
			d = fmt.Sprint(gpio.GetPinMode(logical))
			v = fmt.Sprint(gpio.ReadPin(logical))
		}
		table.Append([]string{
			fmt.Sprintf("%v", pin), l, d, v,
		})
	}

	table.Render()
}
