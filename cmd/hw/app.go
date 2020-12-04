package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.Publisher
	gopi.Platform
	gopi.GPIO
	gopi.I2C
	gopi.FontManager
	gopi.Command

	fontdir *string
	i2cbus  *uint
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *app) Define(cfg gopi.Config) error {
	// Define flags
	this.fontdir = cfg.FlagString("fontdir", "", "Font directory")
	this.i2cbus = cfg.FlagUint("i2c.bus", 0, "I2C Bus")

	// Define commands
	cfg.Command("version", "Return information about the command", func(context.Context) error {
		if err := this.PrintVersion(cfg); err != nil {
			return err
		} else {
			return gopi.ErrHelp
		}
	})

	cfg.Command("hw", "Return hardware platform information", this.RunHardware)
	cfg.Command("spi", "Return SPI interface parameters", nil)
	cfg.Command("i2c", "Return I2C interface parameters", this.RunI2C)
	cfg.Command("gpio", "Control GPIO interface", this.RunGPIO)
	cfg.Command("fonts", "Return Font faces", this.RunFonts) // Not yet implemented

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	// Set the command to run
	if this.Command = cfg.GetCommand(nil); this.Command == nil {
		return gopi.ErrHelp
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) RunHardware(context.Context) error {
	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.Append([]string{
		"Product", this.Platform.Product(), fmt.Sprint(this.Platform.Type()),
	})
	table.Append([]string{
		"Serial Number", "", this.Platform.SerialNumber(),
	})
	table.Append([]string{
		"Uptime", "", this.Platform.Uptime().Truncate(time.Second).String(),
	})
	if l1, l5, l15 := this.Platform.LoadAverages(); l1 != 0 && l5 != 0 && l15 != 0 {
		table.AppendBulk([][]string{
			{"Load Averages", "1m", fmt.Sprintf("%.2f", l1)},
			{"Load Averages", "5m", fmt.Sprintf("%.2f", l5)},
			{"Load Averages", "15m", fmt.Sprintf("%.2f", l15)},
		})
	}
	if zones := this.Platform.TemperatureZones(); len(zones) > 0 {
		for k, v := range zones {
			table.Append([]string{
				"Temperature Zones", k, fmt.Sprintf("%.2fC", v),
			})
		}
	}
	table.Append([]string{
		"Number of Displays", "", fmt.Sprint(this.Platform.NumberOfDisplays()),
	})
	table.Append([]string{
		"Attached Displays", "", fmt.Sprint(this.Platform.AttachedDisplays()),
	})
	table.Render()

	// Return success
	return nil
}

func (this *app) RunFonts(context.Context) error {
	if this.fontdir == nil || *this.fontdir == "" {
		return gopi.ErrBadParameter.WithPrefix("Missing -fontdir flag")
	} else if stat, err := os.Stat(*this.fontdir); os.IsNotExist(err) || stat.IsDir() == false {
		return gopi.ErrBadParameter.WithPrefix("Invalid -fontdir flag")
	}

	manager := this.FontManager
	if err := manager.OpenFacesAtPath(*this.fontdir, nil); err != nil {
		return err
	}
	if families := manager.Families(); len(families) == 0 {
		return fmt.Errorf("No fonts found")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)

	table.SetHeader([]string{"Family", "Name", "Style", "Glyphs"})
	for _, family := range manager.Families() {
		for _, face := range manager.Faces(family, gopi.FONT_FLAGS_STYLE_ANY) {
			table.Append([]string{
				family,
				face.Name(),
				face.Style(),
				fmt.Sprint(face.NumGlyphs()),
			})
		}
	}

	table.Render()

	// Return success
	return nil
}
