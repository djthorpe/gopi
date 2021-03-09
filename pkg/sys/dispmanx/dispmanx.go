// +build dispmanx

package dispmanx

import (
	"fmt"
	"reflect"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: bcm_host
#include <bcm_host.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Display     C.DISPMANX_DISPLAY_HANDLE_T
	DisplayInfo C.DISPMANX_MODEINFO_T
	Element     C.DISPMANX_ELEMENT_HANDLE_T
	Resource    C.DISPMANX_RESOURCE_HANDLE_T
	Update      C.DISPMANX_UPDATE_HANDLE_T
	PixFormat   C.VC_IMAGE_TYPE_T
	Transform   C.DISPMANX_TRANSFORM_T
	Rect        C.VC_RECT_T
	Protection  C.DISPMANX_PROTECTION_T
	Alpha       C.VC_DISPMANX_ALPHA_T
	AlphaFlag   C.DISPMANX_FLAGS_ALPHA_T
	Clamp       C.DISPMANX_CLAMP_T
)

type Data struct {
	buf   uintptr
	cap   uint32
	data8 []byte
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DISPMANX_NO_ROTATE  Transform = C.DISPMANX_NO_ROTATE
	DISPMANX_ROTATE_90  Transform = C.DISPMANX_ROTATE_90
	DISPMANX_ROTATE_180 Transform = C.DISPMANX_ROTATE_180
	DISPMANX_ROTATE_270 Transform = C.DISPMANX_ROTATE_270
)

