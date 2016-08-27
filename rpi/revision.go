/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

////////////////////////////////////////////////////////////////////////////////

type Product uint

////////////////////////////////////////////////////////////////////////////////

const (
	RPI_REVISION_WARRANTY_MASK     = 0x03000000
	RPI_REVISION_ENCODING_MASK     = 0x00800000
	RPI_REVISION_PCB_MASK          = 0x0000000F
	RPI_REVISION_PRODUCT_MASK      = 0x00000FF0
	RPI_REVISION_PROCESSOR_MASK    = 0x0000F000
	RPI_REVISION_MANUFACTURER_MASK = 0x000F0000
	RPI_REVISION_MEMORY_MASK       = 0x00700000
)

const (
	RPI_MODEL_UNKNOWN Product = iota
	RPI_MODEL_A
	RPI_MODEL_B
	RPI_MODEL_A_PLUS
	RPI_MODEL_B_PLUS
	RPI_MODEL_B_PI_2
	RPI_MODEL_B_PI_3
	RPI_MODEL_ALPHA
	RPI_MODEL_COMPUTE_MODULE
	RPI_MODEL_ZERO
)

////////////////////////////////////////////////////////////////////////////////

var productmap1 = map[uint32]Product{
	0x00: RPI_MODEL_UNKNOWN,
	0x01: RPI_MODEL_UNKNOWN,
	0x02: RPI_MODEL_B,
	0x03: RPI_MODEL_B,
	0x04: RPI_MODEL_B,
	0x05: RPI_MODEL_B,
	0x06: RPI_MODEL_B,
	0x07: RPI_MODEL_A,
	0x08: RPI_MODEL_A,
	0x09: RPI_MODEL_A,
	0x0A: RPI_MODEL_UNKNOWN,
	0x0B: RPI_MODEL_UNKNOWN,
	0x0C: RPI_MODEL_UNKNOWN,
	0x0D: RPI_MODEL_B,
	0x0E: RPI_MODEL_B,
	0x0F: RPI_MODEL_B,
	0x10: RPI_MODEL_B_PLUS,
	0x11: RPI_COMPUTE_MODULE,
	0x12: RPI_MODEL_A_PLUS,
	0x13: RPI_MODEL_B_PLUS,
	0x14: RPI_COMPUTE_MODULE,
	0x15: RPI_MODEL_A_PLUS,
}

var productmap2 = map[uint32]Product{
	0 << 4: RPI_MODEL_A,
	1 << 4: RPI_MODEL_B,
	2 << 4: RPI_MODEL_A_PLUS,
	3 << 4: RPI_MODEL_B_PLUS,
	4 << 4: RPI_MODEL_B_PI_2,
	5 << 4: RPI_MODEL_ALPHA,
	6 << 4: RPI_COMPUTE_MODULE,
	7 << 4: RPI_MODEL_UNKNOWN,
	8 << 4: RPI_MODEL_B_PI_3,
	9 << 4: RPI_MODEL_ZERO,
}

var pcbmap1 = map[uint32]PCBRevision{
	0x02: 1,
	0x03: 1,
	0x04: 2,
	0x05: 2,
	0x06: 2,
	0x07: 2,
	0x08: 2,
	0x09: 2,
	0x0D: 2,
	0x0E: 2,
	0x0F: 2,
	0x10: 1,
	0x11: 1,
	0x12: 1,
	0x13: 1,
	0x14: 1,
	0x15: 1,
}

var memorymap1 = map[uint32]MemoryMB{
	0x02: RPI_MEMORY_256MB,
	0x03: RPI_MEMORY_256MB,
	0x04: RPI_MEMORY_256MB,
	0x05: RPI_MEMORY_256MB,
	0x06: RPI_MEMORY_256MB,
	0x07: RPI_MEMORY_256MB,
	0x08: RPI_MEMORY_256MB,
	0x09: RPI_MEMORY_256MB,
	0x0D: RPI_MEMORY_512MB,
	0x0E: RPI_MEMORY_512MB,
	0x0F: RPI_MEMORY_512MB,
	0x10: RPI_MEMORY_512MB,
	0x11: RPI_MEMORY_512MB,
	0x12: RPI_MEMORY_256MB,
	0x13: RPI_MEMORY_512MB,
	0x14: RPI_MEMORY_512MB,
	0x15: RPI_MEMORY_256MB,
}

var processormap1 = map[uint32]Processor{
	0x02: RPI_PROCESSOR_BCM2835,
	0x03: RPI_PROCESSOR_BCM2835,
	0x04: RPI_PROCESSOR_BCM2835,
	0x05: RPI_PROCESSOR_BCM2835,
	0x06: RPI_PROCESSOR_BCM2835,
	0x07: RPI_PROCESSOR_BCM2835,
	0x08: RPI_PROCESSOR_BCM2835,
	0x09: RPI_PROCESSOR_BCM2835,
	0x0D: RPI_PROCESSOR_BCM2835,
	0x0E: RPI_PROCESSOR_BCM2835,
	0x0F: RPI_PROCESSOR_BCM2835,
	0x10: RPI_PROCESSOR_BCM2835,
	0x11: RPI_PROCESSOR_BCM2835,
	0x12: RPI_PROCESSOR_BCM2835,
	0x13: RPI_PROCESSOR_BCM2835,
	0x14: RPI_PROCESSOR_BCM2835,
	0x15: RPI_PROCESSOR_BCM2835,
}

