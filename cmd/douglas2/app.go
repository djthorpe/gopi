package main

import (
	"context"

	// Modules
	"github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.Logger
	*IT8951
}

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Define(cfg gopi.Config) error {
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.IT8951)
	return nil
}

func (this *app) Run(ctx context.Context) error {
	return nil
}
