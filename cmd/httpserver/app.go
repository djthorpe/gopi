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
	gopi.HttpStatic
	gopi.HttpLogger
	gopi.Publisher
	gopi.Logger
	gopi.Command
	gopi.MetricWriter
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("serve", "Serve static files", this.Serve)
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
	// Add static serving
	if folder, err := this.getFolderRoot(); err != nil {
		return err
	} else if err := this.HttpStatic.ServeFolder("/", folder); err != nil {
		return err
	} else if err := this.HttpLogger.Log("httpserver"); err != nil {
		return err
	}

	// Wait for interrupt, print out metrics
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
