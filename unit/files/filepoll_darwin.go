// +build darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files

import (

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type filepoll struct {
	base.Unit
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *filepoll) Init(config FilePoll) error {
	return gopi.ErrNotImplemented
}
