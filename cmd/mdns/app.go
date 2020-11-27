package main

import (
	"context"
	"fmt"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/mdns"
)

type app struct {
	gopi.Unit
	*mdns.Listener
	*mdns.Discovery
	gopi.Publisher
	gopi.Command
}

func (this *app) Define(cfg gopi.Config) error {
	cfg.Command("listen", "Listen for mDNS messages", this.Listen)
	cfg.Command("discovery", "Discover mDNS services", this.Discover)
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if this.Command = cfg.GetCommand(nil); this.Command == nil {
		return gopi.ErrHelp
	}
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return this.Command.Run(ctx)
}

func (this *app) Listen(ctx context.Context) error {
	fmt.Println("Waiting for CTRL+C")
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)
	for {
		select {
		case <-ctx.Done():
			return nil
		case evt := <-ch:
			fmt.Println(evt)
		}
	}
}

func (this *app) Discover(ctx context.Context) error {
	other, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return this.Discovery.EnumerateServices(other)
}
