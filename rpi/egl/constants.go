/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package egl

/*
	#cgo CFLAGS: -I/opt/vc/include
	#include "EGL/egl.h"
*/
import "C"

const (
	EGL_FALSE	= 0
	EGL_TRUE	= 1
)

