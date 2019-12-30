/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Error uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ErrNone           Error = iota // No error condition
	ErrNotImplemented              // Method or feature not implemented
	ErrBadParameter                // Error with parameter passed to method
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY AND WRAP

func (this Error) Error() string {
	switch this {
	case ErrNotImplemented:
		return "Not Implemented"
	case ErrBadParameter:
		return "Bad Parameter"
	default:
		return "[?? Invalid Error value]"
	}
}

func (this Error) WithPrefix(prefix string) error {
	return fmt.Errorf("%s: %w", prefix, this)
}
