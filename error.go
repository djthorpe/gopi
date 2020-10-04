package gopi

import "fmt"

///////////////////////////////////////////////////////////////////////////////
// Types

type Error uint

///////////////////////////////////////////////////////////////////////////////
// Globals

const (
	ErrNone Error = iota
	ErrBadParameter
	ErrNotImplemented
	ErrUnexpectedResponse
	ErrHelp
	ErrInternalAppError
)

///////////////////////////////////////////////////////////////////////////////
// Implementation

func (e Error) Error() string {
	switch e {
	case ErrNone:
		return "No Error"
	case ErrBadParameter:
		return "Bad Parameter"
	case ErrNotImplemented:
		return "Not Implemented"
	case ErrHelp:
		return "Help Requested"
	case ErrUnexpectedResponse:
		return "Unexpected Response"
	case ErrInternalAppError:
		return "Internal Application Error"
	default:
		return "[?? Invalid Error]"
	}
}

func (e Error) WithPrefix(prefix string) error {
	return fmt.Errorf("%s: %w", prefix, e)
}
