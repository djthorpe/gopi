// +build linux

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
	"fmt"
	"os"
	"sync"
	"syscall"
	"time"

	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type FilePollMode int

type FilePoll struct {
	// Delta is the time between invocations of epoll_wait, defaults to 200ms
	Delta time.Duration

	// Maximum number of events we can handle, defaults to 64
	Events uint
}

// Callback that is called when an event occurs
type FilePollCallback func(*os.File, FilePollMode)

// FilePoll interface
type FilePollInterface interface {
	gopi.Driver

	Watch(*os.File, FilePollMode, FilePollCallback) error
	Unwatch(*os.File) error
}

// private driver
type filepoll struct {
	log      gopi.Logger
	handle   int        // Poll file handle
	lock     sync.Mutex // Exclusive lock
	watchers map[int]*filepoll_watcher
	events   []syscall.EpollEvent
	delta    time.Duration
	ctlEvent syscall.EpollEvent // An event buffer for EpollCtl
	ctx      context.Context
	cancel   context.CancelFunc
	done     chan struct{}
}

// private watcher
type filepoll_watcher struct {
	events   uint32
	handle   *os.File
	callback FilePollCallback
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	FILEPOLL_MODE_READ   FilePollMode = syscall.EPOLLIN
	FILEPOLL_MODE_WRITE  FilePollMode = syscall.EPOLLOUT
	FILEPOLL_MODE_EDGE   FilePollMode = syscall.EPOLLET
	FILEPOLL_MODE_HANGUP FilePollMode = syscall.EPOLLHUP
	FILEPOLL_MODE_ERROR  FilePollMode = syscall.EPOLLERR
)

var (
	filepoll_delta  = time.Millisecond * 200
	filepoll_events = 64
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func (config FilePoll) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.linux.filepoll.Open>{ }")

	this := new(filepoll)
	this.log = log

	// Delta between epoll_wait invocations
	if config.Delta == 0 {
		this.delta = filepoll_delta
	} else {
		this.delta = config.Delta
	}

	// Array of events
	if config.Events == 0 {
		this.events = make([]syscall.EpollEvent, filepoll_events)
	} else {
		this.events = make([]syscall.EpollEvent, config.Events)
	}

	// Create poll
	if handle, err := syscall.EpollCreate1(syscall.EPOLL_CLOEXEC); err != nil {
		return nil, err
	} else {
		this.handle = handle
	}

	// Setup
	this.watchers = make(map[int]*filepoll_watcher)
	this.done = make(chan struct{})

	// Start the background event processor
	this.ctx, this.cancel = context.WithCancel(context.Background())
	go this.epollwait(this.ctx)

	return this, nil
}

