/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

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

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Open a driver - opens the concrete version given the config method
func Open(config Config, log Logger) (Driver, error) {
	if driver, err := config.Open(log); err != nil {
		return nil, err
	} else {
		return driver, nil
	}
}
