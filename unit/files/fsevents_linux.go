// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package files

import (
	"os"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type fsevents struct {
	fh       *os.File
	filepoll gopi.FilePoll
	watch    map[uint32]string

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// INIT AND CLOSE

func (this *fsevents) Init(config FSEvents) error {
	if config.FilePoll == nil {
		return gopi.ErrBadParameter.WithPrefix("filepoll")
	} else {
		this.filepoll = config.FilePoll
	}
	if fh, err := linux.InotifyInit(); err != nil {
		return err
	} else if err := this.filepoll.Watch(fh.Fd(), gopi.FILEPOLL_FLAG_READ, this.FilepollRead); err != nil {
		fh.Close()
		return err
	} else {
		this.fh = fh
		this.watch = make(map[uint32]string)
	}

	return nil
}

func (this *fsevents) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Remove Inotify watchers
	errs := gopi.NewCompoundError()
	for watch := range this.watch {
		errs.Add(linux.InotifyRemoveWatch(this.fh.Fd(), watch))
	}
	if errs.ErrorOrSelf() != nil {
		return errs
	}

	// Stop watching Inotify
	if err := this.filepoll.Unwatch(this.fh.Fd()); err != nil {
		return err
	}

	// Close filehandle
	if this.fh != nil {
		if err := this.fh.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.watch = nil
	this.fh = nil
	this.filepoll = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.FSEvents

func (this *fsevents) Watch(path string, flags gopi.FSEventFlags, cb gopi.FSEventFunc) (uint32, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check incoming parameters
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, gopi.ErrNotFound.WithPrefix("path")
	} else if err != nil {
		return 0, err
	}
	if cb == nil {
		return 0, gopi.ErrBadParameter.WithPrefix("FSEventFunc")
	}
	if flags == gopi.FSEVENT_FLAG_NONE {
		flags = gopi.FSEVENT_FLAG_ALL
	}

	// Start watching
	if watch, err := linux.InotifyAddWatch(this.fh.Fd(), path, flagsToInotifyMask(flags)); err != nil {
		return 0, err
	} else {
		this.watch[watch] = path
		return watch, nil
	}
}

func (this *fsevents) Unwatch(watch uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if _, exists := this.watch[watch]; exists == false {
		return gopi.ErrNotFound.WithPrefix("Unwatch")
	} else if err := linux.InotifyRemoveWatch(this.fh.Fd(), watch); err != nil {
		return err
	} else {
		delete(this.watch, watch)
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *fsevents) FilepollRead(fd uintptr, flags gopi.FilePollFlags) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if fd == this.fh.Fd() && flags&gopi.FILEPOLL_FLAG_READ > 0 {
		if evt, err := linux.InotifyRead(fd); err != nil {
			this.Log.Error(err)
		} else {
			// TODO
			this.Log.Debug("TODO", evt)
		}
	}
}

func flagsToInotifyMask(flags gopi.FSEventFlags) linux.InotifyMask {
	mask := linux.IN_NONE
	if flags&gopi.FSEVENT_FLAG_ATTRIB > 0 {
		mask |= linux.IN_ATTRIB
	}
	if flags&gopi.FSEVENT_FLAG_CREATE > 0 {
		mask |= linux.IN_CREATE
	}
	if flags&gopi.FSEVENT_FLAG_DELETE > 0 {
		mask |= linux.IN_DELETE | linux.IN_DELETE_SELF
	}
	if flags&gopi.FSEVENT_FLAG_MODIFY > 0 {
		mask |= linux.IN_MODIFY
	}
	if flags&gopi.FSEVENT_FLAG_MOVE > 0 {
		mask |= linux.IN_MOVE_SELF | linux.IN_MOVED_FROM | linux.IN_MOVED_TO
	}
	if flags&gopi.FSEVENT_FLAG_UNMOUNT > 0 {
		mask |= linux.IN_UNMOUNT
	}
	return mask
}

/*
func (this *epoller) Watch(t *testing.T, stop <-chan struct{}) {
	defer this.WaitGroup.Done()
FOR_LOOP:
	for {
		select {
		case <-stop:
			break FOR_LOOP
		default:
			if evts, err := linux.EpollWait(this.fd, time.Millisecond*100, 10); err != nil {
				t.Error(err)
			} else {
				for _, evt := range evts {
					if uintptr(evt.Fd) == this.inotify && evt.Flags() == linux.EPOLL_MODE_READ {
						this.ReadNotify(t, this.inotify)
					}
				}
			}
		}
	}
}

func (this *epoller) ReadNotify(t *testing.T, fd uintptr) {
	if evt, err := linux.InotifyRead(fd); err != nil {
		t.Error(err)
	} else {
		t.Log(evt)
	}
}
*/
