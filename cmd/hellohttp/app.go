package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.Server
	gopi.HttpStatic
	gopi.Logger
}

func (this *app) Run(ctx context.Context) error {
	if this.Server == nil {
		return gopi.ErrInternalAppError.WithPrefix("Server")
	}

	if err := this.Server.StartInBackground("tcp", ":0"); err != nil {
		return err
	} else if err := this.HttpStatic.ServeFolder("/", ""); err != nil {
		return err
	}

	fmt.Println("Started server, http://localhost" + this.Server.Addr() + "/")
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()

	// Return success
	return nil
}
