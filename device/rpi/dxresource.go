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
	"image"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

import (
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host
    #cgo LDFLAGS:  -L/opt/vc/lib -lbcm_host
	#include "vc_dispmanx.h"
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	dxResourceHandle uint32
)

type DXResource struct {
	handle dxResourceHandle
	model  DXColorModel    // color model, which should be RGBA32 (4 bytes per pixel)
	size   khronos.EGLSize // size of the bitmap
	stride uint32          // number of bytes per row rounded up to 16-byte boundaries
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DX_RESOURCE_NONE    dxResourceHandle = 0
	DX_RESOURCE_SUCCESS                  = DX_SUCCESS
)

////////////////////////////////////////////////////////////////////////////////

func (this *DXDisplay) CreateResource(model DXColorModel, size khronos.EGLSize) (*DXResource, error) {
	resource := new(DXResource)
	resource.size = size
	resource.model = model
	resource.stride = dxAlignUp(uint32(size.Width), uint32(16)) * 4
	this.log.Debug2("<rpi.DX>CreateResource model=%v size=%v stride=%v", model, size, resource.stride)
	resource.handle = dxResourceCreate(model, size.Width, size.Height)
	if resource.handle == DX_RESOURCE_NONE {
		return nil, this.log.Error("dxResourceCreate failed")
	}
	return resource, nil
}

func (this *DXDisplay) DestroyResource(resource *DXResource) error {
	this.log.Debug2("<rpi.DX>DestroyResource")
	if dxResourceDelete(resource.handle) != true {
		return this.log.Error("dxResourceDelete failed")
	}
	resource.handle = DX_RESOURCE_NONE
	return nil
}

func (this *DXResource) String() string {
	return fmt.Sprintf("<rpi.DXResource>{ handle=%v model=%v size=%v stride=%v }", this.handle, this.model, this.size, this.stride)
}

func (h dxResourceHandle) String() string {
	return fmt.Sprintf("<rpi.DXResourceHandle>{%08X}", uint32(h))
}

func (this *DXResource) GetSize() khronos.EGLSize {
	return this.size
}

func (this *DXResource) GetFrame() khronos.EGLFrame {
	return khronos.EGLFrame{khronos.EGLPoint{0, 0}, khronos.EGLSize{this.size.Width, this.size.Height}}
}

func (this *DXResource) GetHandle() dxResourceHandle {
	return this.handle
}

func (this *DXResource) ClearToColor(color khronos.EGLColorRGBA32) error {
	data, err := dxReadBitmap(this, false)
	if err != nil {
		return err
	}
	value := color.Uint32()
	for i := 0; i < len(data); i++ {
		data[i] = value
	}

	// Write bitmap
	if err := dxWriteBitmap(this, data); err != nil {
		return err
	}

	return nil
}

func (this *DXResource) PaintImage(pt khronos.EGLPoint, bitmap image.Image) error {
	data, err := dxReadBitmap(this, false)
	if err != nil {
		return err
	}
	bounds := bitmap.Bounds()
	for i := uint(0); i < uint(len(data)); i++ {
		dx := i % uint(this.stride>>2)
		dy := i / uint(this.stride>>2)
		if dx >= this.size.Width || dy >= this.size.Height {
			continue
		}
		sx := int(dx)
		sy := int(dy)
		if sx > bounds.Dx() || sy > bounds.Dy() {
			continue
		}
		r, g, b, a := bitmap.At(int(sx), int(sy)).RGBA()
		data[i] = ((r & 0xFF00) >> 8) | (g & 0xFF00) | ((b & 0xFF00) << 8) | ((a & 0xFF00) << 16)
	}

	// Write bitmap
	if err := dxWriteBitmap(this, data); err != nil {
		return err
	}

	return nil
}

