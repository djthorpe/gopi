package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.Server
	gopi.HttpStatic
	gopi.Logger
	gopi.Command

	addr *string
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("serve", "Serve static files", this.Serve)
	this.addr = cfg.FlagString("addr", ":0", "Address for server", "serve")
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
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

////////////////////////////////////////////////////////////////////////////////
// COMMANDS

func (this *app) Serve(ctx context.Context) error {
	// Check parameters
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("Server")
	}

	// Start server and add services
	if folder, err := this.getFolderRoot(); err != nil {
		return err
	} else if addr, err := this.getAddr(); err != nil {
		return err
	} else if err := this.Server.StartInBackground("tcp", addr); err != nil {
		return err
	} else if err := this.HttpStatic.ServeFolder("/", folder); err != nil {
		return err
	}

	// Wait for interrupt
	fmt.Println("Started server, http://localhost" + this.Server.Addr() + "/")
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *app) getFolderRoot() (string, error) {
	args := this.Args()
	if len(args) != 1 {
		return "", gopi.ErrBadParameter.WithPrefix("Missing folder")
	} else if stat, err := os.Stat(args[0]); err != nil {
		return "", err
	} else if stat.IsDir() == false {
		return "", gopi.ErrBadParameter.WithPrefix("Missing folder")
	} else {
		return filepath.Clean(args[0]), nil
	}
}

func (this *app) getAddr() (string, error) {
	return *this.addr, nil
}
