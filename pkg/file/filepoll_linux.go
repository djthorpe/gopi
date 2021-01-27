// +build linux

package file

import (
	"context"
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

type filepoll struct {
	gopi.Unit
	gopi.Logger
	sync.RWMutex
	pipe

	handle uintptr
	cap    uint
	funcs  map[uintptr]gopi.FilePollFunc
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	epollDefaultCapacity = 10
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *filepoll) New(gopi.Config) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if handle, err := linux.EpollCreate(); err != nil {
		return err
	} else {
		this.handle = handle
		this.cap = epollDefaultCapacity
		this.funcs = make(map[uintptr]gopi.FilePollFunc)
	}

	// Initialize pipe
	if err := this.pipe.Init(); err != nil {
		linux.EpollClose(this.handle)
		return err
	} else if err := linux.EpollAdd(this.handle, this.pipe.ReadFd(), linux.EPOLL_MODE_READ); err != nil {
		linux.EpollClose(this.handle)
		return err
	}

	// Return success
	return nil
}

func (this *filepoll) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Wake up pipe
	if err := this.pipe.Wake(); err != nil {
		return err
	}

	// Close polling
	if this.handle != 0 {
		if err := linux.EpollClose(this.handle); err != nil {
			return err
		}
	}

	// Close pipes
	if err := this.pipe.Close(); err != nil {
		return err
	}

	// Release resources
	this.handle = 0

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *filepoll) String() string {
	str := "<filepoll"
	if this.handle != 0 {
		str += " handle=" + fmt.Sprint(this.handle)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *filepoll) Watch(fd uintptr, mode gopi.FilePollFlags, handler gopi.FilePollFunc) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Convert flags
	flags := linux.EpollMode(0)
	if mode&gopi.FILEPOLL_FLAG_READ == gopi.FILEPOLL_FLAG_READ {
		flags |= linux.EPOLL_MODE_READ
	}
	if mode&gopi.FILEPOLL_FLAG_WRITE == gopi.FILEPOLL_FLAG_WRITE {
		flags |= linux.EPOLL_MODE_WRITE
	}

	// Add watcher and start background goroutine
	if _, exists := this.funcs[fd]; exists || fd == 0 || handler == nil || flags == 0 {
		return gopi.ErrBadParameter.WithPrefix("Watch")
	} else if err := linux.EpollAdd(this.handle, fd, flags); err != nil {
		return err
	} else {
		this.funcs[fd] = handler
	}

	// Success
	return nil
}

func (this *filepoll) Unwatch(fd uintptr) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if _, exists := this.funcs[fd]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("fd")
	} else if err := linux.EpollDelete(this.handle, fd); err != nil {
		return err
	} else {
		delete(this.funcs, fd)
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *filepoll) Run(ctx context.Context) error {
	stop := false

	go func() {
		<-ctx.Done()
		stop = true
		this.pipe.Wake()
	}()

	for stop == false {
		evts, err := linux.EpollWait(this.handle, 0, this.cap)
		if err != nil {
			this.Print("FilePoll: ", err)
		} else {
			for _, evt := range evts {
				// events need to be called in the right order
				this.call(uintptr(evt.Fd), maskToFlags(evt.Flags()))
			}
		}
	}

	return nil
}

func (this *filepoll) call(fd uintptr, flags gopi.FilePollFlags) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if fd == this.pipe.ReadFd() && flags == gopi.FILEPOLL_FLAG_READ {
		if err := this.pipe.Clear(); err != nil {
			this.Print("FilePoll: ", err)
		}
	} else if handler, exists := this.funcs[fd]; exists {
		handler(fd, flags)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

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
