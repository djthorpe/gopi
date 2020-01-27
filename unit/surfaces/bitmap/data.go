// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap

import (
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPE

type Data struct {
	data *rpi.DXData
	size uint
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// NEW AND FREE

// Set capacity for the data
func (this *Data) SetCapacity(size uint) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	return this.setCapacity(size)
}

func (this *Data) setCapacity(size uint) error {
	if size == 0 || (this.data != nil && this.data.Cap() < size) {
		if this.data != nil {
			this.data.Free()
		}
		this.data = nil
		this.size = 0
	}
	if this.data == nil && size > 0 {
		this.data = rpi.DXNewData(size)
		if this.data == nil {
			return gopi.ErrInternalAppError.WithPrefix("Cap")
		}
	}

	// Set new size
	this.size = size

	// Success
	return nil
}

func (this *Data) Bytes() []byte {
	if this.data == nil {
		return nil
	} else {
		return this.data.Bytes(this.size)
	}
}

func (this *Data) Uint32() []uint32 {
	if this.data == nil {
		return nil
	} else {
		return this.data.Uint32()
	}
}

func (this *Data) Ptr() uintptr {
	if this.data == nil {
		return 0
	} else {
		return this.data.Ptr()
	}
}

func (this *Data) FillByte(value byte) {
	bytes := this.Bytes()
	for i := range bytes {
		bytes[i] = value
	}
}

func (this *Data) FillUint32(value uint32) {
	uints := this.Uint32()
	for i := range uints {
		uints[i] = value
	}
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE ROWS

// Read will read bitmap data from GPU memory assuming number of bytes
// as width ('stride')
func (this *Data) Read(handle rpi.DXResource, offset, length uint, stride uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check parameters
	if handle == rpi.DX_NO_HANDLE || length == 0 || stride == 0 {
		return gopi.ErrBadParameter.WithPrefix("Read")
	}
	if stride != rpi.DXAlignUp(stride, 16) {
		// stride should be on 16-byte boundary
		return gopi.ErrBadParameter.WithPrefix("stride")
	}
	// Extend capacity of buffer as necessary
	if err := this.setCapacity(length * uint(stride)); err != nil {
		return err
	}
	// Read data
	ptr := this.data.Ptr() - uintptr(offset*uint(stride))
	rect := rpi.DXNewRect(0, int32(offset), uint32(stride), uint32(length))
	return rpi.DXResourceReadData(handle, rect, ptr, stride)
}

func (this *Data) Write(handle rpi.DXResource, dxmode rpi.DXImageType, offset, length uint, stride uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check parameters
	if handle == rpi.DX_NO_HANDLE || length == 0 || stride == 0 {
		return gopi.ErrBadParameter.WithPrefix("Write")
	}
	if stride != rpi.DXAlignUp(stride, 16) {
		// stride should be on 16-byte boundary
		return gopi.ErrBadParameter.WithPrefix("stride")
	}
	if this.data == nil {
		return gopi.ErrInternalAppError.WithPrefix("Write")
	}

	// Set the pointer to the strip and move y forward and ptr back for each strip
	ptr := this.data.Ptr() - uintptr(offset*uint(stride))
	rect := rpi.DXNewRect(0, int32(offset), uint32(stride), uint32(length))
	return rpi.DXResourceWriteData(handle, dxmode, stride, ptr, rect)
}
