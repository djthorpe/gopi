// +build darwin

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
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type lirc struct {
	devin  *os.File // device in
	devout *os.File // device out

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *lirc) Init(config LIRC) error {
	return gopi.ErrNotImplemented
}
