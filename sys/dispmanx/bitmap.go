// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package dispmanx

import (
	"fmt"
	"image/color"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type ImageType rpi.DXImageType

type Bitmap interface {
	// Bitmap properties
	Bounds() rpi.DXRect
	Size() rpi.DXSize
	Centre() rpi.DXPoint
	NorthWest() rpi.DXPoint
	SouthWest() rpi.DXPoint
	NorthEast() rpi.DXPoint
	SouthEast() rpi.DXPoint

	// ClearToColor clears the screen
	ClearToColor(color.Color)

	// PaintPixel changes a single pixel
	PaintPixel(color.Color, rpi.DXPoint)

	// PaintCircle paints a filled circle with origin and radius
	PaintCircle(color.Color, rpi.DXPoint, uint32)

	// PaintLine paints a line
	PaintLine(color.Color, rpi.DXPoint, rpi.DXPoint)
}

type bitmap struct {
	handle        rpi.DXResource
	imageType     ImageType
	bounds        rpi.DXRect
	bytesPerPixel uint32
	stride        uint32
	data          *rpi.DXData
	dirty         rpi.DXRect

	RetainCount
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Supported image types
	IMAGE_TYPE_NONE   = ImageType(rpi.DX_IMAGE_TYPE_NONE)
	IMAGE_TYPE_RGB565 = ImageType(rpi.DX_IMAGE_TYPE_RGB565)
	IMAGE_TYPE_RGB888 = ImageType(rpi.DX_IMAGE_TYPE_RGB888)
	IMAGE_TYPE_RGBA32 = ImageType(rpi.DX_IMAGE_TYPE_RGBA32)
)

////////////////////////////////////////////////////////////////////////////////
// NEW AND CLOSE

func NewBitmap(imageType ImageType, w, h uint32) (Bitmap, error) {
	this := new(bitmap)
	if imageType == IMAGE_TYPE_NONE {
		return nil, gopi.ErrBadParameter.WithPrefix("imageType")
	} else if w == 0 || h == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("size")
	}
	if handle, err := rpi.DXResourceCreate(rpi.DXImageType(imageType), rpi.DXSize{w, h}); err != nil {
		return nil, err
	} else {
		this.handle = handle
		this.imageType = imageType
		this.bounds = rpi.DXNewRect(0, 0, w, h)
		this.bytesPerPixel = bytesPerPixelForImageType(imageType)
		this.stride = rpi.DXAlignUp(w, 16) * this.bytesPerPixel
		this.dirty = nil
	}

	// Success
	return this, nil
}

func (this *bitmap) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.data != nil {
		this.data.Free()
	}
	if this.handle == rpi.DX_NO_HANDLE {
		return nil
	}
	err := rpi.DXResourceDelete(this.handle)
	this.handle = rpi.DX_NO_HANDLE
	this.data = nil
	return err
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *bitmap) Bounds() rpi.DXRect {
	return this.bounds
}

func (this *bitmap) Size() rpi.DXSize {
	return rpi.DXRectSize(this.bounds)
}

func (this *bitmap) Centre() rpi.DXPoint {
	size := rpi.DXRectSize(this.bounds)
	return rpi.DXPoint{int32(size.W >> 1), int32(size.H >> 1)}
}

func (this *bitmap) NorthWest() rpi.DXPoint {
	return rpi.DXRectOrigin(this.bounds)
}

func (this *bitmap) SouthEast() rpi.DXPoint {
	origin := rpi.DXRectOrigin(this.bounds)
	size := rpi.DXRectSize(this.bounds)
	return origin.Add(rpi.DXPoint{int32(size.W) - 1, int32(size.H) - 1})
}

func (this *bitmap) SouthWest() rpi.DXPoint {
	origin := rpi.DXRectOrigin(this.bounds)
	size := rpi.DXRectSize(this.bounds)
	return origin.Add(rpi.DXPoint{0, int32(size.H) - 1})
}

func (this *bitmap) NorthEast() rpi.DXPoint {
	origin := rpi.DXRectOrigin(this.bounds)
	size := rpi.DXRectSize(this.bounds)
	return origin.Add(rpi.DXPoint{int32(size.W) - 1, 0})
}

// Retain increments the counter and returns the resource handle
func (this *bitmap) Retain() rpi.DXResource {
	this.RetainCount.Inc()
	return this.handle
}

