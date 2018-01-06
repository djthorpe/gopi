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
	major, minor int
}

// Raspberry-pi specific interface for SurfaceManager
type SurfaceManager interface {
	gopi.SurfaceManager

	// Return a list of extensions the GPU provides
	Extensions() []string
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config EGL) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.surface.rpi.SurfaceManager.Open>{ Display=%v }", config.Display)
	this := new(egl)
	this.log = log

	// Check display
	this.display = config.Display
	if this.display == nil {
		return nil, gopi.ErrBadParameter
	}

	// Initialize EGL
	n := to_eglNativeDisplayType(this.display.Display())
	if handle, err := eglGetDisplay(n); err != EGL_SUCCESS {
		return nil, os.NewSyscallError("eglGetDisplay", err)
	} else {
		this.handle = handle
	}
	if major, minor, err := eglInitialize(this.handle); err != EGL_SUCCESS {
		return nil, os.NewSyscallError("eglInitialize", err)
	} else {
		this.major = int(major)
		this.minor = int(minor)
	}

	// Get configurations
	if configs, err := eglGetConfigs(this.handle); err != EGL_SUCCESS {
		return nil, os.NewSyscallError("eglGetConfigs", err)
	} else {
		for i, config := range configs {
			if a, err := eglGetConfigAttribs(this.handle, config); err != EGL_SUCCESS {
				return nil, os.NewSyscallError("eglGetConfigAttribs", err)
			} else {
				fmt.Println(i, a)
			}
		}
	}

	return this, nil
}

func (this *egl) Close() error {
	this.log.Debug("<sys.surface.rpi.SurfaceManager.Close>{ Display=%v }", this.display)
	if this.display == nil {
		return nil
	}
	if err := eglTerminate(this.handle); err != EGL_SUCCESS {
		return os.NewSyscallError("Close", err)
	} else {
		this.display = nil
		this.handle = eglDisplay(EGL_NO_DISPLAY)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *egl) String() string {
	if this.display == nil {
		return fmt.Sprintf("<sys.surface.rpi.SurfaceManager>{ nil }")
	} else {
		return fmt.Sprintf("<sys.surface.rpi.SurfaceManager>{ handle=%v name=%v version={ %v,%v } extensions=%v client_apis=%v display=%v }", this.handle, this.Name(), this.major, this.minor, this.Extensions(), this.ClientAPIs(), this.display)
	}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

func (this *egl) Display() gopi.Display {
	return this.display
}

func (this *egl) Name() string {
	return fmt.Sprintf("%v %v", eglQueryString(this.handle, EGL_VENDOR), eglQueryString(this.handle, EGL_VERSION))
}

func (this *egl) Extensions() string {
	return eglQueryString(this.handle, EGL_EXTENSIONS)
}

func (this *egl) ClientAPIs() string {
	return eglQueryString(this.handle, EGL_CLIENT_APIS)
}
