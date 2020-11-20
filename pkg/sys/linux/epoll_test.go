// +build linux

package linux_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

func Test_Epoll_000(t *testing.T) {
	if handle, err := linux.EpollCreate(); err != nil {
		t.Error(err)
	} else if err := linux.EpollClose(handle); err != nil {
		t.Error(err)
	}
}
