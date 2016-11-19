/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package khronos /* import "github.com/djthorpe/gopi/khronos" */

import (
	"fmt"
	"io"
)

import (
	gopi "github.com/djthorpe/gopi"
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

// EGLColorRGBA32
type EGLColorRGBA32 struct {
	R uint8
	G uint8
	B uint8
	A uint8
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

	// Create Background & Surfaces
	CreateBackground(api string, opacity float32) (EGLSurface, error)
	CreateSurface(api string, size EGLSize, origin EGLPoint, layer uint16, opacity float32) (EGLSurface, error)
	CreateSurfaceWithBitmap(bitmap EGLBitmap, origin EGLPoint, layer uint16, opacity float32) (EGLSurface, error)
	
	// Destroy Surface
	DestroySurface(surface EGLSurface) error

	// Create bitmap resource
	CreateImage(r io.Reader) (EGLBitmap, error)

	// Destroy bitmap resource
	DestroyImage(bitmap EGLBitmap) error

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
// Variables

var (
	EGLZeroPoint = EGLPoint{ 0, 0 }
	EGLWhiteColor = EGLColorRGBA32{ 0xFF, 0xFF, 0xFF, 0xFF }
	EGLRedColor = EGLColorRGBA32{ 0xFF, 0x00, 0x00, 0xFF }
	EGLGreenColor = EGLColorRGBA32{ 0x00, 0xFF, 0x00, 0xFF }
	EGLBlueColor = EGLColorRGBA32{ 0x00, 0x00, 0xFF, 0xFF }
	EGLBlackColor = EGLColorRGBA32{ 0x00, 0x00, 0x00, 0xFF }
	EGLGreyColor = EGLColorRGBA32{ 0x80, 0x80, 0x80, 0xFF }
)

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

func (this EGLColorRGBA32) String() string {
	return fmt.Sprintf("<EGLColorRGBA32>{r=%02X g=%02X b=%02X a=%02X}",this.R,this.G,this.B,this.A)
}

////////////////////////////////////////////////////////////////////////////////
// Point, Size and Frame methods

// Return the result of adding two points
func (this EGLPoint) Add(that EGLPoint) EGLPoint {
	return EGLPoint{this.X + that.X, this.Y + that.Y}
}

// Return if equals
func (this EGLPoint) Equals(that EGLPoint) bool {
	if this.X != that.X {
		return false
	}
	if this.Y != that.Y {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// Color methods

func (color EGLColorRGBA32) Uint32() uint32 {
	return uint32(color.A) << 24 + uint32(color.B) << 16 + uint32(color.G) << 8 + uint32(color.R)
}

