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
	}

	// Success
	return nil
}

func (this *manager) Close() error {
	if this.handle == 0 {
		return nil
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
// IMPLEMENTATION gopi.SurfaceManager surface methods

func (this *manager) CreateSurfaceWithBitmap(bitmap gopi.Bitmap, flags gopi.SurfaceFlags, opacity float32, layer uint16, origin gopi.Point, size gopi.Size) (gopi.Surface, error) {
	flags = gopi.SURFACE_FLAG_BITMAP | bitmap.Type() | flags.Mod()
	if opacity < 0.0 || opacity > 1.0 {
		return nil, gopi.ErrBadParameter.WithPrefix("opacity")
	} else if layer < gopi.SURFACE_LAYER_DEFAULT || layer > gopi.SURFACE_LAYER_MAX {
		return nil, gopi.ErrBadParameter.WithPrefix("layer")
	} else if bitmap == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("bitmap")
	} else if size = size_from_bitmap(bitmap, size); size == gopi.ZeroSize {
		return nil, gopi.ErrBadParameter.WithPrefix("size")
	} else if native_surface, err := NewSurface(this.update, bitmap, flags, opacity, layer, origin, size); err != nil {
		return nil, err
	} else {
		// Return the surface
		s := &surface{
			log:     this.log,
			flags:   flags,
			opacity: opacity,
			layer:   layer,
			native:  native_surface,
			bitmap:  bitmap,
		}
		this.surfaces = append(this.surfaces, s)
		return s, nil
	}
}

func (this *manager) DestroySurface(gopi.Surface) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.SurfaceManager update methods

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
// IMPLEMENTATION gopi.SurfaceManager bitmap methods

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
