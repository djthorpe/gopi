/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gopi

type Error uint

const (
	ErrNone Error = iota
	ErrNotImplemented
	ErrBadParameter
)

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