func (this *DXResource) PaintText(text string, face khronos.VGFace, color khronos.EGLColorRGBA32, origin khronos.EGLPoint, size float32) error {
	// Get frame & color
	frame := this.GetFrame()
	color_uint32 := color.Uint32()

	// Get bitmap
	data, err := dxReadBitmap(this, false)
	if err != nil {
		return err
	}

	// Set font size
	if err := face.(*vgfFace).SetSize(size); err != nil {
		return err
	}

	// Draw
	for i, w := 0, 0; i < len(text); i += w {

		// Load bitmap for rune
		char, width := utf8.DecodeRuneInString(text[i:])
		buffer, pixel_mode, size, advance, stride, err := face.(*vgfFace).LoadBitmapForRune(char)
		if err != nil {
			return err
		}

		// Check pixel mode to ensure it's supported
		if pixel_mode != VG_FONT_FT_PIXEL_MODE_GRAY {
			// We only support greyscale for the moment
			return errors.New("Glyph pixel mode unsupported")
		}

		// convert buffer to []byte (supports 8-bit greyscales)
		data_src := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{buffer, int(stride * size.Height), int(stride * size.Height)}))
		offset_src := uint(stride * size.Height)

		// paint pixels
		for y := uint(0); y < size.Height; y++ {
			offset_src -= stride
			for x := uint(0); x < size.Width; x++ {
				point := khronos.EGLPoint{origin.X + int(x), origin.Y - int(y)}
				if point.InFrame(frame) {
					offset_dst := point.X + point.Y*int(this.stride>>2)
					data[offset_dst] = dxPixelOver(data[offset_dst], color_uint32, data_src[offset_src+x])
				}
			}
		}
		origin = origin.Offset(advance)
		w = width
	}

	// Write bitmap
	if err := dxWriteBitmap(this, data); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods

// OVER pixel function
func dxPixelOver(src uint32, dst uint32, alpha byte) uint32 {
	sr := float32(src & 0x000000FF)
	sg := float32((src & 0x0000FF00) >> 8)
	sb := float32((src & 0x00FF0000) >> 16)
	dr := float32(dst & 0x000000FF)
	dg := float32((dst & 0x0000FF00) >> 8)
	db := float32((dst & 0x00FF0000) >> 16)
	da := (dst & 0xFF000000) >> 24
	pa := float32(alpha / 255.0)
	na := 1.0 - pa
	dr = dr*pa + na*sr
	dg = dg*pa + na*sg
	db = db*pa + na*sb
	return da<<24 | uint32(byte(db))<<16 | uint32(byte(dg))<<8 | uint32(byte(dr))
}

// Create a bitmap buffer and optionally read the data from the resource
func dxReadBitmap(resource *DXResource, readData bool) ([]uint32, error) {
	data := make([]uint32, uint(resource.stride/4)*uint(resource.size.Height))
	frame := DXFrame{DXPoint{int32(0), int32(0)}, DXSize{uint32(resource.size.Width), uint32(resource.size.Height)}}
	if success := dxResourceReadData(resource.handle, &frame, unsafe.Pointer(&data[0]), resource.stride); success == false {
		return nil, EGLErrorInvalidParameter
	}
	return data, nil
}

// Write bitmap
func dxWriteBitmap(resource *DXResource, data []uint32) error {
	frame := DXFrame{DXPoint{int32(0), int32(0)}, DXSize{uint32(resource.size.Width), uint32(resource.size.Height)}}
	if success := dxResourceWriteData(resource.handle, resource.model, resource.stride, unsafe.Pointer(&data[0]), &frame); success == false {
		return EGLErrorInvalidParameter
	}
	return nil
}

func dxResourceCreate(model DXColorModel, w, h uint) dxResourceHandle {
	var dummy C.uint32_t
	return dxResourceHandle(C.vc_dispmanx_resource_create(C.VC_IMAGE_TYPE_T(model), C.uint32_t(w), C.uint32_t(h), (*C.uint32_t)(unsafe.Pointer(&dummy))))
}

func dxResourceDelete(handle dxResourceHandle) bool {
	return C.vc_dispmanx_resource_delete(C.DISPMANX_RESOURCE_HANDLE_T(handle)) == DX_RESOURCE_SUCCESS
}

func dxResourceWriteData(handle dxResourceHandle, model DXColorModel, src_pitch uint32, src_buffer unsafe.Pointer, dst_rect *DXFrame) bool {
	return C.vc_dispmanx_resource_write_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), C.VC_IMAGE_TYPE_T(model), C.int(src_pitch), src_buffer, (*C.VC_RECT_T)(unsafe.Pointer(dst_rect))) == DX_RESOURCE_SUCCESS
}

func dxResourceReadData(handle dxResourceHandle, src_rect *DXFrame, dst_buffer unsafe.Pointer, dst_pitch uint32) bool {
	return C.vc_dispmanx_resource_read_data(C.DISPMANX_RESOURCE_HANDLE_T(handle), (*C.VC_RECT_T)(unsafe.Pointer(src_rect)), unsafe.Pointer(dst_buffer), C.uint32_t(dst_pitch)) == DX_RESOURCE_SUCCESS
}

func dxAlignUp(value, alignment uint32) uint32 {
	return ((value - 1) & ^(alignment - 1)) + alignment
}
