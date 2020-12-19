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

	fontdir, name *string
	i2cbus, port  *uint
	timeout       *time.Duration
}

type header struct {
	string
}

func (h header) Format() (string, table.Alignment, table.Color) {
	return h.string, table.Auto, table.Bold
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *app) Define(cfg gopi.Config) error {
	// Define flags
	//this.fontdir = cfg.FlagString("fontdir", "", "Font directory", "fonts")
	this.i2cbus = cfg.FlagUint("bus", 0, "I2C Bus", "i2c")
	this.timeout = cfg.FlagDuration("timeout", time.Second, "Discovery timeout", "mdns")
	this.port = cfg.FlagUint("port", 0, "Service Port", "mdns serve")
	this.name = cfg.FlagString("name", "", "Service Name", "mdns serve")

	// Define commands
	cfg.Command("info", "Hardware information", this.RunInfo)
	cfg.Command("version", "Version information", func(context.Context) error {
		return this.PrintVersion(cfg)
	})
	cfg.Command("mdns", "mDNS Service Discovery", this.RunDiscovery)
	cfg.Command("mdns serve", "Serve mDNS service record for this host", this.RunDiscoveryServe)
	cfg.Command("i2c", "Detect I2C devices", this.RunI2C)

	/*
			// TODO Define other commands
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
