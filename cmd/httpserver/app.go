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
	gopi.HttpStatic
	gopi.HttpLogger
	gopi.HttpTemplate
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Run(ctx context.Context) error {
	if err := this.HttpStatic.ServeStatic("/"); err != nil {
		return err
	} else if err := this.HttpTemplate.ServeTemplate("/", "page.tmpl"); err != nil {
		return err
	}

	// Wait for interrupt, print out metrics
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS
/*
func (this *app) getFolderRoot() (string, error) {
	args := this.Args()
	if len(args) != 1 {
		return "", gopi.ErrBadParameter.WithPrefix("Missing folder")
	} else if stat, err := os.Stat(args[0]); err != nil {
		return "", err
	} else if stat.IsDir() == false {
		return "", gopi.ErrBadParameter.WithPrefix("Missing folder")
	} else {
		return filepath.Clean(args[0]), nil
	}
}
*/
