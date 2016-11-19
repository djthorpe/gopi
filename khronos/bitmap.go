/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package khronos /* import "github.com/djthorpe/gopi/khronos" */

import (
	"image"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract 2D bitmap surface
type EGLBitmap interface {

	// Return window size
	GetSize() EGLSize

	// Paint an image into the bitmap
	PaintImage(pt EGLPoint,bitmap image.Image) error

	// Clear bitmap to one color
	ClearToColor(color EGLColorRGBA32) error
}
