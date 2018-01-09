// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host
#cgo LDFLAGS:  -L/opt/vc/lib -lbcm_host
#include "vc_dispmanx.h"
*/
import "C"
import "unsafe"
import "fmt"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type dxRect C.VC_RECT_T
type dxResourceHandle C.DISPMANX_RESOURCE_HANDLE_T
type dxUpdateHandle C.DISPMANX_UPDATE_HANDLE_T
type dxElementHandle C.DISPMANX_ELEMENT_HANDLE_T
type dxDisplayHandle C.DISPMANX_DISPLAY_HANDLE_T
type dxImageType int
type dxError int
type dxProtection uint32
type dxAlphaFlags int
type dxTransformFlags int

type dxAlpha struct {
	Flags   dxAlphaFlags
	Opacity uint32
	Mask    dxResourceHandle
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// dxImageType values
	// We only list defaults for supported color models on the Raspberry Pi
	// From /opt/vc/include/interface/vctypes/vc_image_types.h
	DX_IMAGETYPE_RGB565    dxImageType = 1  // 16 bits per pixel
	DX_IMAGETYPE_YUV420    dxImageType = 3  // 16 bits per pixel
	DX_IMAGETYPE_RGB888    dxImageType = 5  // 24 bits per pixel
	DX_IMAGETYPE_4BPP      dxImageType = 7  // 4bpp palettised image
	DX_IMAGETYPE_RGBA32    dxImageType = 15 /* 32 bits per pixel - RGB888 0xAABBGGRR */
	DX_IMAGETYPE_YUV422    dxImageType = 16 /* a line of Y (32-byte padded), a line of U (16-byte padded), and a line of V (16-byte padded) */
	DX_IMAGETYPE_RGBA565   dxImageType = 17 /* RGB565 with a transparent patch */
	DX_IMAGETYPE_RGBA16    dxImageType = 18 /* 16 bits per pixel - Compressed (4444) version of RGBA32 */
	DX_IMAGETYPE_YUV_UV    dxImageType = 19 /* VCIII codec format */
	DX_IMAGETYPE_TF_RGBA32 dxImageType = 20 /* VCIII T-format RGBA8888 */
	DX_IMAGETYPE_TF_RGBX32 dxImageType = 21 /* VCIII T-format RGBx8888 */
	DX_IMAGETYPE_TF_RGBA16 dxImageType = 23 /* VCIII T-format RGBA4444 */
	DX_IMAGETYPE_TF_RGB565 dxImageType = 25 /* VCIII T-format RGB565 */
)

const (
	DX_SUCCESS dxError = iota
	DX_RESOURCE_ERROR
	DX_UPDATE_ERROR
	DX_ELEMENT_ERROR
)

const (
	DX_NO_RESOURCE = C.DISPMANX_RESOURCE_HANDLE_T(0)
	DX_NO_UPDATE   = C.DISPMANX_UPDATE_HANDLE_T(0)
	DX_NO_ELEMENT  = C.DISPMANX_ELEMENT_HANDLE_T(0)
)

const (
	// From /opt/vc/include/interface/vmcs_host/vc_dispmanx_types.h
	DX_PROTECTION_NONE dxProtection = 0x000000000
	DX_PROTECTION_HDCP dxProtection = 0x00000000B
	DX_PROTECTION_MAX  dxProtection = 0x00000000F
)

const (
	// From /opt/vc/include/interface/vmcs_host/vc_dispmanx_types.h
	/* Bottom 2 bits sets the alpha mode */
	DX_ALPHA_FROM_SOURCE       dxAlphaFlags = 0
	DX_ALPHA_FIXED_ALL_PIXELS  dxAlphaFlags = 1
	DX_ALPHA_FIXED_NON_ZERO    dxAlphaFlags = 2
	DX_ALPHA_FIXED_EXCEED_0X07 dxAlphaFlags = 3
	DX_ALPHA_PREMULT           dxAlphaFlags = 1 << 16
	DX_ALPHA_MIX               dxAlphaFlags = 1 << 17
)

