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

type Point struct {
	X int32
	Y int32
}

type Size struct {
	Width  uint32
	Height uint32
}

type Rectangle struct {
	Point
	Size
}

type (
	DisplayHandle  uint32
	UpdateHandle   uint32
	ElementHandle  uint32
	ResourceHandle uint32
	UpdatePriority int32
	InputFormat    uint32
	Opacity        uint8
	Protection     uint32
	Transform      int
	ImageType      int
	ClampMode      int
)

type ModeInfo struct {
	Size          Size
	Transform     Transform
	InputFormat   InputFormat
	Handle        DisplayHandle
}

type Color struct {
	Red, Green, Blue uint8
}

type Resource struct {
	handle ResourceHandle
	size Size
	buffer *byte
}

type Element struct {
	handle ElementHandle
	frame *Rectangle
	layer int32
}

type Alpha struct {
	Flags   uint32
	Opacity uint32
	Mask    ResourceHandle
}

type Clamp struct {
	Mode    ClampMode
	Flags   int
	Opacity uint32
	Mask    ResourceHandle
}

type VideoCore struct {
	display uint16
	size    Size
	handle  DisplayHandle
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	/* Success and failure conditions */
	DISPMANX_SUCCESS   = 0
	DISPMANX_INVALID   = -1
	DISPMANX_NO_HANDLE = 0
)

