package main

import (
	"context"
	"strings"
	"time"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.ConnPool
	gopi.Logger
	gopi.ServiceDiscovery
	gopi.Command
	Chromecast
	Rotel

	service *string
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	KeyStub = "Stub"
	KeyArgs = "Args"
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Define(cfg gopi.Config) error {
	this.Chromecast.Define(cfg)
	this.Rotel.Define(cfg)

	// Global flags
	this.service = cfg.FlagString("srv", "", "name, service:name or host:port")

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.ConnPool, this.Logger, this.ServiceDiscovery)

	if cmd, err := cfg.GetCommand(cfg.Args()); err != nil {
		return gopi.ErrHelp
	} else if cmd == nil {
		return gopi.ErrHelp
	} else {
		this.Command = cmd
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	name := this.Command.Name()
	switch {
	case strings.HasPrefix(name, "rotel"):
		name = "gopi.rotel.Manager"
	case strings.HasPrefix(name, "cast"):
		name = "gopi.chromecast.Manager"
	}
	if stub, err := this.GetStub(name); err != nil {
		return err
	} else {
		ctx = context.WithValue(ctx, KeyStub, stub)
		ctx = context.WithValue(ctx, KeyArgs, this.Command.Args())
		return this.Command.Run(ctx)
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET STUB

func (this *app) GetStub(name string) (gopi.ServiceStub, error) {
	// Timeout for lookup after 500ms
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	service := "grpc"
	if *this.service != "" {
		if strings.Contains(*this.service, ":") == false {
			service = service + ":" + *this.service
		} else {
			service = *this.service
		}
	}

	if conn, err := this.ConnPool.ConnectService(ctx, "tcp", service, 0); err != nil {
		return nil, err
	} else if stub := conn.NewStub(name); stub == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("Cannot create stub: ", name)
	} else {
		return stub, nil
	}
}

/*

	// Set flags for cast functions
	this.castId = cfg.FlagString("id", "", "Chromecast Id", "cast", "cast app", "cast load", "cast seek", "cast pause", "cast vol")

	// Set watch flag
	this.watch = cfg.FlagBool("watch", false, "Watch for events", "cast")
*/

/*
func (this *app) GetPingStub() (gopi.PingStub, error) {
	if stub, err := this.GetStub("gopi.ping.Ping"); err != nil {
		return nil, err
	} else {
		return stub.(gopi.PingStub), nil
	}
}

func (this *app) GetMetricsStub() (gopi.MetricsStub, error) {
	if stub, err := this.GetStub("gopi.metrics.Metrics"); err != nil {
		return nil, err
	} else {
		return stub.(gopi.MetricsStub), nil
	}
}

func (this *app) GetServeAddress() (string, string, error) {
	var network, addr string

	args := this.Args()
	switch {
	case len(args) == 0:
		network = "tcp"
		addr = ":0"
	case len(args) == 1:
		if _, _, err := net.SplitHostPort(args[0]); err != nil {
			return "", "", err
		} else {
			network = "tcp"
			addr = args[0]
		}
	default:
		return "", "", gopi.ErrBadParameter.WithPrefix(this.Command.Name())
	}

	// Return success
	return network, addr, nil
}
*/
/*
	cfg.Command("server", "Start ping service", func(ctx context.Context) error {
		if network, addr, err := this.GetServeAddress(); err != nil {
			return err
		} else {
			return this.RunServer(ctx, network, addr)
		}
	})
	cfg.Command("version", "Display server version information", func(ctx context.Context) error {
		if stub, err := this.GetPingStub(); err != nil {
			return err
		} else {
			return this.RunVersion(ctx, stub)
		}
	})
	cfg.Command("ping", "Perform ping to server", func(ctx context.Context) error {
		if stub, err := this.GetPingStub(); err != nil {
			return err
		} else {
			return this.RunPing(ctx, stub)
		}
	})
	cfg.Command("metrics", "Retrieve metrics from server", func(ctx context.Context) error {
		if stub, err := this.GetMetricsStub(); err != nil {
			return err
		} else {
			return this.RunMetrics(ctx, stub)
		}
	})
*/
