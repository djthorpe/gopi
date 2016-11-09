/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"io"
	"image"
)

import (
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
)

import (
	khronos "github.com/djthorpe/gopi/khronos"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC FUNCTIONS

// Create an image resource
func (this *eglDriver) CreateImage(r io.Reader) (khronos.EGLBitmap,error) {
	// Decode the bitmap
	bitmap, _, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	bounds := bitmap.Bounds()
	resource, err := this.dx.CreateResource(DX_IMAGE_RGBA16,khronos.EGLSize{ uint(bounds.Dx()), uint(bounds.Dy()) })
	if err != nil {
		return nil, err
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r,g,b,a := bitmap.At(x, y).RGBA()
			resource.SetPixel(khronos.EGLPoint{ x, y },uint16(r),uint16(g),uint16(b),uint16(a))
		}
	}

	return resource, nil
}

// Destroy an image resource
func (this *eglDriver) DestroyImage(bitmap khronos.EGLBitmap) error {
	return this.dx.DestroyResource(bitmap.(*DXResource))
}

