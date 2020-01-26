/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package fonts

import (
	"image"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

type Face interface {
	// Return a bitmap for a rune at pixel size
	BitmapForRunePixels(rune, uint) (image.Image, error)

	// Return a bitmap for a rune at point size (with dpi)
	BitmapForRunePoints(rune, float32, uint) (image.Image, error)

	gopi.FontFace
}
