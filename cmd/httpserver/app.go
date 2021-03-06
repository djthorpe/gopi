package main

import (
	"context"
	"fmt"
	"os"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/http/renderer"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.HttpTemplate

	// Renderers
	*renderer.HttpIndexRenderer
	*renderer.HttpTextRenderer

	// Document Root
	docroot string
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.HttpTemplate, this.HttpIndexRenderer, this.HttpTextRenderer)

	if docroot, err := docRoot(cfg.Args()); err != nil {
		return err
	} else {
		this.docroot = docroot
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {

	// Serve templates under "/"
	if err := this.HttpTemplate.Serve("/", this.docroot); err != nil {
		return err
	}

	// Wait for interrupt, print out metrics
	fmt.Println("Press CTRL+C to end")
	<-ctx.Done()

	// Return success
	return nil
}

func docRoot(args []string) (string, error) {
	var docroot string
	if len(args) == 0 {
		if wd, err := os.Getwd(); err != nil {
			return "", err
		} else {
			docroot = wd
		}
	} else if len(args) == 1 {
		docroot = args[0]
	} else {
		return "", gopi.ErrBadParameter.WithPrefix("Too many arguments")
	}

	return docroot, nil
}
