package main

import (
	"context"
	"encoding/csv"
	"io"
	"os"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.Metrics
	gopi.Publisher
	gopi.Command
}

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("scan", "Scan a CSV file and display definition", this.Scan)
	cfg.Command("dump", "Display CSV file", this.Dump)
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

func (this *app) Scan(ctx context.Context) error {
	// Check argument
	args := this.Command.Args()
	if len(args) != 1 {
		return gopi.ErrBadParameter.WithPrefix(this.Command.Name())
	} else if stat, err := os.Stat(args[0]); os.IsNotExist(err) {
		return err
	} else if stat.Mode().IsRegular() == false {
		return gopi.ErrBadParameter.WithPrefix(args[0])
	} else if err != nil {
		return err
	}

	// Create a table for scanning
	table := NewTable(true)

	// Open file for CSV parsing
	fh, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fh.Close()
	csv := csv.NewReader(fh)
	for {
		record, err := csv.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		table.Scan(record)
	}

	// Write out the table
	table.Schema().Write(os.Stdout)

	// Return success
	return nil
}

func (this *app) Dump(ctx context.Context) error {
	// Check argument
	args := this.Command.Args()
	if len(args) != 1 {
		return gopi.ErrBadParameter.WithPrefix(this.Command.Name())
	} else if stat, err := os.Stat(args[0]); os.IsNotExist(err) {
		return err
	} else if stat.Mode().IsRegular() == false {
		return gopi.ErrBadParameter.WithPrefix(args[0])
	} else if err != nil {
		return err
	}

	// Create a table for scanning
	table := NewTable(true)

	// Open file for CSV parsing
	fh, err := os.Open(args[0])
	if err != nil {
		return err
	}
	defer fh.Close()
	csv := csv.NewReader(fh)
	for {
		record, err := csv.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		table.Append(record)
	}

	// Write out the table
	table.Write(os.Stdout)

	// Return success
	return nil
}
