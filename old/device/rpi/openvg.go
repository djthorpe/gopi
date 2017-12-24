/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	khronos "github.com/djthorpe/gopi/khronos"
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
	log   gopi.Logger
	egl   *eglDriver
	lock  sync.Mutex
	path  map[C.VGPath]*vgPath
	paint map[C.VGPaint]*vgPaint
}

// Paths
type vgPath struct {
	handle C.VGPath
	log    gopi.Logger
}

// Paints
type vgPaint struct {
	handle       C.VGPaint
	line_width   float32
	join_style   khronos.VGStrokeJoinStyle
	cap_style    khronos.VGStrokeCapStyle
	dash_pattern []float32
	fill_rule    khronos.VGFillRule
	miter_limit  float32
	log          gopi.Logger
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
	VG_FILL_RULE               uint16  = 0x1101
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

// Path segment commands
const (
	VG_SEGMENT_CLOSE_PATH C.VGPathSegment = (0 << 1)
	VG_SEGMENT_MOVE_TO    C.VGPathSegment = (1 << 1)
	VG_SEGMENT_LINE_TO    C.VGPathSegment = (2 << 1)
	VG_SEGMENT_HLINE_TO   C.VGPathSegment = (3 << 1)
	VG_SEGMENT_VLINE_TO   C.VGPathSegment = (4 << 1)
	VG_SEGMENT_QUAD_TO    C.VGPathSegment = (5 << 1)
	VG_SEGMENT_CUBIC_TO   C.VGPathSegment = (6 << 1)
	VG_SEGMENT_SQUAD_TO   C.VGPathSegment = (7 << 1)
	VG_SEGMENT_SCUBIC_TO  C.VGPathSegment = (8 << 1)
	VG_SEGMENT_SCCWARC_TO C.VGPathSegment = (9 << 1)
	VG_SEGMENT_SCWARC_TO  C.VGPathSegment = (10 << 1)
	VG_SEGMENT_LCCWARC_TO C.VGPathSegment = (11 << 1)
	VG_SEGMENT_LCWARC_TO  C.VGPathSegment = (12 << 1)
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

// Open
func (config OpenVG) Open(log gopi.Logger) (gopi.Driver, error) {
	this := new(vgDriver)
	this.log = log
	this.log.Debug2("<rpi.OpenVG>Open")

	// EGL driver
	egl, ok := config.EGL.(*eglDriver)
	if egl == nil || ok != true {
		return nil, this.log.Error("Invalid configuration parameter: EGL")
	}
	this.egl = egl

	// Create maps for paths and paints
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
	return fmt.Sprintf("<rpi.OpenVG>{ egl=%v }", this.egl)
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
// FLUSH DRAWING

func (this *vgDriver) Do(surface khronos.EGLSurface, callback func() error) error {
	this.log.Debug2("<rpi.OpenVG>Do surface=%v", surface)

	// Lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Set Current surface
	if err := this.egl.SetCurrentContext(surface); err != nil {
		return err
	}

	// Set identity matrix
	C.vgLoadIdentity()

	// Draw
	if err := callback(); err != nil {
		return err
	}

	// Flush
	C.vgFlush()
	if err := this.egl.FlushSurface(surface); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// TRANSLATION FUNCTIONS

func (this *vgDriver) Translate(offset khronos.VGPoint) error {
	C.vgTranslate(C.VGfloat(offset.X), C.VGfloat(offset.Y))
	return vgGetError()
}

func (this *vgDriver) Scale(x, y float32) error {
	C.vgScale(C.VGfloat(x), C.VGfloat(y))
	return vgGetError()
}

func (this *vgDriver) Shear(x, y float32) error {
	C.vgShear(C.VGfloat(x), C.VGfloat(y))
	return vgGetError()
}

func (this *vgDriver) Rotate(r float32) error {
	C.vgRotate(C.VGfloat(r))
	return vgGetError()
}

func (this *vgDriver) LoadIdentity() error {
	C.vgLoadIdentity()
	return vgGetError()
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
		handle:     handle,
		log:        this.log,
		line_width: 1.0,
		cap_style:  khronos.VG_STYLE_CAP_NONE,
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

func (this *vgPaint) SetStrokeWidth(width float32) error {
	if this.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	if width < 0.0 {
		return vgError(VG_ILLEGAL_ARGUMENT_ERROR)
	}
	this.line_width = width
	return nil
}

func (this *vgPaint) SetStrokeStyle(join khronos.VGStrokeJoinStyle, cap khronos.VGStrokeCapStyle) error {
	if this.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	if cap != khronos.VG_STYLE_CAP_NONE {
		this.cap_style = cap
	}
	if join != khronos.VG_STYLE_JOIN_NONE {
		this.join_style = join
	}
	return nil
}

func (this *vgPaint) SetStrokeDash(pattern ...float32) error {
	if this.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	this.dash_pattern = pattern
	return nil
}

func (this *vgPaint) SetFillRule(style khronos.VGFillRule) error {
	if this.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	this.fill_rule = style
	return nil
}

func (this *vgPaint) SetMiterLimit(value float32) error {
	if this.handle == VG_PAINT_HANDLE_NONE {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	this.miter_limit = value
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PATH METHODS

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
	C.vgDestroyPath(obj.handle)
	delete(this.path, obj.handle)
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

	// STROKE
	stroke_obj, ok := stroke.(*vgPaint)
	if stroke_obj != nil && ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	if stroke_obj != nil && stroke_obj.handle != VG_PAINT_HANDLE_NONE {
		flags |= VG_STROKE_PATH
		if stroke_obj.cap_style != khronos.VG_STYLE_CAP_NONE {
			C.vgSeti(C.VGParamType(VG_STROKE_CAP_STYLE), C.VGint(stroke_obj.cap_style))
		}
		if stroke_obj.join_style != khronos.VG_STYLE_JOIN_NONE {
			C.vgSeti(C.VGParamType(VG_STROKE_JOIN_STYLE), C.VGint(stroke_obj.join_style))
		}
		if stroke_obj.miter_limit != 0.0 {
			C.vgSetf(C.VGParamType(VG_STROKE_MITER_LIMIT), C.VGfloat(stroke_obj.miter_limit))
		}
		if stroke_obj.dash_pattern != nil && len(stroke_obj.dash_pattern) > 0 {
			header := *(*reflect.SliceHeader)(unsafe.Pointer(&stroke_obj.dash_pattern))
			C.vgSetfv(C.VGParamType(VG_STROKE_DASH_PATTERN), C.VGint(header.Len), (*C.VGfloat)(unsafe.Pointer(header.Data)))
		} else {
			C.vgSetfv(C.VGParamType(VG_STROKE_DASH_PATTERN), 0, nil)
		}
		C.vgSetPaint(stroke_obj.handle, VG_STROKE_PATH)
		C.vgSetf(C.VGParamType(VG_STROKE_LINE_WIDTH), C.VGfloat(stroke_obj.line_width))
	}

	// FILL
	fill_obj, ok := fill.(*vgPaint)
	if fill_obj != nil && ok == false {
		return vgError(VG_BAD_HANDLE_ERROR)
	}
	if fill_obj != nil && fill_obj.handle != VG_PAINT_HANDLE_NONE {
		flags |= VG_FILL_PATH
		if fill_obj.fill_rule != khronos.VG_STYLE_FILL_NONE {
			C.vgSeti(C.VGParamType(VG_FILL_RULE), C.VGint(fill_obj.fill_rule))
		}
		C.vgSetPaint(fill_obj.handle, VG_FILL_PATH)

	}

	// DRAW PATH
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

// Close Path
func (this *vgPath) Close() error {
	cmd := C.VGubyte(VG_SEGMENT_CLOSE_PATH)
	C.vgAppendPathData(this.handle, C.VGint(1), &cmd, unsafe.Pointer(&khronos.VGZeroPoint))
	return vgGetError()
}

// Move To
func (this *vgPath) MoveTo(point khronos.VGPoint) error {
	cmd := C.VGubyte(VG_SEGMENT_MOVE_TO)
	C.vgAppendPathData(this.handle, C.VGint(1), &cmd, unsafe.Pointer(&point))
	return vgGetError()
}

// Line To
func (this *vgPath) LineTo(points ...khronos.VGPoint) error {
	// Create an array of LINE_TO commands
	cmd := make([]C.VGubyte, len(points))
	for i := range points {
		cmd[i] = C.VGubyte(VG_SEGMENT_LINE_TO)
	}
	// Append path data
	cmd_hdr := (*reflect.SliceHeader)(unsafe.Pointer(&cmd))
	points_hdr := (*reflect.SliceHeader)(unsafe.Pointer(&points))
	C.vgAppendPathData(this.handle, C.VGint(len(points)), (*C.VGubyte)(unsafe.Pointer(cmd_hdr.Data)), unsafe.Pointer(points_hdr.Data))
	return vgGetError()
}

// Quad To
func (this *vgPath) QuadTo(p1, p2 khronos.VGPoint) error {
	cmd := C.VGubyte(VG_SEGMENT_QUAD_TO)
	points := []float32{p1.X, p1.Y, p2.X, p2.Y}
	C.vgAppendPathData(this.handle, C.VGint(1), &cmd, unsafe.Pointer(&points[0]))
	return vgGetError()
}

// Cubic To
func (this *vgPath) CubicTo(p1, p2, p3 khronos.VGPoint) error {
	cmd := C.VGubyte(VG_SEGMENT_QUAD_TO)
	points := []float32{p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y}
	C.vgAppendPathData(this.handle, C.VGint(1), &cmd, unsafe.Pointer(&points[0]))
	return vgGetError()
}

// Clear surface to color
func (this *vgDriver) Clear(surface khronos.EGLSurface, color khronos.VGColor) error {
	size := surface.GetSize()
	C.vgSetfv(C.VGParamType(VG_CLEAR_COLOR), C.VGint(4), (*C.VGfloat)(unsafe.Pointer(&color)))
	C.vgClear(C.VGint(0), C.VGint(0), C.VGint(size.Width), C.VGint(size.Height))
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// ERROR HANDLING

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

// Get error from OpenVG
func vgGetError() error {
	err := vgErrorType(C.vgGetError())
	if err == VG_NO_ERROR {
		return nil
	}
	//	fmt.Println("******", vgError(err), "********") // REMOVE THIS CODE
	panic(vgError(err))
	return vgError(err)
}

// Return error based on an OpenVG error code
func vgError(err vgErrorType) error {
	if err == VG_NO_ERROR {
		return nil
	}
	return errors.New(err.String())
}
