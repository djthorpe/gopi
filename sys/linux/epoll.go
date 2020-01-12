// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	EpollOp   int
	EpollEvt  syscall.EpollEvent
	EpollMode uint32
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EPOLL_OP_MOD EpollOp = syscall.EPOLL_CTL_MOD
	EPOLL_OP_ADD EpollOp = syscall.EPOLL_CTL_ADD
	EPOLL_OP_DEL EpollOp = syscall.EPOLL_CTL_DEL
)

const (
	EPOLL_MODE_EDGE   EpollMode = syscall.EPOLLPRI
	EPOLL_MODE_READ   EpollMode = syscall.EPOLLIN
	EPOLL_MODE_WRITE  EpollMode = syscall.EPOLLOUT
	EPOLL_MODE_HANGUP EpollMode = syscall.EPOLLHUP
	EPOLL_MODE_ERROR  EpollMode = syscall.EPOLLERR
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func EpollCreate() (int, error) {
	if fd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC); err != nil {
		return 0, os.NewSyscallError("EpollCreate1", err)
	} else {
		return fd, nil
	}
}

func EpollClose(handle int) error {
	if err := syscall.Close(handle); err != nil {
		return os.NewSyscallError("Close", err)
	} else {
		return nil
	}
}

func EpollAdd(handle, fd int, mode EpollMode) error {
	event := new(EpollEvt)
	event.Fd = int32(fd)
	event.Events = uint32(mode)
	if err := syscall.EpollCtl(handle, int(EPOLL_OP_ADD), fd, (*syscall.EpollEvent)(event)); err != nil {
		return os.NewSyscallError("EpollAdd epoll_ctl", err)
	} else {
		return nil
	}
}

func EpollDelete(handle, fd int) error {
	if err := syscall.EpollCtl(handle, int(EPOLL_OP_DEL), fd, nil); err != nil {
		return os.NewSyscallError("EpollDelete epoll_ctl", err)
	} else {
		return nil
	}
}

func EpollWait(handle int, timeout time.Duration, cap uint) ([]EpollEvt, error) {
	if cap == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("cap")
	}
	events := make([]syscall.EpollEvent, cap)
	if n, err := syscall.EpollWait(handle, events[:], int(timeout.Milliseconds())); err == syscall.EAGAIN || err == syscall.EINTR {
		return nil, nil
	} else if err != nil {
		return nil, os.NewSyscallError("epoll_wait", err)
	} else {
		evts := make([]EpollEvt, n)
		for i := range evts {
			evts[i] = EpollEvt(events[i])
		}
		return evts, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// EPOLL EVENT

func (evt EpollEvt) Flags() EpollMode {
	return EpollMode(evt.Events)
}

func (evt EpollEvt) String() string {
	return "<epoll.event fd=" + fmt.Sprint(evt.Fd) + " flags=" + fmt.Sprint(evt.Flags()) + ">"
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v EpollMode) String() string {
	if v == 0 {
		return "EPOLL_MODE_NONE"
	}
	str := ""
	if v&EPOLL_MODE_READ == EPOLL_MODE_READ {
		str += "EPOLL_MODE_READ" + "|"
	}
	if v&EPOLL_MODE_WRITE == EPOLL_MODE_WRITE {
		str += "EPOLL_MODE_WRITE" + "|"
	}
	if v&EPOLL_MODE_EDGE == EPOLL_MODE_EDGE {
		str += "EPOLL_MODE_EDGE" + "|"
	}
	if v&EPOLL_MODE_HANGUP == EPOLL_MODE_HANGUP {
		str += "EPOLL_MODE_HANGUP" + "|"
	}
	if v&EPOLL_MODE_ERROR == EPOLL_MODE_ERROR {
		str += "EPOLL_MODE_ERROR" + "|"
	}
	return strings.TrimSuffix(str, "|")
}