const (
	// From /opt/vc/include/interface/vmcs_host/vc_dispmanx_types.h
	/* Bottom 2 bits sets the orientation */
	DX_TRANSFORM_NO_ROTATE  dxTransformFlags = 0
	DX_TRANSFORM_ROTATE_90  dxTransformFlags = 1
	DX_TRANSFORM_ROTATE_180 dxTransformFlags = 2
	DX_TRANSFORM_ROTATE_270 dxTransformFlags = 3

	DX_TRANSFORM_FLIP_HRIZ dxTransformFlags = 1 << 16
	DX_TRANSFORM_FLIP_VERT dxTransformFlags = 1 << 17

	/* invert left/right images */
	DX_TRANSFORM_STEREOSCOPIC_INVERT dxTransformFlags = 1 << 19

	/* extra flags for controlling 3d duplication behaviour */
	DX_TRANSFORM_STEREOSCOPIC_NONE dxTransformFlags = 0 << 20
	DX_TRANSFORM_STEREOSCOPIC_MONO dxTransformFlags = 1 << 20
	DX_TRANSFORM_STEREOSCOPIC_SBS  dxTransformFlags = 2 << 20
	DX_TRANSFORM_STEREOSCOPIC_TB   dxTransformFlags = 3 << 20
	DX_TRANSFORM_STEREOSCOPIC_MASK dxTransformFlags = 15 << 20

	/* extra flags for controlling snapshot behaviour */
	DX_TRANSFORM_SNAPSHOT_NO_YUV        dxTransformFlags = 1 << 24
	DX_TRANSFORM_SNAPSHOT_NO_RGB        dxTransformFlags = 1 << 25
	DX_TRANSFORM_SNAPSHOT_FILL          dxTransformFlags = 1 << 26
	DX_TRANSFORM_SNAPSHOT_SWAP_RED_BLUE dxTransformFlags = 1 << 27
	DX_TRANSFORM_SNAPSHOT_PACK          dxTransformFlags = 1 << 28
)

var (
	DX_NULL = unsafe.Pointer(uintptr(0))
)

////////////////////////////////////////////////////////////////////////////////
// RECT

func dxRectSet(x, y, width, height uint32) dxRect {
	var rect C.VC_RECT_T
	C.vc_dispmanx_rect_set(&rect, C.uint32_t(x), C.uint32_t(y), C.uint32_t(width), C.uint32_t(height))
	return dxRect(rect)
}

////////////////////////////////////////////////////////////////////////////////
// RESOURCES

// Create a new resource
func dxResourceCreate(image_type dxImageType, width, height uint32) (dxResourceHandle, dxError) {
	var native_image_handle C.uint32_t
	if handle := C.vc_dispmanx_resource_create(C.VC_IMAGE_TYPE_T(image_type), C.uint32_t(width), C.uint32_t(height), &native_image_handle); handle == DX_NO_RESOURCE {
		return dxResourceHandle(DX_NO_RESOURCE), DX_RESOURCE_ERROR
	} else {
		return dxResourceHandle(handle), DX_SUCCESS
	}
}

// Delete a resource
func dxResourceDelete(handle dxResourceHandle) dxError {
	if err := C.vc_dispmanx_resource_delete(C.DISPMANX_RESOURCE_HANDLE_T(handle)); dxError(err) != DX_SUCCESS {
		return DX_RESOURCE_ERROR
	} else {
		return DX_SUCCESS
	}
}

// Write the bitmap data to VideoCore memory
func dxResourceWriteDataHandle(handle dxResourceHandle, src_type dxImageType, src_pitch int, source unsafe.Pointer, rect *dxRect) dxError {
	if err := C.vc_dispmanx_resource_write_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), C.VC_IMAGE_TYPE_T(src_type), C.int(src_pitch), unsafe.Pointer(source), (*C.VC_RECT_T)(rect)); dxError(err) != DX_SUCCESS {
		return DX_RESOURCE_ERROR
	} else {
		return DX_SUCCESS
	}
}

