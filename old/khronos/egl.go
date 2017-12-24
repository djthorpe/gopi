/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// EGL
//
// This package defines abstract surface, which can be bitmap, 2D or 3D vector.
// It also defines types for points, sizes and frames and associated calculations
//
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

// Frame alignment
type EGLFrameAlignFlag uint8

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
	CreateCursor() (EGLSurface, error)
	CreateSurface(api string, size EGLSize, origin EGLPoint, layer uint16, opacity float32) (EGLSurface, error)
	CreateSurfaceWithBitmap(bitmap EGLBitmap, origin EGLPoint, layer uint16, opacity float32) (EGLSurface, error)

	// Destroy Surface
	DestroySurface(surface EGLSurface) error

	// Create Bitmap resource from an image
	CreateImage(r io.Reader) (EGLBitmap, error)

	// Write Bitmap out to stream as PNG
	WriteImagePNG(w io.Writer, bitmap EGLBitmap) error

	// Destroy Bitmap resource
	DestroyImage(bitmap EGLBitmap) error

	// Create a bitmap with a copy of the screen contents
	SnapshotImage() (EGLBitmap, error)

	// Move surface origin relative to current origin
	MoveSurfaceOriginBy(surface EGLSurface, rel EGLPoint) error

	// Move surface origin absolutely
	MoveSurfaceOriginTo(surface EGLSurface, abs EGLPoint) error

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
// CONSTANTS

const (
	EGL_ALIGN_VCENTER EGLFrameAlignFlag = 1 << iota
	EGL_ALIGN_TOP                       = 1 << iota
	EGL_ALIGN_BOTTOM                    = 1 << iota
	EGL_ALIGN_HCENTER                   = 1 << iota
	EGL_ALIGN_LEFT                      = 1 << iota
	EGL_ALIGN_RIGHT                     = 1 << iota
	EGL_ALIGN_CENTER                    = EGL_ALIGN_VCENTER | EGL_ALIGN_HCENTER
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	EGLZeroPoint  = EGLPoint{0, 0}
	EGLZeroSize   = EGLSize{0, 0}
	EGLWhiteColor = EGLColorRGBA32{0xFF, 0xFF, 0xFF, 0xFF}
	EGLRedColor   = EGLColorRGBA32{0xFF, 0x00, 0x00, 0xFF}
	EGLGreenColor = EGLColorRGBA32{0x00, 0xFF, 0x00, 0xFF}
	EGLBlueColor  = EGLColorRGBA32{0x00, 0x00, 0xFF, 0xFF}
	EGLBlackColor = EGLColorRGBA32{0x00, 0x00, 0x00, 0xFF}
	EGLGreyColor  = EGLColorRGBA32{0x80, 0x80, 0x80, 0xFF}
)

////////////////////////////////////////////////////////////////////////////////
// METHODS String() methods

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
	return fmt.Sprintf("<EGLColorRGBA32>{r=%02X g=%02X b=%02X a=%02X}", this.R, this.G, this.B, this.A)
}

////////////////////////////////////////////////////////////////////////////////
// METHODS Point methods

// Return the result of adding two points
func (this EGLPoint) Add(that EGLPoint) EGLPoint {
	return EGLPoint{this.X + that.X, this.Y + that.Y}
}

// Return the result of adding a size to a point
func (this EGLPoint) Offset(that EGLSize) EGLPoint {
	return EGLPoint{this.X + int(that.Width), this.Y + int(that.Height)}
}

// Return boolean value that determines if point is within a frame
func (this EGLPoint) InFrame(that EGLFrame) bool {
	if this.X < that.X || this.Y < that.Y {
		return false
	}
	if this.X >= that.X+int(that.Width) {
		return false
	}
	if this.Y >= that.Y+int(that.Height) {
		return false
	}
	return true
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
// METHODS Frame methods

// Return origin of a frame
func (this EGLFrame) Origin() EGLPoint {
	return EGLPoint{this.X, this.Y}
}

// Return size of a frame
func (this EGLFrame) Size() EGLSize {
	return EGLSize{this.Width, this.Height}
}

// Aligns origin of rectangle to another frame, and returns new frame. The other
// frame is passed as a pointer, so no need to copy the frame.
func (this EGLFrame) AlignTo(that *EGLFrame, flags EGLFrameAlignFlag) EGLFrame {
	// Vertical
	switch {
	case (flags & EGL_ALIGN_VCENTER) == EGL_ALIGN_VCENTER:
		this.Y = that.Y + ((int(that.Height) - int(this.Height)) / 2)
		break
	case (flags & EGL_ALIGN_TOP) == EGL_ALIGN_TOP:
		this.Y = that.Y
		break
	case (flags & EGL_ALIGN_BOTTOM) == EGL_ALIGN_BOTTOM:
		this.Y = that.Y + int(that.Height) - int(this.Height)
		break
	}
	// Horizontal
	switch {
	case (flags & EGL_ALIGN_HCENTER) == EGL_ALIGN_HCENTER:
		this.X = that.X + ((int(that.Width) - int(this.Width)) / 2)
		break
	case (flags & EGL_ALIGN_LEFT) == EGL_ALIGN_LEFT:
		this.X = that.X
		break
	case (flags & EGL_ALIGN_RIGHT) == EGL_ALIGN_RIGHT:
		this.X = that.X + int(that.Width) - int(this.Width)
		break
	}
	return this
}

////////////////////////////////////////////////////////////////////////////////
// METHODS Color methods

func (color EGLColorRGBA32) Uint32() uint32 {
	return uint32(color.A)<<24 + uint32(color.B)<<16 + uint32(color.G)<<8 + uint32(color.R)
}
