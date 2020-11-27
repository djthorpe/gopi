package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.Server
	gopi.PingService
	gopi.Logger
	gopi.Command
}

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("version", "Return version information", this.Version(cfg.Version()))
	cfg.Command("serve", "Start server", this.Serve)
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if this.Command = cfg.GetCommand(cfg.Args()); this.Command == nil {
		return gopi.ErrHelp
	}
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) Version(v gopi.Version) func(context.Context) error {
	return func(context.Context) error {
		if name := v.Name(); name != "" {
			fmt.Printf("%-10s %s\n", "Name", v.Name())
		}
		tag, branch, hash := v.Version()
		if tag != "" {
			fmt.Printf("%-10s %s\n", "Tag", tag)
		}
		if branch != "" {
			fmt.Printf("%-10s %s\n", "Branch", branch)
		}
		if hash != "" {
			fmt.Printf("%-10s %s\n", "Hash", hash)
		}
		if buildtime := v.BuildTime(); buildtime.IsZero() == false {
			fmt.Printf("%-10s %s\n", "BuildTime", buildtime.Format(time.RFC3339))
		}
		if version := v.GoVersion(); version != "" {
			fmt.Printf("%-10s %s\n", "GoVersion", version)
		}
		return nil
	}
}

func (this *app) Serve(ctx context.Context) error {
	var network, addr string

	args := this.Command.Args()
	switch {
	case len(args) == 0:
		network = "tcp"
		addr = ":0"
	case len(args) == 1:
		if _, _, err := net.SplitHostPort(args[0]); err != nil {
			return err
		} else {
			network = "tcp"
			addr = args[0]
		}
	default:
		return gopi.ErrBadParameter.WithPrefix(this.Command.Name())
	}

	if err := this.Server.StartInBackground(network, addr); err != nil {
		return err
	}

	fmt.Println("Started server, ", this.Server)
	fmt.Println("Press CTRL+C to end")

	// Wait until done
	<-ctx.Done()

	// Close gracefully
	return this.Server.Stop(false)
}
