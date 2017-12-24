/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package mmal /* import "github.com/djthorpe/gopi/device/rpi/mmal" */

import (
	"errors"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type status int

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_SUCCESS   status = iota
	MMAL_ENOMEM           // Out of memory
	MMAL_ENOSPC           // Out of resources (other than memory)
	MMAL_EINVAL           // Argument is invalid
	MMAL_ENOSYS           // Function not implemented
	MMAL_ENOENT           // No such file or directory
	MMAL_ENXIO            // No such device or address
	MMAL_EIO              // I/O error
	MMAL_ESPIPE           // Illegal seek
	MMAL_ECORRUPT         // Data is corrupt
	MMAL_ENOTREADY        // Component is not ready
	MMAL_ECONFIG          // Component is not configured
	MMAL_EISCONN          // Port is already connected
	MMAL_ENOTCONN         // Port is disconnected
	MMAL_EAGAIN           // Resource temporarily unavailable. Try again later
	MMAL_EFAULT           // Bad address
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (s status) Error() error {
	switch s {
	case MMAL_SUCCESS:
		return nil
	case MMAL_ENOMEM:
		return errors.New("MMAL_ENOMEM")
	case MMAL_ENOSPC:
		return errors.New("MMAL_ENOSPC")
	case MMAL_EINVAL:
		return errors.New("MMAL_EINVAL")
	case MMAL_ENOSYS:
		return errors.New("MMAL_ENOSYS")
	case MMAL_ENOENT:
		return errors.New("MMAL_ENOENT")
	case MMAL_ENXIO:
		return errors.New("MMAL_ENXIO")
	case MMAL_EIO:
		return errors.New("MMAL_EIO")
	case MMAL_ESPIPE:
		return errors.New("MMAL_ESPIPE")
	case MMAL_ECORRUPT:
		return errors.New("MMAL_ECORRUPT")
	case MMAL_ENOTREADY:
		return errors.New("MMAL_ENOTREADY")
	case MMAL_ECONFIG:
		return errors.New("MMAL_ECONFIG")
	case MMAL_EISCONN:
		return errors.New("MMAL_EISCONN")
	case MMAL_ENOTCONN:
		return errors.New("MMAL_ENOTCONN")
	case MMAL_EAGAIN:
		return errors.New("MMAL_EAGAIN")
	case MMAL_EFAULT:
		return errors.New("MMAL_EFAULT")
	default:
		return errors.New("[?? Invalid status value]")
	}
}
