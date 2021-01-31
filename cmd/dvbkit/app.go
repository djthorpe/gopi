package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type app struct {
	gopi.Unit
	gopi.Logger
	gopi.DVBManager
	gopi.DVBTuner

	timeout *time.Duration
	tuner   *uint
	params  []gopi.DVBTunerParams
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *app) Define(cfg gopi.Config) error {
	this.timeout = cfg.FlagDuration("timeout", 2*time.Second, "Tune timeout")
	this.tuner = cfg.FlagUint("tuner", 0, "Tuner identifier")
	return nil
}

func (this *app) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.DVBManager)

	args := cfg.Args()
	if len(args) != 1 {
		return gopi.ErrBadParameter.WithPrefix("file")
	}

	// Set the tuner
	tuners := this.DVBManager.Tuners()
	for _, tuner := range tuners {
		if tuner.Id() == *this.tuner {
			this.DVBTuner = tuner
		}
	}
	if this.DVBTuner == nil {
		return gopi.ErrNotFound.WithPrefix("Tuner")
	}

	// Parse tuner params
	if fh, err := os.Open(args[0]); err != nil {
		return err
	} else if params, err := this.DVBManager.ParseTunerParams(fh); err != nil {
		return err
	} else {
		this.params = params
	}

	// Return success
	return nil
}

func (this *app) Run(ctx context.Context) error {
	// Tune all channels
	for _, param := range this.params {
		if ctx.Err() != nil {
			break
		}

		this.Print("Tune:", this.DVBTuner, param.Name())
		tunectx, cancel := context.WithTimeout(ctx, *this.timeout)
		defer cancel()

		if err := this.Tune(tunectx, this.DVBTuner, param, this.Tuned); err != nil {
			this.Print("  Error:", err)
		} else {
			// Wait until tuning timeout
			<-tunectx.Done()
		}
	}

	// Wait for interrupt
	//fmt.Println("Press CTRL+C to end")
	//<-ctx.Done()

	// Return success
	return nil
}

func (this *app) Tuned(ctx gopi.DVBContext) {
	fmt.Println("Tuned=", ctx)
}
