/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package input

import (
	"encoding/binary"
	"errors"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
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
	slots  []*TouchEvent
	slot     uint32
	position image.Point
}

type FT5406Callback func(syscall.EpollEvent)

////////////////////////////////////////////////////////////////////////////////

const (
	PATH_INPUT_DEVICES   = "/sys/class/input/event*"
	MAX_POLL_EVENTS      = 32
	MAX_EVENT_SIZE_BYTES = 1024
)

////////////////////////////////////////////////////////////////////////////////

var (
	REGEXP_DEVICENAME = regexp.MustCompile("^FT5406")
)

////////////////////////////////////////////////////////////////////////////////
// Create Touchscreen

func NewFT5406(slots uint) (*FT5406, error) {
	var err error

	// check for inputs
	if slots == 0 {
		return nil, errors.New("Invalid slots value")
	}

	this := new(FT5406)
	this.device, err = getDeviceName()
	if err != nil {
		return nil, err
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
	if err = syscall.EpollCtl(this.poll, syscall.EPOLL_CTL_ADD, int(this.event.Fd), &this.event); err != nil {
		syscall.Close(this.poll)
		this.file.Close()
		return nil, err
	}

	// set position and slot to zero, plus create the slots
	this.position.X = 0
	this.position.Y = 0
	this.slot = 0
	this.slots = make([]*TouchEvent, slots)

	// set up the slot values
	for i, _ := range this.slots {
		this.slots[i] = &TouchEvent{ Slot: uint32(i) }
	}

	return this, nil
}

func (this *FT5406) Close() error {
	syscall.Close(this.poll)
	this.file.Close()
	return nil
}

func (this *FT5406) waitForEvents(callback FT5406Callback) error {
	events := make([]syscall.EpollEvent, MAX_POLL_EVENTS)
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
			var event rawEvent
			err := binary.Read(this.file, binary.LittleEndian, &event)
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Decode error:", err)
				continue
			}
			err = this.processEvent(&event)
			if err != nil {
				fmt.Println("Decode error:", err)
			}
		}
	})
}

func (this *FT5406) processEvent(event *rawEvent) error {
	switch {
	case event.Type == EV_SYN:
		this.slots[this.slot].Timestamp = time.Duration(time.Duration(event.Second)*time.Second + time.Duration(event.Microsecond)*time.Microsecond)
		fmt.Println("SYNC:",this.slots[this.slot])
		return nil
	case event.Type == EV_KEY && event.Code == BTN_TOUCH && event.Value == BTN_TOUCH_PRESS:
		fmt.Println("BTN_TOUCH_PRESS")
		return nil
	case event.Type == EV_KEY && event.Code == BTN_TOUCH && event.Value == BTN_TOUCH_RELEASE:
		fmt.Println("BTN_TOUCH_RELEASE")
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_SLOT:
		if event.Value >= uint32(len(this.slots)) {
			return errors.New("Invalid slot value")
		}
		this.slot = event.Value
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_POSITION_X:
		this.slots[this.slot].Point.X = int(event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_POSITION_Y:
		this.slots[this.slot].Point.Y = int(event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_MT_TRACKING_ID:
		// Identifier is a 16 bit value which we turn into an int
		id := int(int16(event.Value))
		if id == -1 {
			fmt.Println("RELEASE:",this.slots[this.slot])
		} else {
			this.slots[this.slot].Identifier = id
		}
		return nil
	case event.Type == EV_ABS && event.Code == ABS_X:
		this.position.X = int(event.Value)
		return nil
	case event.Type == EV_ABS && event.Code == ABS_Y:
		this.position.Y = int(event.Value)
		return nil
	default:
		return errors.New(fmt.Sprintf("Invalid event: %v", event))
	}
}

func (this *FT5406) GetPosition() image.Point {
	return this.position
}

////////////////////////////////////////////////////////////////////////////////
// Private Methods

func getDeviceName() (string, error) {
	files, err := filepath.Glob(PATH_INPUT_DEVICES)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		buf, err := ioutil.ReadFile(path.Join(file, "device", "name"))
		if err != nil {
			continue
		}
		if REGEXP_DEVICENAME.Match(buf) {
			return path.Join("/", "dev", "input", path.Base(file)), nil
		}
	}
	return "", nil
}
