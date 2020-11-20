// +build linux

package linux

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v3"
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

func EpollCreate() (uintptr, error) {
	if fd, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC); err != nil {
		return 0, os.NewSyscallError("EpollCreate1", err)
	} else {
		return uintptr(fd), nil
	}
}

func EpollClose(handle uintptr) error {
	if err := syscall.Close(int(handle)); err != nil {
		return os.NewSyscallError("Close", err)
	} else {
		return nil
	}
}

func EpollAdd(handle, fd uintptr, mode EpollMode) error {
	event := new(EpollEvt)
	event.Fd = int32(fd)
	event.Events = uint32(mode)
	if err := syscall.EpollCtl(int(handle), int(EPOLL_OP_ADD), int(fd), (*syscall.EpollEvent)(event)); err != nil {
		return os.NewSyscallError("EpollAdd epoll_ctl", err)
	} else {
		return nil
	}
}

func EpollDelete(handle, fd uintptr) error {
	if err := syscall.EpollCtl(int(handle), int(EPOLL_OP_DEL), int(fd), nil); err != nil {
		return os.NewSyscallError("EpollDelete epoll_ctl", err)
	} else {
		return nil
	}
}

func EpollWait(handle uintptr, timeout time.Duration, cap uint) ([]EpollEvt, error) {
	if cap == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("cap")
	}
	events := make([]syscall.EpollEvent, cap)
	ms := int(timeout.Milliseconds())
	// Block if time is zero
	if ms == 0 {
		ms = -1
	}
	if n, err := syscall.EpollWait(int(handle), events[:], ms); err == syscall.EAGAIN || err == syscall.EINTR {
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
