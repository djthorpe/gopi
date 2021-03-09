package main

import (
	"context"
	"fmt"
	"image/color"

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
	this.SurfaceManager.Do(func(ctx gopi.GraphicsContext) error {
		if surface, err := this.SurfaceManager.CreateSurface(ctx, 0, 1.0, 100, gopi.Point{500, 500}, gopi.Size{100, 100}); err != nil {
			return err
		} else {
			bitmap := surface.Bitmap()
			bitmap.ClearToColor(color.Gray{0x80})
			this.Print(surface)
		}
		return nil
	})

	// Wait for interrupt
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()

	// Return success
	return nil
}