const (
	VC_IMAGE_RGB565       PixFormat = C.VC_IMAGE_RGB565
	VC_IMAGE_1BPP         PixFormat = C.VC_IMAGE_1BPP
	VC_IMAGE_YUV420       PixFormat = C.VC_IMAGE_YUV420
	VC_IMAGE_48BPP        PixFormat = C.VC_IMAGE_48BPP
	VC_IMAGE_RGB888       PixFormat = C.VC_IMAGE_RGB888
	VC_IMAGE_8BPP         PixFormat = C.VC_IMAGE_8BPP
	VC_IMAGE_4BPP         PixFormat = C.VC_IMAGE_4BPP        // 4bpp palettised image
	VC_IMAGE_3D32         PixFormat = C.VC_IMAGE_3D32        /* A separated format of 16 colour/light shorts followed by 16 z values */
	VC_IMAGE_3D32B        PixFormat = C.VC_IMAGE_3D32B       /* 16 colours followed by 16 z values */
	VC_IMAGE_3D32MAT      PixFormat = C.VC_IMAGE_3D32MAT     /* A separated format of 16 material/colour/light shorts followed by 16 z values */
	VC_IMAGE_RGB2X9       PixFormat = C.VC_IMAGE_RGB2X9      /* 32 bit format containing 18 bits of 6.6.6 RGB, 9 bits per short */
	VC_IMAGE_RGB666       PixFormat = C.VC_IMAGE_RGB666      /* 32-bit format holding 18 bits of 6.6.6 RGB */
	VC_IMAGE_RGBA32       PixFormat = C.VC_IMAGE_RGBA32      /* RGB888 with an alpha byte after each pixel */ /* xxx: isn't it BEFORE each pixel? */
	VC_IMAGE_YUV422       PixFormat = C.VC_IMAGE_YUV422      /* a line of Y (32-byte padded), a line of U (16-byte padded), and a line of V (16-byte padded) */
	VC_IMAGE_RGBA565      PixFormat = C.VC_IMAGE_RGBA565     /* RGB565 with a transparent patch */
	VC_IMAGE_RGBA16       PixFormat = C.VC_IMAGE_RGBA16      /* Compressed (4444) version of RGBA32 */
	VC_IMAGE_YUV_UV       PixFormat = C.VC_IMAGE_YUV_UV      /* VCIII codec format */
	VC_IMAGE_TF_RGBA32    PixFormat = C.VC_IMAGE_TF_RGBA32   /* VCIII T-format RGBA8888 */
	VC_IMAGE_TF_RGBX32    PixFormat = C.VC_IMAGE_TF_RGBX32   /* VCIII T-format RGBx8888 */
	VC_IMAGE_TF_FLOAT     PixFormat = C.VC_IMAGE_TF_FLOAT    /* VCIII T-format float */
	VC_IMAGE_TF_RGBA16    PixFormat = C.VC_IMAGE_TF_RGBA16   /* VCIII T-format RGBA4444 */
	VC_IMAGE_TF_RGBA5551  PixFormat = C.VC_IMAGE_TF_RGBA5551 /* VCIII T-format RGB5551 */
	VC_IMAGE_TF_RGB565    PixFormat = C.VC_IMAGE_TF_RGB565   /* VCIII T-format RGB565 */
	VC_IMAGE_TF_YA88      PixFormat = C.VC_IMAGE_TF_YA88     /* VCIII T-format 8-bit luma and 8-bit alpha */
	VC_IMAGE_TF_BYTE      PixFormat = C.VC_IMAGE_TF_BYTE     /* VCIII T-format 8 bit generic sample */
	VC_IMAGE_TF_PAL8      PixFormat = C.VC_IMAGE_TF_PAL8     /* VCIII T-format 8-bit palette */
	VC_IMAGE_TF_PAL4      PixFormat = C.VC_IMAGE_TF_PAL4     /* VCIII T-format 4-bit palette */
	VC_IMAGE_TF_ETC1      PixFormat = C.VC_IMAGE_TF_ETC1     /* VCIII T-format Ericsson Texture Compressed */
	VC_IMAGE_BGR888       PixFormat = C.VC_IMAGE_BGR888      /* RGB888 with R & B swapped */
	VC_IMAGE_BGR888_NP    PixFormat = C.VC_IMAGE_BGR888_NP   /* RGB888 with R & B swapped, but with no pitch, i.e. no padding after each row of pixels */
	VC_IMAGE_BAYER        PixFormat = C.VC_IMAGE_BAYER       /* Bayer image, extra defines which variant is being used */
	VC_IMAGE_CODEC        PixFormat = C.VC_IMAGE_CODEC       /* General wrapper for codec images e.g. JPEG from camera */
	VC_IMAGE_YUV_UV32     PixFormat = C.VC_IMAGE_YUV_UV32    /* VCIII codec format */
	VC_IMAGE_TF_Y8        PixFormat = C.VC_IMAGE_TF_Y8       /* VCIII T-format 8-bit luma */
	VC_IMAGE_TF_A8        PixFormat = C.VC_IMAGE_TF_A8       /* VCIII T-format 8-bit alpha */
	VC_IMAGE_TF_SHORT     PixFormat = C.VC_IMAGE_TF_SHORT    /* VCIII T-format 16-bit generic sample */
	VC_IMAGE_TF_1BPP      PixFormat = C.VC_IMAGE_TF_1BPP     /* VCIII T-format 1bpp black/white */
	VC_IMAGE_OPENGL       PixFormat = C.VC_IMAGE_OPENGL
	VC_IMAGE_YUV444I      PixFormat = C.VC_IMAGE_YUV444I      /* VCIII-B0 HVS YUV 4:4:4 interleaved samples */
	VC_IMAGE_YUV422PLANAR PixFormat = C.VC_IMAGE_YUV422PLANAR /* Y, U, & V planes separately (VC_IMAGE_YUV422 has them interleaved on a per line basis) */
	VC_IMAGE_ARGB8888     PixFormat = C.VC_IMAGE_ARGB8888     /* 32bpp with 8bit alpha at MS byte, with R, G, B (LS byte) */
	VC_IMAGE_XRGB8888     PixFormat = C.VC_IMAGE_XRGB8888     /* 32bpp with 8bit unused at MS byte, with R, G, B (LS byte) */
	VC_IMAGE_YUV422YUYV   PixFormat = C.VC_IMAGE_YUV422YUYV   /* interleaved 8 bit samples of Y, U, Y, V */
	VC_IMAGE_YUV422YVYU   PixFormat = C.VC_IMAGE_YUV422YVYU   /* interleaved 8 bit samples of Y, V, Y, U */
	VC_IMAGE_YUV422UYVY   PixFormat = C.VC_IMAGE_YUV422UYVY   /* interleaved 8 bit samples of U, Y, V, Y */
	VC_IMAGE_YUV422VYUY   PixFormat = C.VC_IMAGE_YUV422VYUY   /* interleaved 8 bit samples of V, Y, U, Y */
	VC_IMAGE_RGBX32       PixFormat = C.VC_IMAGE_RGBX32       /* 32bpp like RGBA32 but with unused alpha */
	VC_IMAGE_RGBX8888     PixFormat = C.VC_IMAGE_RGBX8888     /* 32bpp, corresponding to RGBA with unused alpha */
	VC_IMAGE_BGRX8888     PixFormat = C.VC_IMAGE_BGRX8888     /* 32bpp, corresponding to BGRA with unused alpha */
	VC_IMAGE_YUV420SP     PixFormat = C.VC_IMAGE_YUV420SP     /* Y as a plane, then UV byte interleaved in plane with with same pitch, half height */
	VC_IMAGE_YUV444PLANAR PixFormat = C.VC_IMAGE_YUV444PLANAR /* Y, U, & V planes separately 4:4:4 */
	VC_IMAGE_TF_U8        PixFormat = C.VC_IMAGE_TF_U8        /* T-format 8-bit U - same as TF_Y8 buf from U plane */
	VC_IMAGE_TF_V8        PixFormat = C.VC_IMAGE_TF_V8        /* T-format 8-bit U - same as TF_Y8 buf from V plane */
	VC_IMAGE_YUV420_16    PixFormat = C.VC_IMAGE_YUV420_16    /* YUV4:2:0 planar, 16bit values */
	VC_IMAGE_YUV_UV_16    PixFormat = C.VC_IMAGE_YUV_UV_16    /* YUV4:2:0 codec format, 16bit values */
	VC_IMAGE_YUV420_S     PixFormat = C.VC_IMAGE_YUV420_S     /* YUV4:2:0 with U,V in side-by-side format */
	VC_IMAGE_YUV10COL     PixFormat = C.VC_IMAGE_YUV10COL     /* 10-bit YUV 420 column image format */
	VC_IMAGE_RGBA1010102  PixFormat = C.VC_IMAGE_RGBA1010102  /* 32-bpp, 10-bit R/G/B, 2-bit Alpha */
)

