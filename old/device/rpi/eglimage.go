/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"image"
	"io"
)

import (
	_ "image/gif"
	_ "image/jpeg"
	png "image/png"
)

import (
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

// Create an image resource
func (this *eglDriver) CreateImage(r io.Reader) (khronos.EGLBitmap, error) {
	// Decode the bitmap
	bitmap, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	bounds := bitmap.Bounds()
	resource, err := this.dx.CreateResource(DX_IMAGE_RGBA32, khronos.EGLSize{uint(bounds.Dx()), uint(bounds.Dy())})
	if err != nil {
		return nil, err
	}

	if err := resource.PaintImage(khronos.EGLZeroPoint, bitmap); err != nil {
		this.dx.DestroyResource(resource)
		return nil, err
	}

	return resource, nil
}

// Destroy an image resource
func (this *eglDriver) DestroyImage(bitmap khronos.EGLBitmap) error {
	return this.dx.DestroyResource(bitmap.(*DXResource))
}

// Snapshot screen
func (this *eglDriver) SnapshotImage() (khronos.EGLBitmap, error) {
	resource, err := this.dx.CreateSnapshot()
	if err != nil {
		return nil, err
	}
	return resource, nil
}

// Write image as a PNG
func (this *eglDriver) WriteImagePNG(w io.Writer, bitmap khronos.EGLBitmap) error {
	size := bitmap.GetSize()
	image := image.NewRGBA(image.Rect(0, 0, int(size.Width), int(size.Height)))
	data, err := dxReadBitmap(bitmap.(*DXResource), true)
	stride := (bitmap.(*DXResource).stride >> 2)
	if err != nil {
		return err
	}
	for y := uint32(0); y < uint32(size.Height); y++ {
		for x := uint32(0); x < uint32(size.Width); x++ {
			source_pixel := data[x+y*stride]
			destination_offset := image.PixOffset(int(x), int(y))
			image.Pix[destination_offset+0] = uint8((source_pixel & 0x000000FF) >> 0)  // R
			image.Pix[destination_offset+1] = uint8((source_pixel & 0x0000FF00) >> 8)  // G
			image.Pix[destination_offset+2] = uint8((source_pixel & 0x00FF0000) >> 16) // B
			image.Pix[destination_offset+3] = uint8((source_pixel & 0xFF000000) >> 24) // A
		}
	}

	return png.Encode(w, image)
}
