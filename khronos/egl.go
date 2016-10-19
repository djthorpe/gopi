/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package khronos

import (
	"fmt"
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

// Store state for the non-abstract input driver
type EGLState struct {
	driver EGLDriver
}

// Abstract driver interface
type EGLDriver interface {
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

	// Close closes the driver and frees the underlying resources
	Close() error
}

// Abstract configuration which is used to open and return the
// concrete driver
type EGLConfig interface {
	// Opens the driver from configuration, or returns error
	Open() (EGLDriver, error)
}

// Abstract window
type EGLWindow interface {

}

////////////////////////////////////////////////////////////////////////////////
// Opener interface

// Open opens a connection EGL
func Open(config EGLConfig) (EGLDriver, error) {
	driver, err := config.Open()
	if err != nil {
		return nil, err
	}
	return &EGLState{ driver }, nil
}

////////////////////////////////////////////////////////////////////////////////
// Driver interface

// Closes the device and frees the resources
func (this *EGLState) Close() error {
	return this.driver.Close()
}

// Return EGL version number
func (this *EGLState) GetVersion() (int, int) {
	return this.driver.GetVersion()
}

// Return Vendor information
func (this *EGLState) GetVendorString() string {
	return this.driver.GetVendorString()
}

// Return Extension information
func (this *EGLState) GetExtensions() []string {
	return this.driver.GetExtensions()
}

// Return Client API information
func (this *EGLState) GetSupportedClientAPIs() []string {
	return this.driver.GetSupportedClientAPIs()
}

// Bind to an API
func (this *EGLState) BindAPI(api string) error {
	return this.driver.BindAPI(api)
}

// Bound API
func (this *EGLState) QueryAPI() (string, error) {
	return this.driver.QueryAPI()
}

// Return size of display
func (this *EGLState) GetFrame() *EGLFrame {
	return this.driver.GetFrame()
}

// Create Window
func (this *EGLState) CreateWindow(api string,size *EGLSize,origin *EGLPoint) (EGLWindow,error) {
	return this.driver.CreateWindow(api,size,origin)
}

// Create Background
func (this *EGLState) CreateBackground(api string) (EGLWindow,error) {
	return this.driver.CreateBackground(api)
}

// Close a window
func (this *EGLState) CloseWindow(window EGLWindow) error {
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

