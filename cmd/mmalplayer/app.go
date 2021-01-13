package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	mmal "github.com/djthorpe/gopi/v3/pkg/media/mmal"
	multierror "github.com/hashicorp/go-multierror"
)

type app struct {
	gopi.Unit
	gopi.Logger
	*mmal.Manager

	r       io.ReadCloser
	w       io.Writer
	decoder mmal.ImageComponent
}

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.Manager)

	args := cfg.Args()
	if len(args) != 1 {
		return errors.New("Require input filename")
	}
	if r, err := os.Open(args[0]); err != nil {
		return err
	} else {
		this.r = r
	}

	// Create output buffer
	this.w = new(bytes.Buffer)

	// Create decoder
	if decoder, err := this.Manager.ImageDecoder(); err != nil {
		return err
	} else {
		this.decoder = decoder
	}

	// Set format
	if err := this.decoder.SetInputFormatJPEG(); err != nil {
		return err
	}

	// Reader -> Decoder -> Writer
	if _, err := this.Manager.CreateReaderForComponent(this.r, this.decoder, 0); err != nil {
		return err
	} else if _, err := this.Manager.CreateWriterForComponent(this.w, this.decoder, 0); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *app) Dispose() error {
	var result error

	// Close reader
	if err := this.r.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	// Release resources
	this.decoder = nil
	this.r = nil
	this.w = nil

	return result

}

func (this *app) Run(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := this.Manager.Exec(ctx); err != nil {
		return err
	} else {
		return nil
	}
}
