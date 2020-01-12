// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package linux_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/sys/linux"
)

func Test_Epoll_000(t *testing.T) {
	if handle, err := linux.EpollCreate(); err != nil {
		t.Error(err)
	} else if err := linux.EpollClose(handle); err != nil {
		t.Error(err)
	}
}
