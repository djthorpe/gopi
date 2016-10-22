/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

import (
	_ "fmt"
)

////////////////////////////////////////////////////////////////////////////////

type Product uint32
type Processor uint32

////////////////////////////////////////////////////////////////////////////////

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

const (
	RPI_PROCESSOR_UNKNOWN Processor = iota
	RPI_PROCESSOR_BCM2835
	RPI_PROCESSOR_BCM2836
	RPI_PROCESSOR_BCM2837
)

const (
	RPI_BROADCOM_2835_PERIPHERAL_BASE uint32 = 0x20000000
	RPI_BROADCOM_2836_PERIPHERAL_BASE uint32 = 0x3F000000
	RPI_BROADCOM_2837_PERIPHERAL_BASE uint32 = 0x3F000000
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
		RPI_MODEL_A:              "A",
		RPI_MODEL_B:              "B",
		RPI_MODEL_A_PLUS:         "A+",
		RPI_MODEL_B_PLUS:         "B+",
		RPI_MODEL_B_PI_2:         "B2",
		RPI_MODEL_B_PI_3:         "B3",
		RPI_MODEL_ALPHA:          "alpha",
		RPI_MODEL_COMPUTE_MODULE: "compute",
		RPI_MODEL_ZERO:           "zero",
		RPI_MODEL_UNKNOWN:        "unknown",
	}

	processormap2 = map[uint32]Processor{
		0 << 12: RPI_PROCESSOR_BCM2835,
		1 << 12: RPI_PROCESSOR_BCM2836,
		2 << 12: RPI_PROCESSOR_BCM2837,
	}

	processorstringmap = map[Processor]string{
		RPI_PROCESSOR_UNKNOWN: "unknown",
		RPI_PROCESSOR_BCM2835: "BCM2835",
		RPI_PROCESSOR_BCM2836: "BCM2836",
		RPI_PROCESSOR_BCM2837: "BCM2837",
	}

	peripheralbasemap = map[Processor]uint32{
		RPI_PROCESSOR_UNKNOWN: 0,
		RPI_PROCESSOR_BCM2835: RPI_BROADCOM_2835_PERIPHERAL_BASE,
		RPI_PROCESSOR_BCM2836: RPI_BROADCOM_2836_PERIPHERAL_BASE,
		RPI_PROCESSOR_BCM2837: RPI_BROADCOM_2837_PERIPHERAL_BASE,
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
	return ((revision & w) != 0), nil
}

func (this *RaspberryPi) Product() (Product, error) {

	// Get revision
	revision, err := this.GetRevision()
	if err != nil {
		return RPI_MODEL_UNKNOWN, err
	}

	// Decode differently depending on the format
	var product Product
	var ok bool

	revision = revision & ^RPI_REVISION_WARRANTY_MASK
	if (revision & RPI_REVISION_ENCODING_MASK) != 0 {
		// Raspberry Pi 2 style revision coding
		product, ok = productmap2[revision&RPI_REVISION_PRODUCT_MASK]
	} else {
		// Raspberry Pi 1 style revision coding
		product, ok = productmap1[revision]
	}

	if ok == false {
		return RPI_MODEL_UNKNOWN, nil
	}

	return product, nil
}

func (this *RaspberryPi) Processor() (Processor, error) {

	// Get revision
	revision, err := this.GetRevision()
	if err != nil {
		return RPI_PROCESSOR_UNKNOWN, err
	}

	// Decode differently depending on the format
	var processor Processor
	var ok bool

	revision = revision & ^RPI_REVISION_WARRANTY_MASK
	if (revision & RPI_REVISION_ENCODING_MASK) != 0 {
		// Raspberry Pi 2 style revision coding
		processor, ok = processormap2[revision&RPI_REVISION_PROCESSOR_MASK]
	} else {
		// Raspberry Pi 1 style revision coding
		processor = RPI_PROCESSOR_BCM2835
		ok = true
	}

	if ok == false {
		return RPI_PROCESSOR_UNKNOWN, nil
	}

	return processor, nil
}

func (this *RaspberryPi) ProductName() (string, error) {
	product, err := this.Product()
	if err != nil {
		return "", err
	}
	name, ok := productstringmap[product]
	if ok == false {
		return productstringmap[RPI_MODEL_UNKNOWN], nil
	}
	return name, nil
}

func (this *RaspberryPi) ProcessorName() (string, error) {
	processor, err := this.Processor()
	if err != nil {
		return "", err
	}
	name, ok := processorstringmap[processor]
	if ok == false {
		return processorstringmap[RPI_PROCESSOR_UNKNOWN], nil
	}
	return name, nil
}

func (this *RaspberryPi) PeripheralBase() (uint32, error) {
	processor, err := this.Processor()
	if err != nil {
		return 0, err
	}
	return peripheralbasemap[processor], nil
}
