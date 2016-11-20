/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
)

import (
	gopi "github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Capability Tuple
type Tuple struct {
	Key  gopi.Capability
	Func tupleCallback
}

// Capability callback returns a string (for the moment)
type tupleCallback func(key gopi.Capability) string

////////////////////////////////////////////////////////////////////////////////
// Hardware Capabilities

// Return all capabilities
func (this *Device) makeCapabilities() []gopi.Tuple {
	tuples := make([]gopi.Tuple, 8)

	tuples[0] = &Tuple{Key: gopi.CAP_HW_SERIAL, Func: this.getCapSerial}
	tuples[1] = &Tuple{Key: gopi.CAP_HW_PLATFORM, Func: this.getCapPlatform}
	tuples[2] = &Tuple{Key: gopi.CAP_HW_MODEL, Func: this.getCapModel}
	tuples[3] = &Tuple{Key: gopi.CAP_HW_REVISION, Func: this.getCapRevision}
	tuples[4] = &Tuple{Key: gopi.CAP_HW_PCB, Func: this.getCapPCB}
	tuples[5] = &Tuple{Key: gopi.CAP_HW_WARRANTY, Func: this.getCapWarranty}
	tuples[6] = &Tuple{Key: gopi.CAP_HW_PROCESSOR_NAME, Func: this.getCapProcessor}
	tuples[7] = &Tuple{Key: gopi.CAP_HW_PROCESSOR_TEMP, Func: this.getCapCoreTemperature}

	return tuples
}

func (tuple *Tuple) Capability() gopi.Capability {
	return tuple.Key
}

func (tuple *Tuple) String() string {
	return fmt.Sprint(tuple.Func(tuple.Key))
}

func (this *Device) getCapPlatform(key gopi.Capability) string {
	return "RPI"
}

func (this *Device) getCapSerial(key gopi.Capability) string {
	serial, _ := this.GetSerialNumber()
	return fmt.Sprintf("%016X", serial)
}

func (this *Device) getCapModel(key gopi.Capability) string {
	model, _, _ := this.GetModel()
	return fmt.Sprintf("%s", model)
}

func (this *Device) getCapPCB(key gopi.Capability) string {
	_, pcb, _ := this.GetModel()
	return fmt.Sprintf("%s", pcb)
}

func (this *Device) getCapRevision(key gopi.Capability) string {
	revision, _ := this.GetRevision()
	return fmt.Sprintf("%s", revision)
}

func (this *Device) getCapProcessor(key gopi.Capability) string {
	processor, _ := this.GetProcessor()
	return fmt.Sprintf("%s", processor)
}

func (this *Device) getCapWarranty(key gopi.Capability) string {
	warranty, _ := this.GetWarrantyBit()
	return fmt.Sprintf("%s", warranty)
}

func (this *Device) getCapCoreTemperature(key gopi.Capability) string {
	temp, _ := this.GetCoreTemperatureCelcius()
	return fmt.Sprintf("%s", temp)
}
