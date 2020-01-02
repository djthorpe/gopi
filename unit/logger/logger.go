/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package logger

import (
	"fmt"
	"io"
	"strconv"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type Log struct {
	Writer  io.Writer
	Unit    string
	Debug   bool
	Verbose bool
}

type log struct {
	verbose bool

	base.Logger
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Log) Name() string { return "gopi.Logger" }

func (config Log) New(gopi.Logger) (gopi.Unit, error) {
	this := new(log)
	if err := this.Logger.Init(config.Writer, config.Unit, config.Debug); err != nil {
		return nil, err
	} else {
		this.verbose = config.Verbose
	}
	return this, nil
}

func (this *log) String() string {
	return fmt.Sprintf("<gopi.Logger name=%s debug=%v verbose=%v>", strconv.Quote(this.Name()), this.IsDebug(), this.verbose)
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Logger

func (this *log) Error(err error) error {
	this.Lock()
	defer this.Unlock()
	return this.Logger.Error(err)
}

func (this *log) Debug(args ...interface{}) {
	this.Lock()
	defer this.Unlock()
	this.Logger.Debug(args...)
}

func (this *log) Clone(name string) gopi.Logger {
	that := new(log)
	if base := this.Logger.Clone(name).(*base.Logger); base == nil {
		return nil
	} else {
		that.Logger = *base
		that.verbose = this.verbose
		return that
	}
}
