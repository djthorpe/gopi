/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"fmt"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	// Error represents a gopi error
	Error uint

	// CompoundError represents a set of errors
	CompoundError struct {
		errs []error
	}
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ErrNone             Error = iota // No error condition
	ErrNotImplemented                // Method or feature not implemented
	ErrBadParameter                  // Error with parameter passed to method
	ErrNotFound                      // Missing object
	ErrHelp                          // Help requested from command line
	ErrInternalAppError              // Internal application error
	ErrSignalCaught                  // Signal caught
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY AND WRAP

func (this Error) Error() string {
	switch this {
	case ErrNotImplemented:
		return "Not Implemented"
	case ErrBadParameter:
		return "Bad Parameter"
	case ErrNotFound:
		return "Not Found"
	case ErrHelp:
		return "Help Requested"
	case ErrInternalAppError:
		return "Internal Application Error"
	case ErrSignalCaught:
		return "Signal caught"
	default:
		return "[?? Invalid Error value]"
	}
}

func (this Error) WithPrefix(prefix string) error {
	return fmt.Errorf("%s: %w", prefix, this)
}

func (this *CompoundError) Add(err error) error {
	if err == nil {
		return this
	}
	if this.errs == nil {
		this.errs = make([]error, 0, 1)
	}
	this.errs = append(this.errs, err)
	return this
}

func (this *CompoundError) Error() string {
	if len(this.errs) == 0 {
		return ErrNone.Error()
	} else if len(this.errs) == 1 {
		return this.errs[0].Error()
	} else {
		str := ""
		for _, err := range this.errs {
			str += err.Error() + ","
		}
		return strings.TrimSuffix(str, ",")
	}
}

func (this *CompoundError) ErrorOrSelf() error {
	if len(this.errs) == 0 {
		return nil
	} else if len(this.errs) == 1 {
		return this.errs[0]
	} else {
		return this
	}
}
