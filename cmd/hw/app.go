package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
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
	gopi.LIRC
	gopi.ServiceDiscovery
	gopi.FontManager
	gopi.Command

	fontdir *string
	i2cbus  *uint
	timeout *time.Duration
}

type header struct {
	string
}

func (h header) Format() (string, table.Alignment, table.Color) {
	return "[" + h.string + "]", table.Auto, table.White | table.Inverse
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *app) Define(cfg gopi.Config) error {
	// Define flags
	this.fontdir = cfg.FlagString("fontdir", "", "Font directory", "fonts")
	this.i2cbus = cfg.FlagUint("bus", 0, "I2C Bus", "i2c")
	this.timeout = cfg.FlagDuration("timeout", time.Second, "Discovery timeout", "mdns")

	// Define commands
	cfg.Command("version", "Return information about the command", func(context.Context) error {
		if err := this.PrintVersion(cfg); err != nil {
			return err
		} else {
			return gopi.ErrHelp
		}
	})
	cfg.Command("lirc", "IR sending and receiving control", func(ctx context.Context) error {
		return this.RunLIRC(ctx, cfg)
	})
	cfg.Command("lirc print", "Print LIRC Parameters", func(ctx context.Context) error {
		return nil
	})

	/*
		// Define mDNS command
		cfg.Command("mdns", "Service discovery", this.RunDiscovery)

		// Define other commands
		cfg.Command("hw", "Return hardware platform information", this.RunHardware)
		cfg.Command("i2c", "Return I2C interface parameters", this.RunI2C)
		cfg.Command("gpio", "Control GPIO interface", this.RunGPIO)
		cfg.Command("fonts", "Return Font faces", this.RunFonts) // Not yet implemented
	*/

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	// Set the command to run
	if cmd, err := cfg.GetCommand(nil); err != nil {
		return err
	} else if cmd == nil {
		return gopi.ErrHelp
	} else {
		this.Command = cmd
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) RunHardware(context.Context) error {
	// Display platform information
	table := table.New(table.WithHeader(false), table.WithMergeCells())

	table.Append(header{"Product"}, this.Platform.Product(), fmt.Sprint(this.Platform.Type()))
	table.Append("Serial Number", "", this.Platform.SerialNumber())
	table.Append("Uptime", "", this.Platform.Uptime().Truncate(time.Second).String())
	if l1, l5, l15 := this.Platform.LoadAverages(); l1 != 0 && l5 != 0 && l15 != 0 {
		table.Append("Load Averages", "1m", fmt.Sprintf("%.2f", l1))
		table.Append("Load Averages", "5m", fmt.Sprintf("%.2f", l5))
		table.Append("Load Averages", "15m", fmt.Sprintf("%.2f", l15))
	}
	if zones := this.Platform.TemperatureZones(); len(zones) > 0 {
		for k, v := range zones {
			table.Append("Temperature Zones", k, fmt.Sprintf("%.2fC", v))
		}
	}
	table.Append("Number of Displays", "", fmt.Sprint(this.Platform.NumberOfDisplays()))
	table.Append("Attached Displays", "", fmt.Sprint(this.Platform.AttachedDisplays()))
	table.Render(os.Stdout)

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
