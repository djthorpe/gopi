// +build openvg

package openvg

import (
	"reflect"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

/*
  #cgo pkg-config: brcmvg brcmglesv2
  #include <VG/openvg.h>
  #include <VG/vgu.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Error          C.VGErrorCode
	Param          C.VGParamType
	Handle         C.VGHandle
	MaskOperation  C.VGMaskOperation
	Path           C.VGPath
	PaintMode      C.VGbitfield
	MaskLayer      C.VGMaskLayer
	PathType       C.VGPathDatatype
	PathCapability C.VGPathCapabilities
	QueryType      C.VGHardwareQueryType
	QueryResult    C.VGHardwareQueryResult
	QueryString    C.VGStringID
)

// See /opt/vc/include/KHR/khrplatform.h
//typedef khronos_float_t  VGfloat;
//typedef khronos_int8_t   VGbyte;
//typedef khronos_uint8_t  VGubyte;
//typedef khronos_int16_t  VGshort;
//typedef khronos_int32_t  VGint;
//typedef khronos_uint32_t VGuint;
//typedef khronos_uint32_t VGbitfield;

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	VG_NO_ERROR                       Error = C.VG_NO_ERROR
	VG_BAD_HANDLE_ERROR               Error = C.VG_BAD_HANDLE_ERROR
	VG_ILLEGAL_ARGUMENT_ERROR         Error = C.VG_ILLEGAL_ARGUMENT_ERROR
	VG_OUT_OF_MEMORY_ERROR            Error = C.VG_OUT_OF_MEMORY_ERROR
	VG_PATH_CAPABILITY_ERROR          Error = C.VG_PATH_CAPABILITY_ERROR
	VG_UNSUPPORTED_IMAGE_FORMAT_ERROR Error = C.VG_UNSUPPORTED_IMAGE_FORMAT_ERROR
	VG_UNSUPPORTED_PATH_FORMAT_ERROR  Error = C.VG_UNSUPPORTED_PATH_FORMAT_ERROR
	VG_IMAGE_IN_USE_ERROR             Error = C.VG_IMAGE_IN_USE_ERROR
	VG_NO_CONTEXT_ERROR               Error = C.VG_NO_CONTEXT_ERROR
)

const (
	VG_MATRIX_MODE                 Param = C.VG_MATRIX_MODE /* Mode settings */
	VG_FILL_RULE                   Param = C.VG_FILL_RULE
	VG_IMAGE_QUALITY               Param = C.VG_IMAGE_QUALITY
	VG_RENDERING_QUALITY           Param = C.VG_RENDERING_QUALITY
	VG_BLEND_MODE                  Param = C.VG_BLEND_MODE
	VG_IMAGE_MODE                  Param = C.VG_IMAGE_MODE
	VG_SCISSOR_RECTS               Param = C.VG_SCISSOR_RECTS   /* Scissoring rectangles */
	VG_COLOR_TRANSFORM             Param = C.VG_COLOR_TRANSFORM /* Color Transformation */
	VG_COLOR_TRANSFORM_VALUES      Param = C.VG_COLOR_TRANSFORM_VALUES
	VG_STROKE_LINE_WIDTH           Param = C.VG_STROKE_LINE_WIDTH /* Stroke parameters */
	VG_STROKE_CAP_STYLE            Param = C.VG_STROKE_CAP_STYLE
	VG_STROKE_JOIN_STYLE           Param = C.VG_STROKE_JOIN_STYLE
	VG_STROKE_MITER_LIMIT          Param = C.VG_STROKE_MITER_LIMIT
	VG_STROKE_DASH_PATTERN         Param = C.VG_STROKE_DASH_PATTERN
	VG_STROKE_DASH_PHASE           Param = C.VG_STROKE_DASH_PHASE
	VG_STROKE_DASH_PHASE_RESET     Param = C.VG_STROKE_DASH_PHASE_RESET
	VG_TILE_FILL_COLOR             Param = C.VG_TILE_FILL_COLOR /* Edge fill color for VG_TILE_FILL tiling mode */
	VG_CLEAR_COLOR                 Param = C.VG_CLEAR_COLOR     /* Color for vgClear */
	VG_GLYPH_ORIGIN                Param = C.VG_GLYPH_ORIGIN    /* Glyph origin */
	VG_MASKING                     Param = C.VG_MASKING         /* Enable/disable alpha masking and scissoring */
	VG_SCISSORING                  Param = C.VG_SCISSORING
	VG_PIXEL_LAYOUT                Param = C.VG_PIXEL_LAYOUT /* Pixel layout information */
	VG_SCREEN_LAYOUT               Param = C.VG_SCREEN_LAYOUT
	VG_FILTER_FORMAT_LINEAR        Param = C.VG_FILTER_FORMAT_LINEAR /* Source format selection for image filters */
	VG_FILTER_FORMAT_PREMULTIPLIED Param = C.VG_FILTER_FORMAT_PREMULTIPLIED
	VG_FILTER_CHANNEL_MASK         Param = C.VG_FILTER_CHANNEL_MASK /* Destination write enable mask for image filters */
	VG_MAX_SCISSOR_RECTS           Param = C.VG_MAX_SCISSOR_RECTS   /* Implementation limits (read-only) */
	VG_MAX_DASH_COUNT              Param = C.VG_MAX_DASH_COUNT
	VG_MAX_KERNEL_SIZE             Param = C.VG_MAX_KERNEL_SIZE
	VG_MAX_SEPARABLE_KERNEL_SIZE   Param = C.VG_MAX_SEPARABLE_KERNEL_SIZE
	VG_MAX_COLOR_RAMP_STOPS        Param = C.VG_MAX_COLOR_RAMP_STOPS
	VG_MAX_IMAGE_WIDTH             Param = C.VG_MAX_IMAGE_WIDTH
	VG_MAX_IMAGE_HEIGHT            Param = C.VG_MAX_IMAGE_HEIGHT
	VG_MAX_IMAGE_PIXELS            Param = C.VG_MAX_IMAGE_PIXELS
	VG_MAX_IMAGE_BYTES             Param = C.VG_MAX_IMAGE_BYTES
	VG_MAX_FLOAT                   Param = C.VG_MAX_FLOAT
	VG_MAX_GAUSSIAN_STD_DEVIATION  Param = C.VG_MAX_GAUSSIAN_STD_DEVIATION
)

