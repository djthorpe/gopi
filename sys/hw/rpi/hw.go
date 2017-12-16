/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import (
	"fmt"
	"strings"

	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
   #cgo CFLAGS: -I/opt/vc/include
   #cgo LDFLAGS: -L/opt/vc/lib -lbcm_host
   #include "bcm_host.h"
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Hardware struct{}

type hardware struct {
	log      gopi.Logger
	service  int
	serial   uint64
	revision uint32
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open
func (config Hardware) Open(logger gopi.Logger) (gopi.Driver, error) {
	logger.Debug("sys.rpi.Hardware.Open{  }")

	// Initialise
	if err := bcmHostInit(); err != nil {
		return nil, err
	}

	this := new(hardware)
	this.log = logger
	this.service = GENCMD_SERVICE_NONE
	this.serial = GENCMD_SERIAL_NONE
	this.revision = GENCMD_REVISION_NONE

	// Success
	return this, nil
}

// Close
func (this *hardware) Close() error {
	this.log.Debug("sys.rpi.Hardware.Close{ }")

	// vcgencmd interface
	if this.service != GENCMD_SERVICE_NONE {
		if err := vcGencmdTerminate(); err != nil {
			bcmHostTerminate()
			return err
		}
		this.service = GENCMD_SERVICE_NONE
	}

	// host terminate
	if err := bcmHostTerminate(); err != nil {
		return err
	}

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetName returns the name of the hardware
func (this *hardware) Name() string {
	if product, err := this.GetProduct(); err != nil {
		this.log.Error("Error fetching serial number: %v", err)
		return ""
	} else {
		return fmt.Sprintf("%v (revision %v)", product.model, product.revision)
	}
}

// SerialNumber returns the serial number of the hardware, if available
func (this *hardware) SerialNumber() string {
	if serial, err := this.GetSerialNumberUint64(); err != nil {
		this.log.Error("Error fetching serial number: %v", err)
		return ""
	} else {
		return fmt.Sprintf("%X", serial)
	}
}

// Return the number of displays which can be opened
func (this *hardware) NumberOfDisplays() uint {
	return uint(DX_ID_MAX) + 1
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *hardware) String() string {
	// Cache serial and revision
	this.GetSerialNumberUint64()
	this.GetRevisionUint32()
	product, _ := this.GetProduct()
	params := []string{
		fmt.Sprintf("name=%v", this.Name()),
		fmt.Sprintf("serial=0x%X", this.serial),
		fmt.Sprintf("revision=0x%08X", this.revision),
		fmt.Sprintf("product=%v", product),
		fmt.Sprintf("displays=%v", this.NumberOfDisplays()),
		fmt.Sprintf("peripheral_addr=0x%08X", bcmHostGetPeripheralAddress()),
		fmt.Sprintf("peripheral_size=0x%08X", bcmHostGetPeripheralSize()),
	}
	return fmt.Sprintf("sys.rpi.Hardware{ %v }", strings.Join(params, " "))
}

////////////////////////////////////////////////////////////////////////////////
// BCMHOST

func bcmHostInit() error {
	C.bcm_host_init()
	return nil
}

func bcmHostTerminate() error {
	C.bcm_host_deinit()
	return nil
}

func bcmHostGetPeripheralAddress() uint32 {
	return uint32(C.bcm_host_get_peripheral_address())
}

func bcmHostGetPeripheralSize() uint32 {
	return uint32(C.bcm_host_get_peripheral_size())
}

func bcmHostGetSDRAMAddress() uint32 {
	return uint32(C.bcm_host_get_sdram_address())
}
