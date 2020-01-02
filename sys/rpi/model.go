// +build rpi
// +build !darwin

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package rpi

import "fmt"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	Model     uint
	Processor uint
	Revision  uint
)

type ProductInfo struct {
	Model       Model
	Processor   Processor
	Revision    Revision
	WarrantyBit bool
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RPI_WARRANTY_MASK     uint32 = 0x03000000
	RPI_ENCODING_MASK     uint32 = 0x00800000
	RPI_REVISION_MASK     uint32 = 0x0000000F
	RPI_MODEL_MASK        uint32 = 0x00000FF0
	RPI_PROCESSOR_MASK    uint32 = 0x0000F000
	RPI_MANUFACTURER_MASK uint32 = 0x000F0000
	RPI_MEMORY_MASK       uint32 = 0x00700000
)

const (
	RPI_MODEL_A                    Model = (0x00 << 4)
	RPI_MODEL_B                    Model = (0x01 << 4)
	RPI_MODEL_A_PLUS               Model = (0x02 << 4)
	RPI_MODEL_B_PLUS               Model = (0x03 << 4)
	RPI_MODEL_B_2                  Model = (0x04 << 4)
	RPI_MODEL_ALPHA                Model = (0x05 << 4)
	RPI_MODEL_COMPUTE_MODULE       Model = (0x06 << 4)
	RPI_MODEL_B_3                  Model = (0x08 << 4)
	RPI_MODEL_ZERO                 Model = (0x09 << 4)
	RPI_MODEL_COMPUTE_MODULE_3     Model = (0x0A << 4)
	RPI_MODEL_ZERO_W               Model = (0x0C << 4)
	RPI_MODEL_B_3PLUS              Model = (0x0D << 4)
	RPI_MODEL_A_3PLUS              Model = (0x0E << 4)
	RPI_MODEL_UNKNOWN              Model = (0x0F << 4)
	RPI_MODEL_COMPUTE_MODULE_3PLUS Model = (0x10 << 4)
	RPI_MODEL_B_4                  Model = (0x11 << 4)
)

const (
	RPI_PROCESSOR_UNKNOWN Processor = 0xFFFFFFFF
	RPI_PROCESSOR_BCM2835 Processor = (0 << 12)
	RPI_PROCESSOR_BCM2836 Processor = (1 << 12)
	RPI_PROCESSOR_BCM2837 Processor = (2 << 12)
	RPI_PROCESSOR_BCM2838 Processor = (3 << 12)
)

const (
	RPI_REVISION_UNKNOWN Revision = 0
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
	pcbmap1 = map[uint32]Revision{
		0x02: Revision(1),
		0x03: Revision(1),
		0x04: Revision(2),
		0x05: Revision(2),
		0x06: Revision(2),
		0x07: Revision(2),
		0x08: Revision(2),
		0x09: Revision(2),
		0x0D: Revision(2),
		0x0E: Revision(2),
		0x0F: Revision(2),
		0x10: Revision(1),
		0x11: Revision(1),
		0x12: Revision(1),
		0x13: Revision(1),
		0x14: Revision(1),
		0x15: Revision(1),
	}
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewProductInfo(product uint32) ProductInfo {
	if product == 0 {
		return ProductInfo{}
	} else {
		model, revision := modelRevision(product)
		return ProductInfo{
			model,
			processor(product),
			revision,
			warrantyBit(product),
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// warrantyBit returns the value of the warranty bit
func warrantyBit(revision uint32) bool {
	return ((revision & uint32(RPI_WARRANTY_MASK)) != 0)
}

// modelRevision returns model and revision information
func modelRevision(revision uint32) (Model, Revision) {
	revision = revision & ^RPI_WARRANTY_MASK
	if (revision & RPI_ENCODING_MASK) != 0 {
		// Raspberry Pi 2 style revision coding
		return Model(revision & RPI_MODEL_MASK), Revision(revision & RPI_REVISION_MASK)
	} else {
		// Raspberry Pi 1 style revision coding
		if model, ok := productmap1[revision]; !ok {
			return RPI_MODEL_UNKNOWN, RPI_REVISION_UNKNOWN
		} else if pcb, ok := pcbmap1[revision]; !ok {
			return RPI_MODEL_UNKNOWN, RPI_REVISION_UNKNOWN
		} else {
			return model, pcb
		}
	}
}

// processor returns processor information
func processor(revision uint32) Processor {
	revision = revision & ^RPI_WARRANTY_MASK
	if (revision & RPI_ENCODING_MASK) != 0 {
		// Raspberry Pi 2 style revision coding
		return Processor(revision & RPI_PROCESSOR_MASK)
	} else {
		// Raspberry Pi 1 style revision coding
		return RPI_PROCESSOR_BCM2835
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p ProductInfo) String() string {
	return fmt.Sprintf("rpi.ProductInfo{ model=%v processor=%v revision=%v warrantyBit=%v }", p.Model, p.Processor, p.Revision, p.WarrantyBit)
}

func (m Model) String() string {
	switch m {
	case RPI_MODEL_A:
		return "RPI_MODEL_A"
	case RPI_MODEL_B:
		return "RPI_MODEL_B"
	case RPI_MODEL_A_PLUS:
		return "RPI_MODEL_A_PLUS"
	case RPI_MODEL_B_PLUS:
		return "RPI_MODEL_B_PLUS"
	case RPI_MODEL_B_2:
		return "RPI_MODEL_B_2"
	case RPI_MODEL_B_3:
		return "RPI_MODEL_B_3"
	case RPI_MODEL_ALPHA:
		return "RPI_MODEL_ALPHA"
	case RPI_MODEL_COMPUTE_MODULE:
		return "RPI_MODEL_COMPUTE_MODULE"
	case RPI_MODEL_COMPUTE_MODULE_3:
		return "RPI_MODEL_COMPUTE_MODULE_3"
	case RPI_MODEL_ZERO:
		return "RPI_MODEL_ZERO"
	case RPI_MODEL_ZERO_W:
		return "RPI_MODEL_ZERO_W"
	case RPI_MODEL_B_3PLUS:
		return "RPI_MODEL_B_3PLUS"
	case RPI_MODEL_A_3PLUS:
		return "RPI_MODEL_A_3PLUS"
	case RPI_MODEL_COMPUTE_MODULE_3PLUS:
		return "RPI_MODEL_COMPUTE_MODULE_3PLUS"
	case RPI_MODEL_B_4:
		return "RPI_MODEL_B_4"
	default:
		return fmt.Sprintf("[?? Unknown Model value 0x%02X]", uint32(m))
	}
}

func (p Processor) String() string {
	switch p {
	case RPI_PROCESSOR_BCM2835:
		return "RPI_PROCESSOR_BCM2835"
	case RPI_PROCESSOR_BCM2836:
		return "RPI_PROCESSOR_BCM2836"
	case RPI_PROCESSOR_BCM2837:
		return "RPI_PROCESSOR_BCM2837"
	case RPI_PROCESSOR_BCM2838:
		return "RPI_PROCESSOR_BCM2838"
	default:
		return "[?? Unknown Processor value]"
	}
}

func (p Revision) String() string {
	if p == RPI_REVISION_UNKNOWN {
		return "[?? Unknown Revision value]"
	}
	return fmt.Sprintf("RPI_REVISION_V%v", uint(p))
}
