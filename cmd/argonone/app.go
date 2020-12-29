package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.Publisher

	// Control fan
	gopi.ArgonOne

	// Registered services
	gopi.PingService
	gopi.InputService

	// Emit LIRC codes
	gopi.LIRC
	gopi.LIRCKeycodeManager

	// Write metrics
	gopi.MetricWriter
}

func (this *app) Run(ctx context.Context) error {
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()
	return nil
}
