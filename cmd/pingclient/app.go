package main

import (
	"context"
	"fmt"
	"time"

	"github.com/djthorpe/gopi/v3"
)

type app struct {
	gopi.Unit
	gopi.ConnPool

	conn gopi.Conn // Connection to server
}

func (this *app) New(cfg gopi.Config) error {
	if args := cfg.Args(); len(args) != 1 {
		return gopi.ErrBadParameter
	} else if conn, err := this.ConnPool.Connect("tcp", args[0]); err != nil {
		return err
	} else {
		this.conn = conn
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	if stub := this.conn.NewStub("gopi.ping.Ping").(gopi.PingStub); stub == nil {
		return gopi.ErrInternalAppError
	} else {
		timer := time.NewTicker(time.Second)
	FOR_LOOP:
		for {
			select {
			case <-timer.C:
				fmt.Println("ping")
				if err := stub.Ping(ctx); err != nil {
					fmt.Println(err)
					break FOR_LOOP
				}
			case <-ctx.Done():
				break FOR_LOOP
			}
		}
		timer.Stop()
	}

	// Return success
	return nil
}