const (
	VG_CLEAR_MASK     MaskOperation = C.VG_CLEAR_MASK
	VG_FILL_MASK      MaskOperation = C.VG_FILL_MASK
	VG_SET_MASK       MaskOperation = C.VG_SET_MASK
	VG_UNION_MASK     MaskOperation = C.VG_UNION_MASK
	VG_INTERSECT_MASK MaskOperation = C.VG_INTERSECT_MASK
	VG_SUBTRACT_MASK  MaskOperation = C.VG_SUBTRACT_MASK
)

const (
	VG_STROKE_PATH PaintMode = C.VG_STROKE_PATH
	VG_FILL_PATH   PaintMode = C.VG_FILL_PATH
)

const (
	VG_PATH_DATATYPE_S_8  PathType = C.VG_PATH_DATATYPE_S_8
	VG_PATH_DATATYPE_S_16 PathType = C.VG_PATH_DATATYPE_S_16
	VG_PATH_DATATYPE_S_32 PathType = C.VG_PATH_DATATYPE_S_32
	VG_PATH_DATATYPE_F    PathType = C.VG_PATH_DATATYPE_F
)

const (
	VG_PATH_CAPABILITY_APPEND_FROM             PathCapability = C.VG_PATH_CAPABILITY_APPEND_FROM
	VG_PATH_CAPABILITY_APPEND_TO               PathCapability = C.VG_PATH_CAPABILITY_APPEND_TO
	VG_PATH_CAPABILITY_MODIFY                  PathCapability = C.VG_PATH_CAPABILITY_MODIFY
	VG_PATH_CAPABILITY_TRANSFORM_FROM          PathCapability = C.VG_PATH_CAPABILITY_TRANSFORM_FROM
	VG_PATH_CAPABILITY_TRANSFORM_TO            PathCapability = C.VG_PATH_CAPABILITY_TRANSFORM_TO
	VG_PATH_CAPABILITY_INTERPOLATE_FROM        PathCapability = C.VG_PATH_CAPABILITY_INTERPOLATE_FROM
	VG_PATH_CAPABILITY_INTERPOLATE_TO          PathCapability = C.VG_PATH_CAPABILITY_INTERPOLATE_TO
	VG_PATH_CAPABILITY_PATH_LENGTH             PathCapability = C.VG_PATH_CAPABILITY_PATH_LENGTH
	VG_PATH_CAPABILITY_POINT_ALONG_PATH        PathCapability = C.VG_PATH_CAPABILITY_POINT_ALONG_PATH
	VG_PATH_CAPABILITY_TANGENT_ALONG_PATH      PathCapability = C.VG_PATH_CAPABILITY_TANGENT_ALONG_PATH
	VG_PATH_CAPABILITY_PATH_BOUNDS             PathCapability = C.VG_PATH_CAPABILITY_PATH_BOUNDS
	VG_PATH_CAPABILITY_PATH_TRANSFORMED_BOUNDS PathCapability = C.VG_PATH_CAPABILITY_PATH_TRANSFORMED_BOUNDS
	VG_PATH_CAPABILITY_ALL                     PathCapability = C.VG_PATH_CAPABILITY_ALL
)

