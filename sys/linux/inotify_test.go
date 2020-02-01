// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package linux_test

import (
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2/sys/linux"
)

func Test_INotify_000(t *testing.T) {
	t.Log("Test_INotify_000")
}

func Test_INotify_001(t *testing.T) {
	if fh, err := linux.InotifyInit(); err != nil {
		t.Error(err)
	} else if err := fh.Close(); err != nil {
		t.Error(err)
	}
}

func Test_INotify_002(t *testing.T) {
	fh, err := linux.InotifyInit()
	if err != nil {
		t.Fatal(err)
	}
	poll := NewEventPoller(t, fh.Fd())
	if poll == nil {
		t.Fatal("NewEventPoller failed")
	}
	time.Sleep(time.Second)
	if err := poll.Close(t); err != nil {
		t.Fatal(err)
	}
	if err := fh.Close(); err != nil {
		t.Fatal(err)
	}
}

func Test_INotify_003(t *testing.T) {
	fh, err := linux.InotifyInit()
	if err != nil {
		t.Fatal(err)
	}
	poll := NewEventPoller(t, fh.Fd())
	if poll == nil {
		t.Fatal("NewEventPoller failed")
	}
	t.Log("Watching", os.TempDir())
	if watch, err := linux.InotifyAddWatch(fh.Fd(), os.TempDir(), linux.IN_DEFAULT); err != nil {
		t.Error(err)
	} else {
		time.Sleep(time.Second * 1)
		t.Log("Creating tmp file")
		if tmpfile, err := ioutil.TempFile("", "inotify"); err != nil {
			t.Error(err)
		} else {
			t.Log("Removing tmp file", tmpfile.Name())
			tmpfile.Close()
			os.Remove(tmpfile.Name())
		}
		time.Sleep(time.Second * 1)
		if err := linux.InotifyRemoveWatch(fh.Fd(), watch); err != nil {
			t.Error(err)
		}
	}
	if err := poll.Close(t); err != nil {
		t.Fatal(err)
	}
	if err := fh.Close(); err != nil {
		t.Fatal(err)
	}
}

////////////////////////////////////////////////////////////////////////////////
// EPOLL HARNESS

type epoller struct {
	fd      uintptr
	stop    chan struct{}
	inotify uintptr

	sync.WaitGroup
}

func NewEventPoller(t *testing.T, inotify uintptr) *epoller {
	this := new(epoller)
	if fd, err := linux.EpollCreate(); err != nil {
		t.Fatal(err)
		return nil
	} else if err := linux.EpollAdd(fd, inotify, linux.EPOLL_MODE_READ); err != nil {
		linux.EpollClose(fd)
		t.Error(err)
		return nil
	} else {
		this.fd = fd
		this.inotify = inotify
		this.stop = make(chan struct{})
		this.WaitGroup.Add(1)
		go this.Watch(t, this.stop)
		return this
	}
}

func (this *epoller) Close(t *testing.T) error {

	// Wait for go routine to end
	this.stop <- struct{}{}
	this.WaitGroup.Wait()

	// Close epoll
	if err := linux.EpollClose(this.fd); err != nil {
		t.Error(err)
		return err
	} else {
		return nil
	}
}

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
