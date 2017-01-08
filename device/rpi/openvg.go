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

	path map[C.VGPath]*vgPath
	paint map[C.VGPaint]*vgPaint
}

// Paths
type vgPath struct {
	handle C.VGPath
}

// Paints
type vgPaint struct {
	handle C.VGPaint
}

// Errors
type vgErrorType uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_PATH_HANDLE_NONE C.VGPath = C.VGPath(0)
	VG_PATH_CAPACITY = 100
	VG_PAINT_HANDLE_NONE C.VGPaint = C.VGPaint(0)
	VG_PAINT_CAPACITY = 100
)

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

	// EGL driver
	egl, ok := config.EGL.(*eglDriver)
	if egl == nil || ok != true {
		return nil, this.log.Error("Invalid configuration parameter: EGL")
	}
	this.egl = egl

	// Create slice for paths and paints
	this.path = make(map[C.VGPath]*vgPath,VG_PATH_CAPACITY)
	this.paint = make(map[C.VGPaint]*vgPaint,VG_PAINT_CAPACITY)

	// Success
	return this, nil
}

// Close the driver
func (this *vgDriver) Close() error {

	// Close path and paint objects
	for _,path := range this.path {
		this.DestroyPath(path)
	}
	for _,paint := range this.paint {
		this.DestroyPaint(paint)
	}

	this.log.Debug2("<rpi.OpenVG>Close")
	return nil
}

// Return human-readable form of driver
func (this *vgDriver) String() string {
	return fmt.Sprintf("<rpi.OpenVG>{ egl=%v surface=%v }", this.egl, this.surface)
}

// Return human-readable form of path object
func (this *vgPath) String() string {
	return fmt.Sprintf("<rpi.VGPath>{ handle=0x%08X}", this.handle)
}

// Return human-readable form of paint object
func (this *vgPaint) String() string {
	return fmt.Sprintf("<rpi.VGPaint>{ handle=0x%08X}", this.handle)
}

////////////////////////////////////////////////////////////////////////////////
// BEGIN AND END

func (this *vgDriver) Begin(surface khronos.EGLSurface) error {
	this.log.Debug2("<rpi.OpenVG>Begin surface=%v", surface)
	if this.surface != nil {
		this.log.Warn("Begin() called before Flush()")
		if err := this.Flush(); err != nil {
			return err
		}
	}
	if err := this.egl.SetCurrentContext(surface); err != nil {
		return err
	}

	// Set identity matrix
	this.surface = surface
	C.vgLoadIdentity();

	return nil
}

