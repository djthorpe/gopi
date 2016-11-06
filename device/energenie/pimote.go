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
package pimote /* import "github.com/djthorpe/gopi/device/entergenie/pimote" */

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
	hw   "github.com/djthorpe/gopi/hw"
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
	log *util.LoggerDevice
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	PIMOTE_K0 = hw.GPIOPin(17)
	PIMOTE_K1 = hw.GPIOPin(22)
	PIMOTE_K2 = hw.GPIOPin(23)
	PIMOTE_K3 = hw.GPIOPin(27)
	PIMOTE_MOD_SEL = hw.GPIOPin(24)
	PIMOTE_CE_EN = hw.GPIOPin(25)
)

////////////////////////////////////////////////////////////////////////////////

func (config Pimote) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug2("<energenie.pimote>Open")

	this := new(PimoteDriver)
	this.gpio = config.GPIO
	this.log = log

	return nil
}

func (this *PimoteDriver) Close() error {
	log.Debug2("<energenie.pimote>Close")
	// do nothing
}

