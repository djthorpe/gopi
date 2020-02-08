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
	"strings"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	display "github.com/djthorpe/gopi/v2/unit/display"
	tablewriter "github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////

type (
	CommandFunc func(app gopi.App) error
	Command     struct {
		Name        string
		Description string
		Command     CommandFunc
	}
)

var (
	Commands = []Command{
		Command{"platform", "Return information about the hardware platform", PlatformCommand},
		Command{"displays", "Return information about the displays", DisplaysCommand},
	}
)

////////////////////////////////////////////////////////////////////////////////

func Main(app gopi.App, args []string) error {
	if len(args) == 0 {
		return Commands[0].Command(app)
	} else {
		for _, command := range args {
			if err := RunCommand(app, command); err != nil {
				return err
			}
		}
	}

	// Return success
	return nil
}

func RunCommand(app gopi.App, name string) error {
	for _, cmd := range Commands {
		if cmd.Name == strings.ToLower(name) {
			return cmd.Command(app)
		}
	}

	// Return not found
	return gopi.ErrNotFound.WithPrefix(name)
}

func PlatformCommand(app gopi.App) error {
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

func DisplaysCommand(app gopi.App) error {
	platform := app.Platform()
	if platform.NumberOfDisplays() == 0 {
		return fmt.Errorf("No displays found")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, i := range platform.AttachedDisplays() {
		if display_, err := gopi.New(display.Display{
			Id:       i,
			Platform: app.Platform(),
		}, app.Log()); err != nil {
			return fmt.Errorf("Display %v: %w", i, err)
		} else {
			display := display_.(gopi.Display)
			defer display.Close()

			w, h := display.Size()
			ppi := display.PixelsPerInch()
			ppiStr := fmt.Sprint(ppi)
			if ppi == 0 {
				ppiStr = "-"
			}
			table.Append([]string{
				fmt.Sprint(display.DisplayId()),
				display.Name(),
				fmt.Sprintf("%vx%v", w, h),
				ppiStr,
			})
		}
	}
	table.Render()

	// Return success
	return nil
}
