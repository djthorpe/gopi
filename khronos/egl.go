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
	Width uint
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
	GetVersion() (int,int)

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
	CreateBackground(api string) (EGLWindow,error)

	// Create Window
	CreateWindow(api string,size EGLSize,origin EGLPoint,layer uint16) (EGLWindow,error)

	// Close window
	CloseWindow(window EGLWindow) error
}

// Abstract window
type EGLWindow interface {
	// Put in window functions here
}

////////////////////////////////////////////////////////////////////////////////
// String() methods

func (this EGLSize) String() string {
	return fmt.Sprintf("<EGLSize>{%v,%v}",this.Width,this.Height)
}

func (this EGLPoint) String() string {
	return fmt.Sprintf("<EGLPoint>{%v,%v}",this.X,this.Y)
}

func (this EGLFrame) String() string {
	return fmt.Sprintf("<EGLFrame>{%v,%v}",this.EGLPoint,this.EGLSize)
}

