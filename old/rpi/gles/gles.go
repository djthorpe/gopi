/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package gles

/*
	#cgo CFLAGS:   -I/opt/vc/include
	#cgo LDFLAGS:  -L/opt/vc/lib -lGLESv2
	#include <GLES2/gl2.h>
	#include <GLES2/gl2ext.h>
	#include <GLES2/gl2platform.h>
*/
import "C"
