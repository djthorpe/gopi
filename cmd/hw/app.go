package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/hw/platform"
	"github.com/djthorpe/gopi/v3/pkg/hw/spi"
	"github.com/djthorpe/gopi/v3/pkg/log"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	*log.Log
	*platform.Platform
	*spi.Devices
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *app) Run(context.Context) error {
	this.PlatformTable()

	if err := this.SPITable(); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PLATFORM

func (this *app) PlatformTable() {
	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.Append([]string{
		"Product", this.Platform.Product(),
	})
	table.Append([]string{
		"Type", fmt.Sprint(this.Platform.Type()),
	})
	table.Append([]string{
		"Serial Number", this.Platform.SerialNumber(),
	})
	table.Append([]string{
		"Uptime", this.Platform.Uptime().Truncate(time.Second).String(),
	})
	table.Append([]string{
		"Load Averages", fmt.Sprint(this.Platform.LoadAverages()),
	})
	table.Append([]string{
		"Number of Displays", fmt.Sprint(this.Platform.NumberOfDisplays()),
	})
	table.Append([]string{
		"Attached Displays", fmt.Sprint(this.Platform.AttachedDisplays()),
	})
	table.Render()
}

////////////////////////////////////////////////////////////////////////////////
// SPI

func (this *app) SPITable() error {
	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, dev := range this.Devices.Enumerate() {
		spi, err := this.Devices.Open(dev, 0)
		if err != nil {
			return err
		}
		table.Append([]string{
			dev.String(),
			fmt.Sprint(spi),
		})
	}
	table.Render()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *app) String() string {
	str := "<app"
	str += " platform=" + fmt.Sprint(this.Platform)
	str += " spi=" + fmt.Sprint(this.Devices)
	return str + ">"
}