const (
	DISPMANX_FLAGS_ALPHA_FROM_SOURCE          AlphaFlag = C.DISPMANX_FLAGS_ALPHA_FROM_SOURCE
	DISPMANX_FLAGS_ALPHA_FIXED_ALL_PIXELS     AlphaFlag = C.DISPMANX_FLAGS_ALPHA_FIXED_ALL_PIXELS
	DISPMANX_FLAGS_ALPHA_FIXED_NON_ZERO       AlphaFlag = C.DISPMANX_FLAGS_ALPHA_FIXED_NON_ZERO
	DISPMANX_FLAGS_ALPHA_FIXED_EXCEED_0X07    AlphaFlag = C.DISPMANX_FLAGS_ALPHA_FIXED_EXCEED_0X07
	DISPMANX_FLAGS_ALPHA_PREMULT              AlphaFlag = C.DISPMANX_FLAGS_ALPHA_PREMULT
	DISPMANX_FLAGS_ALPHA_MIX                  AlphaFlag = C.DISPMANX_FLAGS_ALPHA_MIX
	DISPMANX_FLAGS_ALPHA_DISCARD_LOWER_LAYERS AlphaFlag = C.DISPMANX_FLAGS_ALPHA_DISCARD_LOWER_LAYERS
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - RECT

func NewRect(x, y int32, w, h uint32) *Rect {
	r := new(Rect)
	C.vc_dispmanx_rect_set((*C.VC_RECT_T)(r), C.uint32_t(x), C.uint32_t(y), C.uint32_t(w), C.uint32_t(h))
	return r
}

func (r *Rect) Origin() (int32, int32) {
	return int32(r.x), int32(r.y)
}

func (r *Rect) Size() (uint32, uint32) {
	return uint32(r.width), uint32(r.height)
}

func (r *Rect) String() string {
	str := "<rect"
	x, y := r.Origin()
	str += fmt.Sprintf(" origin={%d,%d}", x, y)
	w, h := r.Size()
	str += fmt.Sprintf(" size={%d,%d}", w, h)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - UPDATES

func UpdateStart(priority int32) (Update, error) {
	if ctx := C.vc_dispmanx_update_start(C.int32_t(priority)); ctx == 0 {
		return 0, gopi.ErrBadParameter
	} else {
		return Update(ctx), nil
	}
}

/*
func UpdateSubmit(cb callback, userInfo uintptr) error {
	if err := C.vc_dispmanx_update_submit(C.DISPMANX_UPDATE_HANDLE_T(ctx)); err != 0 {
		return gopi.ErrUnexpectedResponse
	} else {
		return nil
	}
}
*/

func UpdateSubmitSync(ctx Update) error {
	if err := C.vc_dispmanx_update_submit_sync(C.DISPMANX_UPDATE_HANDLE_T(ctx)); err != 0 {
		return gopi.ErrUnexpectedResponse
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - ALPHA

func NewAlphaFromSource(opacity uint8) *Alpha {
	this := new(Alpha)
	this.flags = (C.DISPMANX_FLAGS_ALPHA_T)(DISPMANX_FLAGS_ALPHA_FIXED_ALL_PIXELS | DISPMANX_FLAGS_ALPHA_FROM_SOURCE)
	this.opacity = C.uint32_t(opacity)
	return this
}

func NewAlphaFixed(opacity uint8) *Alpha {
	this := new(Alpha)
	this.flags = (C.DISPMANX_FLAGS_ALPHA_T)(DISPMANX_FLAGS_ALPHA_FIXED_ALL_PIXELS)
	this.opacity = C.uint32_t(opacity)
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - ELEMENTS

func ElementAdd(ctx Update, display Display, layer uint16, destrect *Rect, src Resource, srcrect *Rect, protection Protection, alpha *Alpha, clamp *Clamp, transform Transform) (Element, error) {
	if element := C.vc_dispmanx_element_add(
		C.DISPMANX_UPDATE_HANDLE_T(ctx),
		C.DISPMANX_DISPLAY_HANDLE_T(display),
		C.int32_t(layer),
		(*C.VC_RECT_T)(destrect),
		C.DISPMANX_RESOURCE_HANDLE_T(src),
		(*C.VC_RECT_T)(srcrect),
		C.DISPMANX_PROTECTION_T(protection),
		(*C.VC_DISPMANX_ALPHA_T)(alpha),
		(*C.DISPMANX_CLAMP_T)(clamp),
		C.DISPMANX_TRANSFORM_T(transform),
	); element == 0 {
		return 0, gopi.ErrBadParameter
	} else {
		return Element(element), nil
	}
}

func ElementRemove(ctx Update, element Element) error {
	if err := C.vc_dispmanx_element_remove(C.DISPMANX_UPDATE_HANDLE_T(ctx), C.DISPMANX_ELEMENT_HANDLE_T(element)); err != 0 {
		return gopi.ErrBadParameter
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - RESOURCES

func ResourceCreate(f PixFormat, w, h uint32) (Resource, error) {
	var dummy C.uint32_t
	handle := C.vc_dispmanx_resource_create(C.VC_IMAGE_TYPE_T(f), C.uint32_t(w), C.uint32_t(h), (*C.uint32_t)(unsafe.Pointer(&dummy)))
	if handle == 0 {
		return 0, gopi.ErrBadParameter
	} else {
		return Resource(handle), nil
	}
}

func ResourceDelete(handle Resource) error {
	if C.vc_dispmanx_resource_delete(C.DISPMANX_RESOURCE_HANDLE_T(handle)) == 0 {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

func ResourceRead(handle Resource, r *Rect, dest uintptr, pitch uint32) error {
	if C.vc_dispmanx_resource_read_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), (*C.VC_RECT_T)(r), unsafe.Pointer(dest), C.uint32_t(pitch)) == 0 {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

func ResourceWrite(handle Resource, f PixFormat, pitch uint32, source uintptr, r *Rect) error {
	if C.vc_dispmanx_resource_write_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), C.VC_IMAGE_TYPE_T(f), C.int(pitch), unsafe.Pointer(source), (*C.VC_RECT_T)(r)) == 0 {
		return nil
	} else {
		return gopi.ErrBadParameter
	}
}

func ResourceStride(w uint32) uint32 {
	return AlignUp(w, 16)
}

func AlignUp(value, alignment uint32) uint32 {
	return ((value - 1) & ^(alignment - 1)) + alignment
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - DISPLAYS

func DisplayOpen(display uint32) (Display, error) {
	if handle := Display(C.vc_dispmanx_display_open(C.uint32_t(display))); handle != 0 {
		return handle, nil
	} else {
		return 0, gopi.ErrBadParameter
	}
}

func DisplayOpenOffscreen(handle Resource, transform Transform) (Display, error) {
	if handle := Display(C.vc_dispmanx_display_open_offscreen(C.DISPMANX_RESOURCE_HANDLE_T(handle), C.DISPMANX_TRANSFORM_T(transform))); handle != 0 {
		return handle, nil
	} else {
		return 0, gopi.ErrBadParameter
	}
}

func DisplayClose(display Display) error {
	if C.vc_dispmanx_display_close(C.DISPMANX_DISPLAY_HANDLE_T(display)) == 0 {
		return nil
	} else {
		return gopi.ErrUnexpectedResponse
	}
}

func DisplayGetInfo(display Display) (DisplayInfo, error) {
	var info DisplayInfo
	if C.vc_dispmanx_display_get_info(C.DISPMANX_DISPLAY_HANDLE_T(display), (*C.DISPMANX_MODEINFO_T)(&info)) == 0 {
		return info, nil
	} else {
		return info, gopi.ErrUnexpectedResponse
	}
}

func DisplaySnapshot(display Display) (Resource, error) {
	if info, err := DisplayGetInfo(display); err != nil {
		return 0, err
	} else if bitmap, err := ResourceCreate(info.PixFormat(), info.Width(), info.Height()); err != nil {
		return 0, err
	} else if err := C.vc_dispmanx_snapshot(C.DISPMANX_DISPLAY_HANDLE_T(display), C.DISPMANX_RESOURCE_HANDLE_T(bitmap), 0); err != 0 {
		return 0, gopi.ErrUnexpectedResponse
	} else {
		return bitmap, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - DISPLAY INFO

func (d DisplayInfo) Width() uint32 {
	return uint32(d.width)
}

func (d DisplayInfo) Height() uint32 {
	return uint32(d.height)
}

func (d DisplayInfo) Transform() Transform {
	return Transform(d.transform & 0x03)
}

func (d DisplayInfo) PixFormat() PixFormat {
	switch d.input_format {
	case C.VCOS_DISPLAY_INPUT_FORMAT_RGB888:
		return VC_IMAGE_RGB888
	case C.VCOS_DISPLAY_INPUT_FORMAT_RGB565:
		return VC_IMAGE_RGB565
	default:
		return 0
	}
}

func (d DisplayInfo) Num() uint32 {
	return uint32(d.display_num)
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: BYTE BUFFER

func NewData(size uint32) *Data {
	// Align to 4-byte boundaries
	this := new(Data)
	this.cap = AlignUp(size, 4)
	this.buf = uintptr(C.malloc(C.uint(this.cap)))
	if this.cap == 0 || this.buf == 0 {
		return nil
	}

	// Set data8 slice to point to allocated bytes
	h8 := (*reflect.SliceHeader)(unsafe.Pointer(&this.data8))
	h8.Data = this.buf
	h8.Len = int(size)
	h8.Cap = int(this.cap)

	return this
}

func (this *Data) Dispose() {
	C.free(unsafe.Pointer(this.buf))
	this.buf = 0
	this.cap = 0
}

func (this *Data) Stride() uint32 {
	return this.cap
}

func (this *Data) Bytes() []byte {
	return this.data8
}

func (this *Data) Ptr() uintptr {
	return this.buf
}

func (this *Data) PtrMinusOffset(offset uint32) uintptr {
	return this.buf - uintptr(offset)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t Transform) String() string {
	switch t & 0x04 {
	case DISPMANX_NO_ROTATE:
		return "DISPMANX_NO_ROTATE"
	case DISPMANX_ROTATE_90:
		return "DISPMANX_ROTATE_90"
	case DISPMANX_ROTATE_180:
		return "DISPMANX_ROTATE_180"
	case DISPMANX_ROTATE_270:
		return "DISPMANX_ROTATE_270"
	default:
		return "[?? Invalid Transform value]"
	}
}

func (f PixFormat) String() string {
	switch f {
	case VC_IMAGE_RGB565:
		return "VC_IMAGE_RGB565"
	case VC_IMAGE_1BPP:
		return "VC_IMAGE_1BPP"
	case VC_IMAGE_YUV420:
		return "VC_IMAGE_YUV420"
	case VC_IMAGE_48BPP:
		return "VC_IMAGE_48BPP"
	case VC_IMAGE_RGB888:
		return "VC_IMAGE_RGB888"
	case VC_IMAGE_8BPP:
		return "VC_IMAGE_8BPP"
	case VC_IMAGE_4BPP:
		return "VC_IMAGE_4BPP"
	case VC_IMAGE_3D32:
		return "VC_IMAGE_3D32"
	case VC_IMAGE_3D32B:
		return "VC_IMAGE_3D32B"
	case VC_IMAGE_3D32MAT:
		return "VC_IMAGE_3D32MAT"
	case VC_IMAGE_RGB2X9:
		return "VC_IMAGE_RGB2X9"
	case VC_IMAGE_RGB666:
		return "VC_IMAGE_RGB666"
	case VC_IMAGE_RGBA32:
		return "VC_IMAGE_RGBA32"
	case VC_IMAGE_YUV422:
		return "VC_IMAGE_YUV422"
	case VC_IMAGE_RGBA565:
		return "VC_IMAGE_RGBA565"
	case VC_IMAGE_RGBA16:
		return "VC_IMAGE_RGBA16"
	case VC_IMAGE_YUV_UV:
		return "VC_IMAGE_YUV_UV"
	case VC_IMAGE_TF_RGBA32:
		return "VC_IMAGE_TF_RGBA32"
	case VC_IMAGE_TF_RGBX32:
		return "VC_IMAGE_TF_RGBX32"
	case VC_IMAGE_TF_FLOAT:
		return "VC_IMAGE_TF_FLOAT"
	case VC_IMAGE_TF_RGBA16:
		return "VC_IMAGE_TF_RGBA16"
	case VC_IMAGE_TF_RGBA5551:
		return "VC_IMAGE_TF_RGBA5551"
	case VC_IMAGE_TF_RGB565:
		return "VC_IMAGE_TF_RGB565"
	case VC_IMAGE_TF_YA88:
		return "VC_IMAGE_TF_YA88"
	case VC_IMAGE_TF_BYTE:
		return "VC_IMAGE_TF_BYTE"
	case VC_IMAGE_TF_PAL8:
		return "VC_IMAGE_TF_PAL8"
	case VC_IMAGE_TF_PAL4:
		return "VC_IMAGE_TF_PAL4"
	case VC_IMAGE_TF_ETC1:
		return "VC_IMAGE_TF_ETC1"
	case VC_IMAGE_BGR888:
		return "VC_IMAGE_BGR888"
	case VC_IMAGE_BGR888_NP:
		return "VC_IMAGE_BGR888_NP"
	case VC_IMAGE_BAYER:
		return "VC_IMAGE_BAYER"
	case VC_IMAGE_CODEC:
		return "VC_IMAGE_CODEC"
	case VC_IMAGE_YUV_UV32:
		return "VC_IMAGE_YUV_UV32"
	case VC_IMAGE_TF_Y8:
		return "VC_IMAGE_TF_Y8"
	case VC_IMAGE_TF_A8:
		return "VC_IMAGE_TF_A8"
	case VC_IMAGE_TF_SHORT:
		return "VC_IMAGE_TF_SHORT"
	case VC_IMAGE_TF_1BPP:
		return "VC_IMAGE_TF_1BPP"
	case VC_IMAGE_OPENGL:
		return "VC_IMAGE_OPENGL"
	case VC_IMAGE_YUV444I:
		return "VC_IMAGE_YUV444I"
	case VC_IMAGE_YUV422PLANAR:
		return "VC_IMAGE_YUV422PLANAR"
	case VC_IMAGE_ARGB8888:
		return "VC_IMAGE_ARGB8888"
	case VC_IMAGE_XRGB8888:
		return "VC_IMAGE_XRGB8888"
	case VC_IMAGE_YUV422YUYV:
		return "VC_IMAGE_YUV422YUYV"
	case VC_IMAGE_YUV422YVYU:
		return "VC_IMAGE_YUV422YVYU"
	case VC_IMAGE_YUV422UYVY:
		return "VC_IMAGE_YUV422UYVY"
	case VC_IMAGE_YUV422VYUY:
		return "VC_IMAGE_YUV422VYUY"
	case VC_IMAGE_RGBX32:
		return "VC_IMAGE_RGBX32"
	case VC_IMAGE_RGBX8888:
		return "VC_IMAGE_RGBX8888"
	case VC_IMAGE_BGRX8888:
		return "VC_IMAGE_BGRX8888"
	case VC_IMAGE_YUV420SP:
		return "VC_IMAGE_YUV420SP"
	case VC_IMAGE_YUV444PLANAR:
		return "VC_IMAGE_YUV444PLANAR"
	case VC_IMAGE_TF_U8:
		return "VC_IMAGE_TF_U8"
	case VC_IMAGE_TF_V8:
		return "VC_IMAGE_TF_V8"
	case VC_IMAGE_YUV420_16:
		return "VC_IMAGE_YUV420_16"
	case VC_IMAGE_YUV_UV_16:
		return "VC_IMAGE_YUV_UV_16"
	case VC_IMAGE_YUV420_S:
		return "VC_IMAGE_YUV420_S"
	case VC_IMAGE_YUV10COL:
		return "VC_IMAGE_YUV10COL"
	case VC_IMAGE_RGBA1010102:
		return "VC_IMAGE_RGBA1010102"
	default:
		return "[?? Invalid PixFormat value]"
	}
}
