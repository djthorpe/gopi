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
	_ "image/png"
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