// Release returns true if the bitmap should be released
func (this *bitmap) Release() bool {
	return this.RetainCount.Dec()
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE ROWS

func (this *bitmap) ReadRows(offset, height uint32, read bool) (*rpi.DXData, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Get size of bitmap
	size := rpi.DXRectSize(this.bounds)

	// Check parameters
	if height == 0 || height > size.H {
		return nil, gopi.ErrBadParameter.WithPrefix("height")
	}
	if offset+height > size.H {
		return nil, gopi.ErrBadParameter.WithPrefix("offset")
	}

	// Create a data buffer or use existing one
	cap := uint(height * this.stride)
	if this.data == nil || this.data.Cap() < cap {
		if this.data != nil {
			this.data.Free()
		}
		if this.data = rpi.DXNewData(cap); this.data == nil {
			return nil, gopi.ErrInternalAppError.WithPrefix("DXNewData")
		}
	}

	// Read data - we offset using ptr and set rect to 0
	if read {
		ptr := this.data.Ptr() - uintptr(offset*this.stride)
		rect := rpi.DXNewRect(0, int32(offset), size.W, height)
		if err := rpi.DXResourceReadData(this.handle, rect, ptr, this.stride); err != nil {
			return nil, err
		}
	}

	// Return success
	return this.data, nil
}

func (this *bitmap) WriteRows(data *rpi.DXData, offset uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Get size of bitmap
	size := rpi.DXRectSize(this.bounds)

	// Check parameters
	if data == nil {
		return gopi.ErrBadParameter.WithPrefix("data")
	}
	if offset >= size.H {
		return gopi.ErrBadParameter.WithPrefix("offset")
	}
	if uint32(data.Cap())%this.stride != 0 {
		return gopi.ErrBadParameter.WithPrefix("stride")
	}

	// determine how many rows to write
	height := uint32(data.Cap()) / this.stride
	if (offset + height) > size.H {
		return gopi.ErrBadParameter.WithPrefix("offset")
	}

	// Set the pointer to the strip and move y forward and ptr back for each strip
	ptr := data.Ptr() - uintptr(offset*this.stride)
	rect := rpi.DXNewRect(0, int32(offset), size.W, height)

	// Write data
	if err := rpi.DXResourceWriteData(this.handle, rpi.DXImageType(this.imageType), this.stride, ptr, rect); err != nil {
		return err
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *bitmap) String() string {
	return "<bitmap" +
		" image_type=" + fmt.Sprint(rpi.DXImageType(this.imageType)) +
		" size=" + fmt.Sprint(rpi.DXRectSize(this.bounds)) +
		" stride=" + fmt.Sprint(this.stride) +
		">"
}

////////////////////////////////////////////////////////////////////////////////
// DIRTY

func (this *bitmap) setDirty(rect rpi.DXRect) {
	this.Mutex.Lock()
	this.dirty = nil
	this.Mutex.Unlock()
	this.addDirty(rect)
}

func (this *bitmap) addDirty(rect rpi.DXRect) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	fmt.Println("-> dirty=>", rpi.DXRectString(this.dirty), " add=>", rpi.DXRectString(rect))
	if this.dirty == nil || rect == nil {
		this.dirty = rect
	} else {
		this.dirty = rpi.DXRectUnion(this.dirty, rect)
	}
	// Clip to bitmap size
	if this.dirty != nil {
		this.dirty = rpi.DXRectIntersection(this.bounds, this.dirty)
	}
	fmt.Println("<- dirty=>", rpi.DXRectString(this.dirty))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func bytesPerPixelForImageType(imageType ImageType) uint32 {
	switch imageType {
	case IMAGE_TYPE_RGBA32:
		return 4
	case IMAGE_TYPE_RGB888:
		return 3
	case IMAGE_TYPE_RGB565:
		return 2
	default:
		return 0
	}
}

func colorToBytes(t ImageType, c color.Color) []byte {
	// Returns color 0000 <= v <= FFFF
	r, g, b, a := c.RGBA()
	// Convert to []byte
	switch t {
	case IMAGE_TYPE_RGB888:
		return []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8)}
	case IMAGE_TYPE_RGB565:
		r := uint16(r>>(8+3)) << (5 + 6)
		g := uint16(g>>(8+2)) << 5
		b := uint16(b >> (8 + 3))
		v := r | g | b
		return []byte{byte(v), byte(v >> 8)}
	case IMAGE_TYPE_RGBA32:
		return []byte{byte(r >> 8), byte(g >> 8), byte(b >> 8), byte(a >> 8)}
	default:
		return nil
	}
}
