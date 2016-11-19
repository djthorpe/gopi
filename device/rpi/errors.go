/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"errors"
)

var (
	eglError = map[uint16]error{
		EGL_SUCCESS:             nil,
		EGL_NOT_INITIALIZED:     errors.New("EGL_NOT_INITIALIZED"),
		EGL_BAD_ACCESS:          errors.New("EGL_BAD_ACCESS"),
		EGL_BAD_ALLOC:           errors.New("EGL_BAD_ALLOC"),
		EGL_BAD_ATTRIBUTE:       errors.New("EGL_BAD_ATTRIBUTE"),
		EGL_BAD_CONFIG:          errors.New("EGL_BAD_CONFIG"),
		EGL_BAD_CONTEXT:         errors.New("EGL_BAD_CONTEXT"),
		EGL_BAD_CURRENT_SURFACE: errors.New("EGL_BAD_CURRENT_SURFACE"),
		EGL_BAD_DISPLAY:         errors.New("EGL_BAD_DISPLAY"),
		EGL_BAD_MATCH:           errors.New("EGL_BAD_MATCH"),
		EGL_BAD_NATIVE_PIXMAP:   errors.New("EGL_BAD_NATIVE_PIXMAP"),
		EGL_BAD_NATIVE_WINDOW:   errors.New("EGL_BAD_NATIVE_WINDOW"),
		EGL_BAD_PARAMETER:       errors.New("EGL_BAD_PARAMETER"),
		EGL_BAD_SURFACE:         errors.New("EGL_BAD_SURFACE"),
		EGL_CONTEXT_LOST:        errors.New("EGL_CONTEXT_LOST"),
	}
	EGLErrorUnknown                  = errors.New("Unknown EGL error")
	EGLErrorInvalidDisplayDriver     = errors.New("Invalid display driver parameter")
	EGLErrorInvalidAPIBind           = errors.New("Invalid EGL API binding parameter")
	EGLErrorInvalidFrameBufferConfig = errors.New("Invalid EGL framebuffer parameter")
	EGLErrorNoBitmap                 = errors.New("No bitmap")
	EGLErrorInvalidParameter         = errors.New("Invalid parameter")
	ErrInvalidMasterParam = errors.New("Invalid I2C master number")
	ErrNotImplemented = errors.New("Not Implemented")
)
