// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2018
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"errors"
	"fmt"
	"os"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Resource struct {
	ImageType dxImageType
	Width     uint32
	Height    uint32
}

type resource struct {
	log        gopi.Logger
	handle     dxResourceHandle
	image_type dxImageType
	width      uint32
	height     uint32
	stride     uint32 // number of bytes per row rounded up to 16-byte boundaries
}

////////////////////////////////////////////////////////////////////////////////
// ERRORS

var (
	ErrUnsupportedImageType = errors.New("Unsupported Image Type")
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Resource) Open(log gopi.Logger) (*resource, error) {
	log.Debug("<sys.surface.rpi.Bitmap.Open>{ image_type=%v size={ %v,%v } }", config.ImageType, config.Width, config.Height)

	// Check configuration parameters
	if config.Width == 0 || config.Height == 0 {
		return nil, gopi.ErrBadParameter
	}
	if config.ImageType != DX_IMAGETYPE_RGBA32 {
		return nil, ErrUnsupportedImageType
	}

	// Create resource
	this := new(resource)
	this.log = log
	this.image_type = config.ImageType
	this.width = config.Width
	this.height = config.Height
	this.stride = dxAlignUp(uint32(config.Width), uint32(16)) * 4 // uint32 to bytes

	if handle, err := dxResourceCreate(this.image_type, this.width, this.height); err != DX_SUCCESS {
		return nil, os.NewSyscallError("dxResourceCreate", err)
	} else {
		this.handle = handle
	}

	return this, nil
}

func (this *resource) Close() error {
	// If already closed
	if this.handle == dxResourceHandle(DX_NO_RESOURCE) {
		return nil
	}
	if err := dxResourceDelete(this.handle); err != DX_SUCCESS {
		return os.NewSyscallError("dxResourceDelete", err)
	} else {
		this.handle = dxResourceHandle(DX_NO_RESOURCE)
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *resource) String() string {
	if this.handle == dxResourceHandle(DX_NO_RESOURCE) {
		return "<sys.surface.rpi.Bitmap.Open>{ nil }"
	} else {
		return fmt.Sprintf("<sys.surface.rpi.Bitmap.Open>{ size={ %v,%v } image_type=%v stride=%v handle=%v }", this.width, this.height, this.image_type, this.stride, this.handle)
	}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func dxAlignUp(value, alignment uint32) uint32 {
	return ((value - 1) & ^(alignment - 1)) + alignment
}
