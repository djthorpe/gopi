/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package input

import (
	"os"
	"io"
	"path"
	"regexp"
	"syscall"
	"io/ioutil"
    "path/filepath"
	"encoding/binary"
)

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////

type FT5406 struct {
	device string
	file   *os.File
	poll   int
	event  syscall.EpollEvent
}

type FT5406Callback func (syscall.EpollEvent)

type FT5406Event struct {
	Second uint32
	Microsecond uint32
	Type uint16
	Code uint16
	Value uint32
}

////////////////////////////////////////////////////////////////////////////////

const (
	PATH_INPUT_DEVICES = "/sys/class/input/event*"
	MAX_POLL_EVENTS = 32
	MAX_EVENT_SIZE_BYTES = 1024
)

const (

)


////////////////////////////////////////////////////////////////////////////////

var (
	REGEXP_DEVICENAME = regexp.MustCompile("^FT5406")
)

////////////////////////////////////////////////////////////////////////////////
// Create Touchscreen

func NewFT5406() (*FT5406, error) {
	var err error

	this := new(FT5406)
	this.device, err = getDeviceName()
	if err != nil {
		return nil,err
	}

	// open driver
	this.file, err = os.Open(this.device)
	if err != nil {
		return nil, err
	}

	// create a poll
	this.poll, err = syscall.EpollCreate1(0)
	if err != nil {
		this.file.Close()
		return nil, err
	}

	// register the poll with the device
	this.event.Events = syscall.EPOLLIN
	this.event.Fd = int32(this.file.Fd())
	if err = syscall.EpollCtl(this.poll,syscall.EPOLL_CTL_ADD,int(this.event.Fd),&this.event); err != nil {
		syscall.Close(this.poll)
		this.file.Close()
		return nil,err
	}

	return this, nil
}

func (this *FT5406) Close() error {
	syscall.Close(this.poll)
	this.file.Close()
	return nil
}

func (this *FT5406) waitForEvents(callback FT5406Callback) error {
	events := make([]syscall.EpollEvent,MAX_POLL_EVENTS)
	for {
		n, err := syscall.EpollWait(this.poll, events, -1)
		if err != nil {
			return err
		}
		if n <= 0 {
			continue
		}
		for _, event := range events[:n] {
			if event.Fd != int32(this.file.Fd()) {
				continue
			}
			callback(event)
		}
	}
}

func (this *FT5406) ProcessEvents() {
	this.waitForEvents(func(event syscall.EpollEvent) {
		for {
			var event FT5406Event
			err := binary.Read(this.file,binary.LittleEndian,&event)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Decode error:",err)
			} else {
				fmt.Println(event)
			}
		}
	})
}

////////////////////////////////////////////////////////////////////////////////
// Private Methods

func getDeviceName() (string,error) {
	files, err := filepath.Glob(PATH_INPUT_DEVICES)
	if err != nil {
		return "",err
	}
	for _, file := range files {
		buf, err := ioutil.ReadFile(path.Join(file,"device","name"))
		if err != nil {
			continue
		}
		if REGEXP_DEVICENAME.Match(buf) {
			return path.Join("/","dev","input",path.Base(file)),nil
		}
	}
	return "",nil
}
