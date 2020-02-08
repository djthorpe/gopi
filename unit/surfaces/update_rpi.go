// +build rpi

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
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Update struct {
	update rpi.DXUpdate
	sync.RWMutex
}

func (this *Update) Start(priority int32) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.update != rpi.DX_NO_HANDLE {
		return gopi.ErrOutOfOrder.WithPrefix("update")
	}
	if update, err := rpi.DXUpdateStart(priority); err != nil {
		return err
	} else {
		this.update = update
	}
	// Success
	return nil
}

func (this *Update) Submit() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.update == rpi.DX_NO_HANDLE {
		return gopi.ErrOutOfOrder.WithPrefix("update")
	}
	err := rpi.DXUpdateSubmitSync(this.update)
	this.update = 0
	return err
}

func (this *Update) Update() rpi.DXUpdate {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.update
}
