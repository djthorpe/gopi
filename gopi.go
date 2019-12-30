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
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Unit configuration interface
type Config interface {
	Name() string             // Returns name of the unit
	New(Logger) (Unit, error) // Opens the driver from configuration, or returns error
}

// Unit interface
type Unit interface {
	Close() error   // Close closes the driver and frees the underlying resources
	String() string // String returns a string representation of the unit
}

// Abstract logging interface
type Logger interface {
	Unit

	Name() string              // Return unit name
	Error(error) error         // Output logging messages
	Debug(args ...interface{}) // Debug output
	IsDebug() bool             // Return IsDebug flag
}

// UnitBase is the struct for any unit
type UnitBase struct {
	Log Logger
}

// LoggerBase is the struct for any logger
type LoggerBase struct {
	UnitBase
	name   string
	writer io.Writer
	debug  bool
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Config

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
// IMPLEMENTATION gopi.Unit

func (this *UnitBase) Init(log Logger) error {
	this.Log = log
	return nil
}

func (this *UnitBase) Close() error {
	this.Log = nil
	return nil
}

func (this *UnitBase) String() string {
	if this.Log != nil {
		return "<" + this.Log.Name() + ">"
	} else {
		return "<gopi.UnitBase>"
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Logger

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

func (this *LoggerBase) Debug(args ...interface{}) {
	if this.debug {
		fmt.Fprintln(this.writer, args...)
	}
}

func (this *LoggerBase) IsDebug() bool {
	return this.debug
}

func (this *LoggerBase) Name() string {
	return this.name
}
