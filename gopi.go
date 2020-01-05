/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"log"
	"sync"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract driver interface
type Driver interface {
	// Close closes the driver and frees the underlying resources
	Close() error
}

// Abstract configuration which is used to open and return the
// concrete driver
type Config interface {
	// Opens the driver from configuration, or returns error
	Open(Logger) (Driver, error)
}

// Abstract logging interface
type Logger interface {
	Driver

	// Output logging messages
	Fatal(format string, v ...interface{}) error
	Error(format string, v ...interface{}) error
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Debug2(format string, v ...interface{})

	// Return IsDebug flag
	IsDebug() bool
}

// Concrete basic logger
type logger struct {
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Open a driver - opens the concrete version given the config method
func Open(config Config, log Logger) (Driver, error) {
	if log == nil {
		log = new(logger)
	}
	if driver, err := config.Open(log); err != nil {
		return nil, err
	} else {
		return driver, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *logger) Close() error {
	this.Lock()
	defer this.Unlock()
	return nil
}
func (this *logger) Fatal(format string, v ...interface{}) error {
	this.Lock()
	defer this.Unlock()
	err := fmt.Errorf(format, v...)
	log.Printf("Fatal: %v", err.Error())
	return err
}
func (this *logger) Error(format string, v ...interface{}) error {
	this.Lock()
	defer this.Unlock()
	err := fmt.Errorf(format, v...)
	log.Printf("Error: %v", err.Error())
	return err
}
func (this *logger) Warn(format string, v ...interface{}) {
	this.Lock()
	defer this.Unlock()
	log.Printf("Warn: %v", fmt.Sprintf(format, v...))
}
func (this *logger) Info(format string, v ...interface{}) {
	this.Lock()
	defer this.Unlock()
	log.Printf("Info: %v", fmt.Sprintf(format, v...))
}
func (this *logger) Debug(format string, v ...interface{}) {
	this.Lock()
	defer this.Unlock()
	log.Printf("Debug: %v", fmt.Sprintf(format, v...))
}
func (this *logger) Debug2(format string, v ...interface{}) {
	this.Lock()
	defer this.Unlock()
	log.Printf("Debug: %v", fmt.Sprintf(format, v...))
}
func (this *logger) IsDebug() bool {
	return true
}
