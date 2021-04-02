package main

import (
	"context"
	"fmt"
	"net"
	"strconv"

	gopi "github.com/djthorpe/gopi/v3"
)

func (this *app) RunServer(ctx context.Context, network, addr string) error {
	var port uint16

	// Start the server and discover the port served on
	if err := this.Server.StartInBackground(network, addr); err != nil {
		return err
	} else if _, port_, err := net.SplitHostPort(this.Server.Addr()); err != nil {
		return err
	} else if port__, err := strconv.ParseUint(port_, 0, 32); err != nil {
		return err
	} else {
		port = uint16(port__)
	}

	// Serve the ServiceRecord until parent context is done
	if r, err := this.ServiceDiscovery.NewServiceRecord("_gopi._tcp", "rpcping", port, nil, gopi.SERVICE_FLAG_IP4); err != nil {
		return err
	} else {
		child, cancel := context.WithCancel(ctx)
		defer cancel()
		go func(ctx context.Context) {
			if err := this.ServiceDiscovery.Serve(ctx, []gopi.ServiceRecord{r}); err != nil {
				this.Print(err)
			}
		}(child)
	}

	fmt.Println("Started server, ", this.Server)
	fmt.Println("Press CTRL+C to end")

	// Wait until done
	<-ctx.Done()

	// Close gracefully
	return this.Server.Stop(false)
}
