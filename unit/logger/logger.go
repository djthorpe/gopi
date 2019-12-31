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
)

type Log struct {
	Writer  io.Writer
	Unit    string
	Debug   bool
	Verbose bool
}

type log struct {
	gopi.LoggerBase
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Log) Name() string { return "gopi.Logger" }

func (config Log) New(gopi.Logger) (gopi.Unit, error) {
	this := new(log)
	if err := this.LoggerBase.Init(config.Writer, config.Unit, config.Debug); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *log) String() string {
	return fmt.Sprintf("<gopi.Logger unit=%s debug=%v>", strconv.Quote(this.Name()), this.IsDebug())
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Logger

func (this *log) Error(err error) error {
	this.Lock()
	defer this.Unlock()
	return this.LoggerBase.Error(err)
}

func (this *log) Debug(args ...interface{}) {
	this.Lock()
	defer this.Unlock()
	this.LoggerBase.Debug(args...)
}
