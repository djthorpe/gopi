package main

import (
	"context"
	"fmt"
	"os"
	"time"

	// Modules

	"github.com/djthorpe/data"
	table "github.com/djthorpe/data/pkg/table"
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

	cfg.Command("cast list", "Return list of Chromecasts", func(ctx context.Context) error {
		stub := this.GetStub(ctx)
		casts, err := stub.List(ctx, time.Second)
		if err != nil {
			return err
		}
		table := table.NewTable("Id", "Name", "Model", "Service", "State")
		for _, cast := range casts {
			table.Append(cast.Id(), cast.Name(), cast.Model(), cast.Service(), cast.State())
		}
		return table.Write(os.Stdout, table.OptHeader(), table.OptAscii(80, data.BorderLines))
	})
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Chromecast) GetStub(ctx context.Context) gopi.CastStub {
	return ctx.Value(KeyStub).(gopi.CastStub)
}