// Read bitmap data from VideoCore memory
func dxResourceReadData(handle dxResourceHandle, rect *dxRect, dst_pitch uint32) ([]byte, dxError) {
	// TODO: Calculate the size of the buffer we need
	destination := make([]byte, 0)
	if err := C.vc_dispmanx_resource_read_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), (*C.VC_RECT_T)(rect), unsafe.Pointer(&destination[0]), C.uint32_t(dst_pitch)); dxError(err) != DX_SUCCESS {
		return nil, DX_RESOURCE_ERROR
	} else {
		return destination, DX_SUCCESS
	}
}

////////////////////////////////////////////////////////////////////////////////
// UPDATES

// Start an update
func dxUpdateStart(priority int32) (dxUpdateHandle, dxError) {
	if handle := C.vc_dispmanx_update_start(C.int32_t(priority)); handle == DX_NO_UPDATE {
		return dxUpdateHandle(DX_NO_UPDATE), DX_UPDATE_ERROR
	} else {
		return dxUpdateHandle(handle), DX_SUCCESS
	}
}

// End an update and wait for it to complete
func dxUpdateSubmitSync(handle dxUpdateHandle) dxError {
	if err := C.vc_dispmanx_update_submit_sync(C.DISPMANX_UPDATE_HANDLE_T(handle)); dxError(err) != DX_SUCCESS {
		return DX_UPDATE_ERROR
	} else {
		return DX_SUCCESS
	}
}

// Ends an update and callback
// int vc_dispmanx_update_submit( DISPMANX_UPDATE_HANDLE_T update, DISPMANX_CALLBACK_FUNC_T cb_func, void *cb_arg );

////////////////////////////////////////////////////////////////////////////////
// ELEMENTS

// Add an elment to a display as part of an update
func dxElementAdd(update dxUpdateHandle, display dxDisplayHandle, layer int32, dest_rect *dxRect, resource dxResourceHandle, src_rect *dxRect, protection dxProtection, alpha dxAlpha, transform dxTransformFlags) (dxElementHandle, dxError) {
	if element := C.vc_dispmanx_element_add(
		C.DISPMANX_UPDATE_HANDLE_T(update),
		C.DISPMANX_DISPLAY_HANDLE_T(display),
		C.int32_t(layer),
		(*C.VC_RECT_T)(dest_rect),
		C.DISPMANX_RESOURCE_HANDLE_T(resource),
		(*C.VC_RECT_T)(src_rect),
		C.DISPMANX_PROTECTION_T(protection),
		(*C.VC_DISPMANX_ALPHA_T)(unsafe.Pointer(&alpha)),
		(*C.DISPMANX_CLAMP_T)(DX_NULL), // TODO: We ignore clamp for the moment
		C.DISPMANX_TRANSFORM_T(transform),
	); element == DX_NO_ELEMENT {
		return dxElementHandle(DX_NO_ELEMENT), DX_ELEMENT_ERROR
	} else {
		return dxElementHandle(element), DX_SUCCESS
	}
}

// Remove a display element from its display
func dxElementRemove(update dxUpdateHandle, element dxElementHandle) dxError {
	if err := C.vc_dispmanx_element_remove(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element)); dxError(err) != DX_SUCCESS {
		return DX_ELEMENT_ERROR
	} else {
		return DX_SUCCESS
	}
}

// Signal that a region of the bitmap has been modified
func dxElementModified(update dxUpdateHandle, element dxElementHandle, rect *dxRect) dxError {
	if err := C.vc_dispmanx_element_modified(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element), (*C.VC_RECT_T)(rect)); dxError(err) != DX_SUCCESS {
		return DX_ELEMENT_ERROR
	} else {
		return DX_SUCCESS
	}
}

// Change the source image of a display element
func dxElementChangeSource(update dxUpdateHandle, element dxElementHandle, src dxResourceHandle) dxError {
	if err := C.vc_dispmanx_element_change_source(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element), C.DISPMANX_RESOURCE_HANDLE_T(src)); dxError(err) != DX_SUCCESS {
		return DX_ELEMENT_ERROR
	} else {
		return DX_SUCCESS
	}
}

// Change the layer number of a display element
func dxElementChangeLayer(update dxUpdateHandle, element dxElementHandle, layer int32) dxError {
	if err := C.vc_dispmanx_element_change_layer(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element), C.int32_t(layer)); dxError(err) != DX_SUCCESS {
		return DX_ELEMENT_ERROR
	} else {
		return DX_SUCCESS
	}
}

