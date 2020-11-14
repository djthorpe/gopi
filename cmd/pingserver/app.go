package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.Server
	gopi.PingService
}

func (this *app) Run(ctx context.Context) error {
	if err := this.Server.StartInBackground("tcp", ":0"); err != nil {
		return err
	}

	fmt.Println("Started Ping Server, ", this.Server)
	fmt.Println("Press CTRL+C to end")

	// Wait until done
	<-ctx.Done()

	// Close gracefully
	return this.Server.Stop(false)
}
