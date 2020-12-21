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
	gopi.DisplayManager
	gopi.SurfaceManager
	gopi.Display
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *app) Define(cfg gopi.Config) error {
	// Return success
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	if this.DisplayManager == nil {
		return gopi.ErrInternalAppError.WithPrefix("Invalid DisplayManager")
	}
	if this.SurfaceManager == nil {
		return gopi.ErrInternalAppError.WithPrefix("Invalid SurfaceManager")
	}

	if display := this.DisplayManager.PrimaryDisplay(); display == nil {
		return fmt.Errorf("No connected display")
	} else {
		this.Display = display
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	if bg, err := this.SurfaceManager.CreateBackground(this.Display, gopi.SURFACE_FLAG_BITMAP); err != nil {
		return err
	} else {
		this.Print(bg)
	}
	return nil
}
