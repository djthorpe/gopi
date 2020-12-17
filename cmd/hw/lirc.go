package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *app) RunLIRC(ctx context.Context, cfg gopi.Config) error {
	if this.LIRC == nil {
		return fmt.Errorf("LIRC is not enabled")
	}

	// Return sucess
	return nil
}
