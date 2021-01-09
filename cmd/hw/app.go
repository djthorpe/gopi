package main

import (
	"context"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/table"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	sync.WaitGroup
	gopi.Unit
	gopi.CastManager
	gopi.Logger
	gopi.Publisher
	gopi.Platform
	gopi.GPIO
	gopi.I2C
	gopi.LIRC
	gopi.ServiceDiscovery
	gopi.FontManager
	gopi.Command

	name         *string
	i2cbus, port *uint
	timeout      *time.Duration
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
	this.i2cbus = cfg.FlagUint("bus", 0, "I2C Bus", "i2c")
	this.timeout = cfg.FlagDuration("timeout", time.Second, "Discovery timeout", "mdns")
	this.name = cfg.FlagString("name", "", "Service", "mdns serve")
	this.port = cfg.FlagUint("port", 0, "Service Port", "mdns serve")

	// Define commands
	cfg.Command("info", "Hardware information", this.RunInfo)
	cfg.Command("version", "Version information", func(context.Context) error {
		return this.PrintVersion(cfg)
	})
	cfg.Command("mdns", "mDNS Service Discovery", this.RunDiscovery)
	cfg.Command("mdns serve", "Serve mDNS service record for this host", this.RunDiscoveryServe)

	cfg.Command("i2c", "Detect I2C devices", this.RunI2C)

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