// New function added to VCHI to change attributes, set_opacity does not work there.
func dxElementChangeAttributes(update dxUpdateHandle, element dxElementHandle, change_flags uint32, layer int32, opacity uint8, dest_rect, src_rect *dxRect, mask dxResourceHandle, transform dxTransformFlags) dxError {
	if err := C.vc_dispmanx_element_change_attributes(C.DISPMANX_UPDATE_HANDLE_T(update), C.DISPMANX_ELEMENT_HANDLE_T(element), C.uint32_t(change_flags), C.int32_t(layer), C.uint8_t(opacity), (*C.VC_RECT_T)(dest_rect), (*C.VC_RECT_T)(src_rect), C.DISPMANX_RESOURCE_HANDLE_T(mask), C.DISPMANX_TRANSFORM_T(transform)); dxError(err) != DX_SUCCESS {
		return DX_ELEMENT_ERROR
	} else {
		return DX_SUCCESS
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (e dxError) Error() string {
	switch e {
	case DX_SUCCESS:
		return "DX_SUCCESS"
	case DX_RESOURCE_ERROR:
		return "DX_RESOURCE_ERROR"
	case DX_UPDATE_ERROR:
		return "DX_UPDATE_ERROR"
	case DX_ELEMENT_ERROR:
		return "DX_ELEMENT_ERROR"
	default:
		return "[?? Invalid dxError value]"
	}
}

func (i dxImageType) String() string {
	switch i {
	case DX_IMAGETYPE_RGB565:
		return "DX_IMAGETYPE_RGB565"
	case DX_IMAGETYPE_YUV420:
		return "DX_IMAGETYPE_YUV420"
	case DX_IMAGETYPE_RGB888:
		return "DX_IMAGETYPE_RGB888"
	case DX_IMAGETYPE_4BPP:
		return "DX_IMAGETYPE_4BPP"
	case DX_IMAGETYPE_RGBA32:
		return "DX_IMAGETYPE_RGBA32"
	case DX_IMAGETYPE_YUV422:
		return "DX_IMAGETYPE_YUV422"
	case DX_IMAGETYPE_RGBA565:
		return "DX_IMAGETYPE_RGBA565"
	case DX_IMAGETYPE_RGBA16:
		return "DX_IMAGETYPE_RGBA16"
	case DX_IMAGETYPE_YUV_UV:
		return "DX_IMAGETYPE_YUV_UV"
	case DX_IMAGETYPE_TF_RGBA32:
		return "DX_IMAGETYPE_TF_RGBA32"
	case DX_IMAGETYPE_TF_RGBX32:
		return "DX_IMAGETYPE_TF_RGBX32"
	case DX_IMAGETYPE_TF_RGBA16:
		return "DX_IMAGETYPE_TF_RGBA16"
	case DX_IMAGETYPE_TF_RGB565:
		return "DX_IMAGETYPE_TF_RGB565"
	default:
		return "[?? Invalid dxImageType value]"
	}
}

func (p dxProtection) String() string {
	switch p {
	case DX_PROTECTION_NONE:
		return "DX_PROTECTION_NONE"
	case DX_PROTECTION_HDCP:
		return "DX_PROTECTION_HDCP"
	default:
		return "[?? Invalid dxProtection value]"

	}
}

func (h dxResourceHandle) String() string {
	if h == dxResourceHandle(DX_NO_RESOURCE) {
		return "dxResource{nil}"
	} else {
		return fmt.Sprintf("dxResource{0x%08X}", uint32(h))
	}
}

func (h dxElementHandle) String() string {
	if h == dxElementHandle(DX_NO_ELEMENT) {
		return "dxElement{nil}"
	} else {
		return fmt.Sprintf("dxElement{0x%08X}", uint32(h))
	}
}

func (h dxUpdateHandle) String() string {
	if h == dxUpdateHandle(DX_NO_UPDATE) {
		return "dxUpdate{nil}"
	} else {
		return fmt.Sprintf("dxUpdate{0x%08X}", uint32(h))
	}
}
