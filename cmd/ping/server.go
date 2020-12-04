package main

import (
	"context"
	"fmt"
	"net"

	"github.com/djthorpe/gopi/v3"
)

func (this *app) RunServer(ctx context.Context) error {
	var network, addr string

	args := this.Command.Args()
	switch {
	case len(args) == 0:
		network = "tcp"
		addr = ":0"
	case len(args) == 1:
		if _, _, err := net.SplitHostPort(args[0]); err != nil {
			return err
		} else {
			network = "tcp"
			addr = args[0]
		}
	default:
		return gopi.ErrBadParameter.WithPrefix(this.Command.Name())
	}

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
