/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi // import "github.com/djthorpe/gopi"

import (
	"errors"
	"flag"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// ErrNoTasks is an error returned when no tasks are to be run
	ErrNoTasks = errors.New("No tasks to run")
	// ErrAppError is a general application error
	ErrAppError = errors.New("General application error")
	// ErrHelp is returned when -help was called for on the command line
	ErrHelp = flag.ErrHelp
	// ErrNotImplemented is returned when a feature is not supported
	ErrNotImplemented = errors.New("Feature not implemented")
	// ErrBadParameter is returned when a supplied parameter is invalid
	ErrBadParameter = errors.New("Bad Parameter")
)
