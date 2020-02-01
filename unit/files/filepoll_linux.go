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
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

type filepoll struct {
	handle  uintptr
	timeout time.Duration
	cap     uint
	stop    map[uintptr]chan struct{}

	base.Unit
	sync.Mutex
	sync.WaitGroup
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	EPOLL_DEFAULT_CAPCITY = 10
	EPOLL_DEFAULT_TIMEOUT = 100 * time.Millisecond
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
		this.stop = make(map[uintptr]chan struct{})
		this.cap = EPOLL_DEFAULT_CAPCITY
		this.timeout = EPOLL_DEFAULT_TIMEOUT
	}

	// Return success
	return nil
}

func (this *filepoll) Close() error {
	// Unwatch all events
	for k := range this.stop {
		if err := this.Unwatch(k); err != nil {
			return err
		}
	}

	// Wait for all gooutines to end
	this.WaitGroup.Wait()

	// Close polling
	if this.handle != 0 {
		if err := linux.EpollClose(this.handle); err != nil {
			return err
		}
	}

	// Release resources
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
	} else if _, exists := this.stop[fd]; exists {
		return gopi.ErrDuplicateItem.WithPrefix("fd")
	} else if err := linux.EpollAdd(this.handle, fd, flags); err != nil {
		return err
	} else {
		this.stop[fd] = make(chan struct{})
		this.WaitGroup.Add(1)
		go this.watch(fd, handler,this.stop[fd])
	}

	// Success
	return nil
}

func (this *filepoll) Unwatch(fd uintptr) error {
	this.Lock()
	defer this.Unlock()

	if fd == 0 {
		return gopi.ErrBadParameter.WithPrefix("fd")
	} else if stop, exists := this.stop[fd]; exists == false {
		return gopi.ErrNotFound.WithPrefix("fd")
	} else {
		// Send stop
		stop <- struct{}{}
		<-stop
		delete(this.stop, fd)
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// BACKGROUND METHODS

func (this *filepoll) watch(fd uintptr, handler gopi.FilePollFunc,stop chan struct{}) {
	defer this.WaitGroup.Done()

	// Continue receiving epoll events until stop channel is closed
	timer := time.NewTicker(this.timeout)
FOR_LOOP:
	for {
		select {
		case <-timer.C:
			if evts, err := linux.EpollWait(this.handle, this.timeout, this.cap); err != nil {
				this.Log.Error(err)
			} else {
				for _, evt := range evts {
					if evt.Flags()&linux.EPOLL_MODE_READ == linux.EPOLL_MODE_READ {
						handler(uintptr(evt.Fd), gopi.FILEPOLL_FLAG_READ)
					}
					if evt.Flags()&linux.EPOLL_MODE_WRITE == linux.EPOLL_MODE_WRITE {
						handler(uintptr(evt.Fd), gopi.FILEPOLL_FLAG_WRITE)
					}
				}
			}
		case <-stop:
			timer.Stop()
			break FOR_LOOP
		}
	}

	// Delete the watcher
	if err := linux.EpollDelete(this.handle, fd); err != nil {
		this.Log.Error(err)
	}

	// Close stop
	close(stop)
}
