// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type EGL struct {
	Display gopi.Display
}

type egl struct {
	log          gopi.Logger
	display      gopi.Display
	handle       eglDisplay
	major, minor eglInt
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config EGL) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<surface.rpi.Open>{ Display=%v }", config.Display)
	this := new(egl)
	this.log = log

	// Check display
	this.display = config.Display
	if this.display == nil {
		return nil, gopi.ErrBadParameter
	}

	// Initialize EGL
	this.handle = to_eglDisplay(config.Display.Display())
	if major, minor, err := eglInitialize(this.handle); err != EGL_SUCCESS {
		return nil, os.NewSyscallError("Open", err)
	} else {
		this.major = major
		this.minor = minor
	}

	return this, nil
}

func (this *egl) Close() error {
	this.log.Debug("<surface.rpi.Close>{ Display=%v }", this.display)
	if this.display == nil {
		return nil
	}
	if err := eglTerminate(this.handle); err != EGL_SUCCESS {
		return os.NewSyscallError("Close", err)
	} else {
		this.display = nil
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *egl) String() string {
	if this.display == nil {
		return fmt.Sprintf("<surface.rpi>{ nil }")
	} else {
		return fmt.Sprintf("<surface.rpi>{ version={ %v,%v } display=%v }", this.major, this.minor, this.display)
	}
}
