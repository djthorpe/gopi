/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved
    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"os"
	"syscall"
	"encoding/binary"
	"io"
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// PRIVATE TYPES

type evPoll struct {
	handle int
	device *os.File
	event syscall.EpollEvent
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// Maximum number of events for polling
	INPUT_MAX_POLL_EVENTS = 32
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

func evNewPoll(device *os.File) (*evPoll, error) {
	var err error

	this := new(evPoll)
	this.device = device
	if this.handle, err = syscall.EpollCreate1(syscall.EPOLL_CLOEXEC); err != nil {
		return nil, err
	}

	// register the poll with the device
	this.event.Events = syscall.EPOLLIN
	this.event.Fd = int32(device.Fd())
	if err = syscall.EpollCtl(this.handle, syscall.EPOLL_CTL_ADD, int(this.event.Fd), &this.event); err != nil {
		return nil, err
	}

	// Success
	return this, nil
}

func (this *evPoll) Close() error {
	// deregister polling
	if err := syscall.EpollCtl(this.handle, syscall.EPOLL_CTL_DEL, int(this.event.Fd), &this.event); err != nil {
		return err
	}

	// Reset
	this.handle = 0
	this.device = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// WATCH FOR EVENTS

func (this *evPoll) evPoll(callback func (*evEvent)) error {
	var raw_event evEvent

	events := make([]syscall.EpollEvent,INPUT_MAX_POLL_EVENTS)
	for {
		n, err := syscall.EpollWait(this.handle, events, -1)
		if err != nil {
			if err == syscall.EINTR {
				continue
			} else {
				return err
			}
		}
		if n <= 0 {
			continue
		}
		for _, event := range events[:n] {
			if event.Fd != int32(event.Fd) {
				continue
			}
			err := binary.Read(this.device, binary.LittleEndian, &raw_event)
			if err == io.EOF {
				fmt.Println("EOF")
				return nil
			}
			if err != nil {
				return err
			}
			// process the event
			callback(&raw_event)
		}
	}
	// we should never get here
	return nil
}

