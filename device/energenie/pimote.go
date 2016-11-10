/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

// This package implements an Energenie Pimote device. Here is the pinouts
// for the pins on the pimote, see the following document for more information:
// https://energenie4u.co.uk/res/pdfs/ENER314%20UM.pdf
//
//  K0 GPIO17 / pin 11
//  K1 GPIO22 / pin 15
//  K2 GPIO23 / pin 16
//  K3 GPIO27 / pin 13
//  MODSEL GPIO24 / pin 18 (low OOK high FSK)
//  CE MODULATOR ENABLE GPIO25 / pin 22 (low off high on)
//
package energenie /* import "github.com/djthorpe/gopi/device/energenie" */

import (
	"time"
)

import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////

// Pimote Configuration
type Pimote struct {
	// the gpio interface
	GPIO hw.GPIODriver
}

// Pimote Driver
type PimoteDriver struct {
	gpio hw.GPIODriver
	log  *util.LoggerDevice
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	PIMOTE_K0         = hw.GPIOPin(17)
	PIMOTE_K1         = hw.GPIOPin(22)
	PIMOTE_K2         = hw.GPIOPin(23)
	PIMOTE_K3         = hw.GPIOPin(27)
	PIMOTE_MOD_SEL    = hw.GPIOPin(24)
	PIMOTE_MOD_EN     = hw.GPIOPin(25)
	PIMOTE_SOCKET_MIN = 1
	PIMOTE_SOCKET_MAX = 4
)

////////////////////////////////////////////////////////////////////////////////
// Open and close

func (config Pimote) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug2("<energenie.Pimote>Open")

	this := new(PimoteDriver)
	this.gpio = config.GPIO
	this.log = log

	// set output pins low
	for _, pin := range []hw.GPIOPin{PIMOTE_K0, PIMOTE_K1, PIMOTE_K2, PIMOTE_K3, PIMOTE_MOD_SEL, PIMOTE_MOD_EN} {
		this.gpio.SetPinMode(pin, hw.GPIO_OUTPUT)
		this.gpio.WritePin(pin, hw.GPIO_LOW)
	}

	// Return success
	return this, nil
}

func (this *PimoteDriver) Close() error {
	this.log.Debug2("<energenie.Pimote>Close")

	// set output pins low
	for _, pin := range []hw.GPIOPin{PIMOTE_K0, PIMOTE_K1, PIMOTE_K2, PIMOTE_K3, PIMOTE_MOD_SEL, PIMOTE_MOD_EN} {
		this.gpio.WritePin(pin, hw.GPIO_LOW)
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// send information

func (this *PimoteDriver) write(reg byte, state bool) {
	if state {
		reg = reg | 8
	}
	// output to pins K0
	if (reg & 0x01) == 0 {
		this.gpio.WritePin(PIMOTE_K0, hw.GPIO_LOW)
	} else {
		this.gpio.WritePin(PIMOTE_K0, hw.GPIO_HIGH)
	}
	// output to pins K1
	if (reg & 0x02) == 0 {
		this.gpio.WritePin(PIMOTE_K1, hw.GPIO_LOW)
	} else {
		this.gpio.WritePin(PIMOTE_K1, hw.GPIO_HIGH)
	}
	// output to pins K2
	if (reg & 0x04) == 0 {
		this.gpio.WritePin(PIMOTE_K2, hw.GPIO_LOW)
	} else {
		this.gpio.WritePin(PIMOTE_K2, hw.GPIO_HIGH)
	}
	// output to pins K3
	if (reg & 0x08) == 0 {
		this.gpio.WritePin(PIMOTE_K3, hw.GPIO_LOW)
	} else {
		this.gpio.WritePin(PIMOTE_K3, hw.GPIO_HIGH)
	}

	// Let it settle, encoder requires this
	time.Sleep(100 * time.Millisecond)

	// Enable the modulator
	this.gpio.WritePin(PIMOTE_MOD_EN, hw.GPIO_HIGH)

	// Keep enabled for a period
	time.Sleep(250 * time.Millisecond)

	// Disable the modulator
	this.gpio.WritePin(PIMOTE_MOD_EN, hw.GPIO_LOW)

	// Let it settle
	time.Sleep(100 * time.Millisecond)
}

func (this *PimoteDriver) send(socket uint, state bool) error {
	/* TODO: Ultimately we should do this in a go routine, and use mutex locks */
	switch {
	case socket == 0:
		this.write(0x3, state)
		break
	case socket == 1:
		this.write(0x7, state)
		break
	case socket == 2:
		this.write(0x6, state)
		break
	case socket == 3:
		this.write(0x5, state)
		break
	case socket == 4:
		this.write(0x4, state)
		break
	default:
		return this.log.Error("Socket out of range: %v", socket)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// On and Off

func (this *PimoteDriver) On(sockets ...uint) error {
	if len(sockets) == 0 {
		return this.send(0, true)
	}
	// Write socket
	for _, socket := range sockets {
		if socket < PIMOTE_SOCKET_MIN || socket > PIMOTE_SOCKET_MAX {
			return this.log.Error("Socket out of range: %v", socket)
		}
		if err := this.send(socket, true); err != nil {
			return err
		}
	}
	// Success
	return nil
}

func (this *PimoteDriver) Off(sockets ...uint) error {
	if len(sockets) == 0 {
		return this.send(0, false)
	}
	// Write socket
	for _, socket := range sockets {
		if socket < PIMOTE_SOCKET_MIN || socket > PIMOTE_SOCKET_MAX {
			return this.log.Error("Socket out of range: %v", socket)
		}
		if err := this.send(socket, false); err != nil {
			return err
		}
	}
	// Success
	return nil
}
