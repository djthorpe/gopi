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

	// Return frame with origin at (0,0)
	GetFrame() EGLFrame

	// Paint an image into the bitmap
	PaintImage(pt EGLPoint, bitmap image.Image) error

	// Clear bitmap to one color
	ClearToColor(color EGLColorRGBA32) error

	// Draw text starting at an origin as (bottom,left) with a particular
	// font face and a point size
	PaintText(text string, face VGFace, color EGLColorRGBA32, origin EGLPoint, size float32) error
}
