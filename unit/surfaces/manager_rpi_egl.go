// +build rpi,egl

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces

import (
	"strconv"
	"strings"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	egl "github.com/djthorpe/gopi/v2/sys/egl"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	display "github.com/djthorpe/gopi/v2/unit/display"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
	element "github.com/djthorpe/gopi/v2/unit/surfaces/element"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Implementation struct {
	egldisplay egl.EGLDisplay
	dxdisplay  rpi.DXDisplayHandle
	bitmap     map[bitmap.Bitmap]bool
	element    map[element.Element]bool

	Update
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *surfacemanager) Init(config SurfaceManager) error {
	if egldisplay := egl.EGLGetDisplay(config.Display.DisplayId()); egldisplay == 0 {
		return gopi.ErrInternalAppError
	} else {
		this.egldisplay = egldisplay
	}

	if dxdisplay, ok := config.Display.(display.NativeDisplay); ok == false {
		return gopi.ErrBadParameter.WithPrefix("display")
	} else {
		this.dxdisplay = dxdisplay.Handle()
	}

	// Initialize EGL
	if _, _, err := egl.EGLInitialize(this.egldisplay); err != nil {
		return err
	}

	// Set bitmaps and elements
	this.bitmap = make(map[bitmap.Bitmap]bool)
	this.element = make(map[element.Element]bool)

	// Return success
	return nil
}

func (this *surfacemanager) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Remove all elements and bitmaps
	if err := this.removeAllElements(); err != nil {
		return err
	}
	if err := this.removeAllBitmaps(); err != nil {
		return err
	}

	// Close EGL
	if err := egl.EGLTerminate(this.egldisplay); err != nil {
		return err
	}

	// Release resources
	this.element = nil
	this.bitmap = nil

	// Success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *surfacemanager) EGLVendor() string {
	return egl.EGLQueryString(this.egldisplay, egl.EGL_QUERY_VENDOR)
}

func (this *surfacemanager) EGLVersion() string {
	return egl.EGLQueryString(this.egldisplay, egl.EGL_QUERY_VERSION)
}

func (this *surfacemanager) EGLExtensions() string {
	return egl.EGLQueryString(this.egldisplay, egl.EGL_QUERY_EXTENSIONS)
}

func (this *surfacemanager) EGLClientAPI() string {
	return egl.EGLQueryString(this.egldisplay, egl.EGL_QUERY_CLIENT_APIS)
}

func (this *surfacemanager) Name() string {
	return this.EGLVendor()
}

