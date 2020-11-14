package main

import (
	"context"
	"fmt"

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

	if services, err := this.conn.ListServices(ctx); err != nil {
		return err
	} else {
		fmt.Println(this.conn)
		for i, service := range services {
			fmt.Println(i, service)
		}
	}

	// Return success
	return nil
}
