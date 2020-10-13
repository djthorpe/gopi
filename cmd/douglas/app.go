package main

import (
	"context"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/dev/waveshare"
)

type app struct {
	gopi.Unit
	gopi.Logger
	*waveshare.EPD
}

func (this *app) Run(ctx context.Context) error {
	this.Debug("epd=", this.EPD)
	return nil
}
