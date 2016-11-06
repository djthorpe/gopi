/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This package file implements the abstract hardware interfaces for GPIO
package hw // import "github.com/djthorpe/gopi/hw"

import (
	"fmt"
	"errors"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////

// LED Configuration
type LED struct {
	// the gpio interface
	GPIO GPIODriver

	// the array of pins which are connected to LED's
	Pins []GPIOPin
}

// LED interface
type LEDDriver interface {
	// Enforces general driver
	gopi.Driver

	// Return array of available LED pins
	Pins() []GPIOPin

	// Start LED. When no arguments are given, all pins are switched to on state
	On(index ...uint) error

	// Stop LED. When no arguments are given, all pins are switched to off state
	Off(index ...uint) error
}

// LED state
type LEDDevice struct {
	gpio           GPIODriver
	pins           []GPIOPin
	log            *util.LoggerDevice
	finish_channel chan bool
	done_channel   chan bool
	on_channel     chan []GPIOPin
	off_channel    chan []GPIOPin
}

////////////////////////////////////////////////////////////////////////////////
// Open and close LED

func (config LED) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug2("<hw.LED>Open")

	// If no pins, then return error
	if len(config.Pins) == 0 {
		return nil, log.Error("No pins specified")
	}

	// Create the driver
	device := new(LEDDevice)
	device.gpio = config.GPIO
	device.pins = config.Pins
	device.log = log
	device.finish_channel = make(chan bool)
	device.done_channel = make(chan bool)
	device.on_channel = make(chan []GPIOPin)
	device.off_channel = make(chan []GPIOPin)

	// Set pins to OUTPUT mode and set to OFF
	for _, pin := range device.pins {
		device.gpio.SetPinMode(pin, GPIO_OUTPUT)
		device.gpio.WritePin(pin, GPIO_LOW)
	}

	// Start background goroutine
	go device.runLoop()

	// Success
	return device, nil
}

func (this *LEDDevice) Close() error {
	this.log.Debug2("<hw.LED>Close")

	// Quit - send signal to finish and then get back done signal
	this.finish_channel <- true
	<-this.done_channel

	// Switch off
	for _, pin := range this.pins {
		this.gpio.WritePin(pin, GPIO_LOW)
	}

	return nil
}

func (this *LEDDevice) String() string {
	return fmt.Sprintf("<hw.LED>{ Pins: %v }", this.pins)
}

func (this *LEDDevice) Pins() []GPIOPin {
	return this.pins
}

////////////////////////////////////////////////////////////////////////////////
// Switch on and off

func (this *LEDDevice) On(index ...uint) error {
	pins, err := this.getPins(index)
	if err != nil {
		return err
	}
	this.on_channel <- pins
	return nil
}

func (this *LEDDevice) Off(index ...uint) error {
	pins, err := this.getPins(index)
	if err != nil {
		return err
	}
	this.off_channel <- pins
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Private methods

func (this *LEDDevice) runLoop() {
	done := false
	for done == false {
		select {
		case <-this.finish_channel:
			// If we receive a finish signal, then break
			done = true
			break
		case pins := <-this.on_channel:
			for _,pin := range(pins) {
				this.gpio.WritePin(pin,GPIO_HIGH)
			}
		case pins := <-this.off_channel:
			for _,pin := range(pins) {
				this.gpio.WritePin(pin,GPIO_LOW)
			}
		}
	}
	this.done_channel <- done
}

func (this *LEDDevice) getPins(pins []uint) ([]GPIOPin, error) {
	return nil,errors.New("Invalid pins")
}

