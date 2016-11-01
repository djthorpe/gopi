/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package khronos /* import "github.com/djthorpe/gopi/khronos" */

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract 2D bitmap surface
type EGLBitmap interface {

	// Return window size
	GetSize() EGLSize

	// Set Pixel
	SetPixel(pt EGLPoint) error
}
