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
	watch        *bool
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
	this.timeout = cfg.FlagDuration("timeout", time.Second, "Discovery timeout", "mdns", "cast", "cast app", "cast vol", "cast mute", "cast unmute")
	this.name = cfg.FlagString("name", "", "Service or Chromecast Name", "mdns serve", "cast app", "cast vol", "cast mute", "cast unmute")
	this.port = cfg.FlagUint("port", 0, "Service Port", "mdns serve")
	this.watch = cfg.FlagBool("watch", false, "Watch for events", "cast app", "cast vol", "cast mute", "cast unmute", "cast load")

	// Define commands
	cfg.Command("info", "Hardware information", this.RunInfo)
	cfg.Command("version", "Version information", func(context.Context) error {
		return this.PrintVersion(cfg)
	})
	cfg.Command("mdns", "mDNS Service Discovery", this.RunDiscovery)
	cfg.Command("mdns serve", "Serve mDNS service record for this host", this.RunDiscoveryServe)

	cfg.Command("i2c", "Detect I2C devices", this.RunI2C)

	cfg.Command("cast", "List Chromecast devices", this.RunCast)
	cfg.Command("cast app", "Launch Application", this.RunCastApp)
	cfg.Command("cast vol", "Set Chromecast volume", this.RunCastVol)
	cfg.Command("cast mute", "Mute Chromecast volume", this.RunCastMute)
	cfg.Command("cast unmute", "Unmute Chromecast volume", this.RunCastUnmute)
	cfg.Command("cast load", "Play media on Chromecast", this.RunCastLoad)
	//	cfg.Command("cast play", "Set play state on Chromecast", this.RunCastPlay)
	//	cfg.Command("cast pause", "Set pause state on Chromecast", this.RunCastPause)
	//	cfg.Command("cast stop", "Set stop state on Chromecast", this.RunCastStop)

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
