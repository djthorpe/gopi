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
	gopi.CastService
	gopi.PingService
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Run(ctx context.Context) error {

	// Wait for interrupt
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()

	return ctx.Err()
}
