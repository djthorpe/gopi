/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package mock

import (
	"context"
	"fmt"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Input struct{}

type input struct {
	log         gopi.Logger
	subscribers []chan gopi.Event
	devices     []gopi.InputDevice
}

type event struct {
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Input) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.mock.Input.Open{  }")

	this := new(input)
	this.log = logger
	this.subscribers = make([]chan gopi.Event, 0)
	this.devices = make([]gopi.InputDevice, 0)

	// Success
	return this, nil
}

// Close
func (this *input) Close() error {
	this.log.Debug("sys.mock.Input.Close{ }")

	for _, device := range this.devices {
		if device == nil {
			// Ignore closed devices
		} else if err := device.Close(); err != nil {
			// Error closing device - report the error
			this.log.Error("%v", err)
		}
	}

	this.subscribers = nil
	this.devices = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// INPUT INTERFACE

// Open Devices by name, type and bus
func (this *input) OpenDevicesByName(name string, flags gopi.InputDeviceType, bus gopi.InputDeviceBus) ([]gopi.InputDevice, error) {
	return nil, nil
}

// Add Device
func (this *input) AddDevice(device gopi.InputDevice) error {
	// Check for nil argument
	if device == nil {
		return gopi.ErrBadParameter
	}
	// TODO
	return gopi.ErrNotImplemented
}

// Close Device
func (this *input) CloseDevice(device gopi.InputDevice) error {
	for i := range this.devices {
		if this.devices[i] == nil {
			// ignore closed devices
		} else if this.devices[i] != device {
			// ignore non-matching device
		} else if err := device.Close(); err != nil {
			// error closing device
			return err
		} else {
			// set device to nil
			this.devices[i] = nil
		}
	}
	return nil
}

// Return a list of open devices
func (this *input) GetOpenDevices() []gopi.InputDevice {
	opened := make([]gopi.InputDevice, 0, len(this.devices))
	for i := range this.devices {
		if this.devices[i] != nil {
			opened = append(opened, this.devices[i])
		}
	}
	return opened
}

// Watch for events with context
func (this *input) Watch(ctx context.Context) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// PUBLISHER INTERFACE

func (this *input) Subscribe() chan gopi.Event {
	subscriber := make(chan gopi.Event)
	this.subscribers = append(this.subscribers, subscriber)
	return subscriber
}

// Unsubscribe from events emitted
func (this *input) Unsubscribe(subscriber chan gopi.Event) {
	for i, s := range this.subscribers {
		if subscriber == s {
			this.subscribers[i] = nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// EVENT INTERFACE

// Source of the event
func (this *event) Source() gopi.Driver {
	return nil
}

// Name of the event
func (this *event) Name() string {
	return "sys.mock.Input.Event"
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *input) String() string {
	return fmt.Sprintf("<sys.mock.Input>{ devices=%v }", this.devices)
}
