/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
	"math"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	khronos "github.com/djthorpe/gopi/khronos"
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
	// The number of the display to open
	Display uint16

	// The physical inches on the diagnol for the display or zero if unknown
	PhysicalInches float64

	// Hardware board driver
	Device gopi.HardwareDriver
}

type DXDisplay struct {
	display uint16
	width   uint32
	height  uint32
	ppi     uint32
	handle  dxDisplayHandle
	log     gopi.Logger
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
func (config DXDisplayConfig) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<rpi.DXDisplay>Open display=%v", config.Display)

	// create new display object
	d := new(DXDisplay)

	// Set logging
	d.log = log

	// get the display size
	d.display = config.Display
	d.width, d.height = config.Device.GetDisplaySize(d.display)

	// set the pixels-per-inch value
	if config.PhysicalInches > 0.0 {
		pixels := math.Sqrt(math.Pow(float64(d.width), 2.0) + math.Pow(float64(d.height), 2.0))
		d.ppi = uint32(math.Floor((pixels / float64(config.PhysicalInches)) + 0.5))
	} else {
		d.ppi = 0
	}

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

// Return pixel density in pixels per inch.
// Returns 0 if no PPI value has been set.
func (this *DXDisplay) GetPixelsPerInch() uint32 {
	return this.ppi
}

// Return the size of the display in pixels
func (this *DXDisplay) GetDisplaySize() (uint32, uint32) {
	return this.width, this.height
}

// Return the display number
func (this *DXDisplay) GetDisplay() uint16 {
	return this.display
}

// Return mode info
func (this *DXDisplay) GetModeInfo() (*DXModeInfo, error) {
	var modeInfo DXModeInfo
	if dxDisplayGetInfo(this.handle, &modeInfo) != true {
		return nil, this.log.Error("dxDisplayGetInfo error")
	}
	return &modeInfo, nil
}

// Create snapshot of screen
func (this *DXDisplay) CreateSnapshot() (*DXResource, error) {
	// create a resource
	resource, err := this.CreateResource(DX_IMAGE_RGBA32, khronos.EGLSize{uint(this.width), uint(this.height)})
	if err != nil {
		return nil, err
	}
	if ret := dxDisplaySnapshot(this.handle, resource.handle, DX_NO_ROTATE); ret == false {
		return nil, EGLErrorSnapshot
	}
	return resource, nil
}

// Human-readable version of the display
func (this *DXDisplay) String() string {
	return fmt.Sprintf("<rpi.DXDisplay>{ handle=%v display=%v size=%v ppi=%v", this.handle, this.display, this.GetSize(), this.GetPixelsPerInch())
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

func dxDisplaySnapshot(display dxDisplayHandle, resource dxResourceHandle, transform DXTransform) bool {
	return C.vc_dispmanx_snapshot(C.DISPMANX_DISPLAY_HANDLE_T(display), C.DISPMANX_RESOURCE_HANDLE_T(resource), C.DISPMANX_TRANSFORM_T(transform)) == DX_SUCCESS
}
