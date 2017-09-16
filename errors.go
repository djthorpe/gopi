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
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	// ErrNoTasks is an error returned when no tasks are to be run
	ErrNoTasks = errors.New("No tasks to run")
	// ErrAppError is a general application error
	ErrAppError = errors.New("General application error")
	// ErrModuleNotFound is an error when module cannot be found by name or type
	ErrModuleNotFound = errors.New("Module not found")
)
