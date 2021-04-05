package bitmap

import (
	"fmt"
	"sync"

	// Modules
	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACE

type BitmapFactory interface {
	New(fmt gopi.SurfaceFormat, w, h uint32) (gopi.Bitmap, error)
	Dispose(gopi.Bitmap) error
}

type Manager struct {
	gopi.Unit
	gopi.Platform
	sync.Mutex

	bitmaps []gopi.Bitmap
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	factories = make(map[gopi.SurfaceFormat]BitmapFactory)
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Platform)
	return nil
}

func (this *Manager) Dispose() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Release all bitmaps
	var result error
	for _, bitmap := range this.bitmaps {
		if bitmap == nil {
			// NOOP
		} else if err := this.disposeBitmap(bitmap); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.bitmaps = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) NewBitmap(format gopi.SurfaceFormat, w, h uint32) (gopi.Bitmap, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if factory, exists := factories[format]; exists == false {
		return nil, gopi.ErrNotFound.WithPrefix(format)
	} else if bitmap, err := factory.New(format, w, h); err != nil {
		return nil, err
	} else {
		this.bitmaps = append(this.bitmaps, bitmap)
		return bitmap, nil
	}
}

func (this *Manager) DisposeBitmap(bitmap gopi.Bitmap) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	for i := range this.bitmaps {
		if this.bitmaps[i] != bitmap {
			continue
		}
		err := this.disposeBitmap(bitmap)
		this.bitmaps[i] = nil
		return err
	}
	return gopi.ErrNotFound.WithPrefix("DisposeBitmap")
}

func (this *Manager) disposeBitmap(bitmap gopi.Bitmap) error {
	if factory, exists := factories[bitmap.Format()]; exists == false {
		return gopi.ErrInternalAppError.WithPrefix("Dispose")
	} else {
		return factory.Dispose(bitmap)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func RegisterFactory(factory BitmapFactory, formats ...gopi.SurfaceFormat) {
	for _, format := range formats {
		if _, exists := factories[format]; exists {
			panic("RegisterFactory: Duplicate: " + fmt.Sprint(format))
		} else {
			factories[format] = factory
		}
	}
}

func AlignUp(v uint32, a uint32) uint32 {
	// Align a value on a byte-bounrary
	return ((v - 1) & ^(a - 1)) + a
}
