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
	gopi "../.."      /* import "github.com/djthorpe/gopi" */
	util "../../util" /* import "github.com/djthorpe/gopi/util" */
)

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host
    #cgo LDFLAGS:  -L/opt/vc/lib -lbcm_host
	#include "vc_dispmanx.h"
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type DXDisplayConfig struct {
	Display uint16
	Device  gopi.HardwareDriver
}

type DXDisplay struct {
	display uint16
	width   uint32
	height  uint32
	handle  dxDisplayHandle
	log     *util.LoggerDevice
}

type DXModeInfo struct {
	Size        DXSize
	Transform   DXTransform
	InputFormat DXInputFormat
	handle      dxDisplayHandle
}

type (
	dxDisplayHandle uint32
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DX_DISPLAY_NONE dxDisplayHandle = 0
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new Display object, returns error if not possible
func (config DXDisplayConfig) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<rpi.DXDisplay>Open display=%v", config.Display)

	// create new display object
	d := new(DXDisplay)

	// Set logging
	d.log = log

	// get the display size
	d.display = config.Display
	d.width, d.height = config.Device.GetDisplaySize(d.display)

	// open the display
	d.handle = dxDisplayOpen(d.display)
	if d.handle == DX_DISPLAY_NONE {
		return nil, d.log.Error("Cannot open display %v", d.display)
	}

	// success
	return d, nil
}

// Close the display
func (this *DXDisplay) Close() error {
	// Close display
	if dxDisplayClose(this.handle) != true {
		return this.log.Error("dxDisplayClose error")
	}
	// Return success
	this.log.Debug("<rpi.DXDisplay>Close display=%v", this.display)
	return nil
}

// Return display size
func (this *DXDisplay) GetSize() DXSize {
	return DXSize{this.width, this.height}
}

// Return mode info
func (this *DXDisplay) GetModeInfo() (*DXModeInfo, error) {
	var modeInfo DXModeInfo
	if dxDisplayGetInfo(this.handle, &modeInfo) != true {
		return nil, this.log.Error("dxDisplayGetInfo error")
	}
	return &modeInfo, nil
}

// Return logging object
func (this *DXDisplay) Log() *util.LoggerDevice {
	return this.log
}

// Human-readable version of the display
func (this *DXDisplay) String() string {
	return fmt.Sprintf("<rpi.DXDisplay>{ handle=%v display=%v size=%v", this.handle, this.display, this.GetSize())
}

// Human-readable version of the modeInfo
func (this *DXModeInfo) String() string {
	return fmt.Sprintf("<rpi.DXModeInfo>{ size=%v transform=%v inputformat=%v }", this.Size, this.Transform, this.InputFormat)
}

// Human-readable version of the dxDisplayHandle
func (d dxDisplayHandle) String() string {
	return fmt.Sprintf("<rpi.dxDisplayHandle>{%08X}", uint32(d))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func dxDisplayOpen(display uint16) dxDisplayHandle {
	return dxDisplayHandle(C.vc_dispmanx_display_open(C.uint32_t(display)))
}

func dxDisplayClose(display dxDisplayHandle) bool {
	return C.vc_dispmanx_display_close(C.DISPMANX_DISPLAY_HANDLE_T(display)) == DX_SUCCESS
}

func dxDisplayGetInfo(display dxDisplayHandle, info *DXModeInfo) bool {
	return C.vc_dispmanx_display_get_info(C.DISPMANX_DISPLAY_HANDLE_T(display), (*C.DISPMANX_MODEINFO_T)(unsafe.Pointer(info))) == DX_SUCCESS
}