const (
	VG_IMAGE_FORMAT_QUERY  QueryType = C.VG_IMAGE_FORMAT_QUERY
	VG_PATH_DATATYPE_QUERY QueryType = C.VG_PATH_DATATYPE_QUERY
)

const (
	VG_HARDWARE_ACCELERATED   QueryResult = C.VG_HARDWARE_ACCELERATED
	VG_HARDWARE_UNACCELERATED QueryResult = C.VG_HARDWARE_UNACCELERATED
)

const (
	VG_VENDOR     QueryString = C.VG_VENDOR
	VG_RENDERER   QueryString = C.VG_RENDERER
	VG_VERSION    QueryString = C.VG_VERSION
	VG_EXTENSIONS QueryString = C.VG_EXTENSIONS
)

const (
	VG_MATRIX_SIZE = 9
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func Flush() error {
	C.vgFlush()
	return vgGetError()
}

func Finish() error {
	C.vgFinish()
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - GET AND SET

func Setf(p Param, f float32) error {
	C.vgSetf(C.VGParamType(p), C.VGfloat(f))
	return vgGetError()
}

func Seti(p Param, i int32) error {
	C.vgSeti(C.VGParamType(p), C.VGint(i))
	return vgGetError()
}

func Setfv(p Param, f []float32) error {
	count := C.VGint(len(f))
	C.vgSetfv(C.VGParamType(p), count, (*C.VGfloat)(unsafe.Pointer(&f[0])))
	return vgGetError()
}

func Setiv(p Param, i []int32) error {
	count := C.VGint(len(i))
	C.vgSetiv(C.VGParamType(p), count, (*C.VGint)(unsafe.Pointer(&i[0])))
	return vgGetError()
}

func Getf(p Param) (float32, error) {
	return float32(C.vgGetf(C.VGParamType(p))), vgGetError()
}

func Geti(p Param) (int32, error) {
	return int32(C.vgGeti(C.VGParamType(p))), vgGetError()
}

func GetVectorSize(p Param) (int32, error) {
	return int32(C.vgGetVectorSize(C.VGParamType(p))), vgGetError()
}

func Getfv(p Param, count int32) ([]float32, error) {
	var result []float32

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(count)
	sliceHeader.Len = int(count)

	// Get values
	C.vgGetfv(C.VGParamType(p), C.VGint(count), (*C.VGfloat)(unsafe.Pointer(&sliceHeader.Data)))
	if err := vgGetError(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func Getfi(p Param, count int32) ([]int32, error) {
	var result []int32

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(count)
	sliceHeader.Len = int(count)

	// Get values
	C.vgGetiv(C.VGParamType(p), C.VGint(count), (*C.VGint)(unsafe.Pointer(&sliceHeader.Data)))
	if err := vgGetError(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func SetParameterf(h Handle, p Param, f float32) error {
	C.vgSetParameterf(C.VGHandle(h), C.VGint(p), C.VGfloat(f))
	return vgGetError()
}

func SetParameteri(h Handle, p Param, i int32) error {
	C.vgSetParameteri(C.VGHandle(h), C.VGint(p), C.VGint(i))
	return vgGetError()
}

func SetParameterfv(h Handle, p Param, f []float32) error {
	count := C.VGint(len(f))
	C.vgSetParameterfv(C.VGHandle(h), C.VGint(p), count, (*C.VGfloat)(unsafe.Pointer(&f[0])))
	return vgGetError()
}

func SetParameterfi(h Handle, p Param, i []int32) error {
	count := C.VGint(len(i))
	C.vgSetParameteriv(C.VGHandle(h), C.VGint(p), count, (*C.VGint)(unsafe.Pointer(&i[0])))
	return vgGetError()
}

func GetParameterf(h Handle, p Param) (float32, error) {
	f := C.vgGetParameterf(C.VGHandle(h), C.VGint(p))
	return float32(f), vgGetError()
}

func GetParameteri(h Handle, p Param) (int32, error) {
	i := C.vgGetParameteri(C.VGHandle(h), C.VGint(p))
	return int32(i), vgGetError()
}

func GetParameterVectorSize(h Handle, p Param) (int32, error) {
	i := C.vgGetParameterVectorSize(C.VGHandle(h), C.VGint(p))
	return int32(i), vgGetError()
}

func GetParameterfv(h Handle, p Param, count int32) ([]float32, error) {
	var result []float32

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(count)
	sliceHeader.Len = int(count)

	// Get values
	C.vgGetParameterfv(C.VGHandle(h), C.VGint(p), C.VGint(count), (*C.VGfloat)(unsafe.Pointer(&sliceHeader.Data)))
	if err := vgGetError(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func GetParameteriv(h Handle, p Param, count int32) ([]int32, error) {
	var result []int32

	// Make a fake slice
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&result)))
	sliceHeader.Cap = int(count)
	sliceHeader.Len = int(count)

	// Get values
	C.vgGetParameteriv(C.VGHandle(h), C.VGint(p), C.VGint(count), (*C.VGint)(unsafe.Pointer(&sliceHeader.Data)))
	if err := vgGetError(); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - MATRIX MANIPULATION

func LoadIdentity() error {
	C.vgLoadIdentity()
	return vgGetError()
}

func LoadMatrix(f []float32) error {
	if len(f) != VG_MATRIX_SIZE {
		return VG_ILLEGAL_ARGUMENT_ERROR
	}
	C.vgLoadMatrix((*C.VGfloat)(unsafe.Pointer(&f[0])))
	return vgGetError()
}

func GetMatrix() ([]float32, error) {
	f := make([]float32, VG_MATRIX_SIZE)
	C.vgGetMatrix((*C.VGfloat)(unsafe.Pointer(&f[0])))
	if err := vgGetError(); err != nil {
		return nil, err
	} else {
		return f, nil
	}
}

func MultMatrix(f []float32) error {
	if len(f) != VG_MATRIX_SIZE {
		return VG_ILLEGAL_ARGUMENT_ERROR
	}
	C.vgMultMatrix((*C.VGfloat)(unsafe.Pointer(&f[0])))
	return vgGetError()
}

func Translate(x, y float32) error {
	C.vgTranslate(C.VGfloat(x), C.VGfloat(y))
	return vgGetError()
}

func Scale(x, y float32) error {
	C.vgScale(C.VGfloat(x), C.VGfloat(y))
	return vgGetError()
}

func Shear(x, y float32) error {
	C.vgShear(C.VGfloat(x), C.VGfloat(y))
	return vgGetError()
}

func Rotate(angle float32) error {
	C.vgRotate(C.VGfloat(angle))
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// MASKING AND CLEARING

func Mask(handle Handle, o MaskOperation, x, y, w, h int32) error {
	C.vgMask(C.VGHandle(handle), C.VGMaskOperation(o), C.VGint(x), C.VGint(y), C.VGint(w), C.VGint(h))
	return vgGetError()
}

func RenderToMask(p Path, m PaintMode, o MaskOperation) error {
	C.vgRenderToMask(C.VGPath(p), C.VGbitfield(m), C.VGMaskOperation(o))
	return vgGetError()
}

func CreateMaskLayer(w, h int32) (MaskLayer, error) {
	return MaskLayer(C.vgCreateMaskLayer(C.VGint(w), C.VGint(h))), vgGetError()
}

func DestroyMaskLayer(handle MaskLayer) error {
	C.vgDestroyMaskLayer(C.VGMaskLayer(handle))
	return vgGetError()
}

func FillMaskLayer(handle MaskLayer, x, y, w, h int32, f float32) error {
	C.vgFillMaskLayer(C.VGMaskLayer(handle), C.VGint(x), C.VGint(y), C.VGint(w), C.VGint(h), C.VGfloat(f))
	return vgGetError()
}

func CopyMask(handle MaskLayer, dx, dy, sx, sy, w, h int32) error {
	C.vgCopyMask(C.VGMaskLayer(handle), C.VGint(dx), C.VGint(dy), C.VGint(sx), C.VGint(sy), C.VGint(w), C.VGint(h))
	return vgGetError()
}

func Clear(x, y, w, h int32) error {
	C.vgClear(C.VGint(x), C.VGint(y), C.VGint(w), C.VGint(h))
	return vgGetError()
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - PATHS

func CreatePath(t PathType, scale, bias float32, segments, coords float32, cap PathCapability) (Path, error) {
	path := C.vgCreatePath(C.VG_PATH_FORMAT_STANDARD, C.VGPathDatatype(t), C.VGfloat(scale), C.VGfloat(bias), C.VGint(segments), C.VGint(coords), C.VGbitfield(cap))
	return Path(path), vgGetError()
}

func ClearPath(handle Path, cap PathCapability) error {
	C.vgClearPath(C.VGPath(handle), C.VGbitfield(cap))
	return vgGetError()
}

func DestroyPath(handle Path) error {
	C.vgDestroyPath(C.VGPath(handle))
	return vgGetError()
}

func RemovePathCapabilities(handle Path, cap PathCapability) error {
	C.vgRemovePathCapabilities(C.VGPath(handle), C.VGbitfield(cap))
	return vgGetError()
}

func GetPathCapabilities(handle Path) (PathCapability, error) {
	return PathCapability(C.vgGetPathCapabilities(C.VGPath(handle))), vgGetError()
}

func AppendPath(dest, source Path) error {
	C.vgAppendPath(C.VGPath(dest), C.VGPath(source))
	return vgGetError()
}

func TransformPath(dest, source Path) error {
	C.vgTransformPath(C.VGPath(dest), C.VGPath(source))
	return vgGetError()
}

func InterpolatePath(dest, start, end Path, amount float32) (bool, error) {
	result := C.vgInterpolatePath(C.VGPath(dest), C.VGPath(start), C.VGPath(end), C.VGfloat(amount))
	return result != 0, vgGetError()
}

func PathLength(handle Path, start, count int32) (float32, error) {
	result := C.vgPathLength(C.VGPath(handle), C.VGint(start), C.VGint(count))
	return float32(result), vgGetError()
}

/*
VG_API_CALL void VG_API_ENTRY vgAppendPathData(VGPath dstPath,
    VGint numSegments,
    const VGubyte * pathSegments,
    const void * pathData) VG_API_EXIT;
VG_API_CALL void VG_API_ENTRY vgModifyPathCoords(VGPath dstPath, VGint startIndex,
      VGint numSegments,
      const void * pathData) VG_API_EXIT;
VG_API_CALL void VG_API_ENTRY vgPointAlongPath(VGPath path,
    VGint startSegment, VGint numSegments,
    VGfloat distance,
    VGfloat * x, VGfloat * y,
    VGfloat * tangentX, VGfloat * tangentY) VG_API_EXIT;
VG_API_CALL void VG_API_ENTRY vgPathBounds(VGPath path,
VGfloat * minX, VGfloat * minY,
VGfloat * width, VGfloat * height) VG_API_EXIT;
VG_API_CALL void VG_API_ENTRY vgPathTransformedBounds(VGPath path,
           VGfloat * minX, VGfloat * minY,
           VGfloat * width, VGfloat * height) VG_API_EXIT;
           VG_API_CALL void VG_API_ENTRY vgDrawPath(VGPath path, VGbitfield paintModes) VG_API_EXIT;
*/

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - HARDWARE QUERIES

func HardwareQuery(query QueryType, setting int32) (QueryResult, error) {
	result := QueryResult(C.vgHardwareQuery(C.VGHardwareQueryType(query), C.VGint(setting)))
	return result, vgGetError()
}

func GetString(query QueryString) (string, error) {
	str := (*C.char)(C.vgGetString(C.VGStringID(query)))
	if str == nil {
		return "", VG_NO_CONTEXT_ERROR
	} else if err := vgGetError(); err != nil {
		return "", err
	} else {
		return C.GoString(str), nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// ERRORS

func vgGetError() error {
	if err := Error(C.vgGetError()); err == VG_NO_ERROR {
		return nil
	} else {
		return err
	}
}

func (e Error) Error() string {
	return e.String()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e Error) String() string {
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
	default:
		return "[?? Invalid VGError value]"
	}
}
