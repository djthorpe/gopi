package rpi

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host
    #cgo LDFLAGS:  -L/opt/vc/lib -lbcm_host
	#include "vc_dispmanx.h"
    #include "bcm_host.h"
*/
import "C"

import (
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

type Size struct {
	Width  uint32
	Height uint32
}

type (
	DisplayHandle  uint32
	UpdateHandle   uint32
	ElementHandle  uint32
	ResourceHandle uint32
	UpdatePriority int32
	Transform      uint32
	InputFormat    uint32
	ImageType      int
)

type VideoCore struct {
	display uint16
	size    Size
	handle  DisplayHandle
}

type ModeInfo struct {
	Width, Height int32
	Transform     Transform
	InputFormat   InputFormat
	Handle        DisplayHandle
}

type Color struct {
	Red, Green, Blue uint8
}

type Resource struct {
	handle ResourceHandle
}

////////////////////////////////////////////////////////////////////////////////

const (
	/* Success and failure conditions */
	DISPMANX_SUCCESS   = 0
	DISPMANX_INVALID   = -1
	DISPMANX_NO_HANDLE = 0
)

const (
	/* ImageType */
	_ ImageType = iota
	VC_IMAGE_RGB565
	VC_IMAGE_1BPP
	VC_IMAGE_YUV420
	VC_IMAGE_48BPP
	VC_IMAGE_RGB888
	VC_IMAGE_8BPP
	VC_IMAGE_4BPP          // 4bpp palettised image
	VC_IMAGE_3D32          /* A separated format of 16 colour/light shorts followed by 16 z values */
	VC_IMAGE_3D32B         /* 16 colours followed by 16 z values */
	VC_IMAGE_3D32MAT       /* A separated format of 16 material/colour/light shorts followed by 16 z values */
	VC_IMAGE_RGB2X9        /* 32 bit format containing 18 bits of 6.6.6 RGB, 9 bits per short */
	VC_IMAGE_RGB666        /* 32-bit format holding 18 bits of 6.6.6 RGB */
	VC_IMAGE_PAL4_OBSOLETE // 4bpp palettised image with embedded palette
	VC_IMAGE_PAL8_OBSOLETE // 8bpp palettised image with embedded palette
	VC_IMAGE_RGBA32        /* RGB888 with an alpha byte after each pixel */ /* xxx: isn't it BEFORE each pixel? */
	VC_IMAGE_YUV422        /* a line of Y (32-byte padded), a line of U (16-byte padded), and a line of V (16-byte padded) */
	VC_IMAGE_RGBA565       /* RGB565 with a transparent patch */
	VC_IMAGE_RGBA16        /* Compressed (4444) version of RGBA32 */
	VC_IMAGE_YUV_UV        /* VCIII codec format */
	VC_IMAGE_TF_RGBA32     /* VCIII T-format RGBA8888 */
	VC_IMAGE_TF_RGBX32     /* VCIII T-format RGBx8888 */
	VC_IMAGE_TF_FLOAT      /* VCIII T-format float */
	VC_IMAGE_TF_RGBA16     /* VCIII T-format RGBA4444 */
	VC_IMAGE_TF_RGBA5551   /* VCIII T-format RGB5551 */
	VC_IMAGE_TF_RGB565     /* VCIII T-format RGB565 */
	VC_IMAGE_TF_YA88       /* VCIII T-format 8-bit luma and 8-bit alpha */
	VC_IMAGE_TF_BYTE       /* VCIII T-format 8 bit generic sample */
	VC_IMAGE_TF_PAL8       /* VCIII T-format 8-bit palette */
	VC_IMAGE_TF_PAL4       /* VCIII T-format 4-bit palette */
	VC_IMAGE_TF_ETC1       /* VCIII T-format Ericsson Texture Compressed */
	VC_IMAGE_BGR888        /* RGB888 with R & B swapped */
	VC_IMAGE_BGR888_NP     /* RGB888 with R & B swapped, but with no pitch, i.e. no padding after each row of pixels */
	VC_IMAGE_BAYER         /* Bayer image, extra defines which variant is being used */
	VC_IMAGE_CODEC         /* General wrapper for codec images e.g. JPEG from camera */
	VC_IMAGE_YUV_UV32      /* VCIII codec format */
	VC_IMAGE_TF_Y8         /* VCIII T-format 8-bit luma */
	VC_IMAGE_TF_A8         /* VCIII T-format 8-bit alpha */
	VC_IMAGE_TF_SHORT      /* VCIII T-format 16-bit generic sample */
	VC_IMAGE_TF_1BPP       /* VCIII T-format 1bpp black/white */
	VC_IMAGE_OPENGL
	VC_IMAGE_YUV444I      /* VCIII-B0 HVS YUV 4:4:4 interleaved samples */
	VC_IMAGE_YUV422PLANAR /* Y, U, & V planes separately (VC_IMAGE_YUV422 has them interleaved on a per line basis) */
	VC_IMAGE_ARGB8888     /* 32bpp with 8bit alpha at MS byte, with R, G, B (LS byte) */
	VC_IMAGE_XRGB8888     /* 32bpp with 8bit unused at MS byte, with R, G, B (LS byte) */
	VC_IMAGE_YUV422YUYV   /* interleaved 8 bit samples of Y, U, Y, V */
	VC_IMAGE_YUV422YVYU   /* interleaved 8 bit samples of Y, V, Y, U */
	VC_IMAGE_YUV422UYVY   /* interleaved 8 bit samples of U, Y, V, Y */
	VC_IMAGE_YUV422VYUY   /* interleaved 8 bit samples of V, Y, U, Y */
	VC_IMAGE_RGBX32       /* 32bpp like RGBA32 but with unused alpha */
	VC_IMAGE_RGBX8888     /* 32bpp, corresponding to RGBA with unused alpha */
	VC_IMAGE_BGRX8888     /* 32bpp, corresponding to BGRA with unused alpha */
	VC_IMAGE_YUV420SP     /* Y as a plane, then UV byte interleaved in plane with with same pitch, half height */
	VC_IMAGE_YUV444PLANAR /* Y, U, & V planes separately 4:4:4 */
	VC_IMAGE_TF_U8        /* T-format 8-bit U - same as TF_Y8 buf from U plane */
	VC_IMAGE_TF_V8        /* T-format 8-bit U - same as TF_Y8 buf from V plane */
)

////////////////////////////////////////////////////////////////////////////////

// Create new VideoCore object, returns error if not possible
func (rpi *RaspberryPi) NewVideoCore(display uint16) (*VideoCore, error) {

	// create object
	this := new(VideoCore)

	// get the display size
	this.display = display
	width, height, err := graphicsGetDisplaySize(display)
	if err != nil {
		return nil, err
	}
	this.size.Width = width
	this.size.Height = height

	// open the display
	handle, err := displayOpen(uint32(display))
	if err != nil {
		return nil, err
	}
	this.handle = handle

	// success
	return this, nil
}

// Close unmaps GPIO memory
func (this *VideoCore) Close() error {
	err := displayClose(this.handle)
	return err
}

////////////////////////////////////////////////////////////////////////////////

func (this *VideoCore) GetDisplay() uint16 {
	return this.display
}

func (this *VideoCore) GetSize() Size {
	return this.size
}

func (this *VideoCore) GetModeInfo() (ModeInfo, error) {
	var info ModeInfo
	err := displayGetInfo(this.handle, &info)
	return info, err
}

func (this *VideoCore) SetBackgroundColor(handle UpdateHandle, color Color) error {
	return displaySetBackground(handle, this.handle, color.Red, color.Green, color.Blue)
}

func (this *VideoCore) UpdateBegin() (UpdateHandle, error) {
	return updateStart(UpdatePriority(0))
}

func (this *VideoCore) UpdateSubmit(handle UpdateHandle) error {
	return updateSubmitSync(handle)
}

func (this *VideoCore) CreateResource(format ImageType, size Size) (*Resource, error) {
	var buffer uint32
	handle, err := resourceCreate(format, size.Width, size.Height, &buffer)
	if err != nil {
		return nil, err
	}

	resource := new(Resource)
	resource.handle = handle

	return resource, nil
}

func (this *VideoCore) DeleteResource(resource *Resource) error {
	err := resourceDelete(resource.handle)
	if err != nil {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods - display information

func graphicsGetDisplaySize(display uint16) (uint32, uint32, error) {
	var w, h uint32
	success := C.graphics_get_display_size((C.uint16_t)(display), (*C.uint32_t)(&w), (*C.uint32_t)(&h))
	if success != 0 {
		return 0, 0, ErrorDisplay
	}
	return w, h, nil
}

func displayOpen(display uint32) (DisplayHandle, error) {
	handle := DisplayHandle(C.vc_dispmanx_display_open(C.uint32_t(display)))
	if handle >= DisplayHandle(0) {
		return handle, nil
	} else {
		return DisplayHandle(0), ErrorDisplay
	}
}

func displayClose(display DisplayHandle) error {
	if C.vc_dispmanx_display_close(C.DISPMANX_DISPLAY_HANDLE_T(display)) != DISPMANX_SUCCESS {
		return ErrorDisplay
	}
	return nil
}

func displayGetInfo(display DisplayHandle, info *ModeInfo) error {
	if C.vc_dispmanx_display_get_info(C.DISPMANX_DISPLAY_HANDLE_T(display), (*C.DISPMANX_MODEINFO_T)(unsafe.Pointer(info))) != DISPMANX_SUCCESS {
		return ErrorDisplay
	}
	return nil
}

func displaySetBackground(update UpdateHandle, display DisplayHandle, r, g, b uint8) error {
	if C.vc_dispmanx_display_set_background(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_DISPLAY_HANDLE_T(display), C.uint8_t(r), C.uint8_t(g), C.uint8_t(b)) != DISPMANX_SUCCESS {
		return ErrorDisplay
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods - updates

func updateStart(priority UpdatePriority) (UpdateHandle, error) {
	handle := C.vc_dispmanx_update_start(C.int32_t(priority))
	if handle == DISPMANX_NO_HANDLE {
		return UpdateHandle(0), ErrorUpdate
	}
	return UpdateHandle(handle), nil
}

func updateSubmitSync(handle UpdateHandle) error {
	if C.vc_dispmanx_update_submit_sync(C.DISPMANX_UPDATE_HANDLE_T(handle)) != DISPMANX_SUCCESS {
		return ErrorUpdate
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods - elements

func elementAdd(handle UpdateHandle,display DisplayHandle,layer Layer,dest_rect *Rectangle,src_resource ResourceHandle,src_rect *Rectangle,protection Protection,alpha Alpha,clamp Clamp,transform Transform) (ElementHandle,error) {
	// TODO
}

func elementRemove(handle UpdateHandle,element ElementHandle) error {
	// TODO
}

func elementModified(handle UpdateHandle,element ElementHandle,rect Rectangle) error {
	// TODO
}

func elementChangeLayer(handle UpdateHandle,element ElementHandle,layer Layer) error {
	// TODO
}

func elementChangeSource(handle UpdateHandle,element ElementHandle,src_resource ResourceHandle) error {
	// TODO
}

////////////////////////////////////////////////////////////////////////////////
// Private methods - resources

func resourceCreate(format ImageType, w, h uint32, buffer *uint32) (ResourceHandle, error) {
	handle := C.vc_dispmanx_resource_create(C.VC_IMAGE_TYPE_T(format), C.uint32_t(w), C.uint32_t(h), (*C.uint32_t)(unsafe.Pointer(buffer)))
	if handle == DISPMANX_NO_HANDLE {
		return ResourceHandle(0), ErrorResource
	}
	return ResourceHandle(handle), nil
}

func resourceDelete(handle ResourceHandle) error {
	if C.vc_dispmanx_resource_delete(C.DISPMANX_RESOURCE_HANDLE_T(handle)) != DISPMANX_SUCCESS {
		return ErrorResource
	}
	return nil
}
