// +build rpi
// +build egl

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"unsafe"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	egl "github.com/djthorpe/gopi/v2/sys/egl"
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	display "github.com/djthorpe/gopi/v2/unit/display"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type manager struct {
	display      gopi.Display
	handle       egl.EGLDisplay
	major, minor int
	bitmaps      map[rpi.DXResource]*bitmap
	surfaces     []*surface
	update       rpi.DXUpdate

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// INIT AND CLOSE

func (this *manager) Init(config SurfaceManager) error {
	// Check display
	if config.Display == nil {
		return gopi.ErrBadParameter.WithPrefix("display")
	} else {
		this.display = config.Display
	}

	// Initialize EGL
	if handle := egl.EGLGetDisplay(this.display.DisplayId()); handle == 0 {
		return gopi.ErrBadParameter.WithPrefix("display")
	} else if major, minor, err := egl.EGLInitialize(handle); err != nil {
		return err
	} else {
		this.handle = handle
		this.major = major
		this.minor = minor
		this.bitmaps = make(map[rpi.DXResource]*bitmap)
		this.surfaces = make([]*surface, 0, 4)
	}

	// Success
	return nil
}

func (this *manager) Close() error {
	if this.handle == 0 {
		return nil
	}

	// Free surfaces
	if err := this.Do(func(gopi.SurfaceManager) error {
		for _, surface := range this.surfaces {
			if err := surface.Destroy(this.update); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// Free bitmaps
	for _, bitmap := range this.bitmaps {
		if err := bitmap.Destroy(); err != nil {
			return err
		} else {
			delete(this.bitmaps, bitmap.handle)
		}
	}

	// Terminate EGL
	if err := egl.EGLTerminate(this.handle); err != nil {
		return err
	}

	// Free resources
	this.display = nil
	this.handle = 0
	this.bitmaps = nil
	this.surfaces = nil

	// Success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *manager) String() string {
	if this.Unit.Closed {
		return "<gopi.SurfaceManager>"
	} else {
		return "<gopi.SurfaceManager display=" + fmt.Sprint(this.display) + " " +
			"name=" + strconv.Quote(this.Name()) + " " +
			"api=" + fmt.Sprint(this.Types()) + ">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.SurfaceManager

func (this *manager) Display() gopi.Display {
	return this.display
}

func (this *manager) Name() string {
	if this.handle == 0 {
		return ""
	}
	return fmt.Sprintf("%v EGL %v", egl.EGLQueryString(this.handle, egl.EGL_QUERY_VENDOR), egl.EGLQueryString(this.handle, egl.EGL_QUERY_VERSION))
}

func (this *manager) Extensions() []string {
	if this.handle == 0 {
		return nil
	}
	return strings.Split(egl.EGLQueryString(this.handle, egl.EGL_QUERY_EXTENSIONS), " ")
}

func (this *manager) Types() []gopi.SurfaceFlags {
	if this.handle == 0 {
		return nil
	}
	types := strings.Split(egl.EGLQueryString(this.handle, egl.EGL_QUERY_CLIENT_APIS), " ")
	surface_types := make([]gopi.SurfaceFlags, 0, len(types))
	for _, t := range types {
		if t_, ok := egl.EGLSurfaceTypeMap[t]; ok {
			surface_types = append(surface_types, t_)
		}
	}
	// always include bitmaps
	return append(surface_types, gopi.SURFACE_FLAG_BITMAP)
}

////////////////////////////////////////////////////////////////////////////////
// SURFACE METHODS

func (this *manager) CreateSurfaceWithBitmap(bm gopi.Bitmap, flags gopi.SurfaceFlags, opacity float32, layer uint16, origin gopi.Point, size gopi.Size) (gopi.Surface, error) {
	if layer == 0 {
		layer = gopi.SURFACE_LAYER_DEFAULT
	}
	flags = gopi.SURFACE_FLAG_BITMAP | bm.Type() | flags.Mod()
	if bitmap_, ok := bm.(*bitmap); ok == false || bitmap_ == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if opacity < 0.0 || opacity > 1.0 {
		return nil, gopi.ErrBadParameter.WithPrefix("opacity")
	} else if layer < gopi.SURFACE_LAYER_DEFAULT || layer > gopi.SURFACE_LAYER_MAX {
		return nil, gopi.ErrBadParameter.WithPrefix("layer")
	} else if size = size_from_bitmap(bm, size); size == gopi.ZeroSize {
		return nil, gopi.ErrBadParameter.WithPrefix("size")
	} else if native, err := NewNativeSurface(this.update, bitmap_, rpi_dx_display(this.display), flags, opacity, layer, origin, size); err != nil {
		return nil, err
	} else {
		surface := NewSurface(flags, opacity, layer, native)
		this.surfaces = append(this.surfaces, surface)
		return surface, nil
	}
}

func (this *manager) CreateSurface(flags gopi.SurfaceFlags, opacity float32, layer uint16, origin gopi.Point, size gopi.Size) (gopi.Surface, error) {
	// api
	api := flags.Type()

	// Set layer to default
	if layer == 0 {
		layer = gopi.SURFACE_LAYER_DEFAULT
	}

	// if Bitmap, then create a bitmap
	if api == gopi.SURFACE_FLAG_BITMAP {
		if bitmap, err := this.CreateBitmap(flags, size); err != nil {
			return nil, err
		} else if surface, err := this.CreateSurfaceWithBitmap(bitmap, flags, opacity, layer, origin, size); err != nil {
			return nil, err
		} else {
			return surface, nil
		}
	}

	// Return gopi.ErrNotImplemented
	return nil,gopi.ErrNotImplemented

	// Choose r,g,b,a bits per pixel
	var r, g, b, a uint
	switch flags.Config() {
	case gopi.SURFACE_FLAG_RGB565:
		r = 5
		g = 6
		b = 5
		a = 0
	case gopi.SURFACE_FLAG_RGBA32:
		r = 8
		g = 8
		b = 8
		a = 8
	case gopi.SURFACE_FLAG_RGB888:
		r = 8
		g = 8
		b = 8
		a = 0
	default:
		return nil, gopi.ErrNotImplemented.WithPrefix("flags")
	}

	// Create EGL context
	if api_, exists := egl.EGLAPIMap[api]; exists == false {
		return nil, gopi.ErrBadParameter.WithPrefix("api")
	} else if renderable_, exists := egl.EGLRenderableMap[api]; exists == false {
		return nil, gopi.ErrBadParameter.WithPrefix("api")
	} else if opacity < 0.0 || opacity > 1.0 {
		return nil, gopi.ErrBadParameter.WithPrefix("opacity")
	} else if layer < gopi.SURFACE_LAYER_DEFAULT || layer > gopi.SURFACE_LAYER_MAX {
		return nil, gopi.ErrBadParameter.WithPrefix("layer")
	} else if err := egl.EGLBindAPI(api_); err != nil {
		return nil, err
	} else if config, err := egl.EGLChooseConfig(this.handle, r, g, b, a, egl.EGL_SURFACETYPE_FLAG_WINDOW, renderable_); err != nil {
		return nil, err
	} else if native, err := NewNativeSurface(this.update, nil, rpi_dx_display(this.display), flags, opacity, layer, origin, size); err != nil {
		return nil, err
	} else if handle, err := egl.EGLCreateSurface(this.handle, config, egl_nativewindow(native)); err != nil {
		native.Destroy(this.update)
		return nil, err
	} else if context, err := egl.EGLCreateContext(this.handle, config, nil); err != nil {
		// TODO: Destroy EGL Window
		native.Destroy(this.update)
		return nil, err
	} else if err := egl.EGLMakeCurrent(this.handle, handle, handle, context); err != nil {
		// TODO: Destroy context, surface, window, ...
		native.Destroy(this.update)
		return nil, err
	} else {
		surface := NewSurface(flags, opacity, layer, native)
		this.surfaces = append(this.surfaces, surface)
		return surface, nil
	}
}

func (this *manager) DestroySurface(gopi.Surface) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// UPDATE METHODS

func (this *manager) Do(callback gopi.SurfaceCallback) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.handle == 0 {
		return gopi.ErrBadParameter.WithPrefix("handle")
	}
	if callback == nil {
		return gopi.ErrBadParameter.WithPrefix("callback")
	}
	if this.update != 0 {
		return gopi.ErrInternalAppError.WithPrefix("update")
	}
	// TODO rpi.DX_UPDATE_PRIORITY_DEFAULT
	if update, err := rpi.DXUpdateStart(0); err != nil {
		return err
	} else {
		this.update = update
		defer func() {
			if rpi.DXUpdateSubmitSync(update); err != nil {
				this.Log.Error(err)
			}
			this.update = 0
		}()
		if err := callback(this); err != nil {
			return err
		}
	}
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BITMAP METHODS

func (this *manager) CreateBitmap(flags gopi.SurfaceFlags, size gopi.Size) (gopi.Bitmap, error) {
	flags = gopi.SURFACE_FLAG_BITMAP | flags.Config() | flags.Mod()
	if bitmap, err := NewBitmap(flags, size); err != nil {
		return nil, err
	} else if _, exists := this.bitmaps[bitmap.handle]; exists {
		return nil, gopi.ErrDuplicateItem.WithPrefix("bitmap")
	} else {
		this.bitmaps[bitmap.handle] = bitmap
		return bitmap, nil
	}
}

func (this *manager) CreateSnapshot(flags gopi.SurfaceFlags) (gopi.Bitmap, error) {
	flags = gopi.SURFACE_FLAG_BITMAP | flags.Config() | flags.Mod()
	w, h := this.display.Size()
	if bitmap, err := NewBitmap(flags, gopi.Size{float32(w), float32(h)}); err != nil {
		return nil, err
	} else if _, exists := this.bitmaps[bitmap.handle]; exists {
		return nil, gopi.ErrDuplicateItem.WithPrefix("bitmap")
	} else if err := rpi.DXDisplaySnapshot(rpi_dx_display(this.display), bitmap.handle, rpi.DX_TRANSFORM_NONE); err != nil {
		bitmap.Destroy()
		return nil, err
	} else {
		this.bitmaps[bitmap.handle] = bitmap
		return bitmap, nil
	}
}

func (this *manager) DestroyBitmap(b gopi.Bitmap) error {
	if b == nil {
		return gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if bitmap_, ok := b.(*bitmap); ok == false {
		return gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if _, exists := this.bitmaps[bitmap_.handle]; exists == false {
		return gopi.ErrNotFound.WithPrefix("bitmap")
	} else if err := bitmap_.Destroy(); err != nil {
		return err
	} else {
		delete(this.bitmaps, bitmap_.handle)
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func rpi_dx_display(d gopi.Display) rpi.DXDisplayHandle {
	return d.(display.NativeDisplay).Handle()
}

func opacity_from_float(opacity float32) uint8 {
	if opacity < 0.0 {
		opacity = 0.0
	} else if opacity > 1.0 {
		opacity = 1.0
	}
	// Opacity is between 0 (fully transparent) and 255 (fully opaque)
	return uint8(opacity * float32(0xFF))
}

func size_from_bitmap(bitmap gopi.Bitmap, size gopi.Size) gopi.Size {
	if size == gopi.ZeroSize {
		return bitmap.Size()
	} else {
		return size
	}
}

func egl_nativewindow(window *nativesurface) egl.EGLNativeWindow {
	return egl.EGLNativeWindow(unsafe.Pointer(window))
}
