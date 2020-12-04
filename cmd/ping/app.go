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
	gopi.Command

	//stub gopi.PingStub // Connection to server
	//cmd    // Command to run
}

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("server", "Start ping service", this.RunServer)
	cfg.Command("ping", "Perform ping to server", this.RunPing)
	cfg.Command("version", "Display server version information", this.RunVersion)
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if this.Command := cfg.GetCommand(cfg.Args()); cmd == nil {
		return gopi.ErrHelp
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.cmd.Run(ctx)
}

/*


	else if len(cmd.Args()) != 1 {
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
*/
