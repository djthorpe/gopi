package main

import (
	"context"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////

func (this *app) Remux(ctx context.Context) error {
	args := this.Args()
	if len(args) != 2 {
		return gopi.ErrHelp
	}

	src, err := this.MediaManager.OpenFile(args[0])
	if err != nil {
		return err
	}
	dst, err := this.MediaManager.CreateFile(args[1])
	if err != nil {
		return err
	}

	fmt.Println(src, "=>", dst)

	// Return success
	return nil
}
