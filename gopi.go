/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"io"
	"os"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Unit configuration interface
type Config interface {
	// Returns name of the unit
	Name() string

	// Opens the driver from configuration, or returns error
	New(Logger) (Unit, error)
}

// Unit interface
type Unit interface {
	// Close closes the driver and frees the underlying resources
	Close() error

	// String returns a string representation of the unit
	String() string
}

// Abstract logging interface
type Logger interface {
	Unit

	// Return unit name
	Name() string

	// Output logging messages
	Error(error) error

	// Return IsDebug flag
	IsDebug() bool
}

// Base struct for any unit
type UnitBase struct {
	log Logger
}

// Base struct for any logger
type LoggerBase struct {
	UnitBase
	name   string
	writer io.Writer
	debug  bool
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func New(config Config, log Logger) (Unit, error) {
	if log == nil {
		log_ := new(LoggerBase)
		if err := log_.Init(os.Stderr, config.Name(), false); err != nil {
			return nil, err
		} else {
			log = log_
		}
	}
	if driver, err := config.New(log); err != nil {
		return nil, err
	} else {
		return driver, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.UnitBase

func (this *UnitBase) Init(log Logger) error {
	this.log = log
	return nil
}

func (this *UnitBase) Close() error {
	return gopi.ErrNotImplemented
}

func (this *UnitBase) String() string {
	if this.log != nil {
		return "<" + this.log.Name() + ">"
	} else {
		return "<gopi.UnitBase>"
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.LoggerBase

func (this *LoggerBase) Init(writer io.Writer, name string, debug bool) error {
	if name == "" || writer == nil {
		return ErrBadParameter
	} else {
		this.writer = writer
		this.name = name
		this.debug = debug
		return nil
	}
}

func (this *LoggerBase) Error(err error) error {
	if this.name != "" {
		err = fmt.Errorf("%s: %w", this.name, err)
	}
	fmt.Fprintln(this.writer, err)
	return err
}

func (this *LoggerBase) IsDebug() bool {
	return this.debug
}

func (this *LoggerBase) Name() string {
	return this.name
}
