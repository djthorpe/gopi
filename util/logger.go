/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This package provides utility functions such as logging
package util /* import "github.com/djthorpe/gopi/util" */

import (
	"fmt"
	"os"
	"errors"
)

////////////////////////////////////////////////////////////////////////////////

// The level of the logging required
type LogLevel uint

type LoggerDevice struct {
	driver LoggerInterface
	level LogLevel
}

// Abstract configuration which is used to open and return the
// concrete logging device
type LogConfig interface {
	// Opens the logger from configuration, or returns error
	Open() (LoggerInterface, error)
}

// Abstract interface to the logger, to perform the logging
type LoggerInterface interface {

	// Perform logging
	Log(level LogLevel,message string)

	// Close the logging device
	Close() error
}

// Concrete Stderr Logger Configuration
type StderrLogger struct {

}

// Concrete Stderr Logger Device
type StderrLoggerDevice struct {
	device *os.File
}

////////////////////////////////////////////////////////////////////////////////

const (
	LOG_ANY LogLevel = iota
	LOG_DEBUG2
	LOG_DEBUG
	LOG_INFO
	LOG_WARN
	LOG_ERROR
	LOG_FATAL
)

////////////////////////////////////////////////////////////////////////////////
// Opener interface

// Open opens a connection EGL
func Logger(config LogConfig) (*LoggerDevice, error) {
	logger, err := config.Open()
	if err != nil {
		return nil, err
	}
	return &LoggerDevice{ driver: logger, level: LOG_INFO }, nil
}

////////////////////////////////////////////////////////////////////////////////
// Driver interface

// Closes the device and frees the resources
func (this *LoggerDevice) Close() error {
	return this.driver.Close()
}

////////////////////////////////////////////////////////////////////////////////
// Display log level

func (l LogLevel) String() string {
	switch(l) {
	case LOG_DEBUG2:
		return "DEBUG"
	case LOG_DEBUG:
		return "DEBUG"
	case LOG_INFO:
		return "INFO"
	case LOG_WARN:
		return "WARN"
	case LOG_ERROR:
		return "ERROR"
	case LOG_FATAL:
		return "FATAL"
	default:
		return "OTHER"
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get and set level

func (this* LoggerDevice) GetLevel() LogLevel {
	return this.level
}

func (this* LoggerDevice) SetLevel(level LogLevel) {
	this.level = level
}

////////////////////////////////////////////////////////////////////////////////
// Methods to print out log information

func (this* LoggerDevice) Info(format string,v... interface{}) {
	message := fmt.Sprintf(format,v...)
	if this.level <= LOG_INFO || this.level==LOG_ANY {
		this.driver.Log(LOG_INFO,message)
	}
}

func (this* LoggerDevice) Debug(format string,v... interface{}) {
	message := fmt.Sprintf(format,v...)
	if this.level <= LOG_DEBUG || this.level==LOG_ANY {
		this.driver.Log(LOG_DEBUG,message)
	}
}

func (this* LoggerDevice) Debug2(format string,v... interface{}) {
	message := fmt.Sprintf(format,v...)
	if this.level <= LOG_DEBUG2 || this.level==LOG_ANY {
		this.driver.Log(LOG_DEBUG2,message)
	}
}

func (this* LoggerDevice) Warn(format string,v... interface{}) {
	message := fmt.Sprintf(format,v...)
	if this.level <= LOG_WARN || this.level==LOG_ANY {
		this.driver.Log(LOG_WARN,message)
	}
}

func (this* LoggerDevice) Error(format string,v... interface{}) error {
	message := fmt.Sprintf(format,v...)
	if this.level <= LOG_ERROR || this.level==LOG_ANY {
		this.driver.Log(LOG_ERROR,message)
	}
	return errors.New(message)
}

func (this* LoggerDevice) Fatal(format string,v... interface{})  error {
	message := fmt.Sprintf(format,v...)
	if this.level <= LOG_FATAL || this.level==LOG_ANY {
		this.driver.Log(LOG_FATAL,message)
	}
	return errors.New(message)
}

////////////////////////////////////////////////////////////////////////////////
// Stderr Open and Close logger

// Initialise the EGL interface
func (config StderrLogger) Open() (LoggerInterface, error) {
	this := new(StderrLoggerDevice)
	this.device = os.Stderr
	return this, nil
}

// Close logger
func (this *StderrLoggerDevice) Close() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Output

func (this *StderrLoggerDevice) Log(level LogLevel,message string) {
	fmt.Fprintf(this.device,"[%v] %v\n",level,message)
}





