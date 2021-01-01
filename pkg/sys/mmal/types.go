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
	MMAL_RECT_T     C.struct_MMAL_RECT_T
	MMAL_RATIONAL_T C.struct_MMAL_RATIONAL_T
)
