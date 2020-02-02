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
	handle  uintptr
	cap     uint
	stop    chan struct{}
	handler map[uintptr]gopi.FilePollFunc

	Pipe
	base.Unit
	sync.Mutex
	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EPOLL_DEFAULT_CAP = 10
)

////////////////////////////////////////////////////////////////////////////////
// INIT & CLOSE

func (this *filepoll) Init(config FilePoll) error {
	this.Lock()
	defer this.Unlock()

	if handle, err := linux.EpollCreate(); err != nil {
		return err
	} else {
		this.handle = handle
		this.stop = make(chan struct{})
		this.cap = EPOLL_DEFAULT_CAP
		this.handler = make(map[uintptr]gopi.FilePollFunc)
	}

	// Initialize pipe
	if err := this.Pipe.Init(); err != nil {
		linux.EpollClose(this.handle)
		return err
	} else if err := linux.EpollAdd(this.handle,this.Pipe.ReadFd(),linux.EPOLL_MODE_READ); err != nil {
		linux.EpollClose(this.handle)
		return err
	} 

	// Watch for events in background
	go this.watch(this.stop)

	// Return success
	return nil
}

func (this *filepoll) Close() error {	
	// Unwatch all events
	for fd := range this.handler {
		if err := this.Unwatch(fd); err != nil {
			return err
		}
	}

	// Indicate end to goroutine
	close(this.stop)

	// Wake up go routine
	if err := this.Pipe.Wake(); err != nil {
		return err
	}

	// Wait for gooutine to end
	this.WaitGroup.Wait()

	// Close polling
	if this.handle != 0 {
		if err := linux.EpollClose(this.handle); err != nil {
			return err
		}
	}

	// Close pipes
	if err := this.Pipe.Close(); err != nil {
		return err
	}

	// Release resources
	this.handler = nil
	this.handle = 0
	this.stop = nil

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

////////////////////////////////////////////////////////////////////////////////
// WATCH AND UNWATCH

func (this *filepoll) Watch(fd uintptr, mode gopi.FilePollFlags, handler gopi.FilePollFunc) error {
	this.Lock()
	defer this.Unlock()

	// Convert flags
	flags := linux.EpollMode(0)
	if mode&gopi.FILEPOLL_FLAG_READ == gopi.FILEPOLL_FLAG_READ {
		flags |= linux.EPOLL_MODE_READ
	}
	if mode&gopi.FILEPOLL_FLAG_WRITE == gopi.FILEPOLL_FLAG_WRITE {
		flags |= linux.EPOLL_MODE_WRITE
	}

	// Add watcher and start background goroutine
	if fd == 0 {
		return gopi.ErrBadParameter.WithPrefix("fd")
	} else if _, exists := this.handler[fd]; exists {
		return gopi.ErrDuplicateItem.WithPrefix("fd")
	} else if err := linux.EpollAdd(this.handle, fd, flags); err != nil {
		return err
	} else {
		this.handler[fd] = handler
	}

	// Success
	return nil
}

func (this *filepoll) Unwatch(fd uintptr) error {
	this.Lock()
	defer this.Unlock()

	if fd == 0 {
		return gopi.ErrBadParameter.WithPrefix("fd")
	} else if _, exists := this.handler[fd]; exists == false {
		return gopi.ErrNotFound.WithPrefix("fd")
	} else if err := linux.EpollDelete(this.handle,fd); err != nil {
		return err
	} else {
		delete(this.handler,fd)
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *filepoll) watch(stop <-chan struct{}) {
	this.WaitGroup.Add(1)
	defer this.WaitGroup.Done()

FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		default:
			if evts, err := linux.EpollWait(this.handle,0, this.cap); err != nil {
				this.Log.Error(err)
			} else {
				for _, evt := range evts {
					this.call(uintptr(evt.Fd),maskToFlags(evt.Flags()))
				}
			}
		}
	}
}

func (this *filepoll) call(fd uintptr,flags gopi.FilePollFlags) {
	this.Lock()
	defer this.Unlock()
	if fd == this.Pipe.ReadFd() && flags == gopi.FILEPOLL_FLAG_READ {
		if err := this.Pipe.Clear(); err != nil {
			this.Log.Error(err)
		}
	} else if handler,exists := this.handler[fd]; exists {
		handler(fd,flags)
	} else {
		this.Log.Warn("Filepoll: Unable to handle fd=",fd," flags=",flags)
	}
}

func maskToFlags(mask linux.EpollMode) gopi.FilePollFlags {
	flags := gopi.FILEPOLL_FLAG_NONE
	if mask&linux.EPOLL_MODE_READ == linux.EPOLL_MODE_READ {
		flags |= gopi.FILEPOLL_FLAG_READ
	}
	if mask&linux.EPOLL_MODE_WRITE == linux.EPOLL_MODE_WRITE {
		flags |= gopi.FILEPOLL_FLAG_WRITE
	}
	return flags
}

