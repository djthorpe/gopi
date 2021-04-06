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
	New(ColorModel, uint32, uint32) (gopi.Bitmap, error)
	Dispose(gopi.Bitmap) error
}

type Bitmaps struct {
	gopi.Unit
	gopi.Platform
	sync.Mutex

	bitmaps []gopi.Bitmap
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	factories   = make(map[gopi.SurfaceFormat]BitmapFactory)
	colormodels = make(map[gopi.SurfaceFormat]ColorModel)
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Bitmaps) New(gopi.Config) error {
	this.Require(this.Platform)
	return nil
}

func (this *Bitmaps) Dispose() error {
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
// STRINGIFY

func (this *Bitmaps) String() string {
	str := "<bitmaps models=["
	for _, model := range colormodels {
		str += fmt.Sprint(" ", model)
	}
	return str + " ]>"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Bitmaps) NewBitmap(format gopi.SurfaceFormat, w, h uint32) (gopi.Bitmap, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Get color model
	model, exists := colormodels[format]
	if exists == false {
		return nil, gopi.ErrNotFound.WithPrefix(format)
	}

	// Create bitmap with factory
	if factory, exists := factories[format]; exists == false {
		return nil, gopi.ErrNotFound.WithPrefix(format)
	} else if bitmap, err := factory.New(model, w, h); err != nil {
		return nil, err
	} else {
		this.bitmaps = append(this.bitmaps, bitmap)
		return bitmap, nil
	}
}

func (this *Bitmaps) DisposeBitmap(bitmap gopi.Bitmap) error {
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

func (this *Bitmaps) disposeBitmap(bitmap gopi.Bitmap) error {
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

func RegisterColorModel(format gopi.SurfaceFormat, model ColorModel) {
	if _, exists := colormodels[format]; exists {
		panic("RegisterColorModel: Duplicate: " + fmt.Sprint(format))
	} else {
		colormodels[format] = model
	}
}

func GetColorModel(fmt gopi.SurfaceFormat) ColorModel {
	if model, exists := colormodels[fmt]; exists {
		return model
	} else {
		return nil
	}
}

func AlignUp(v uint32, a uint32) uint32 {
	// Align a value on a byte-bounrary
	return ((v - 1) & ^(a - 1)) + a
}
