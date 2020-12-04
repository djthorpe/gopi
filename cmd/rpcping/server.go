package main

import (
	"context"
	"fmt"
)

func (this *app) RunServer(ctx context.Context, network, addr string) error {
	if err := this.Server.StartInBackground(network, addr); err != nil {
		return err
	}

	fmt.Println("Started server, ", this.Server)
	fmt.Println("Press CTRL+C to end")

	// Wait until done
	<-ctx.Done()

	// Close gracefully
	return this.Server.Stop(false)
}
