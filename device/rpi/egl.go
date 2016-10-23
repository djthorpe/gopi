/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
)

import (
	gopi "../.." /* import "github.com/djthorpe/gopi" */
	util "../../util" /* import "github.com/djthorpe/gopi/util" */
	khronos "../../khronos" /* import "github.com/djthorpe/gopi/khronos" */
)

////////////////////////////////////////////////////////////////////////////////

/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2
  #include <EGL/egl.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////

// Configuration when creating the EGL driver
type EGL struct {
	Display gopi.DisplayDriver
}

// Display handle
type eglDisplay uintptr

// Context handle
type eglContext uintptr

// Surface handle
type eglSurface uintptr

// Configuration handle
type eglConfig uintptr

// EGL driver
type eglDriver struct {
	major, minor int
	dx           *DXDisplay
	display      eglDisplay
	log          *util.LoggerDevice
}

////////////////////////////////////////////////////////////////////////////////
// Open and close device

func (config EGL) Open(log *util.LoggerDevice) (khronos.EGLDriver, error) {
	log.Debug2("<rpi.EGL> Open")

	this := new(eglDriver)
	this.log = log
	return this, nil
}

func (this *eglDriver) Close() error {
	this.log.Debug2("<rpi.EGL> Close")
	return nil
}

func (this *eglDriver) String() string {
	return fmt.Sprintf("<rpi.EGL>{ TODO }")
}

func (this *eglDriver) Log() *util.LoggerDevice {
	return this.log
}

