/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package base

import (
	"fmt"
	"io"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Logger is the base struct for any logger
type Logger struct {
	name   string
	writer io.Writer
	debug  bool
	Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Logger

func (this *Logger) Init(writer io.Writer, name string, debug bool) error {
	if name == "" {
		return gopi.ErrBadParameter.WithPrefix("name")
	} else if writer == nil {
		return gopi.ErrBadParameter.WithPrefix("writer")
	} else if err := this.Unit.Init(nil); err != nil {
		return err
	} else {
		this.writer = writer
		this.name = name
		this.debug = debug
		return nil
	}
}

func (this *Logger) Error(err error) error {
	if this.name != "" {
		err = fmt.Errorf("%s: %w", this.name, err)
	}
	fmt.Fprintln(this.writer, err)
	return err
}

func (this *Logger) Debug(args ...interface{}) {
	if this.debug {
		fmt.Fprintln(this.writer, args...)
	}
}

func (this *Logger) IsDebug() bool {
	return this.debug
}

func (this *Logger) Name() string {
	return this.name
}
