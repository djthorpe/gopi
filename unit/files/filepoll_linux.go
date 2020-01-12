// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files

import (
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

type filepoll struct {
	handle int

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// INIT & CLOSE

func (this *filepoll) Init(config FilePoll) error {
	this.Lock()
	defer this.Unlock()

	if handle, err := linux.EpollCreate(); err != nil {
		return err
	} else {
		this.handle = handle
	}

	// Return success
	return nil
}

func (this *filepoll) Close() error {
	this.Lock()
	defer this.Unlock()

	if this.handle != 0 {
		if err := linux.EpollClose(this.handle); err != nil {
			return err
		}
	}

	// Release resources
	this.handle = 0

	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *filepoll) String() string {
	if this.Unit.Closed {
		return fmt.Sprintf("<%v>", this.Log.Name())
	} else {
		return fmt.Sprintf("<%v handle=%v>", this.Log.Name(), this.handle)
	}
}

func (this *filepoll) Watch(fd uintptr) error {
	return gopi.ErrNotImplemented
}
