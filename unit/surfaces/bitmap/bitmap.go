// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap

import (
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type bitmap struct {
	mode      gopi.SurfaceFlags
	dxmode    rpi.DXImageType
	size      rpi.DXSize
	bounds    rpi.DXRect
	pixelSize uint32
	stride    uint32
	handle    rpi.DXResource
	dirty 	rpi.DXRect

	Data
	RetainCount
	sync.Mutex
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Config) Name() string { return "gopi/surfaces/bitmap" }

func (config Config) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(bitmap)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *bitmap) Init(config Config) error {
	// Check size
	size := rpi.DXSize{uint32(config.Size.W), uint32(config.Size.H)}
	if size.W == 0 || size.H == 0 {
		return gopi.ErrBadParameter.WithPrefix("size")
	} else {
		this.size = size
	}

	// Check image mode
	switch config.Mode {
	case gopi.SURFACE_FLAG_RGBA32:
		this.mode = config.Mode
		this.dxmode = rpi.DX_IMAGE_TYPE_RGBA32
		this.pixelSize = 4
	case gopi.SURFACE_FLAG_RGB888:
		this.mode = config.Mode
		this.dxmode = rpi.DX_IMAGE_TYPE_RGB888
		this.pixelSize = 3
	case gopi.SURFACE_FLAG_RGB565:
		this.mode = config.Mode
		this.dxmode = rpi.DX_IMAGE_TYPE_RGB565
		this.pixelSize = 2
	default:
		return gopi.ErrBadParameter.WithPrefix("mode")
	}

	// Set stride and bounds
	this.stride = rpi.DXAlignUp(size.W*this.pixelSize, 16)
	this.bounds = rpi.DXNewRect(0,0,uint32(config.Size.W), uint32(config.Size.H))
	this.dirty = nil

	// Create the resource
	if handle, err := rpi.DXResourceCreate(this.dxmode, size); err != nil {
		return err
	} else {
		this.handle = handle
	}

	// Success
	return nil
}

func (this *bitmap) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Free any allocated data
	this.Data.SetCapacity(0)

	// Delete resources
	errs := gopi.NewCompoundError()
	if this.handle != rpi.DX_NO_HANDLE {
		errs.Add(rpi.DXResourceDelete(this.handle))
	}

	// Close
	errs.Add(this.Unit.Close())

	// Release resources
	this.handle = rpi.DX_NO_HANDLE

	// Return sucess
	return errs.ErrorOrSelf()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *bitmap) String() string {
	if this.handle == rpi.DX_NO_HANDLE {
		return "<" + Config{}.Name() + " handle=nil>"
	} else {
		return "<" + Config{}.Name() +
			" handle=" + fmt.Sprint(this.handle) +
			" mode=" + fmt.Sprint(this.dxmode) +
			" size=" + fmt.Sprint(this.size) +
			" bytes_per_pixel=" + fmt.Sprint(this.pixelSize) +
			" stride=" + fmt.Sprint(this.stride) +
			">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *bitmap) DXHandle() rpi.DXResource {
	return this.handle
}

func (this *bitmap) Mode() gopi.SurfaceFlags {
	return this.mode
}

// Origin always returns the gopi.ZeroPoint
func (this *bitmap) Origin() gopi.Point {
	return gopi.ZeroPoint
}

func (this *bitmap) Size() gopi.Size {
	return gopi.Size{float32(this.size.W), float32(this.size.H)}
}

func (this *bitmap) DXSize() rpi.DXSize {
	return this.size
}

func (this *bitmap) DXRect() rpi.DXRect {
	return this.bounds
}

// Bytes returns the bitmap as bytes with bytes per row
func (this *bitmap) Bytes() ([]byte, uint32) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := this.Data.Read(this.handle,0,uint(this.size.H),this.stride); err != nil {
		return nil, 0
	} else {
		return this.Data.Bytes(), this.stride
	}
}

////////////////////////////////////////////////////////////////////////////////
// POINTS

func (this *bitmap) Centre() gopi.Point {
	return gopi.Point{float32(this.size.W) / 2.0, float32(this.size.H) / 2.0}
}

func (this *bitmap) NorthWest() gopi.Point {
	return gopi.Point{0,0}
}

func (this *bitmap) SouthWest() gopi.Point {
	return gopi.Point{0,float32(this.size.H-1)}
}

func (this *bitmap) NorthEast() gopi.Point {
	return gopi.Point{float32(this.size.W-1),0}
}

func (this *bitmap) SouthEast() gopi.Point {
	return gopi.Point{float32(this.size.W-1),float32(this.size.H-1)}
}

////////////////////////////////////////////////////////////////////////////////
// RETAIN AND RELEASE

func (this *bitmap) Retain() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	this.Log.Debug("<" + Config{}.Name() + " handle=" + fmt.Sprint(this.handle) + "> RETAIN")
	this.RetainCount.Inc()
}

func (this *bitmap) Release() bool {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.handle == rpi.DX_NO_HANDLE {
		this.Log.Warn("Should not call Release on closed bitmap")
		return false
	} else {
		release := this.RetainCount.Dec()
		this.Log.Debug("<" + Config{}.Name() + " handle=" + fmt.Sprint(this.handle) + "> RELEASE " + fmt.Sprint(release))
		return release	
	}
}

////////////////////////////////////////////////////////////////////////////////
// DIRTY RECTS

func (this *bitmap) setDirty(rect rpi.DXRect) {
	this.Mutex.Lock()
	this.dirty = nil
	this.Mutex.Unlock()
	this.addDirty(rect)
}

func (this *bitmap) addDirty(rect rpi.DXRect) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	//fmt.Println("-> dirty=>", rpi.DXRectString(this.dirty), " add=>", rpi.DXRectString(rect))
	if this.dirty == nil || rect == nil {
		this.dirty = rect
	} else {
		this.dirty = rpi.DXRectUnion(this.dirty, rect)
	}
	// Clip to bitmap size
	if this.dirty != nil {
		this.dirty = rpi.DXRectIntersection(this.bounds, this.dirty)
	}
	//fmt.Println("<- dirty=>", rpi.DXRectString(this.dirty))
}
