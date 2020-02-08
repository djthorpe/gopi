/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces

import (
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SurfaceManager struct {
	Display gopi.Display
}

type surfacemanager struct {
	display gopi.Display

	Implementation
	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (SurfaceManager) Name() string { return "gopi/surfaces" }

func (config SurfaceManager) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(surfacemanager)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	} else {
		this.display = config.Display
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.SurfaceManager

func (this *surfacemanager) Display() gopi.Display {
	return this.display
}

func (this *surfacemanager) CreateSurface(gopi.SurfaceFlags, float32, uint16, gopi.Point, gopi.Size) (gopi.Surface, error) {
	return nil, gopi.ErrNotImplemented
}

func (this *surfacemanager) CreateSnapshot(flags gopi.SurfaceFlags) (gopi.Bitmap, error) {
	return nil, gopi.ErrNotImplemented
}
