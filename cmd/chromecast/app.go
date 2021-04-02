package main

import (
	"context"
	"fmt"
	"time"

	// Modules
	"github.com/djthorpe/gopi/v3"

	// Dependencies
	_ "github.com/djthorpe/gopi/v3/pkg/dev/chromecast"
	_ "github.com/djthorpe/gopi/v3/pkg/event"
	_ "github.com/djthorpe/gopi/v3/pkg/log"
	_ "github.com/djthorpe/gopi/v3/pkg/mdns"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.CastManager
	gopi.Publisher
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Run(ctx context.Context) error {
	this.Require(this.Logger, this.CastManager)

	// Get devices
	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if devices, err := this.Devices(ctx2); err != nil {
		return err
	} else {
		fmt.Println(devices)
	}

	// Subscribe to chromecast events
	ch := this.Subscribe()
	defer this.Unsubscribe(ch)

	// Wait for interrupt, print events
	fmt.Println("Press CTRL+C to end")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt := <-ch:
			if evt, ok := evt.(gopi.CastEvent); ok {
				fmt.Println(evt)
			}
		}
	}
}
