/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2019
	All Rights Reserved
    Documentation https://gopi.mutablelogic.com/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

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
	// ErrUnexpectedResponse is returned when some response wasn't expected
	ErrUnexpectedResponse = errors.New("Unexpected response")
	// ErrNotFound is returned when something isn't found
	ErrNotFound = errors.New("Not found")
	// ErrOutOfOrder is returned when an operation was executed out of order
	ErrOutOfOrder = errors.New("Operation out of order")
	// ErrDeadlineExceeded is returned when a timeout occurred
	ErrDeadlineExceeded = errors.New("Deadline exceeded")
	// ErrNotModified is returned when a resource is not modified
	ErrNotModified = errors.New("Not modified")
)
