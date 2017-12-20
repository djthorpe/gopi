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
	Fatal(format string, v ...interface{}) Error
	Error(format string, v ...interface{}) Error
	Warn(format string, v ...interface{})
	Info(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Debug2(format string, v ...interface{})

	// Return IsDebug flag
	IsDebug() bool
}

// Error type
type Error struct {
	reason string
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

// Open a driver - opens the concrete version given the config method
// and only returns the driver (or nil). Will return an error as a
// reference.
func Open2(config Config, log Logger, error_ref *Error) Driver {
	var err error
	var driver Driver

	// Create driver
	if err == nil {
		driver, err = config.Open(log)
	}

	// Return error
	if err != nil {
		if error_ref != nil {
			*error_ref = NewError(err)
		}
		return nil
	}

	// Return success
	return driver
}

////////////////////////////////////////////////////////////////////////////////
// ERROR IMPLEMENTATION

// Create a gopi.Error object
func NewError(err error) Error {
	return Error{reason: err.Error()}
}

func (err Error) Error() string {
	return err.reason
}
