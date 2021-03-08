// +build dispmanx

package dispmanx

import (
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
	Resource    C.DISPMANX_RESOURCE_HANDLE_T
	PixFormat   C.VC_IMAGE_TYPE_T
	Transform   C.DISPMANX_TRANSFORM_T
)

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

func ResourceStride(w uint32) uint32 {
	return alignUp(w, 4)
}

func alignUp(value, alignment uint32) uint32 {
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
