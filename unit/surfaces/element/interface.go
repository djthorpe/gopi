/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package element

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

type Config struct {
	Origin  gopi.Point
	Size    gopi.Size
	Layer   uint16
	Opacity float32
	Flags   gopi.SurfaceFlags
	Bitmap  bitmap.Bitmap
	Update  rpi.DXUpdate
	Display rpi.DXDisplayHandle
}

type Element interface {
	// Return bounds
	Origin() gopi.Point
	Size() gopi.Size

	// Return bitmap
	Bitmap() bitmap.Bitmap

	// Set element properties
	SetOrigin(rpi.DXUpdate, gopi.Point) error
	SetSize(rpi.DXUpdate, gopi.Size) error
	SetLayer(rpi.DXUpdate, uint16) error
	SetOpacity(rpi.DXUpdate, float32) error
	SetBitmap(rpi.DXUpdate, bitmap.Bitmap) error

	// Remove element
	RemoveElement(rpi.DXUpdate) error

	// Implements Unit
	gopi.Unit
}
