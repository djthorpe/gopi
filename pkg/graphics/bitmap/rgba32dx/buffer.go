// +build dispmanx

package rgba32dx

import (
	"fmt"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Buffer struct {
	data   *dx.Data
	pixels []uint32
	rect   *dx.Rect
	stride uint32
	y      uint32
	dirty  bool
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Buffer) Init(stride uint32) error {
	if data := dx.NewData(stride); data == nil {
		return gopi.ErrInternalAppError
	} else {
		this.data = data
		this.pixels = data.Uint32(0)
	}

	// Set buffer parameters
	this.stride = stride
	this.y = 0
	this.rect = dx.NewRect(0, 0, this.stride, 1)

	// Return success
	return nil
}

func (this *Buffer) Dispose() {
	this.data.Dispose()
	this.data = nil
	this.pixels = nil
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Buffer) Fill(v uint32) {
	if this.data != nil {
		this.data.SetUint32(0, this.stride>>2, v)
	}
}

func (this *Buffer) GetAt(x uint32) uint32 {
	return this.pixels[x]
}

func (this *Buffer) SetAt(x, c uint32) {
	this.pixels[x] = c
}

func (this *Buffer) WriteRow(resource dx.Resource, y uint32) error {
	if this.data == nil {
		return gopi.ErrInternalAppError
	} else if y == this.y {
		// Don't write if the row is the same as the last one, just mark
		// as dirty instead
		this.dirty = true
		return nil
	}

	if err := this.write(resource, y, 1); err != nil {
		return err
	} else {
		this.y = y
		this.dirty = false
		return nil
	}
}

func (this *Buffer) ReadRow(resource dx.Resource, y uint32) error {
	// Check parameters, don't read the same row
	if this.data == nil {
		return gopi.ErrInternalAppError
	} else if y == this.y {
		return nil
	}

	// we may need to write a row before reading the next
	if this.dirty {
		if err := this.write(resource, this.y, 1); err != nil {
			return err
		}
	}

	// Clear dirty flag
	this.dirty = false

	// Read the row
	if err := this.read(resource, y, 1); err != nil {
		return err
	} else {
		this.y = y
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Write data to GPU memory with the y-axis bounds as y and h
func (this *Buffer) write(dest dx.Resource, y, h uint32) error {
	this.rect.SetY(int32(y), h)
	return dx.ResourceWrite(dest, 0, this.stride, this.data.PtrMinusOffset(y*this.stride), this.rect)
}

// Read data from GPU to buffer
func (this *Buffer) read(src dx.Resource, y, h uint32) error {
	this.rect.SetY(int32(y), h)
	return dx.ResourceRead(src, this.rect, this.data.PtrMinusOffset(y*this.stride), this.stride)
}
