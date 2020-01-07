// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package lirc

import (
	"os"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type lirc struct {
	devin               *os.File          // device in
	devout              *os.File          // device out
	features            linux.LIRCFeature // features
	rcv_mode, send_mode gopi.LIRCMode     // modes

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// LIRC_DEV_IN & LIRC_DEV_OUT are the default device path
	LIRC_DEV_OUT = "/dev/lirc0"
	LIRC_DEV_IN  = "/dev/lirc1"
	// LIRC_CARRIER_FREQUENCY is the default carrier frequency
	LIRC_CARRIER_FREQUENCY = 38000
	// LIRC_DUTY_CYCLE is the default duty cycle
	LIRC_DUTY_CYCLE = 50
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *lirc) Init(config LIRC) error {
	return gopi.ErrNotImolemented
}
