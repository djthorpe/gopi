package main

import (
	"context"
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.Publisher
	gopi.SurfaceManager
}

func (this *app) Run(ctx context.Context) error {
	if this.SurfaceManager == nil {
		return gopi.ErrInternalAppError.WithPrefix("SurfaceManager")
	}

	evts := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(evts)

	this.Print("Waiting for CTRL+C to exit")
	for {
		select {
		case evt := <-evts:
			fmt.Println(evt)
		case <-ctx.Done():
			return nil
		}
	}
}
