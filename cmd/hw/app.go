package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/hw/gpiobcm"
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
	*gpiobcm.GPIO
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *app) Run(context.Context) error {

	// Output platform information
	this.PlatformTable()

	// Output SPI information
	if spi := this.Devices.Enumerate(); len(spi) > 0 {
		if err := this.SPITable(this.Devices.Enumerate()); err != nil {
			return err
		}
	}

	// Output GPIO information
	this.GPIOTable()

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

func (this *app) SPITable(devices []spi.Device) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, dev := range devices {
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
// GPIO

func (this *app) GPIOTable() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeader([]string{"Physical", "Logical", "Direction", "Value"})

	// Physical pins start at index 1
	for pin := uint(1); pin <= this.GPIO.NumberOfPhysicalPins(); pin++ {
		var l, d, v string
		if logical := this.GPIO.PhysicalPin(pin); logical != gopi.GPIO_PIN_NONE {
			l = fmt.Sprint(logical)
			d = fmt.Sprint(this.GPIO.GetPinMode(logical))
			v = fmt.Sprint(this.GPIO.ReadPin(logical))
		}
		table.Append([]string{
			fmt.Sprintf("%v", pin), l, d, v,
		})
	}

	table.Render()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *app) String() string {
	str := "<app"
	str += " platform=" + fmt.Sprint(this.Platform)
	str += " spi=" + fmt.Sprint(this.Devices)
	return str + ">"
}
