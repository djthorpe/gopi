/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"errors"
	"fmt"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	khronos "github.com/djthorpe/gopi/khronos"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////

/*
  #cgo CFLAGS:   -I/opt/vc/include
  #cgo LDFLAGS:  -L/opt/vc/lib -lEGL -lGLESv2 -lOpenVG
  #include <EGL/egl.h>
  #include <VG/openvg.h>
  #include <VG/vgu.h>
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
	log     *util.LoggerDevice
	egl     *eglDriver
	surface khronos.EGLSurface
}

type vgErrorType uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_CLEAR_COLOR          uint16  = 0x1121
	VG_PATH_FORMAT_STANDARD C.VGint = 0
)

// Datatypes
const (
	VG_PATH_DATATYPE_S_8 C.VGPathDatatype = iota
	VG_PATH_DATATYPE_S_16
	VG_PATH_DATATYPE_S_32
	VG_PATH_DATATYPE_F
)

// VG & VGU Error Codes
const (
	VG_NO_ERROR                       vgErrorType = 0x0000
	VG_BAD_HANDLE_ERROR               vgErrorType = 0x1000
	VG_ILLEGAL_ARGUMENT_ERROR         vgErrorType = 0x1001
	VG_OUT_OF_MEMORY_ERROR            vgErrorType = 0x1002
	VG_PATH_CAPABILITY_ERROR          vgErrorType = 0x1003
	VG_UNSUPPORTED_IMAGE_FORMAT_ERROR vgErrorType = 0x1004
	VG_UNSUPPORTED_PATH_FORMAT_ERROR  vgErrorType = 0x1005
	VG_IMAGE_IN_USE_ERROR             vgErrorType = 0x1006
	VG_NO_CONTEXT_ERROR               vgErrorType = 0x1007
	VGU_BAD_HANDLE_ERROR              vgErrorType = 0xF000
	VGU_ILLEGAL_ARGUMENT_ERROR        vgErrorType = 0xF001
	VGU_OUT_OF_MEMORY_ERROR           vgErrorType = 0xF002
	VGU_PATH_CAPABILITY_ERROR         vgErrorType = 0xF003
	VGU_BAD_WARP_ERROR                vgErrorType = 0xF004
)

// Path Capabilities
const (
	VG_PATH_CAPABILITY_APPEND_FROM             uint32 = (1 << 0)
	VG_PATH_CAPABILITY_APPEND_TO               uint32 = (1 << 1)
	VG_PATH_CAPABILITY_MODIFY                  uint32 = (1 << 2)
	VG_PATH_CAPABILITY_TRANSFORM_FROM          uint32 = (1 << 3)
	VG_PATH_CAPABILITY_TRANSFORM_TO            uint32 = (1 << 4)
	VG_PATH_CAPABILITY_INTERPOLATE_FROM        uint32 = (1 << 5)
	VG_PATH_CAPABILITY_INTERPOLATE_TO          uint32 = (1 << 6)
	VG_PATH_CAPABILITY_PATH_LENGTH             uint32 = (1 << 7)
	VG_PATH_CAPABILITY_POINT_ALONG_PATH        uint32 = (1 << 8)
	VG_PATH_CAPABILITY_TANGENT_ALONG_PATH      uint32 = (1 << 9)
	VG_PATH_CAPABILITY_PATH_BOUNDS             uint32 = (1 << 10)
	VG_PATH_CAPABILITY_PATH_TRANSFORMED_BOUNDS uint32 = (1 << 11)
	VG_PATH_CAPABILITY_ALL                     uint32 = (1 << 12) - 1
)

// Path draw mode
const (
	VG_STROKE_PATH C.VGbitfield = (1 << 0)
	VG_FILL_PATH   C.VGbitfield = (1 << 1)
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
	return fmt.Sprintf("<rpi.OpenVG>{ egl=%v surface=%v }", this.egl, this.surface)
}

////////////////////////////////////////////////////////////////////////////////
// BEGIN AND END

func (this *vgDriver) Begin(surface khronos.EGLSurface) error {
	this.log.Debug2("<rpi.OpenVG>Begin surface=%v", surface)
	if this.surface != nil {
		this.log.Warn("Begin() cannot be called without Flush()")
		if err := this.Flush(); err != nil {
			return err
		}
	}
	if err := this.egl.SetCurrentContext(surface); err != nil {
		return err
	}
	this.surface = surface
	return nil
}

func (this *vgDriver) Flush() error {
	this.log.Debug2("<rpi.OpenVG>Flush surface=%v", this.surface)
	if this.surface == nil {
		this.log.Warn("Flush() cannot be called without Begin()")
		return nil
	}
	C.vgFlush()
	if err := this.egl.FlushSurface(this.surface); err != nil {
		return err
	}
	this.surface = nil
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// CREATE PATH

func (this *vgDriver) CreatePath() (khronos.VGPath, error) {
	scale := C.VGfloat(1.0)
	bias := C.VGfloat(0.0)
	segCapacityHint := C.VGint(0)
	coordCapacityHint := C.VGint(0)
	capabilities := VG_PATH_CAPABILITY_ALL
	handle := C.vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_F, C.VGfloat(scale), C.VGfloat(bias), C.VGint(segCapacityHint), C.VGint(coordCapacityHint), C.VGbitfield(capabilities))
	return khronos.VGPath(handle), nil
}

func (this *vgDriver) DrawPath(path khronos.VGPath, stroke bool, fill bool) error {
	var paint_mode C.VGbitfield
	if stroke {
		paint_mode &= VG_STROKE_PATH
	}
	if fill {
		paint_mode &= VG_FILL_PATH
	}
	C.vgDrawPath(C.VGPath(path), C.VGbitfield(paint_mode))
	return nil
}

func (this *vgDriver) DestroyPath(path khronos.VGPath) error {
	C.vgDestroyPath(C.VGPath(path))
	return nil
}

func (this *vgDriver) ClearPath(path khronos.VGPath) error {
	capabilities := VG_PATH_CAPABILITY_ALL
	C.vgClearPath(C.VGPath(path), C.VGbitfield(capabilities))
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// OTHER

// Return string version of vgErrorType
func (e vgErrorType) String() string {
	switch e {
	case VG_NO_ERROR:
		return "VG_NO_ERROR"
	case VG_BAD_HANDLE_ERROR:
		return "VG_BAD_HANDLE_ERROR"
	case VG_ILLEGAL_ARGUMENT_ERROR:
		return "VG_ILLEGAL_ARGUMENT_ERROR"
	case VG_OUT_OF_MEMORY_ERROR:
		return "VG_OUT_OF_MEMORY_ERROR"
	case VG_PATH_CAPABILITY_ERROR:
		return "VG_PATH_CAPABILITY_ERROR"
	case VG_UNSUPPORTED_IMAGE_FORMAT_ERROR:
		return "VG_UNSUPPORTED_IMAGE_FORMAT_ERROR"
	case VG_UNSUPPORTED_PATH_FORMAT_ERROR:
		return "VG_UNSUPPORTED_PATH_FORMAT_ERROR"
	case VG_IMAGE_IN_USE_ERROR:
		return "VG_IMAGE_IN_USE_ERROR"
	case VG_NO_CONTEXT_ERROR:
		return "VG_NO_CONTEXT_ERROR"
	case VGU_BAD_HANDLE_ERROR:
		return "VGU_BAD_HANDLE_ERROR"
	case VGU_ILLEGAL_ARGUMENT_ERROR:
		return "VGU_ILLEGAL_ARGUMENT_ERROR"
	case VGU_OUT_OF_MEMORY_ERROR:
		return "VGU_OUT_OF_MEMORY_ERROR"
	case VGU_PATH_CAPABILITY_ERROR:
		return "VGU_PATH_CAPABILITY_ERROR"
	case VGU_BAD_WARP_ERROR:
		return "VGU_BAD_WARP_ERROR"
	default:
		return "[?? Invalid vgErrorType]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// GRAPHICS PRIMITIVES

func (this *vgDriver) Clear(color khronos.VGColor) error {
	if this.surface == nil {
		this.log.Warn("<rpi.OpenVG> Clear() cannot be called without Begin()")
		return nil
	}
	size := this.surface.GetSize()
	C.vgSetfv(C.VGParamType(VG_CLEAR_COLOR), C.VGint(4), (*C.VGfloat)(unsafe.Pointer(&color)))
	C.vgClear(C.VGint(0), C.VGint(0), C.VGint(size.Width), C.VGint(size.Height))

	return nil
}

func (this *vgDriver) Line(p1 khronos.VGPoint, p2 khronos.VGPoint) error {
	if this.surface == nil {
		this.log.Warn("<rpi.OpenVG> Line() cannot be called without Begin()")
		return nil
	}

	// create a path
	path, err := this.CreatePath()
	if err != nil {
		return err
	}

	// append line to path
	err = vguLine(path, p1, p2)
	if err != nil {
		this.DestroyPath(path)
		return err
	}

	// draw path - stroke but no fill
	if this.DrawPath(path, true, false) != nil {
		this.DestroyPath(path)
		return err
	}

	// destroy path
	return this.DestroyPath(path)
}

////////////////////////////////////////////////////////////////////////////////
// VGU GRAPHICS PRIMITIVES

func vgGetError(err vgErrorType) error {
	if err == VG_NO_ERROR {
		return nil
	}
	return errors.New(err.String())
}

func vguLine(path khronos.VGPath, p1 khronos.VGPoint, p2 khronos.VGPoint) error {
	return vgGetError(vgErrorType(C.vguLine(C.VGPath(path), C.VGfloat(p1.X), C.VGfloat(p1.Y), C.VGfloat(p2.X), C.VGfloat(p2.Y))))
}

/*
func (this *vgDriver) vguPolygon(path khronos.VGPath, points []khronos.VGPoint, closed bool) error {

}

func (this *vgDriver) vguRect(path khronos.VGPath, origin khronos.VGPoint, size khronos.VGSize) error {

}

func (this *vgDriver) vguRoundRect(path khronos.VGPath, origin khronos.VGPoint, size khronos.VGSize, arcWidth khronos.VGFloat, arcHeight khronos.VGFloat) error {

}

func (this *vgDriver) vguEllipse(path khronos.VGPath, origin khronos.VGPoint, size khronos.VGSize) error {

}
*/