func (this *surfacemanager) Types() []gopi.SurfaceFlags {
	types := strings.Split(this.EGLClientAPI(), " ")
	apis := make([]gopi.SurfaceFlags, 1, len(types)+1)
	apis[0] = gopi.SURFACE_FLAG_BITMAP
	for _, api := range types {
		if surfaceType, exists := egl.EGLSurfaceTypeMap[api]; exists {
			apis = append(apis, surfaceType)
		}
	}
	return apis
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *surfacemanager) String() string {
	str := "<" + this.Log.Name()
	str += " name=" + strconv.Quote(this.Name())
	str += " version=" + strconv.Quote(this.EGLVersion())
	str += " extensions=" + strconv.Quote(this.EGLExtensions())
	if types := this.Types(); len(types) > 0 {
		str += " surface_types="
		for _, t := range types {
			str += t.TypeString() + "|"
		}
		str = strings.TrimSuffix(str, "|")
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// CREATE BITMAPS

func (this *surfacemanager) CreateBitmap(flags gopi.SurfaceFlags, size gopi.Size) (gopi.Bitmap, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if bm, err := gopi.New(bitmap.Config{size, flags}, this.Log.Clone(bitmap.Config{}.Name())); err != nil {
		return nil, err
	} else {
		bm_ := bm.(bitmap.Bitmap)
		bm_.Retain()
		this.bitmap[bm_] = true
		return bm_, nil
	}
}

func (this *surfacemanager) DestroyBitmap(bm gopi.Bitmap) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if bm == nil {
		return gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if bm_, ok := bm.(bitmap.Bitmap); ok == false {
		return gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if _, exists := this.bitmap[bm_]; exists == false {
		return gopi.ErrNotFound.WithPrefix("bitmap")
	} else {
		delete(this.bitmap, bm_)
		if bm_.Release() {
			return bm_.Close()
		} else {
			return nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// DO UPDATES

func (this *surfacemanager) Do(cb gopi.SurfaceCallback) error {
	if err := this.Update.Start(0); err != nil {
		return err
	}
	defer this.Update.Submit()
	return cb()
}

////////////////////////////////////////////////////////////////////////////////
// CREATE SURFACES

func (this *surfacemanager) CreateSurfaceWithBitmap(bm gopi.Bitmap, flags gopi.SurfaceFlags, opacity float32, layer uint16, origin gopi.Point, size gopi.Size) (gopi.Surface, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check layer parameter
	if layer < gopi.SURFACE_LAYER_DEFAULT || layer > gopi.SURFACE_LAYER_MAX {
		return nil, gopi.ErrBadParameter.WithPrefix("layer")
	}

	if bm == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if bm_, ok := bm.(bitmap.Bitmap); ok == false {
		return nil, gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if update := this.Update.Update(); update == rpi.DX_NO_HANDLE {
		return nil, gopi.ErrOutOfOrder.WithPrefix("CreateSurfaceWithBitmap")
	} else if em, err := gopi.New(element.Config{
		Origin:  origin,
		Size:    size,
		Bitmap:  bm_,
		Layer:   layer,
		Opacity: opacity,
		Flags:   flags,
		Update:  update,
		Display: this.dxdisplay,
	}, this.Log.Clone(element.Config{}.Name())); err != nil {
		return nil, err
	} else {
		em_ := em.(element.Element)
		this.element[em_] = true
		return em_, nil
	}
}

func (this *surfacemanager) CreateBackground(flags gopi.SurfaceFlags, opacity float32) (gopi.Surface, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	w, h := this.display.Size()
	if w == 0 || h == 0 {
		return nil, gopi.ErrOutOfOrder.WithPrefix("Display")
	} else if update := this.Update.Update(); update == rpi.DX_NO_HANDLE {
		return nil, gopi.ErrOutOfOrder.WithPrefix("CreateBackground")
	} else if em, err := gopi.New(element.Config{
		Origin:  gopi.ZeroPoint,
		Size:    gopi.Size{float32(w), float32(h)},
		Layer:   gopi.SURFACE_LAYER_BACKGROUND,
		Opacity: opacity,
		Flags:   flags,
		Update:  update,
		Display: this.dxdisplay,
	}, this.Log.Clone(element.Config{}.Name())); err != nil {
		return nil, err
	} else {
		em_ := em.(element.Element)
		this.element[em_] = true
		return em_, nil
	}
}

func (this *surfacemanager) DestroySurface(surface gopi.Surface) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if surface == nil {
		return gopi.ErrBadParameter.WithPrefix("surface")
	} else if em, ok := surface.(element.Element); ok == false {
		return gopi.ErrBadParameter.WithPrefix("surface")
	} else if _, exists := this.element[em]; exists == false {
		return gopi.ErrNotFound.WithPrefix("surface")
	} else if update := this.Update.Update(); update == rpi.DX_NO_HANDLE {
		return gopi.ErrOutOfOrder.WithPrefix("DestroySurface")
	} else {
		delete(this.element, em)
		if err := em.RemoveElement(update); err != nil {
			return err
		} else if err := em.Close(); err != nil {
			return err
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *surfacemanager) removeAllElements() error {
	// Ignore if no elements
	if len(this.element) == 0 {
		return nil
	}
	// Start an update
	update, err := rpi.DXUpdateStart(0)
	if err != nil {
		return err
	}
	// Remove elements
	errs := gopi.NewCompoundError()
	for em := range this.element {
		errs.Add(em.RemoveElement(update))
	}
	// Update GPU
	errs.Add(rpi.DXUpdateSubmitSync(update))
	// Delete elements from map
	for em := range this.element {
		delete(this.element, em)
	}
	// Return any errors
	return errs.ErrorOrSelf()
}

func (this *surfacemanager) removeAllBitmaps() error {
	// Ignore if no bitmaps
	if len(this.bitmap) == 0 {
		return nil
	}
	// Release bitmaps
	errs := gopi.NewCompoundError()
	for bm := range this.bitmap {
		errs.Add(bm.Close())
	}
	// Delete bitmaps from map
	for bm := range this.bitmap {
		delete(this.bitmap, bm)
	}
	// Return any errors
	return errs.ErrorOrSelf()
}
