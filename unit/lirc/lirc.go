/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package lirc

import (
	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type LIRC struct {
	// Input and Output Devices
	Dev string

	// Message bus
	Bus gopi.Bus

	// Filepoll interface
	Filepoll gopi.FilePoll
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (LIRC) Name() string { return "gopi.LIRC" }

func (config LIRC) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(lirc)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.LIRC

func (this *lirc) RcvMode() gopi.LIRCMode {
	return 0
}
func (this *lirc) SendMode() gopi.LIRCMode {
	return 0
}

func (this *lirc) SetRcvMode(mode gopi.LIRCMode) error {
	return gopi.ErrNotImplemented
}
func (this *lirc) SetSendMode(mode gopi.LIRCMode) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) GetRcvResolution() (uint32, error) {
	return 0, gopi.ErrNotImplemented
}

func (this *lirc) SetRcvTimeout(micros uint32) error {
	return gopi.ErrNotImplemented
}
func (this *lirc) SetRcvTimeoutReports(enable bool) error {
	return gopi.ErrNotImplemented
}
func (this *lirc) SetRcvCarrierHz(value uint32) error {
	return gopi.ErrNotImplemented
}
func (this *lirc) SetRcvCarrierRangeHz(min uint32, max uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetSendCarrierHz(value uint32) error {
	return gopi.ErrNotImplemented
}
func (this *lirc) SetSendDutyCycle(value uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) PulseSend(values []uint32) error {
	return gopi.ErrNotImplemented
}
