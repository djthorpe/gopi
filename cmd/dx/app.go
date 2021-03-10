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
	gopi.Logger
	gopi.SurfaceManager
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.SurfaceManager)

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	if err := this.SurfaceManager.Do(func(ctx gopi.GraphicsContext) error {
		if surface, err := this.SurfaceManager.CreateSurface(ctx, gopi.SURFACE_FLAG_OPENVG, 1.0, 100, gopi.Point{500, 500}, gopi.Size{100, 100}); err != nil {
			return err
		} else {
			this.Print(surface)
		}
		return nil
	}); err != nil {
		return err
	}

	// Wait for interrupt
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()

	// Return success
	return nil
}
