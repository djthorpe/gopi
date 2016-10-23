/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
	"unsafe"
)

import (
	gopi "../.."            /* import "github.com/djthorpe/gopi" */
	util "../../util"       /* import "github.com/djthorpe/gopi/util" */
	khronos "../../khronos"       /* import "github.com/djthorpe/gopi/khronos" */
)

////////////////////////////////////////////////////////////////////////////////

/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2
  #include <EGL/egl.h>
  #include <VG/openvg.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Configuration when creating the OpenVG driver
type OpenVG struct {
	EGL khronos.EGLDriver
}

// EGL driver
type vgDriver struct {
	log          *util.LoggerDevice
	egl          *eglDriver
	window       khronos.EGLWindow
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_CLEAR_COLOR uint16 = 0x1121
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

// Open
func (config OpenVG) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	this := new(vgDriver)
	this.log = log
	this.log.Debug2("<rpi.OpenVG>Open")

	egl, ok := config.EGL.(*eglDriver)
	if egl == nil || ok != true {
		return nil, this.log.Error("Invalid configuration parameter: EGL")
	}
	this.egl = egl

	// Success
	return this, nil
}

// Close the driver
func (this *vgDriver) Close() error {
	this.log.Debug2("<rpi.OpenVG>Close")
	return nil
}

// Return the logging object
func (this *vgDriver) Log() *util.LoggerDevice {
	return this.log
}

// Return human-readable form of driver
func (this *vgDriver) String() string {
	return fmt.Sprintf("<rpi.OpenVG>{ egl=%v window=%v }",this.egl,this.window)
}

////////////////////////////////////////////////////////////////////////////////
// BEGIN AND END

func (this *vgDriver) Begin(window khronos.EGLWindow) error {
	if this.window != nil {
		this.log.Warn("<rpi.OpenVG> Begin() cannot be called without Flush()")
		if err := this.Flush(); err != nil {
			return err
		}
	}
	this.window = window
	return nil
}

func (this *vgDriver) Flush() error {
	if this.window == nil {
		this.log.Warn("<rpi.OpenVG> Flush() cannot be called without Begin()")
		return nil
	}
	C.vgFlush()
	if err := this.egl.FlushWindow(this.window); err != nil {
		return err
	}
	this.window = nil
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// GRAPHICS PRIMITIVES

func (this *vgDriver) Clear(color khronos.VGColor) {
	if this.window == nil {
		this.log.Warn("<rpi.OpenVG> Clear() cannot be called without Begin()")
		return
	}
	C.vgSetfv(C.VGParamType(VG_CLEAR_COLOR),C.VGint(4),(*C.VGfloat)(unsafe.Pointer(&color)));
	C.vgClear(C.VGint(0), C.VGint(0), C.VGint(100), C.VGint(100));
}

