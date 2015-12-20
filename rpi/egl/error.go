/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package egl

/* TYPES */
type Error struct {
    msg    string // description of error
    code   uint64 // error code
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

/* METHODS */
func (e *Error) Error() string {
	return e.msg
}

func toError(code int32) Error {
	switch code {
		case SUCCESS:
			return Error{ "Success", code }
		case NOT_INITIALIZED:
			return Error{ "Not Initialized", code }
		case BAD_ACCESS:
			return Error{ "Bad Access", code }
		default:
			return Error{ "Other Error", code }
	}
}

