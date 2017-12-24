/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package mmal /* import "github.com/djthorpe/gopi/device/rpi/mmal" */

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/mmal
    #cgo LDFLAGS:  -L/opt/vc/lib -lmmal -lmmal_components -lmmal_core
	#include <mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Type of format (control, audio, video, subpicture or unknown)
type FormatType int

// Encoding and Variant
type FormatEncoding uint32

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_ES_TYPE_UNKNOWN FormatType = iota
	MMAL_ES_TYPE_CONTROL
	MMAL_ES_TYPE_AUDIO
	MMAL_ES_TYPE_VIDEO
	MMAL_ES_TYPE_SUBPICTURE
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t FormatType) String() string {
	switch t {
	case MMAL_ES_TYPE_UNKNOWN:
		return "MMAL_ES_TYPE_UNKNOWN"
	case MMAL_ES_TYPE_CONTROL:
		return "MMAL_ES_TYPE_CONTROL"
	case MMAL_ES_TYPE_AUDIO:
		return "MMAL_ES_TYPE_AUDIO"
	case MMAL_ES_TYPE_VIDEO:
		return "MMAL_ES_TYPE_VIDEO"
	case MMAL_ES_TYPE_SUBPICTURE:
		return "MMAL_ES_TYPE_SUBPICTURE"
	default:
		return "[?? Unknown FormatType value]"
	}
}
