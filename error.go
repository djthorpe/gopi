/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"errors"
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
	ErrNone               Error = iota // No error condition
	ErrNotImplemented                  // Method or feature not implemented
	ErrBadParameter                    // Error with parameter passed to method
	ErrNotFound                        // Missing object
	ErrHelp                            // Help requested from command line
	ErrInternalAppError                // Internal application error
	ErrSignalCaught                    // Signal caught
	ErrUnexpectedResponse              // Unexpected Response
	ErrDuplicateItem                   // Duplicate Item
	ErrMax                = ErrDuplicateItem
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY AND WRAP

func (this Error) Error() string {
	switch this {
	case ErrNone:
		return "No Error"
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
	case ErrUnexpectedResponse:
		return "Unexpected response"
	case ErrDuplicateItem:
		return "Duplicate Item"
	default:
		return "[?? Invalid Error value]"
	}
}

func (this Error) WithPrefix(prefix string) error {
	return fmt.Errorf("%s: %w", prefix, this)
}

func NewCompoundError(errs ...error) *CompoundError {
	compound := &CompoundError{}
	for _, err := range errs {
		compound.Add(err)
	}
	return compound
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

func (this *CompoundError) Is(other error) bool {
	if len(this.errs) == 0 && (other == nil || errors.Is(other, ErrNone)) {
		return true
	} else if len(this.errs) == 1 {
		return errors.Is(this.errs[0], other)
	}
	// If any of the errors match, return true
	for _, err := range this.errs {
		if errors.Is(err, other) == true {
			return true
		}
	}
	// Return false by default
	return false
}
