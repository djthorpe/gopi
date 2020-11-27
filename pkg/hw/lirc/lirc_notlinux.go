// +build !linux

package lirc

import (
	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type lirc struct {
	gopi.Unit
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *lirc) RecvMode() gopi.LIRCMode {
	return gopi.LIRC_MODE_NONE
}

func (this *lirc) SendMode() gopi.LIRCMode {
	return gopi.LIRC_MODE_NONE
}

func (this *lirc) SetRecvMode(gopi.LIRCMode) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetSendMode(gopi.LIRCMode) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) RecvDutyCycle() uint32 {
	return 0
}

func (this *lirc) RecvResolutionMicros() uint32 {
	return 0
}

func (this *lirc) SetRecvTimeoutMs(uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetRecvTimeoutReports(bool) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetRecvCarrierHz(uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SetRecvCarrierRangeHz(min, max uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) SendDutyCycle() uint32 {
	return 0
}
func (this *lirc) SetSendCarrierHz(uint32) error {
	return gopi.ErrNotImplemented

}
func (this *lirc) SetSendDutyCycle(uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lirc) PulseSend([]uint32) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *lirc) String() string {
	str := "<lirc"
	return str + ">"
}
