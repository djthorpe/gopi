package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.HttpStatic
	gopi.HttpLogger
	gopi.HttpTemplate
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Run(ctx context.Context) error {
	// Serve all folders under the current working directory under "/"
	if err := this.HttpStatic.ServeStatic("/"); err != nil {
		return err
	}

	// Wait for interrupt, print out metrics
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()

	// Return success
	return nil
}
