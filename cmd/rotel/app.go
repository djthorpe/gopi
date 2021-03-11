package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.Logger

	// Registered services
	gopi.PingService
	gopi.RotelService
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.PingService, this.RotelService)
	return nil
}

func (this *app) Run(ctx context.Context) error {
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()
	return nil
}
