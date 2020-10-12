package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/hw/display"
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
	*display.Displays
	*spi.Devices
	*gpiobcm.GPIO

	cmd gopi.Command
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *app) Define(cfg gopi.Config) error {
	// Define commands
	cfg.Command("hw", "Return hardware platform information", this.RunHardware)
	cfg.Command("display", "Return display information", this.RunDisplays)
	cfg.Command("spi", "Return SPI interface parameters", this.RunSpi)
	cfg.Command("i2c", "Return I2C interface parameters", nil) // Not yet implemented
	cfg.Command("gpio", "Return GPIO interface parameters", this.RunGpio)

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if cmd := cfg.GetCommand(nil); cmd == nil {
		return gopi.ErrHelp
	} else {
		this.cmd = cmd
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.cmd.Run(ctx)
}

func (this *app) RunHardware(context.Context) error {
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

	// Return success
	return nil
}

func (this *app) RunDisplays(context.Context) error {
	displays := this.Displays.Enumerate()
	if len(displays) == 0 {
		return fmt.Errorf("No Displays found")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	for _, display := range displays {
		w, h := display.Size()
		ppi_ := "-"
		if ppi := display.PixelsPerInch(); ppi != 0 {
			ppi_ = fmt.Sprint(ppi)
		}
		table.Append([]string{
			fmt.Sprint(display.Id()),
			fmt.Sprint(display.Name()),
			fmt.Sprint("{", w, ",", h, "}"),
			ppi_,
		})
	}
	table.Render()

	// Return success
	return nil
}

func (this *app) RunSpi(context.Context) error {
	devices := this.Devices.Enumerate()
	if len(devices) == 0 {
		return fmt.Errorf("No SPI interfaces found")
	}

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

func (this *app) RunGpio(context.Context) error {
	pins := this.GPIO.NumberOfPhysicalPins()
	if pins == 0 {
		return fmt.Errorf("No GPIO interface defined")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.SetHeader([]string{"Physical", "Logical", "Direction", "Value"})

	// Physical pins start at index 1
	for pin := uint(1); pin <= pins; pin++ {
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

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *app) String() string {
	str := "<app"
	str += " platform=" + fmt.Sprint(this.Platform)
	str += " spi=" + fmt.Sprint(this.Devices)
	str += " gpio=" + fmt.Sprint(this.GPIO)
	return str + ">"
}
