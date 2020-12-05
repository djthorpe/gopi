package main

import (
	"context"
	"net"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.ConnPool
	gopi.Command
	gopi.Server
	gopi.PingService
	gopi.Logger
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
		if stub, err := this.GetStub(); err != nil {
			return err
		} else {
			return this.RunVersion(ctx, stub)
		}
	})
	cfg.Command("ping", "Perform ping to server", func(ctx context.Context) error {
		if stub, err := this.GetStub(); err != nil {
			return err
		} else {
			return this.RunPing(ctx, stub)
		}
	})

	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if this.Command = cfg.GetCommand(cfg.Args()); this.Command == nil {
		return gopi.ErrHelp
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) GetStub() (gopi.PingStub, error) {
	args := this.Args()
	if len(args) != 1 {
		return nil, gopi.ErrBadParameter.WithPrefix("Missing server address")
	} else if conn, err := this.ConnPool.Connect("tcp", args[0]); err != nil {
		return nil, err
	} else if stub, _ := conn.NewStub("gopi.ping.Ping").(gopi.PingStub); stub == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("Cannot create stub")
	} else {
		return stub, nil
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
