/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package khronos /* import "github.com/djthorpe/gopi/khronos" */

import (
	"fmt"
)

import (
	gopi ".." /* import "github.com/djthorpe/gopi" */
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Size of something
type EGLSize struct {
	Width  uint
	Height uint
}

// Point on the screen (or off the screen)
type EGLPoint struct {
	X int
	Y int
}

// Frame
type EGLFrame struct {
	EGLPoint
	EGLSize
}

// Abstract driver interface
type EGLDriver interface {
	// Inherit general driver interface
	gopi.Driver

	// Return Major and Minor version of EGL
	GetVersion() (int, int)

	// Return Vendor information
	GetVendorString() string

	// Return list of supported extensions
	GetExtensions() []string

	// Return list of supported Client APIs
	GetSupportedClientAPIs() []string

	// Bind to an API
	BindAPI(api string) error

	// Return Bound API
	QueryAPI() (string, error)

	// Return display size
	GetFrame() EGLFrame

	// Create Background
	CreateBackground(api string, opacity float32) (EGLSurface, error)

	// Create Surface
	CreateSurface(api string, size EGLSize, origin EGLPoint, layer uint16, opacity float32) (EGLSurface, error)

	// Destroy Surface
	DestroySurface(surface EGLSurface) error

	// Move surface origin relative to current origin
	MoveSurfaceOriginBy(surface EGLSurface, rel EGLPoint) error

	// Flush surface updates to screen
	FlushSurface(surface EGLSurface) error

	// Set current surface context
	SetCurrentContext(surface EGLSurface) error
}

// Abstract drawable surface
type EGLSurface interface {
	// Return origin (NW value)
	GetOrigin() EGLPoint

	// Return window size
	GetSize() EGLSize

	// Is the surface the background?
	IsBackgroundSurface() bool

	// Return layer the surface is on
	GetLayer() uint16

	// Return the bitmap associated with the surface (if the
	// surface represents a bitmap resource)
	GetBitmap() (EGLBitmap, error)
}

////////////////////////////////////////////////////////////////////////////////
// String() methods

func (this EGLSize) String() string {
	return fmt.Sprintf("<EGLSize>{%v,%v}", this.Width, this.Height)
}

func (this EGLPoint) String() string {
	return fmt.Sprintf("<EGLPoint>{%v,%v}", this.X, this.Y)
}

func (this EGLFrame) String() string {
	return fmt.Sprintf("<EGLFrame>{%v,%v}", this.EGLPoint, this.EGLSize)
}

////////////////////////////////////////////////////////////////////////////////
// Point, Size and Frame methods

// Return the result of adding two points
func (this EGLPoint) Add(that EGLPoint) EGLPoint {
	return EGLPoint{this.X + that.X, this.Y + that.Y}
}
