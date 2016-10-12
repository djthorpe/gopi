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
	"errors"
	"time"
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

const ( // https://www.kernel.org/doc/Documentation/input/event-codes.txt
	EV_SYN uint16 = 0
	EV_KEY uint16 = 1
	EV_ABS uint16 = 3
	BTN_TOUCH uint16 = 330
	BTN_TOUCH_RELEASE uint32 = 0
	BTN_TOUCH_PRESS uint32 = 1
	ABS_X uint16 = 0
	ABS_Y uint16 = 1
	ABS_MT_SLOT uint16 = 0x2F // 47 MT slot being modified
	ABS_MT_POSITION_X uint16 = 0x35 // 53 Center X of multi touch position
	ABS_MT_POSITION_Y uint16 = 0x36 // 54 Center Y of multi touch position
	ABS_MT_TRACKING_ID uint16 = 0x39 // 57 Unique ID of initiated contact
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
				continue
			}
			err = this.ProcessEvent(&event)
			if err != nil {
				fmt.Println("Decode error:",err)
			}
		}
	})
}

func (this *FT5406) ProcessEvent(event *FT5406Event) error {
	timestamp := time.Duration(time.Duration(event.Second) * time.Second + time.Duration(event.Microsecond) * time.Microsecond)
	switch {
	case event.Type == EV_SYN:
		fmt.Println("SYNC:",timestamp)
		return nil
	case event.Type == EV_KEY && event.Code == BTN_TOUCH && event.Value == BTN_TOUCH_PRESS:
		fmt.Println("BTN_TOUCH_PRESS")
		return nil
	case event.Type == EV_KEY && event.Code == BTN_TOUCH && event.Value == BTN_TOUCH_RELEASE:
		fmt.Println("BTN_TOUCH_RELEASE")
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_SLOT:
		fmt.Println("SLOT:",event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_POSITION_X:
		fmt.Println("X:",event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_POSITION_Y:
		fmt.Println("Y:",event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_TRACKING_ID:
		fmt.Println("ID:",event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_X:
		fmt.Println("ABS X:",event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_Y:
		fmt.Println("ABS X:",event.Value)
		return nil
	default:
		return errors.New(fmt.Sprintf("Invalid event: %v",event))
	}
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
