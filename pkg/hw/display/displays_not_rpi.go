// +build !rpi

package display

import (
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

func (this *display) new(id uint16) error {
	return gopi.ErrNotImplemented
}

func (this *display) close() error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Displays) Enumerate() []gopi.Display {
	return nil
}
