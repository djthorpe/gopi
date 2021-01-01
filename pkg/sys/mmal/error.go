//+build mmal

package mmal

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
#include <interface/mmal/util/mmal_connection.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Error int
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_SUCCESS   Error = C.MMAL_SUCCESS
	MMAL_ENOMEM    Error = C.MMAL_ENOMEM    // Out of memory
	MMAL_ENOSPC    Error = C.MMAL_ENOSPC    // Out of resources (other than memory)
	MMAL_EINVAL    Error = C.MMAL_EINVAL    // Argument is invalid
	MMAL_ENOSYS    Error = C.MMAL_ENOSYS    // Function not implemented
	MMAL_ENOENT    Error = C.MMAL_ENOENT    // No such file or directory
	MMAL_ENXIO     Error = C.MMAL_ENXIO     // No such device or address
	MMAL_EIO       Error = C.MMAL_EIO       // I/O error
	MMAL_ESPIPE    Error = C.MMAL_ESPIPE    // Illegal seek
	MMAL_ECORRUPT  Error = C.MMAL_ECORRUPT  // Data is corrupt
	MMAL_ENOTREADY Error = C.MMAL_ENOTREADY // Component is not ready
	MMAL_ECONFIG   Error = C.MMAL_ECONFIG   // Component is not configured
	MMAL_EISCONN   Error = C.MMAL_EISCONN   // Port is already connected
	MMAL_ENOTCONN  Error = C.MMAL_ENOTCONN  // Port is disconnected
	MMAL_EAGAIN    Error = C.MMAL_EAGAIN    // Resource temporarily unavailable. Try again later
	MMAL_EFAULT    Error = C.MMAL_EFAULT    // Bad address
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e Error) Error() string {
	switch e {
	case MMAL_SUCCESS:
		return "MMAL_SUCCESS"
	case MMAL_ENOMEM:
		return "MMAL_ENOMEM"
	case MMAL_ENOSPC:
		return "MMAL_ENOSPC"
	case MMAL_EINVAL:
		return "MMAL_EINVAL"
	case MMAL_ENOSYS:
		return "MMAL_ENOSYS"
	case MMAL_ENOENT:
		return "MMAL_ENOENT"
	case MMAL_ENXIO:
		return "MMAL_ENXIO"
	case MMAL_EIO:
		return "MMAL_EIO"
	case MMAL_ESPIPE:
		return "MMAL_ESPIPE"
	case MMAL_ECORRUPT:
		return "MMAL_ECORRUPT"
	case MMAL_ENOTREADY:
		return "MMAL_ENOTREADY"
	case MMAL_ECONFIG:
		return "MMAL_ECONFIG"
	case MMAL_EISCONN:
		return "MMAL_EISCONN"
	case MMAL_ENOTCONN:
		return "MMAL_ENOTCONN"
	case MMAL_EAGAIN:
		return "MMAL_EAGAIN"
	case MMAL_EFAULT:
		return "MMAL_EFAULT"
	default:
		return "[?? Invalid Error value]"
	}
}
