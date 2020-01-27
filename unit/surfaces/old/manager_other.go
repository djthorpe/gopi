// +build !rpi
// +build !egl

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package surfaces

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type manager struct {
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// INIT AND CLOSE

func (this *manager) Init(config SurfaceManager) error {
	// Not implemented
	return gopi.ErrNotImplemented
}
