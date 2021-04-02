package main

import (
	"context"
	"fmt"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"

	// Dependencies
	_ "github.com/djthorpe/gopi/v3/pkg/rpc/chromecast"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Chromecast struct {
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Chromecast) Define(cfg gopi.Config) {
	cfg.Command("cast", "Watch for Chromecast Events", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		ch := make(chan gopi.CastEvent)
		go func() {
			fmt.Println("Watching for events, press CTRL+C to end")
			for evt := range ch {
				fmt.Println(evt)
			}
		}()
		stub.Stream(ctx, ch)
		close(ch)
		return nil
	})
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Chromecast) GetStub(ctx context.Context) gopi.CastStub {
	return ctx.Value(KeyStub).(gopi.CastStub)
}
