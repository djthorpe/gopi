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
	"sync"
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
	lock    sync.Mutex
	surface khronos.EGLSurface
	path    map[C.VGPath]*vgPath
	paint   map[C.VGPaint]*vgPaint
}

// Paths
type vgPath struct {
	handle C.VGPath
	log    *util.LoggerDevice
}

// Paints
type vgPaint struct {
	handle     C.VGPaint
	line_width float32
	join_style khronos.VGStrokeJoinStyle
	cap_style  khronos.VGStrokeCapStyle
	log        *util.LoggerDevice
}

// Errors
type vgErrorType uint16

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_PATH_HANDLE_NONE  C.VGPath  = C.VGPath(0)
	VG_PATH_CAPACITY               = 100
	VG_PAINT_HANDLE_NONE C.VGPaint = C.VGPaint(0)
	VG_PAINT_CAPACITY              = 100
)

const (
	VG_CLEAR_COLOR             uint16  = 0x1121
	VG_STROKE_LINE_WIDTH       uint16  = 0x1110
	VG_STROKE_CAP_STYLE        uint16  = 0x1111
	VG_STROKE_JOIN_STYLE       uint16  = 0x1112
	VG_STROKE_MITER_LIMIT      uint16  = 0x1113
	VG_STROKE_DASH_PATTERN     uint16  = 0x1114
	VG_STROKE_DASH_PHASE       uint16  = 0x1115
	VG_STROKE_DASH_PHASE_RESET uint16  = 0x1116
	VG_PAINT_COLOR             uint16  = 0x1A01
	VG_PATH_FORMAT_STANDARD    C.VGint = 0
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
	this.path = make(map[C.VGPath]*vgPath, VG_PATH_CAPACITY)
	this.paint = make(map[C.VGPaint]*vgPaint, VG_PAINT_CAPACITY)

	// Success
	return this, nil
}

