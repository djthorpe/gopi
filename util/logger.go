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
	"strings"
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

// Concrete StderrLogger Configuration
type StderrLogger struct {

}

// Concrete NullLogger Configuration
type NullLogger struct {

}

// Concrete FileLogger Configuration
type FileLogger struct {
	// File filename to write the log to. File will be created if
	// it doesn't already exist
	Filename string

	// Whether to append to the existing file if it already exists
	Append bool
}

// Concrete Logger Device
type logger struct {
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
	LOG_NONE
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
		return "[Invalid LogLevel value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get and set level

// Get logging level
func (this* LoggerDevice) GetLevel() LogLevel {
	return this.level
}

// Set logging level from a LogLevel parameter
func (this* LoggerDevice) SetLevel(level LogLevel) {
	this.level = level
}

// Set logging level from a string parameter. Will return
// an error if the string could not be parsed
func (this* LoggerDevice) SetLevelFromString(level string) error {
	switch(strings.ToLower(strings.TrimSpace(level))) {
	case "debug2":
		this.SetLevel(LOG_DEBUG2)
		return nil
	case "debug":
		this.SetLevel(LOG_DEBUG)
		return nil
	case "info":
		this.SetLevel(LOG_INFO)
		return nil
	case "warn", "warning":
		this.SetLevel(LOG_WARN)
		return nil
	case "error", "err":
		this.SetLevel(LOG_ERROR)
		return nil
	case "fatal":
		this.SetLevel(LOG_FATAL)
		return nil
	case "any":
		this.SetLevel(LOG_ANY)
		return nil
	case "none":
		this.SetLevel(LOG_NONE)
		return nil
	}
	return errors.New("Invalid level")
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
// Open Logger

// Initialise the StderrLogger
func (config StderrLogger) Open() (LoggerInterface, error) {
	this := new(logger)
	this.device = os.Stderr
	return this, nil
}

// Initialise the NullLogger
func (config NullLogger) Open() (LoggerInterface, error) {
	return new(logger), nil
}

// Initialise the FileLogger
func (config FileLogger) Open() (LoggerInterface, error) {
	var err error

	this := new(logger)
	flag := os.O_RDWR | os.O_CREATE
	if config.Append {
		flag |= os.O_APPEND
	}
	this.device, err = os.OpenFile(config.Filename,flag,0666)
	if err != nil {
		return nil, err
	}
	return this, nil
}

// Close logger
func (this *logger) Close() error {
	if this.device != nil {
		return this.device.Close()
	}
	return nil
}

// Output log message to device
func (this *logger) Log(level LogLevel,message string) {
	if this.device != nil {
		fmt.Fprintf(this.device,"[%v] %v\n",level,message)
	}
}

