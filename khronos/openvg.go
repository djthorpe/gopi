/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package khronos /* import "github.com/djthorpe/gopi/khronos" */

import (
	gopi ".." /* import "github.com/djthorpe/gopi" */
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Abstract driver interface
type VGDriver interface {
	// Inherit general driver interface
	gopi.Driver

	// Start drawing
	Begin(window EGLWindow) error

	// Flush
	Flush() error

	// Clear
	Clear(color VGColor)

}

// Color with Alpha value
type VGColor struct {
	R,G,B,A float32
}