// Close the driver
func (this *vgDriver) Close() error {

	// Close path and paint objects
	for _, path := range this.path {
		this.DestroyPath(path)
	}
	for _, paint := range this.paint {
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
	return fmt.Sprintf("<rpi.VGPath>{ handle=0x%08X }", this.handle)
}

// Return human-readable form of paint object
func (this *vgPaint) String() string {
	return fmt.Sprintf("<rpi.VGPaint>{ handle=0x%08X }", this.handle)
}

////////////////////////////////////////////////////////////////////////////////
// BEGIN AND END

func (this *vgDriver) Begin(surface khronos.EGLSurface) error {
	this.log.Debug2("<rpi.OpenVG>Begin surface=%v", surface)

	if this.surface != nil {
		return this.log.Error("Begin() called before Flush()")
	}

	this.lock.Lock()

	if err := this.egl.SetCurrentContext(surface); err != nil {
		this.lock.Unlock()
		return err
	}

	// Set identity matrix
	this.surface = surface
	C.vgLoadIdentity()

	return nil
}

func (this *vgDriver) Flush() error {
	this.log.Debug2("<rpi.OpenVG>Flush surface=%v", this.surface)
	if this.surface == nil {
		this.lock.Unlock()
		return this.log.Error("Flush() called before Begin()")
	}
	C.vgFlush()
	if err := vgGetError(); err != nil {
		this.lock.Unlock()
		return err
	}
	if err := this.egl.FlushSurface(this.surface); err != nil {
		this.lock.Unlock()
		return err
	}
	this.surface = nil
	this.lock.Unlock()
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// GET SURFACE POINTS

// Return point aligned to surface
func (this *vgDriver) GetPoint(flags khronos.EGLFrameAlignFlag) khronos.VGPoint {
	var pt khronos.VGPoint

	if this.surface == nil {
		this.log.Warn("GetPoint() called before Flush()")
		return pt
	}

	size := this.surface.GetSize()

	switch { /* Y */
	case flags&khronos.EGL_ALIGN_VCENTER != 0:
		pt.Y = float32(size.Height >> 1)
	case flags&khronos.EGL_ALIGN_TOP != 0:
		pt.Y = 0
	case flags&khronos.EGL_ALIGN_BOTTOM != 0:
		pt.Y = float32(size.Height - 1)
	}
	switch { /* X */
	case flags&khronos.EGL_ALIGN_HCENTER != 0:
		pt.X = float32(size.Width >> 1)
	case flags&khronos.EGL_ALIGN_LEFT != 0:
		pt.X = 0
	case flags&khronos.EGL_ALIGN_RIGHT != 0:
		pt.X = float32(size.Width - 1)
	}
	return pt
}

////////////////////////////////////////////////////////////////////////////////
// CREATE & DESTROY PAINT

func (this *vgDriver) CreatePaint(color khronos.VGColor) (khronos.VGPaint, error) {
	// Get handle for the paint object
	handle := C.vgCreatePaint()
	if handle == VG_PAINT_HANDLE_NONE {
		return nil, vgGetError()
	}

	// Create the object
	obj := &vgPaint{
		handle: handle,
		log:    this.log,
		line_width: 1.0,
		cap_style: khronos.VG_STYLE_CAP_NONE,
		join_style: khronos.VG_STYLE_JOIN_NONE,
	}
	this.paint[handle] = obj

	// Set the color
	if err := obj.SetColor(color); err != nil {
		this.DestroyPaint(obj)
		return nil, err
	}

	this.log.Debug2("<rpi.OpenVG>CreatePaint{ paint=%v }", obj)

	// success
	return obj, nil
}

func (this *vgDriver) DestroyPaint(paint khronos.VGPaint) error {

	obj, ok := paint.(*vgPaint)
	if ok == false || obj.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	_, ok = this.paint[obj.handle]
	if ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	delete(this.paint, obj.handle)

	this.log.Debug2("<rpi.OpenVG>DestroyPaint{ paint=%v }", paint)

	C.vgDestroyPaint(obj.handle)
	obj.handle = VG_PAINT_HANDLE_NONE
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// PAINT METHODS

func (this *vgPaint) SetColor(color khronos.VGColor) error {
	if this.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	C.vgSetParameterfv(C.VGHandle(this.handle), C.VGint(VG_PAINT_COLOR), 4, (*C.VGfloat)(unsafe.Pointer(&color)))
	return vgGetError()
}

func (this *vgPaint) SetLineWidth(width float32) error {
	if this.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	if width < 0.0 {
		return vgError(VG_ILLEGAL_ARGUMENT_ERROR)
	}
	this.line_width = width
	return nil
}

func (this *vgPaint) SetStrokeStyle(join khronos.VGStrokeJoinStyle, cap VGStrokeCapStyle) error {
	if this.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	this.cap_style = cap
	this.join_style = join
	return nil
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
	obj := &vgPath{
		handle: handle,
		log:    this.log,
	}
	this.path[handle] = obj

	this.log.Debug2("<rpi.OpenVG>CreatePath{ path=%v }", obj)

	return obj, nil
}

func (this *vgDriver) DestroyPath(path khronos.VGPath) error {
	obj, ok := path.(*vgPath)
	if ok == false || obj.handle == VG_PATH_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	_, ok = this.path[obj.handle]
	if ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}

	this.log.Debug2("<rpi.OpenVG>DestroyPath{ path=%v }", obj)

	delete(this.path, obj.handle)
	C.vgDestroyPath(obj.handle)
	obj.handle = VG_PATH_HANDLE_NONE
	return vgGetError()
}

func (this *vgPath) Clear() error {
	capabilities := VG_PATH_CAPABILITY_ALL
	C.vgClearPath(C.VGPath(this.handle), C.VGbitfield(capabilities))
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// DRAW METHODS

func (this *vgPath) Draw(stroke, fill khronos.VGPaint) error {
	var flags C.VGbitfield
	stroke_obj, ok := stroke.(*vgPaint)
	if stroke_obj != nil && ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	fill_obj, ok := fill.(*vgPaint)
	if fill_obj != nil && ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	if stroke_obj != nil && stroke_obj.handle != VG_PAINT_HANDLE_NONE {
		flags |= VG_STROKE_PATH
		C.vgSetPaint(stroke_obj.handle, VG_STROKE_PATH)
		C.vgSetf(C.VGParamType(VG_STROKE_LINE_WIDTH), C.VGfloat(stroke_obj.line_width))
		if stroke_obj.cap_style != VG_STYLE_CAP_NONE {
			C.vgSetf(C.VGParamType(VG_STROKE_CAP_STYLE), C.VGfloat(stroke_obj.cap_style))
		}
		if stroke_obj.cap_style != VG_STYLE_JOIN_NONE {
			C.vgSetf(C.VGParamType(VG_STROKE_JOIN_STYLE), C.VGfloat(stroke_obj.join_style))
		}
	}
	if fill_obj != nil && fill_obj.handle != VG_PAINT_HANDLE_NONE {
		flags |= VG_FILL_PATH
		C.vgSetPaint(fill_obj.handle, VG_FILL_PATH)
	}
	C.vgDrawPath(this.handle, flags)
	return vgGetError()
}

func (this *vgPath) Fill(fill khronos.VGPaint) error {
	return this.Draw(nil, fill)
}

func (this *vgPath) Stroke(stroke khronos.VGPaint) error {
	return this.Draw(stroke, nil)
}

func (this *vgPath) Line(start, end khronos.VGPoint) error {
	err := vgErrorType(C.vguRect(this.handle, C.VGfloat(start.X), C.VGfloat(start.Y), C.VGfloat(end.X), C.VGfloat(end.Y)))
	if err != VG_NO_ERROR {
		return vgError(err)
	}
	return nil
}

func (this *vgPath) Rect(origin, size khronos.VGPoint) error {
	err := vgErrorType(C.vguRect(this.handle, C.VGfloat(origin.X), C.VGfloat(origin.Y), C.VGfloat(size.X), C.VGfloat(size.Y)))
	if err != VG_NO_ERROR {
		return vgError(err)
	}
	return nil
}

func (this *vgPath) Ellipse(center, diameter khronos.VGPoint) error {
	err := vgErrorType(C.vguEllipse(this.handle, C.VGfloat(center.X), C.VGfloat(center.Y), C.VGfloat(diameter.X), C.VGfloat(diameter.Y)))
	if err != VG_NO_ERROR {
		return vgError(err)
	}
	return nil
}

func (this *vgPath) Circle(center khronos.VGPoint, diameter float32) error {
	return this.Ellipse(center, khronos.VGPoint{diameter, diameter})
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

////////////////////////////////////////////////////////////////////////////////
// VGU GRAPHICS PRIMITIVES

func vgGetError() error {
	err := vgErrorType(C.vgGetError())
	if err == VG_NO_ERROR {
		return nil
	}
	fmt.Println("******", vgError(err), "********") // REMOVE THIS CODE
	return vgError(err)
}

func vgError(err vgErrorType) error {
	if err == VG_NO_ERROR {
		return nil
	}
	return errors.New(err.String())
}
