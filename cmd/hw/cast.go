package main

import (
	"context"
	"fmt"
	"time"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *app) RunCast(ctx context.Context) error {
	if this.CastManager == nil {
		return gopi.ErrInternalAppError.WithPrefix("CastManager")
	}

	ctx2, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	if devices, err := this.CastManager.Devices(ctx2); err != nil {
		return err
	} else {
		for _, device := range devices {
			fmt.Println(device)
		}
	}

	// Return success
	return nil
}
