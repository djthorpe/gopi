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
	"image/color"
	"os"
	"sync"

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
	log          gopi.Logger
	lock         sync.Mutex
	handle       dxResourceHandle
	image_type   dxImageType
	width        uint32
	height       uint32
	stride_bytes uint32 // number of bytes per row rounded up to 16-byte boundaries
}

////////////////////////////////////////////////////////////////////////////////
// ERRORS

var (
	ErrUnsupportedImageType = errors.New("Unsupported Image Type")
	ErrInvalidResource      = errors.New("Invalid resource")
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config Resource) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.graphics.rpi.Bitmap.Open>{ image_type=%v size={ %v,%v } }", config.ImageType, config.Width, config.Height)

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
	this.stride_bytes = dxAlignUp(uint32(config.Width), uint32(16)) * 4 // uint32 to bytes

	if handle, err := dxResourceCreate(this.image_type, this.width, this.height); err != DX_SUCCESS {
		return nil, os.NewSyscallError("dxResourceCreate", err)
	} else {
		this.handle = handle
	}

	return this, nil
}

func (this *resource) Close() error {
	this.log.Debug("<sys.graphics.rpi.Bitmap.Close>{ handle=%v }", this.handle)

	// If already closed
	if this.handle == dxResourceHandle(DX_NO_RESOURCE) {
		return nil
	}

	// Lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Delete resource
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
		return "<sys.graphics.rpi.Bitmap>{ nil }"
	} else {
		return fmt.Sprintf("<sys.graphics.rpi.Bitmap>{ size={ %v,%v } image_type=%v stride_bytes=%v handle=%v }", this.width, this.height, this.image_type, this.stride_bytes, this.handle)
	}
}

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

func (this *resource) ClearToColorRGBA(color color.RGBA) error {

	// Checks
	if this.handle == dxResourceHandle(DX_NO_RESOURCE) {
		return ErrInvalidResource
	}

	// Check for correct image type
	if this.image_type != DX_IMAGETYPE_RGBA32 {
		return ErrUnsupportedImageType
	}

	// Lock
	this.lock.Lock()
	defer this.lock.Unlock()

	// Clear buffer to color
	data := make([]uint32, this.height*(this.stride_bytes>>4))
	value := dxToRGBA32(color)
	for i := 0; i < len(data); i++ {
		data[i] = value
	}

	// Write data
	if err := dxResourceWriteDataUint32(this.handle, data); err != DX_SUCCESS {
		return os.NewSyscallError("dxResourceWriteDataUint32", err)
	}

	// Success
	return nil
}

func (this *resource) Type() gopi.SurfaceType {
	// Checks
	if this.handle == dxResourceHandle(DX_NO_RESOURCE) {
		return gopi.SURFACE_TYPE_NONE
	}

	// Return image type
	switch this.image_type {
	case DX_IMAGETYPE_RGBA32:
		return gopi.SURFACE_TYPE_RGBA32
	default:
		return gopi.SURFACE_TYPE_NONE
	}
}

func (this *resource) Size() gopi.Size {
	// Checks
	if this.handle == dxResourceHandle(DX_NO_RESOURCE) {
		return gopi.ZeroSize
	} else {
		return gopi.Size{float32(this.width), float32(this.height)}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Align a width (in 1 pixel per uint32) up to a byte boundary
func dxAlignUp(value, alignment uint32) uint32 {
	return ((value - 1) & ^(alignment - 1)) + alignment
}

// Convert naitive image/color type into uint32
func dxToRGBA32(color color.RGBA) uint32 {
	return uint32(color.A)<<24 | uint32(color.B)<<16 | uint32(color.G)<<8 | uint32(color.R)
}

// dxResourceWriteDataUint32 writes into resource data
func dxResourceWriteDataUint32(resource dxResourceHandle, data []uint32) error {
	return gopi.ErrNotImplemented
}
