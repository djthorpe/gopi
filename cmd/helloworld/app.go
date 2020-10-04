package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/log"
)

type app struct {
	gopi.Unit
	*log.Log

	name *string
}

func (this *app) Define(cfg gopi.Config) error {
	this.name = cfg.FlagString("name", "", "Your name")
	return nil
}

func (this *app) Run(ctx context.Context) error {
	if *this.name != "" {
		this.Print("Hello, ", *this.name)
	} else {
		this.Print("Hello, World!")
	}
	return nil
}

func (this *app) String() string {
	str := "<app"
	str += " log=" + fmt.Sprint(this.Log)
	return str + ">"
}
