package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/olekukonko/tablewriter"
)

type app struct {
	gopi.Unit
	gopi.ConnPool

	stub gopi.PingStub // Connection to server
	cmd  gopi.Command  // Command to run
}

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("ping", "Perform ping to server", this.Ping)
	cfg.Command("version", "Display server version information", this.Version)
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if cmd := cfg.GetCommand(cfg.Args()); cmd == nil {
		return gopi.ErrBadParameter
	} else if len(cmd.Args()) != 1 {
		return gopi.ErrBadParameter.WithPrefix("Missing server address")
	} else if conn, err := this.ConnPool.Connect("tcp", cmd.Args()[0]); err != nil {
		return err
	} else if stub, _ := conn.NewStub("gopi.ping.Ping").(gopi.PingStub); stub == nil {
		return gopi.ErrInternalAppError.WithPrefix("Cannot create stub")
	} else {
		this.stub = stub
		this.cmd = cmd
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.cmd.Run(ctx)
}

func (this *app) Ping(ctx context.Context) error {
	timer := time.NewTicker(time.Second)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			fmt.Println("ping")
			if err := this.stub.Ping(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return nil
		}
	}

	// Return success
	return nil
}

func (this *app) Version(ctx context.Context) error {
	version, err := this.stub.Version(ctx)
	if err != nil {
		return err
	}

	// Display platform information
	table := tablewriter.NewWriter(os.Stdout)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
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

	// Return success
	return nil
}
