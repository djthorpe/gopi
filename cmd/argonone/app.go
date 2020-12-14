package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

type app struct {
	gopi.Unit
	gopi.ArgonOne
	gopi.Command
	gopi.ConnPool
	gopi.InputService
	gopi.LIRC
	gopi.LIRCKeycodeManager
	gopi.Logger
	gopi.MetricWriter
	gopi.PingService
	gopi.Publisher
	gopi.Server
	gopi.ServiceDiscovery
}

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("daemon", "Start daemon", this.RunServe)
	cfg.Command("version", "Server version", this.RunVersion)
	cfg.Command("stream", "Stream events from Argonone", this.RunStream)
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

func (this *app) RunServe(ctx context.Context) error {
	if network, addr, err := this.GetServeAddress(); err != nil {
		return err
	} else if err := this.Server.StartInBackground(network, addr); err != nil {
		return err
	}

	fmt.Println("Started server, ", this.Server)
	fmt.Println("Press CTRL+C to end")

	// Wait until done
	<-ctx.Done()

	// Close gracefully
	return this.Server.Stop(false)
}

func (this *app) RunVersion(ctx context.Context) error {
	if stub, err := this.GetPingStub(); err != nil {
		return err
	} else if version, err := stub.Version(ctx); err != nil {
		return err
	} else {
		PrintVersion(version)
	}

	// Return success
	return nil
}

func (this *app) RunStream(ctx context.Context) error {
	// Make a channel to receive the input events
	ch := make(chan gopi.InputEvent)
	defer close(ch)

	// Receive events in background
	go func(ch <-chan gopi.InputEvent) {
		for evt := range ch {
			fmt.Println(evt)
		}
	}(ch)

	if stub, err := this.GetInputStub(); err != nil {
		return err
	} else if err := stub.Stream(ctx, ch); err != nil {
		return err
	}

	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////

func PrintVersion(version gopi.Version) {
	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCells(true)
	table.Append([]string{
		"Name", version.Name(),
	})
	tag, branch, hash := version.Version()
	if tag != "" {
		table.Append([]string{
			"Tag", tag,
		})
	}
	if branch != "" {
		table.Append([]string{
			"Branch", branch,
		})
	}
	if hash != "" {
		table.Append([]string{
			"Hash", hash,
		})
	}
	table.Append([]string{
		"Go version", version.GoVersion(),
	})
	if t := version.BuildTime(); t.IsZero() == false {
		table.Append([]string{
			"Build time", t.Format(time.RFC3339),
		})
	}
	table.Render()
}

func (this *app) GetPingStub() (gopi.PingStub, error) {
	args := this.Args()
	if len(args) != 1 {
		return nil, gopi.ErrBadParameter.WithPrefix("Missing address")
	} else if conn, err := this.ConnPool.Connect("tcp", args[0]); err != nil {
		return nil, err
	} else if stub, _ := conn.NewStub("gopi.ping.Ping").(gopi.PingStub); stub == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("Cannot create stub")
	} else {
		return stub, nil
	}
}

func (this *app) GetInputStub() (gopi.InputStub, error) {
	args := this.Args()
	if len(args) != 1 {
		return nil, gopi.ErrBadParameter.WithPrefix("Missing address")
	} else if conn, err := this.ConnPool.Connect("tcp", args[0]); err != nil {
		return nil, err
	} else if stub, _ := conn.NewStub("gopi.input.Input").(gopi.InputStub); stub == nil {
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
