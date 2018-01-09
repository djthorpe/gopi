// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"os"
	"sync"

	"github.com/djthorpe/gopi"
)

type Element struct {
	Display gopi.Display
	Type    gopi.SurfaceType
	Origin  gopi.Point
	Size    gopi.Size
	Layer   uint16
	Opacity float32
}

type element struct {
	log    gopi.Logger
	lock   sync.Mutex
	handle dxElementHandle
}

/*
	display dxDisplayHandle
	update dxUpdateHandle
	layer uint16
	opacity uint32
	dst_frame *Rect
	src_bitmap
}

*/

func (config Element) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.graphics.rpi.surface.Open>{ Type=%v Origin=%v Size=%v Layer=%v Opacity=%v }", config.Type, config.Origin, config.Size, config.Layer, config.Opacity)

	this := new(element)
	this.log = log

	if update, err := dxUpdateStart(DX_UPDATE_PRIORITY_DEFAULT); err != DX_SUCCESS {
		return nil, os.NewSyscallError("dxUpdateStart", err)
	} else {
		defer dxUpdateSubmitSync(update)
	}

	// TODO

	return this, nil
}

func (this *element) Close() error {
	this.log.Debug("<sys.graphics.rpi.surface.Close>{ %v }", this.handle)
	if this.handle == dxElementHandle(DX_NO_ELEMENT) {
		return nil
	} else if update, err := dxUpdateStart(DX_UPDATE_PRIORITY_DEFAULT); err != DX_SUCCESS {
		return os.NewSyscallError("dxUpdateStart", err)
	} else {
		defer dxUpdateSubmitSync(update)
		if err := dxElementRemove(update, this.handle); err != DX_SUCCESS {
			this.handle = dxElementHandle(DX_NO_ELEMENT)
			return os.NewSyscallError("dxElementRemove", err)
		} else {
			this.handle = dxElementHandle(DX_NO_ELEMENT)
			return nil
		}
	}
}

/*
func (this *DXDisplay) AddElement(update dxUpdateHandle, layer uint16, opacity uint32, dst_rect *DXFrame, src_resource *DXResource) (*DXElement, error) {

        // destination frame - if nil, then cover whole frame
        if dst_rect == nil {
                size := this.GetSize()
                dst_rect = &DXFrame{DXZeroPoint, size}
        }

        // source frame size
        var src_size DXSize
        if src_resource != nil {
                src_size.Width = uint32(src_resource.GetSize().Width)
                src_size.Height = uint32(src_resource.GetSize().Height)
        } else {
                src_size.Width = dst_rect.Width
                src_size.Height = dst_rect.Height
        }

        // set alpha
        //alpha := dxAlpha{ DX_FLAGS_ALPHA_FROM_SOURCE, opacity, 0 }
        alpha := dxAlpha{DX_FLAGS_ALPHA_FIXED_ALL_PIXELS, opacity, 0}

        // set resource handle
        src_resource_handle := DX_RESOURCE_NONE
        if src_resource != nil {
                src_resource_handle = src_resource.GetHandle()
        }

        // create element structure
        element := new(DXElement)

        // add element
        element.handle = dxElementAdd(update, this.handle, layer, dst_rect, src_resource_handle, src_size, DX_PROTECTION_NONE, &alpha, nil, 0)
        if element.handle == DX_ELEMENT_NONE {
                return nil, this.log.Error("dxElementAdd failed")
        }

        // set other members of the element
        element.layer = layer
        element.frame = dst_rect

        // success
        return element, nil
}

*/
