/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package main

import (
	"context"
	"fmt"
	"io"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
	"github.com/djthorpe/gopi/v2/app"
	"github.com/olekukonko/tablewriter"
)

func DisplayInfo(app gopi.App, writer io.Writer) error {
	lirc := app.LIRC()
	table := tablewriter.NewWriter(writer)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	if recv_mode := lirc.RcvMode(); recv_mode != gopi.LIRC_MODE_NONE {
		table.Append([]string{"RECV_MODE", fmt.Sprint(recv_mode)})
		if duty_cycle := lirc.RcvDutyCycle(); duty_cycle != 0 {
			table.Append([]string{"RECV_DUTY_CYCLE", fmt.Sprintf("%d%%", duty_cycle)})
		}
	}
	if send_mode := lirc.SendMode(); send_mode != gopi.LIRC_MODE_NONE {
		table.Append([]string{"SEND_MODE", fmt.Sprint(send_mode)})
		if duty_cycle := lirc.SendDutyCycle(); duty_cycle != 0 {
			table.Append([]string{"SEND_DUTY_CYCLE", fmt.Sprintf("%d%%", duty_cycle)})
		}
	}

	table.Render()
	return nil
}

func Main(app gopi.App, args []string) error {
	if app.Flags().GetBool("info", gopi.FLAG_NS_DEFAULT) {
		// Display LIRC information
		if err := DisplayInfo(app, os.Stdout); err != nil {
			return err
		}
	} else {
		// Wait for interrupt signal
		fmt.Println("Waiting for CTRL+C")
		if err := app.WaitForSignal(context.Background(), os.Interrupt); err != nil {
			app.Log().Error(err)
		}
	}

	// Return success
	return nil
}

func main() {
	if app, err := app.NewCommandLineTool(Main, Events, "lirc"); err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		app.Flags().FlagBool("info", false, "Show information")
		os.Exit(app.Run())
	}
}
