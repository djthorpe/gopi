/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package manager

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
	element "github.com/djthorpe/gopi/v2/unit/surfaces/element"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

type Config struct {
	Display rpi.DXDisplayHandle
}

type Manager interface {
	// Create a new bitmap
	NewBitmap(gopi.Size, gopi.SurfaceFlags) (bitmap.Bitmap, error)

	// Release a bitmap
	ReleaseBitmap(bitmap.Bitmap) error

	// Create a new element with a size
	AddElementWithSize(gopi.Point, gopi.Size, uint16, float32) (element.Element, error)

	/*


		// Create a new element with a bitmap
		AddElementWithBitmap(gopi.Point, bitmap.Bitmap, uint16, float32) (element.Element, error)

		// Remove an element
		RemoveElement(element.Element) error

		// Perform AddElement, RemoveElement and bitmap operations within Do
		Do(func() error) error
	*/

	// Implements Unit
	gopi.Unit
}
