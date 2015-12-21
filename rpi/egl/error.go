/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package egl

import (
	"fmt"
)

/* TYPES */
type Error struct {
	msg  string // description of error
	code int32  // error code
}

/* CONSTANTS */
const (
	SUCCESS             = 0x3000
	NOT_INITIALIZED     = 0x3001
	BAD_ACCESS          = 0x3002
	BAD_ALLOC           = 0x3003
	BAD_ATTRIBUTE       = 0x3004
	BAD_CONFIG          = 0x3005
	BAD_CONTEXT         = 0x3006
	BAD_CURRENT_SURFACE = 0x3007
	BAD_DISPLAY         = 0x3008
	BAD_MATCH           = 0x3009
	BAD_NATIVE_PIXMAP   = 0x300A
	BAD_NATIVE_WINDOW   = 0x300B
	BAD_PARAMETER       = 0x300C
	BAD_SURFACE         = 0x300D
	CONTEXT_LOST        = 0x300E
)

/* VARIABLES */
var (
	ErrorInvalidGraphicsConfiguration = toError(BAD_CONFIG)
)

/* METHODS */
func (e *Error) Error() string {
	return fmt.Sprintf("%s (egl/0x%04X)", e.msg, e.code)
}

func toError(code int32) *Error {
	switch code {
	case SUCCESS:
		return &Error{"Success", code}
	case NOT_INITIALIZED:
		return &Error{"Not Initialized", code}
	case BAD_ACCESS:
		return &Error{"Bad Access", code}
	case BAD_ALLOC:
		return &Error{"Bad Memory Allocation", code}
	case BAD_ATTRIBUTE:
		return &Error{"Bad Attribute", code}
	case BAD_CONFIG:
		return &Error{"Bad Configuration", code}
	case BAD_CONTEXT:
		return &Error{"Bad Context", code}
	case BAD_CURRENT_SURFACE:
		return &Error{"Bad Current Surface", code}
	case BAD_DISPLAY:
		return &Error{"Bad Display", code}
	case BAD_MATCH:
		return &Error{"Bad Match", code}
	case BAD_NATIVE_PIXMAP:
		return &Error{"Bad Native Pixmap", code}
	case BAD_NATIVE_WINDOW:
		return &Error{"Bad Native Window", code}
	case BAD_PARAMETER:
		return &Error{"Bad Parameter", code}
	case BAD_SURFACE:
		return &Error{"Bad Surface", code}
	case CONTEXT_LOST:
		return &Error{"Context Lost", code}
	default:
		return &Error{"General Error", code}
	}
}
