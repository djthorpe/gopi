/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

    Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Model       uint32
	Processor   uint32
	PCBRevision uint32
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RPI_REVISION_WARRANTY_MASK     uint32 = 0x03000000
	RPI_REVISION_ENCODING_MASK     uint32 = 0x00800000
	RPI_REVISION_PCB_MASK          uint32 = 0x0000000F
	RPI_REVISION_PRODUCT_MASK      uint32 = 0x00000FF0
	RPI_REVISION_PROCESSOR_MASK    uint32 = 0x0000F000
	RPI_REVISION_MANUFACTURER_MASK uint32 = 0x000F0000
	RPI_REVISION_MEMORY_MASK       uint32 = 0x00700000
)

const (
	RPI_MODEL_A              Model = (0 << 4)
	RPI_MODEL_B              Model = (1 << 4)
	RPI_MODEL_A_PLUS         Model = (2 << 4)
	RPI_MODEL_B_PLUS         Model = (3 << 4)
	RPI_MODEL_B_PI_2         Model = (4 << 4)
	RPI_MODEL_ALPHA          Model = (5 << 4)
	RPI_MODEL_COMPUTE_MODULE Model = (6 << 4)
	RPI_MODEL_UNKNOWN        Model = (7 << 4)
	RPI_MODEL_B_PI_3         Model = (8 << 4)
	RPI_MODEL_ZERO           Model = (9 << 4)
)

const (
	RPI_PROCESSOR_UNKNOWN Processor = 0xFFFFFFFF
	RPI_PROCESSOR_BCM2835 Processor = (0 << 12)
	RPI_PROCESSOR_BCM2836 Processor = (1 << 12)
	RPI_PROCESSOR_BCM2837 Processor = (2 << 12)
)

const (
	RPI_PCB_UNKNOWN PCBRevision = 0
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	productmap1 = map[uint32]Model{
		0x02: RPI_MODEL_B,
		0x03: RPI_MODEL_B,
		0x04: RPI_MODEL_B,
		0x05: RPI_MODEL_B,
		0x06: RPI_MODEL_B,
		0x07: RPI_MODEL_A,
		0x08: RPI_MODEL_A,
		0x09: RPI_MODEL_A,
		0x0D: RPI_MODEL_B,
		0x0E: RPI_MODEL_B,
		0x0F: RPI_MODEL_B,
		0x10: RPI_MODEL_B_PLUS,
		0x11: RPI_MODEL_COMPUTE_MODULE,
		0x12: RPI_MODEL_A_PLUS,
		0x13: RPI_MODEL_B_PLUS,
		0x14: RPI_MODEL_COMPUTE_MODULE,
		0x15: RPI_MODEL_A_PLUS,
	}
	pcbmap1 = map[uint32]PCBRevision{
		0x02: PCBRevision(1),
		0x03: PCBRevision(1),
		0x04: PCBRevision(2),
		0x05: PCBRevision(2),
		0x06: PCBRevision(2),
		0x07: PCBRevision(2),
		0x08: PCBRevision(2),
		0x09: PCBRevision(2),
		0x0D: PCBRevision(2),
		0x0E: PCBRevision(2),
		0x0F: PCBRevision(2),
		0x10: PCBRevision(1),
		0x11: PCBRevision(1),
		0x12: PCBRevision(1),
		0x13: PCBRevision(1),
		0x14: PCBRevision(1),
		0x15: PCBRevision(1),
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Returns the value of the warranty bit
func (this *Device) GetWarrantyBit() (bool, error) {
	revision, err := this.GetRevision()
	if err != nil {
		return false, err
	}

	// get warranty bit
	w := uint32(RPI_REVISION_WARRANTY_MASK)
	return ((revision & w) != 0), nil
}

// Return product model and PCB revision information
func (this *Device) GetModel() (Model, PCBRevision, error) {
	revision, err := this.GetRevision()
	if err != nil {
		return RPI_MODEL_UNKNOWN, RPI_PCB_UNKNOWN, err
	}

	// Decode differently depending on the format
	var model Model
	var pcb PCBRevision
	var ok bool

	revision = revision & ^RPI_REVISION_WARRANTY_MASK
	if (revision & RPI_REVISION_ENCODING_MASK) != 0 {
		// Raspberry Pi 2 style revision coding
		model = Model(revision & RPI_REVISION_PRODUCT_MASK)
		pcb = PCBRevision(revision & RPI_REVISION_PCB_MASK)
	} else {
		// Raspberry Pi 1 style revision coding
		model, ok = productmap1[revision]
		if ok == false {
			return RPI_MODEL_UNKNOWN, RPI_PCB_UNKNOWN, nil
		}
		pcb, ok = pcbmap1[revision]
		if ok == false {
			return RPI_MODEL_UNKNOWN, RPI_PCB_UNKNOWN, nil
		}
	}

	return model, pcb, nil
}

func (this *Device) GetProcessor() (Processor, error) {
	revision, err := this.GetRevision()
	if err != nil {
		return RPI_PROCESSOR_UNKNOWN, err
	}

	// Decode differently depending on the format
	revision = revision & ^RPI_REVISION_WARRANTY_MASK
	if (revision & RPI_REVISION_ENCODING_MASK) != 0 {
		// Raspberry Pi 2 style revision coding
		return Processor(revision & RPI_REVISION_PROCESSOR_MASK), nil
	} else {
		// Raspberry Pi 1 style revision coding
		return RPI_PROCESSOR_BCM2835, nil
	}
}
