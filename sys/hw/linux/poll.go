/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"context"
	"errors"
	"os"
	"sync"
	"syscall"
)

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type PollMode int

type PollDriver struct {
	// Poll file handle
	handle int

	// Logging
	log gopi.Logger

	// Exclusive lock
	lock sync.Mutex

	// Events we're waiting for
	events map[int]uint32

	// An event buffer for EpollCtl
	ctlEvent syscall.EpollEvent
}

type PollCallback func(fd int, flags PollMode)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	POLL_EPOLLET    = 1 << 31
	POLL_READFLAGS  = syscall.EPOLLIN | syscall.EPOLLRDHUP
	POLL_EDGEFLAGS  = syscall.EPOLLIN | POLL_EPOLLET | syscall.EPOLLPRI
	POLL_WRITEFLAGS = syscall.EPOLLOUT
)

const (
	POLL_MODE_READ   PollMode = 1 << iota
	POLL_MODE_WRITE  PollMode = 1 << iota
	POLL_MODE_EDGE   PollMode = 1 << iota
	POLL_MODE_HANGUP PollMode = 1 << iota
	POLL_MODE_ERROR  PollMode = 1 << iota
)

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	ErrPollInvalidMode = errors.New("Invalid mode argument")
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func NewPollDriver(log gopi.Logger) (*PollDriver, error) {
	var err error

	log.Debug("<linux.Poll>Open")

	this := new(PollDriver)
	if this.handle, err = syscall.EpollCreate1(syscall.EPOLL_CLOEXEC); err != nil {
		return nil, err
	}
	this.events = make(map[int]uint32)
	this.log = log

	// success
	return this, nil
}

func (this *PollDriver) Close() error {
	this.log.Debug("<linux.Poll>Close")
	return syscall.Close(this.handle)
}

////////////////////////////////////////////////////////////////////////////////
// ADD FILE HANDLE TO WATCH

func (this *PollDriver) Add(fd int, mode PollMode) error {
	var already bool

	// Make this method exclusive
	this.lock.Lock()
	defer this.lock.Unlock()

	// Set up the event structure
	this.ctlEvent.Fd = int32(fd)
	this.ctlEvent.Events, already = this.events[fd]

	switch mode {
	case POLL_MODE_READ:
		this.ctlEvent.Events |= POLL_READFLAGS
	case POLL_MODE_WRITE:
		this.ctlEvent.Events |= POLL_WRITEFLAGS
	default:
		return ErrPollInvalidMode
	}

	// Modify or add poll
	var op int
	if already {
		op = syscall.EPOLL_CTL_MOD
	} else {
		op = syscall.EPOLL_CTL_ADD
	}

	// System call
	if err := syscall.EpollCtl(this.handle, op, fd, &this.ctlEvent); err != nil {
		return os.NewSyscallError("epoll_ctl", err)
	}

	// Record the events we're interested in
	this.events[fd] = this.ctlEvent.Events

	// return success
	return nil
}

func (this *PollDriver) Remove(fd int, mode PollMode) error {

	// Make this method exclusive
	this.lock.Lock()
	defer this.lock.Unlock()

	switch mode {
	case POLL_MODE_READ:
		this.stopWaiting(fd, POLL_READFLAGS)
	case POLL_MODE_WRITE:
		this.stopWaiting(fd, POLL_WRITEFLAGS)
	default:
		return ErrPollInvalidMode
	}

	return nil
}

// Watch for polling events for a specified period
func (this *PollDriver) Watch(ctx context.Context, callback PollCallback) error {
	c := make(chan error, 1)
	go func() { c <- this.watchInner(callback) }()
	select {
	case <-ctx.Done():
		// Cancel
		this.log.Info("TODO: CANCEL WATCH")
		<-c // Wait for transaction to return
		return ctx.Err()
	case err := <-c:
		return err
	}
}

//////////////

func (this *PollDriver) watchInner(callback PollCallback) error {
	// Maximum of 64 events
	events := make([]syscall.EpollEvent, 64)

	// Wait for 200ms each time
	delta := 200

	for {
		this.lock.Lock()
		n, err := syscall.EpollWait(this.handle, events[:], int(delta))
		this.lock.Unlock()

		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EINTR {
				continue
			}
			if err == syscall.EBADF || err == syscall.EINVAL {
				// file descriptor closed so exit successfully
				return nil
			}

			// else return an
			return err
		}
		if n == 0 {
			// No more events to wait for
			continue
		}

		// Process incoming events
		for _, event := range events[:n] {
			var mode PollMode
			// determine modes
			if (event.Events&syscall.EPOLLHUP) != 0 || (event.Events&syscall.EPOLLRDHUP != 0) {
				mode |= POLL_MODE_HANGUP
			}
			if event.Events&syscall.EPOLLERR != 0 {
				mode |= POLL_MODE_ERROR
			}
			if event.Events&syscall.EPOLLIN != 0 {
				mode |= POLL_MODE_READ
			}
			if event.Events&syscall.EPOLLOUT != 0 {
				mode |= POLL_MODE_WRITE
			}
			if event.Events&POLL_EPOLLET != 0 {
				mode |= POLL_MODE_EDGE
			}
			// callback
			callback(int(event.Fd), mode)
		}
	}
	// we never reach here
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (m PollMode) String() string {
	s := ""
	if m == 0 {
		return s
	}
	if m&POLL_MODE_READ != 0 {
		s = s + "POLL_MODE_READ|"
	}
	if m&POLL_MODE_WRITE != 0 {
		s = s + "POLL_MODE_WRITE|"
	}
	if m&POLL_MODE_HANGUP != 0 {
		s = s + "POLL_MODE_HANGUP|"
	}
	if m&POLL_MODE_ERROR != 0 {
		s = s + "POLL_MODE_ERROR|"
	}
	if m&POLL_MODE_EDGE != 0 {
		s = s + "POLL_MODE_EDGE|"
	}
	return s[0 : len(s)-1]
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *PollDriver) stopWaiting(fd int, bits uint32) {
	events, already := this.events[fd]
	if already == false {
		// The fd returned by the kernel may have been cancelled already; return silently.
		return
	}

	// Disable the given bits. If we're still waiting for other events, modify the fd
	// event in the kernel.  Otherwise, delete it.
	events &= ^bits
	if events != 0 {
		this.ctlEvent.Fd = int32(fd)
		this.ctlEvent.Events = events
		if err := syscall.EpollCtl(this.handle, syscall.EPOLL_CTL_MOD, fd, &this.ctlEvent); err != nil {
			this.log.Warn("<linux.Poll> epoll_ctl error: EPOLL_CTL_MOD: %v", err)
		}
		this.events[fd] = events
	} else {
		if err := syscall.EpollCtl(this.handle, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
			this.log.Warn("<linux.Poll> epoll_ctl error: EPOLL_CTL_DEL: %v", err)
		}
		delete(this.events, fd)
	}
}
