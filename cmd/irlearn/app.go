package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/lirc"
	"github.com/olekukonko/tablewriter"
)

type app struct {
	gopi.Unit
	gopi.Command
	gopi.Logger
	gopi.Publisher
	gopi.LIRC
	gopi.LIRCKeycodeManager

	name *string
}

func (this *app) Define(cfg gopi.Config) error {
	// Commands
	cfg.Command("recv", "Receive keycodes", this.Recv)
	cfg.Command("keycodes", "Lookup keycodes", this.Keycodes)
	cfg.Command("learn", "Learn keycodes", this.Learn)

	// Flags
	this.name = cfg.FlagString("name", "default", "Database name")

	// Return success
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

func (this *app) Keycodes(ctx context.Context) error {
	if len(this.Command.Args()) == 0 {
		return fmt.Errorf("Missing argument")
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Keycodes"})
	for _, name := range this.Command.Args() {
		codes := this.LIRCKeycodeManager.Keycode(name)
		if len(codes) == 0 {
			table.Append([]string{name, "nil"})
		} else {
			table.Append([]string{name, FormatCodes(codes)})
		}
	}
	table.Render()
	return nil
}

func (this *app) Learn(ctx context.Context) error {
	if len(this.Command.Args()) == 0 {
		return fmt.Errorf("Missing argument")
	}
	for _, name := range this.Command.Args() {
		codes := this.LIRCKeycodeManager.Keycode(name)
		if len(codes) == 0 {
			fmt.Println("No keys found for", name)
			continue
		}
		for i, code := range codes {
			fmt.Println(i+1, "of", len(codes))
			if err := this.LearnCode(ctx, code, *this.name); err != nil {
				return err
			}
			fmt.Println()
		}
	}
	return nil
}

func (this *app) LearnCode(ctx context.Context, code gopi.KeyCode, name string) error {
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	timer := time.NewTimer(5 * time.Second)
	fmt.Println("Learning", code, "for", strconv.Quote(name))
	fmt.Println("  ...press the key on the remote corresponding to", code, "or CTRL+C to end")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		case evt := <-ch:
			fmt.Println(evt)
		}
	}
}

func (this *app) Recv(ctx context.Context) error {
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	fmt.Println("Receiving IR Events, press CTRL+C to exit")

	for {
		select {
		case <-ctx.Done():
			return nil
		case evt := <-ch:
			if inputevent, ok := evt.(gopi.InputEvent); ok {
				fmt.Println("inputevent=", inputevent)
			}
		}
	}
}