func (this *filepoll) Close() error {
	this.log.Debug("<sys.hw.linux.filepoll.Close>{ }")

	// Unwatch all watched files
	for _, watcher := range this.watchers {
		if err := this.Unwatch(watcher.handle); err != nil {
			this.log.Warn("Unwatch: %v", err)
		}
	}

	// Cancel context
	this.cancel()

	// Wait for done from background thread
	_ = <-this.done

	// Make this method exclusive whilst closing
	this.lock.Lock()
	defer this.lock.Unlock()

	// Close poll
	if err := syscall.Close(this.handle); err != nil {
		return err
	}

	// Clear data structures
	this.watchers = nil
	this.done = nil
	this.events = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// ADD AND REMOVE FILES TO WATCH

// Add a file to watch for certain events
func (this *filepoll) Watch(handle *os.File, mode FilePollMode, callback FilePollCallback) error {
	this.log.Debug2("<sys.hw.linux.filepoll.Watch>{ fd=%v mode=%v }", handle.Fd(), mode)

	// Make this method exclusive
	this.lock.Lock()
	defer this.lock.Unlock()

	fd := int(handle.Fd())
	if fd <= 0 {
		// Bad parameter
		return gopi.ErrBadParameter
	}

	// Set up the event structure
	var already bool
	watcher, already := this.watchers[fd]
	this.ctlEvent.Fd = int32(fd)
	if already {
		this.ctlEvent.Events = watcher.events
	}

	switch mode {
	case FILEPOLL_MODE_EDGE:
		this.ctlEvent.Events |= syscall.EPOLLPRI
	case FILEPOLL_MODE_READ:
		this.ctlEvent.Events |= syscall.EPOLLIN
	case FILEPOLL_MODE_WRITE:
		this.ctlEvent.Events |= syscall.EPOLLOUT
	default:
		return gopi.ErrBadParameter
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
	this.watchers[fd] = &filepoll_watcher{
		events:   this.ctlEvent.Events,
		handle:   handle,
		callback: callback,
	}

	// return success
	return nil
}

func (this *filepoll) Unwatch(handle *os.File) error {
	this.log.Debug2("<sys.hw.linux.filepoll.Unwatch>{ fd=%v }", int(handle.Fd()))

	// Make this method exclusive
	this.lock.Lock()
	defer this.lock.Unlock()

	fd := int(handle.Fd())
	if fd < 0 {
		// file has been closed already, simply remove from map and return
		for k, v := range this.watchers {
			if v.handle == handle {
				delete(this.watchers, k)
			}
		}
		return nil
	}
	// Delete from epoll_ctl
	if _, exists := this.watchers[fd]; exists == false {
		return gopi.ErrBadParameter
	} else if err := syscall.EpollCtl(this.handle, syscall.EPOLL_CTL_DEL, fd, nil); err != nil {
		return os.NewSyscallError("epoll_ctl", err)
	}
	delete(this.watchers, fd)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *filepoll) String() string {
	return fmt.Sprintf("<sys.hw.linux.filepoll>{ handle=%v }", this.handle)
}

func (m FilePollMode) String() string {
	s := ""
	if m == 0 {
		return s
	}
	if m&FILEPOLL_MODE_READ != 0 {
		s = s + "FILEPOLL_MODE_READ|"
	}
	if m&FILEPOLL_MODE_WRITE != 0 {
		s = s + "FILEPOLL_MODE_WRITE|"
	}
	if m&FILEPOLL_MODE_EDGE != 0 {
		s = s + "FILEPOLL_MODE_EDGE|"
	}
	if m&FILEPOLL_MODE_HANGUP != 0 {
		s = s + "FILEPOLL_MODE_HANGUP|"
	}
	if m&FILEPOLL_MODE_ERROR != 0 {
		s = s + "FILEPOLL_MODE_ERROR|"
	}
	return s[0 : len(s)-1]
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *filepoll) epollwait(ctx context.Context) {
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			break FOR_LOOP
		default:
			if err := this.epollwait_inner(this.delta); err != nil {
				this.log.Error("linux.filepoll.epollwait: %v", err)
				break FOR_LOOP
			}
		}
	}
	this.done <- gopi.DONE
}

func (this *filepoll) epollwait_inner(delta time.Duration) error {
	this.lock.Lock()
	n, err := syscall.EpollWait(this.handle, this.events[:], int(delta.Seconds()*1E3))
	this.lock.Unlock()

	if err != nil {
		if err == syscall.EAGAIN || err == syscall.EINTR {
			return nil
		} else {
			return err
		}
	}

	if n == 0 {
		return nil
	}

	this.log.Debug2("<sys.hw.linux.filepoll> got event n=%v", n)

	// Process incoming events
	for _, event := range this.events[:n] {
		// determine modes
		var mode FilePollMode
		if (event.Events&syscall.EPOLLHUP) != 0 || (event.Events&syscall.EPOLLRDHUP != 0) {
			mode |= FILEPOLL_MODE_HANGUP
		}
		if event.Events&syscall.EPOLLERR != 0 {
			mode |= FILEPOLL_MODE_ERROR
		}
		if event.Events&syscall.EPOLLIN != 0 {
			mode |= FILEPOLL_MODE_READ
		}
		if event.Events&syscall.EPOLLOUT != 0 {
			mode |= FILEPOLL_MODE_WRITE
		}
		if event.Events&(syscall.EPOLLET&0xffffffff) != 0 {
			mode |= FILEPOLL_MODE_EDGE
		}
		this.epoll_event(&event, mode)
	}
	return nil
}

func (this *filepoll) epoll_event(event *syscall.EpollEvent, mode FilePollMode) {
	this.log.Debug2("<linux.filepoll.event>{ fd=%v events=%08X mode=%v }", event.Fd, event.Events, mode)
	if watcher, exists := this.watchers[int(event.Fd)]; exists {
		watcher.callback(watcher.handle, mode)
	} else {
		this.log.Warn("<linux.filepoll.event> Missing watcher for file descriptor %v", event.Fd)
	}
}
