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

	if cmd := cfg.GetCommand(this.Args()); cmd == nil {
		return gopi.ErrHelp
	} else {
		fmt.Println("command=", cmd)
	}

	// Return sucess
	return nil
}
