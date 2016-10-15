/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md

	This package provides input mechanisms, including the touchscreen
	interface for the official Raspberry Pi LED.
*/
package input

import (
	"image"
	"time"
)

////////////////////////////////////////////////////////////////////////////////


////////////////////////////////////////////////////////////////////////////////

const (
	// References:

	// Event types
	// https://www.kernel.org/doc/Documentation/input/event-codes.txt
	EV_SYN uint16 = 0x0000
	EV_KEY uint16 = 0x0001
	EV_ABS uint16 = 0x0003

	// Button information
	BTN_TOUCH         uint16 = 0x014A
	BTN_TOUCH_RELEASE uint32 = 0x00000000
	BTN_TOUCH_PRESS   uint32 = 0x00000001

	// Multi-Touch Types
	// https://www.kernel.org/doc/Documentation/input/multi-touch-protocol.txt
	ABS_X              uint16 = 0x0000
	ABS_Y              uint16 = 0x0001
	ABS_MT_SLOT        uint16 = 0x002F // 47 MT slot being modified
	ABS_MT_POSITION_X  uint16 = 0x0035 // 53 Center X of multi touch position
	ABS_MT_POSITION_Y  uint16 = 0x0036 // 54 Center Y of multi touch position
	ABS_MT_TRACKING_ID uint16 = 0x0039 // 57 Unique ID of initiated contact
)


func (this *Device) waitForEvents(callback FT5406Callback) error {
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

func (this *Device) ProcessEvents() {
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

func (this *Device) processEvent(event *rawEvent) error {
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

func (this *Device) GetPosition() image.Point {
	return this.position
}