func (this *vgDriver) Flush() error {
	this.log.Debug2("<rpi.OpenVG>Flush surface=%v", this.surface)
	if this.surface == nil {
		this.log.Warn("Flush() called before Begin()")
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
// CREATE & DESTROY PAINT

func (this *vgDriver) CreatePaint() (khronos.VGPaint, error) {
	handle := C.vgCreatePaint()
	if handle == VG_PAINT_HANDLE_NONE {
		return nil, vgGetError()
	}
	obj := &vgPaint{ handle: handle }
	this.paint[handle] = obj
	return obj, nil
}

func (this *vgDriver) DestroyPaint(paint khronos.VGPaint) error {
	obj, ok := paint.(*vgPaint)
	if ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	_, ok = this.paint[obj.handle]
	if ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	delete(this.paint,obj.handle)
	C.vgDestroyPaint(obj.handle)
	return vgGetError()
}

func (this *vgDriver) SetStroke(paint khronos.VGPaint) error {
	C.vgSetPaint(paint.(*vgPaint).handle,C.VGbitfield(khronos.VG_PAINT_STROKE))
	return vgGetError()
}

func (this *vgDriver) SetFill(paint khronos.VGPaint) error {
	C.vgSetPaint(paint.(*vgPaint).handle,C.VGbitfield(khronos.VG_PAINT_FILL))
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// CREATE & DESTROY PATH

func (this *vgDriver) CreatePath() (khronos.VGPath, error) {
	scale := C.VGfloat(1.0)
	bias := C.VGfloat(0.0)
	segCapacityHint := C.VGint(0)
	coordCapacityHint := C.VGint(0)
	capabilities := VG_PATH_CAPABILITY_ALL
	handle := C.vgCreatePath(VG_PATH_FORMAT_STANDARD, VG_PATH_DATATYPE_F, C.VGfloat(scale), C.VGfloat(bias), C.VGint(segCapacityHint), C.VGint(coordCapacityHint), C.VGbitfield(capabilities))
	if handle == VG_PATH_HANDLE_NONE {
		return nil, vgGetError()
	}
	obj := &vgPath{ handle: handle }
	this.path[handle] = obj
	return obj, nil
}

func (this *vgDriver) DestroyPath(path khronos.VGPath) error {
	obj, ok := path.(*vgPath)
	if ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	_, ok = this.path[obj.handle]
	if ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	delete(this.path,obj.handle)
	C.vgDestroyPath(obj.handle)
	return vgGetError()
}

func (this *vgPath) Clear() error {
	capabilities := VG_PATH_CAPABILITY_ALL
	C.vgClearPath(C.VGPath(this.handle), C.VGbitfield(capabilities))
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// SET PAINT STATE

func (this *vgPaint) SetColor(color khronos.VGColor) error {
	// TODO
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// DRAW METHODS

func (this *vgPath) Draw(flags khronos.VGPaintMode) error {
	C.vgDrawPath(this.handle, C.VGbitfield(flags))
	return vgGetError()
}

func (this *vgPath) Line(start, end khronos.VGPoint) error {
	// TODO
	return nil
}

func (this *vgPath) Rect(origin, size khronos.VGPoint) error {
	// TODO
	return nil
}

func (this *vgPath) Ellipse(center, size khronos.VGPoint) error {
	// TODO
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
		return this.log.Error("<rpi.OpenVG> Clear() cannot be called without Begin()")
	}

	size := this.surface.GetSize()
	C.vgSetfv(C.VGParamType(VG_CLEAR_COLOR), C.VGint(4), (*C.VGfloat)(unsafe.Pointer(&color)))
	C.vgClear(C.VGint(0), C.VGint(0), C.VGint(size.Width), C.VGint(size.Height))

	return vgGetError()
}

/*

func (this *vgDriver) Line(p1 khronos.VGPoint, p2 khronos.VGPoint) error {
	if this.surface == nil {
		return this.log.Error("<rpi.OpenVG> Line() cannot be called without Begin()")
	}

	// create a path
	path, err := this.CreatePath()
	if err != nil {
		return err
	}
	defer this.DestroyPath(path)

	// append line to path
	if err := vguLine(path, p1, p2); err != nil {
		return err
	}

	// draw path - stroke but no fill
	if err := this.DrawPath(path, true, false); err != nil {
		return err
	}

	// success
	return nil
}

func (this *vgDriver) Ellipse(origin khronos.VGPoint,size khronos.VGSize) error {
	if this.surface == nil {
		return this.log.Error("<rpi.OpenVG> Ellipse() cannot be called without Begin()")
	}

	// create a path
	path, err := this.CreatePath()
	if err != nil {
		return err
	}

	// append ellipse to path
	err = vguEllipse(path, origin, size)
	if err != nil {
		this.DestroyPath(path)
		return err
	}

}
*/

////////////////////////////////////////////////////////////////////////////////
// VGU GRAPHICS PRIMITIVES

func vgGetError() error {
	err := vgErrorType(C.vgGetError())
	if err == VG_NO_ERROR {
		return nil
	}
	return vgError(err)
}

func vgError(err vgErrorType) error {
	if err == VG_NO_ERROR {
		return nil
	}
	return errors.New(err.String())
}

/*

func vguLine(path khronos.VGPath, p1 khronos.VGPoint, p2 khronos.VGPoint) error {
	return vgGetError(vgErrorType(C.vguLine(C.VGPath(path), C.VGfloat(p1.X), C.VGfloat(p1.Y), C.VGfloat(p2.X), C.VGfloat(p2.Y))))
}

func vguEllipse(path khronos.VGPath, origin khronos.VGPoint, size khronos.VGSize) error {
	return vgGetError(vgErrorType(C.vguEllipse(C.VGPath(path), C.VGfloat(origin.X), C.VGfloat(origin.Y), C.VGfloat(size.Width), C.VGfloat(size.Height))))
}

func (this *vgDriver) vguPolygon(path khronos.VGPath, points []khronos.VGPoint, closed bool) error {

}

func (this *vgDriver) vguRect(path khronos.VGPath, origin khronos.VGPoint, size khronos.VGSize) error {

}

func (this *vgDriver) vguRoundRect(path khronos.VGPath, origin khronos.VGPoint, size khronos.VGSize, arcWidth khronos.VGFloat, arcHeight khronos.VGFloat) error {

}

func (this *vgDriver) vguEllipse(path khronos.VGPath, origin khronos.VGPoint, size khronos.VGSize) error {

}
*/
