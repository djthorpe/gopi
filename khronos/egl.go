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
	GetFrame() *EGLFrame

	// Create Background
	CreateBackground(api string) (EGLWindow,error)

	// Create Window
	CreateWindow(api string,size *EGLSize,origin *EGLPoint) (EGLWindow,error)

	// Close window
	CloseWindow(window EGLWindow) error
}

// Abstract window
type EGLWindow interface {

}

////////////////////////////////////////////////////////////////////////////////
// Driver interface

// Return EGL version number
func (this *gopi.State) GetVersion() (int, int) {
	return this.driver.GetVersion()
}

// Return Vendor information
func (this *gopi.State) GetVendorString() string {
	return this.driver.GetVendorString()
}

// Return Extension information
func (this *gopi.State) GetExtensions() []string {
	return this.driver.GetExtensions()
}

// Return Client API information
func (this *gopi.State) GetSupportedClientAPIs() []string {
	return this.driver.GetSupportedClientAPIs()
}

// Bind to an API
func (this *gopi.State) BindAPI(api string) error {
	return this.driver.BindAPI(api)
}

// Bound API
func (this *gopi.State) QueryAPI() (string, error) {
	return this.driver.QueryAPI()
}

// Return size of display
func (this *gopi.State) GetFrame() *EGLFrame {
	return this.driver.GetFrame()
}

// Create Window
func (this *gopi.State) CreateWindow(api string,size *EGLSize,origin *EGLPoint) (EGLWindow,error) {
	return this.driver.CreateWindow(api,size,origin)
}

// Create Background
func (this *gopi.State) CreateBackground(api string) (EGLWindow,error) {
	return this.driver.CreateBackground(api)
}

// Close a window
func (this *gopi.State) CloseWindow(window EGLWindow) error {
	return this.driver.CloseWindow(window)
}

////////////////////////////////////////////////////////////////////////////////
// String() methods
// Size of something

func (this EGLSize) String() string {
	return fmt.Sprintf("<EGLSize>{%v,%v}",this.Width,this.Height)
}

func (this EGLPoint) String() string {
	return fmt.Sprintf("<EGLPoint>{%v,%v}",this.X,this.Y)
}

func (this EGLFrame) String() string {
	return fmt.Sprintf("<EGLFrame>{%v,%v}",this.EGLPoint,this.EGLSize)
}

func (this *EGLState) String() string {
	return fmt.Sprintf("<EGLState>{%v}",this.driver)
}