var manufacturermap1 = map[uint32]Manufacturer{
	0x04: RPI_MANUFACTURER_SONY,
	0x05: RPI_MANUFACTURER_QISDA,
	0x06: RPI_MANUFACTURER_EGOMAN,
	0x07: RPI_MANUFACTURER_EGOMAN,
	0x08: RPI_MANUFACTURER_SONY,
	0x09: RPI_MANUFACTURER_QISDA,
	0x0D: RPI_MANUFACTURER_EGOMAN,
	0x0E: RPI_MANUFACTURER_SONY,
	0x0F: RPI_MANUFACTURER_QISDA,
	0x10: RPI_MANUFACTURER_SONY,
	0x11: RPI_MANUFACTURER_SONY,
	0x12: RPI_MANUFACTURER_SONY,
	0x13: RPI_MANUFACTURER_EMBEST,
	0x14: RPI_MANUFACTURER_SONY,
	0x15: RPI_MANUFACTURER_SONY,
}

var manufacturermap2 = map[uint32]Manufacturer{
	0 << 16: RPI_MANUFACTURER_SONY,
	1 << 16: RPI_MANUFACTURER_EGOMAN,
	2 << 16: RPI_MANUFACTURER_EMBEST,
	3 << 16: RPI_MANUFACTURER_UNKNOWN,
	4 << 16: RPI_MANUFACTURER_EMBEST,
}

var processormap2 = map[uint32]Processor{
	0 << 12: RPI_PROCESSOR_BCM2835,
	1 << 12: RPI_PROCESSOR_BCM2836,
}

var memorymap2 = map[uint32]MemoryMB{
	0 << 20: RPI_MEMORY_256MB,
	1 << 20: RPI_MEMORY_512MB,
	2 << 20: RPI_MEMORY_1024MB,
}

var productstringmap = map[Product]string{
	RPI_MODEL_UNKNOWN:  "unknown",
	RPI_MODEL_A:        "A",
	RPI_MODEL_B:        "B",
	RPI_MODEL_A_PLUS:   "A+",
	RPI_MODEL_B_PLUS:   "B+",
	RPI_MODEL_B_PI_2:   "B2",
	RPI_MODEL_ALPHA:    "alpha",
	RPI_COMPUTE_MODULE: "compute",
	RPI_MODEL_ZERO:     "zero",
}

var processorstringmap = map[Processor]string{
	RPI_PROCESSOR_UNKNOWN: "unknown",
	RPI_PROCESSOR_BCM2835: "BCM2835",
	RPI_PROCESSOR_BCM2836: "BCM2836",
}

var peripheralbasemap = map[Processor]PeripheralBase{
	RPI_PROCESSOR_UNKNOWN: RPI_PERIPHERAL_BASE_UNKNOWN,
	RPI_PROCESSOR_BCM2835: RPI_PERIPHERAL_BASE_BCM2835,
	RPI_PROCESSOR_BCM2836: RPI_PERIPHERAL_BASE_BCM2836,
}

var manufacturerstringmap = map[Manufacturer]string{
	RPI_MANUFACTURER_UNKNOWN: "unknown",
	RPI_MANUFACTURER_SONY:    "Sony",
	RPI_MANUFACTURER_EGOMAN:  "Egoman",
	RPI_MANUFACTURER_QISDA:   "Qisda",
	RPI_MANUFACTURER_EMBEST:  "Embest",
}

////////////////////////////////////////////////////////////////////////////////

// Function to return the product name, revision and warranty bit
func (this *State) GetModel() (Product, PCBRevision, Warranty, error) {
	var warranty Warranty
	var product Product
	var pcbrevision PCBRevision

	// set Revision
	revision, err := this.GetRevision()
	if err != nil {
		return Product(0),PCBRevision(0),Warranty(false),err
	}

	// set warranty bit
	w := uint32(RPI_REVISION_WARRANTY_MASK)
	if (revision & w) != 0 {
		warranty = Warranty(true)
	} else {
		warranty = Warranty(false)
	}

	// pare down revision and decode differently depending on the format
	revision = revision & ^w
	if (revision & RPI_REVISION_ENCODING_MASK) != 0 {
		// Raspberry Pi 2 style revision coding
		product = Product(productmap2[revision & RPI_REVISION_PRODUCT_MASK])
		pcbrevision = PCBRevision(revision & RPI_REVISION_PCB_MASK)
	} else {
		// Raspberry Pi 1 style revision coding
		product = Product(productmap1[revision])
		pcbrevision = PCBRevision(pcbmap1[revision])
	}

	return product, pcbrevision, warranty, nil
}
