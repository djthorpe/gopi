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

	// Close closes the driver and frees the underlying resources
	Close() error

	// Do Stuff
	Do() error
}

// Abstract configuration which is used to open and return the
// concrete driver
type EGLConfig interface {
	// Opens the driver from configuration, or returns error
	Open() (EGLDriver, error)
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

// Provides human-readable version
func (this *EGLState) String() string {
	return fmt.Sprintf("<EGLState>{%v}",this.driver)
}

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

// Do stuff
func (this *EGLState) Do() error {
	return this.driver.Do()
}



