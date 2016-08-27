/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

////////////////////////////////////////////////////////////////////////////////

type Product uint32

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

var (
	productmap1 = map[uint32]Product{
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

	productmap2 = map[uint32]Product{
		0 << 4: RPI_MODEL_A,
		1 << 4: RPI_MODEL_B,
		2 << 4: RPI_MODEL_A_PLUS,
		3 << 4: RPI_MODEL_B_PLUS,
		4 << 4: RPI_MODEL_B_PI_2,
		5 << 4: RPI_MODEL_ALPHA,
		6 << 4: RPI_MODEL_COMPUTE_MODULE,
		7 << 4: RPI_MODEL_UNKNOWN,
		8 << 4: RPI_MODEL_B_PI_3,
		9 << 4: RPI_MODEL_ZERO,
	}

	productstringmap = map[Product]string{
		RPI_MODEL_A:        "A",
		RPI_MODEL_B:        "B",
		RPI_MODEL_A_PLUS:   "A+",
		RPI_MODEL_B_PLUS:   "B+",
		RPI_MODEL_B_PI_2:   "B2",
		RPI_MODEL_B_PI_3:   "B3",
		RPI_MODEL_ALPHA:    "alpha",
		RPI_MODEL_COMPUTE_MODULE: "compute",
		RPI_MODEL_ZERO:     "zero",
		RPI_MODEL_UNKNOWN:  "unknown",
	}
)

////////////////////////////////////////////////////////////////////////////////

// Returns the value of the warranty bit
func (this *RaspberryPi) WarrantyBit() (bool, error) {
	revision, err := this.GetRevision()
	if err != nil {
		return false, err
	}

	// get warranty bit
	w := uint32(RPI_REVISION_WARRANTY_MASK)
	return ((revision & w) != 0),nil
}

func (this *RaspberryPi) ProductName() (string, error) {

	// Get revision
	revision, err := this.GetRevision()
	if err != nil {
		return "", err
	}
	revision = revision & ^uint32(RPI_REVISION_WARRANTY_MASK)

	// pare down revision and decode differently depending on the format
	var product Product
	var ok bool
	if (revision & RPI_REVISION_ENCODING_MASK) != 0 {
		// Raspberry Pi 2 style revision coding
		product, ok = productmap2[revision & RPI_REVISION_PRODUCT_MASK]
	} else {
		// Raspberry Pi 1 style revision coding
		product, ok = productmap1[revision]
	}

	if ok == false {
		product = RPI_MODEL_UNKNOWN
	}

	name, ok := productstringmap[product]
	if ok == false {
		return productstringmap[RPI_MODEL_UNKNOWN],nil
	}
	return name, nil
}
