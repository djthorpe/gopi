package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.Logger

	name *string
	wait *bool
}

func (this *app) Define(cfg gopi.Config) error {
	this.name = cfg.FlagString("name", "", "Your name")
	this.wait = cfg.FlagBool("wait", false, "Wait for CTRL+C to exit")
	return nil
}

func (this *app) Run(ctx context.Context) error {
	if *this.name != "" {
		this.Print("Hello, ", *this.name)
	} else {
		this.Print("Hello, World!")
	}

	if *this.wait {
		this.Print("Waiting for CTRL+C to exit")
		<-ctx.Done()
	}

	return nil
}

func (this *app) String() string {
	str := "<app"
	str += " log=" + fmt.Sprint(this.Logger)
	return str + ">"
}