const (
	/* Display ID's */
	DISPMANX_ID_MAIN_LCD uint16 = iota
	DISPMANX_ID_AUX_LCD
	DISPMANX_ID_HDMI
	DISPMANX_ID_SDTV
	DISPMANX_ID_FORCE_LCD
	DISPMANX_ID_FORCE_TV
	DISPMANX_ID_FORCE_OTHER /* non-default display */
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
	VC_IMAGE_RGBA32        /* RGB888 0xAABBGGRR */
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

const (
	/* Alpha flags */
	DISPMANX_FLAGS_ALPHA_FROM_SOURCE       uint32 = 0 /* Bottom 2 bits sets the alpha mode */
	DISPMANX_FLAGS_ALPHA_FIXED_ALL_PIXELS  uint32 = 1
	DISPMANX_FLAGS_ALPHA_FIXED_NON_ZERO    uint32 = 2
	DISPMANX_FLAGS_ALPHA_FIXED_EXCEED_0X07 uint32 = 3
	DISPMANX_FLAGS_ALPHA_PREMULT           uint32 = 1 << 16
	DISPMANX_FLAGS_ALPHA_MIX               uint32 = 1 << 17
)

const (
	/* Clamp values */
	DISPMANX_FLAGS_CLAMP_NONE               ClampMode = 0
	DISPMANX_FLAGS_CLAMP_LUMA_TRANSPARENT   ClampMode = 1
	DISPMANX_FLAGS_CLAMP_TRANSPARENT        ClampMode = 2
	DISPMANX_FLAGS_CLAMP_CHROMA_TRANSPARENT ClampMode = 2
	DISPMANX_FLAGS_CLAMP_REPLACE            ClampMode = 3
)

const (
	/* Protection values */
	DISPMANX_PROTECTION_NONE Protection = 0
	DISPMANX_PROTECTION_HDCP Protection = 11
)

const (
	/* Transform values */
	DISPMANX_NO_ROTATE Transform = iota
	DISPMANX_ROTATE_90
	DISPMANX_ROTATE_180
	DISPMANX_ROTATE_270
)

const (
	ELEMENT_CHANGE_LAYER uint32 = (1<<0)
	ELEMENT_CHANGE_OPACITY uint32 =        (1<<1)
	ELEMENT_CHANGE_DEST_RECT uint32 =      (1<<2)
	ELEMENT_CHANGE_SRC_RECT uint32 =       (1<<3)
	ELEMENT_CHANGE_MASK_RESOURCE uint32 =  (1<<4)
	ELEMENT_CHANGE_TRANSFORM uint32 =      (1<<5)
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	Displays = map[string]uint16{
		"lcd":        DISPMANX_ID_MAIN_LCD,
		"aux":        DISPMANX_ID_AUX_LCD,
		"hdmi":       DISPMANX_ID_HDMI,
		"tv":         DISPMANX_ID_SDTV,
		"forcelcd":   DISPMANX_ID_FORCE_LCD,
		"forcetv":    DISPMANX_ID_FORCE_TV,
		"forceother": DISPMANX_ID_FORCE_OTHER,
	}
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

	// open the display
	handle, err := displayOpen(uint32(display))
	if err != nil {
		return nil, err
	}

	// Populate the structure
	this.size = Size{ width, height }
	this.handle = handle

	// success
	return this, nil
}

// Close unmaps GPIO memory
func (this *VideoCore) Close() error {
	// Close display
	err := displayClose(this.handle)
	return err
}

////////////////////////////////////////////////////////////////////////////////
// Get VideoCore properties

func (this *VideoCore) GetDisplayID() uint16 {
	return this.display
}

func (this *VideoCore) GetSize() Size {
	return this.size
}

func (this *VideoCore) GetFrame() *Rectangle {
	return &Rectangle{ Point{ 0, 0 }, this.GetSize() }
}

func (this *VideoCore) GetModeInfo() (ModeInfo, error) {
	var info ModeInfo
	err := displayGetInfo(this.handle, &info)
	return info, err
}

////////////////////////////////////////////////////////////////////////////////
// UPDATES

func (this *VideoCore) UpdateBegin() (UpdateHandle, error) {
	return updateStart(UpdatePriority(0))
}

func (this *VideoCore) UpdateSubmit(handle UpdateHandle) error {
	return updateSubmitSync(handle)
}

func (this *VideoCore) SetBackgroundColor(handle UpdateHandle, color Color) error {
	return displaySetBackground(handle, this.handle, color.Red, color.Green, color.Blue)
}

////////////////////////////////////////////////////////////////////////////////
// RESOURCES

func (this *VideoCore) CreateResource(format ImageType, size Size) (*Resource, error) {
	handle, err := resourceCreate(format, size.Width, size.Height)
	if err != nil {
		return nil, err
	}

	resource := new(Resource)
	resource.handle = handle
	resource.size = size

	return resource, nil
}

func (this *VideoCore) DeleteResource(resource *Resource) error {
	err := resourceDelete(resource.handle)
	if err != nil {
		return err
	}
	return nil
}

func (this *Resource) GetSize() Size {
	return this.size
}

func (this *Resource) GetFrame() *Rectangle {
	return &Rectangle{ Point{ 0, 0 },this.GetSize() }
}

func (this *Resource) WriteData(format ImageType,src_pitch int,src_buffer []byte,dst_rect *Rectangle) error {
	return resourceWriteData(this.handle,format,src_pitch,&src_buffer[0],dst_rect)
}

////////////////////////////////////////////////////////////////////////////////
// RECTANGLES

func (this *Rectangle) Set(point Point,size Size) {
	C.vc_dispmanx_rect_set((*C.VC_RECT_T)(unsafe.Pointer(this)),C.uint32_t(point.X),C.uint32_t(point.Y),C.uint32_t(size.Width),C.uint32_t(size.Height))
}

////////////////////////////////////////////////////////////////////////////////
// ELEMENTS

func (this *VideoCore) AddElement(update UpdateHandle,layer int32,dst_rect *Rectangle,src_resource *Resource,src_rect *Rectangle) (*Element, error) {
	var src_resource_handle ResourceHandle

	// if there is a source resource, then set the handle
	if src_resource != nil {
		src_resource_handle = src_resource.handle
	}
	// destination frame
	if dst_rect == nil {
		dst_rect = this.GetFrame()
	}
	// source frame
	if src_rect == nil {
		if src_resource == nil {
			return nil, ErrorElement
		}
		src_rect = src_resource.GetFrame()
	}

	// set alpha to 255
	// TODO: Allow Alpha to be set
	alpha := Alpha{ DISPMANX_FLAGS_ALPHA_FIXED_ALL_PIXELS, 255, 0 }

	// add element
	handle, err := elementAdd(update,this.handle,layer,dst_rect,src_resource_handle,src_rect,DISPMANX_PROTECTION_NONE,&alpha,nil,0);
	if err != nil {
		return nil, err
	}

	// create element structure
	element := new(Element)
	element.handle = handle
	element.layer = layer
	element.frame = dst_rect

	return element,nil
}

func (this *VideoCore) RemoveElement(update UpdateHandle,element *Element) error {
	return elementRemove(update,element.handle)
}

func (this *VideoCore) ChangeElementSource(update UpdateHandle,element *Element,resource *Resource) error {
	if element == nil || resource == nil {
		return ErrorElement
	}
	return elementChangeSource(update,element.handle,resource.handle)
}

func (this *VideoCore) ChangeElementLayer(update UpdateHandle,element *Element,layer int32) error {
	if element == nil {
		return ErrorElement
	}
	err := elementChangeLayer(update,element.handle,layer)
	if err != nil {
		return err
	}
	element.layer = layer
	return nil
}

func (this *VideoCore) ChangeElementFrame(update UpdateHandle,element *Element,frame *Rectangle) error {
	if element == nil {
		return ErrorElement
	}
	err := elementChangeDestination(update,element.handle,frame)
	if err != nil {
		return err
	}
	element.frame = frame
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods - display

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
	if handle == DisplayHandle(0) {
		return handle, ErrorDisplay
	}
	return handle, nil
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

func elementAdd(update UpdateHandle, display DisplayHandle, layer int32, dest_rect *Rectangle, src_resource ResourceHandle, src_rect *Rectangle, protection Protection, alpha *Alpha, clamp *Clamp, transform Transform) (ElementHandle, error) {
	handle := C.vc_dispmanx_element_add(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_DISPLAY_HANDLE_T(display), C.int32_t(layer), (*C.VC_RECT_T)(unsafe.Pointer(dest_rect)), C.DISPMANX_RESOURCE_HANDLE_T(src_resource), (*C.VC_RECT_T)(unsafe.Pointer(src_rect)), C.DISPMANX_PROTECTION_T(protection), (*C.VC_DISPMANX_ALPHA_T)(unsafe.Pointer(alpha)), (*C.DISPMANX_CLAMP_T)(unsafe.Pointer(clamp)), C.DISPMANX_TRANSFORM_T(transform))
	if handle == DISPMANX_NO_HANDLE {
		return ElementHandle(0), ErrorElement
	}
	return ElementHandle(handle), nil
}

func elementRemove(update UpdateHandle, element ElementHandle) error {
	if C.vc_dispmanx_element_remove(C.DISPMANX_UPDATE_HANDLE_T(update),C.DISPMANX_ELEMENT_HANDLE_T(element)) != DISPMANX_SUCCESS {
		return ErrorElement
	}
	// success
	return nil
}

func elementModified(update UpdateHandle, element ElementHandle, rect Rectangle) error {
	// TODO
	return ErrorElement
}

func elementChangeLayer(update UpdateHandle, element ElementHandle, layer int32) error {
	if C.vc_dispmanx_element_change_layer(C.DISPMANX_UPDATE_HANDLE_T(update),C.DISPMANX_ELEMENT_HANDLE_T(element),C.int32_t(layer)) != DISPMANX_SUCCESS {
		return ErrorElement
	}
	// success
	return nil
}

func elementChangeSource(update UpdateHandle, element ElementHandle, resource ResourceHandle) error {
	if C.vc_dispmanx_element_change_source(C.DISPMANX_UPDATE_HANDLE_T(update),C.DISPMANX_ELEMENT_HANDLE_T(element),C.DISPMANX_RESOURCE_HANDLE_T(resource)) != DISPMANX_SUCCESS {
		return ErrorElement
	}
	// success
	return nil
}

func elementChangeDestination(update UpdateHandle, element ElementHandle, frame *Rectangle) error {
	if C.vc_dispmanx_element_change_attributes(
		C.DISPMANX_UPDATE_HANDLE_T(update),
		C.DISPMANX_ELEMENT_HANDLE_T(element),
		C.uint32_t(ELEMENT_CHANGE_DEST_RECT),
		C.int32_t(0), // layer
		C.uint8_t(0), // opacity
		(*C.VC_RECT_T)(unsafe.Pointer(frame)), // dest_rect
		(*C.VC_RECT_T)(unsafe.Pointer(nil)), // src_rect
		C.DISPMANX_RESOURCE_HANDLE_T(0), // mask
		C.DISPMANX_TRANSFORM_T(0), // transform
	) != DISPMANX_SUCCESS {
		return ErrorElement
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods - resources

func resourceCreate(format ImageType, w, h uint32) (ResourceHandle, error) {
	var ptr C.uint32_t
	handle := C.vc_dispmanx_resource_create(C.VC_IMAGE_TYPE_T(format), C.uint32_t(w), C.uint32_t(h), (*C.uint32_t)(unsafe.Pointer(&ptr)))
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

func resourceWriteData(handle ResourceHandle,format ImageType,src_pitch int,src_buffer *byte,dst_rect *Rectangle) error {
	if C.vc_dispmanx_resource_write_data(C.DISPMANX_RESOURCE_HANDLE_T(handle),C.VC_IMAGE_TYPE_T(format),C.int(src_pitch),unsafe.Pointer(src_buffer),(*C.VC_RECT_T)(unsafe.Pointer(dst_rect))) != DISPMANX_SUCCESS {
		return ErrorResource
	}
	return nil
}

