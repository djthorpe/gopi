package main

import (
	"context"
	"net"
	"strings"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.ConnPool
	gopi.Command
	gopi.Logger
	gopi.PingService
	gopi.MetricsService
	gopi.Server
	gopi.ServiceDiscovery
	gopi.Unit

	// Flags
	service, castId *string
}

func (this *app) Define(cfg gopi.Config) error {
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
	cfg.Command("cast", "List Google Chromecasts", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCast(ctx, stub)
		}
	})
	cfg.Command("cast app", "Start Chromecast Application", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCastApp(ctx, stub)
		}
	})
	cfg.Command("cast load", "Load media from URL", func(ctx context.Context) error {
		if stub, err := this.GetGoogleCastStub(); err != nil {
			return err
		} else {
			return this.RunCastLoad(ctx, stub)
		}
	})

	// Global flags
	this.service = cfg.FlagString("srv", "", "name, service:name or host:port")

	// Set flags for cast functions
	this.castId = cfg.FlagString("id", "", "Chromecast Id", "cast app", "cast load")

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
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
	return this.Command.Run(ctx)
}

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

func (this *app) GetGoogleCastStub() (gopi.CastStub, error) {
	if stub, err := this.GetStub("gopi.googlecast.Manager"); err != nil {
		return nil, err
	} else {
		return stub.(gopi.CastStub), nil
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
