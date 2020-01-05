// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package spi

import (
	"fmt"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	linux "github.com/djthorpe/gopi/v2/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *spi) Init(config SPI) error {
	// Set SPI parameters
	this.bus = config.Bus
	this.slave = config.Slave
	this.delay_usec = uint16(config.Delay)

	// Open the device
	if dev, err := linux.SPIOpenDevice(this.bus, this.slave); err != nil {
		return err
	} else {
		this.dev = dev
	}

	// Get current mode, speed and bits per word
	if mode, err := linux.SPIMode(this.dev.Fd()); err != nil {
		return err
	} else {
		this.mode = mode
	}
	if speed_hz, err := linux.SPISpeedHz(this.dev.Fd()); err != nil {
		return err
	} else {
		this.speed_hz = speed_hz
	}
	if bits_per_word, err := linux.SPIBitsPerWord(this.dev.Fd()); err != nil {
		return err
	} else {
		this.bits_per_word = bits_per_word
	}

	// Success
	return nil
}

func (this *spi) Close() error {
	if this.dev != nil {
		if err := this.dev.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.dev = nil

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *spi) String() string {
	if this.Unit.Closed {
		return "<gopi.SPI>"
	} else {
		return "<gopi.SPI bus=" + fmt.Sprint(this.bus) + " slave=" + fmt.Sprint(this.slave) + " mode=" + fmt.Sprint(this.mode) + " max_speed=" + fmt.Sprint(this.speed_hz) + "Hz " + " bits_per_word=" + fmt.Sprint(this.bits_per_word) + ">"
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.SPI

func (this *spi) SetMode(mode gopi.SPIMode) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetMode(this.dev.Fd(), mode); err != nil {
		return err
	} else if mode_, err := linux.SPIMode(this.dev.Fd()); err != nil {
		return err
	} else if mode != mode_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("mode")
	} else {
		this.mode = mode
		return nil
	}
}

func (this *spi) SetMaxSpeedHz(speed uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetSpeedHz(this.dev.Fd(), speed); err != nil {
		return err
	} else if speed_, err := linux.SPISpeedHz(this.dev.Fd()); err != nil {
		return err
	} else if speed != speed_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("speed")
	} else {
		this.speed_hz = speed
		return nil
	}
}

func (this *spi) SetBitsPerWord(bits uint8) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetBitsPerWord(this.dev.Fd(), bits); err != nil {
		return err
	} else if bits_, err := linux.SPIBitsPerWord(this.dev.Fd()); err != nil {
		return err
	} else if bits != bits_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("bits")
	} else {
		this.bits_per_word = bits
		return nil
	}
}

func (this *spi) Transfer(send []byte) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(send) == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("send")
	} else if receive, err := linux.SPITransfer(this.dev.Fd(), send, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return nil, err
	} else {
		return receive, nil
	}
}

func (this *spi) Read(len uint32) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("len")
	} else if receive, err := linux.SPIRead(this.dev.Fd(), len, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return nil, err
	} else {
		return receive, nil
	}
}

func (this *spi) Write(send []byte) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(send) == 0 {
		return gopi.ErrBadParameter.WithPrefix("send")
	} else if err := linux.SPIWrite(this.dev.Fd(), send, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return err
	} else {
		return nil
	}
}
