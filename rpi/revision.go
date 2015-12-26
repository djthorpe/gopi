/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

// Code adapted from https://github.com/AndrewFromMelbourne/raspberry_pi_revision
// Please see LICENSE.md for more information

// Memory
const (
	RPI_MEMORY_UNKNOWN = iota
	RPI_MEMORY_256MB
	RPI_MEMORY_512MB
	RPI_MEMORY_1024MB
)

// Processor
const (
	RPI_CPU_UNKNOWN = iota
	RPI_CPU_BROADCOM_2835
	RPI_CPU_BROADCOM_2836
)

// Model
const (
	RPI_MODEL_UNKNOWN = iota
	RPI_MODEL_A
	RPI_MODEL_B
	RPI_MODEL_A_PLUS
	RPI_MODEL_B_PLUS
	RPI_MODEL_B_PI_2
	RPI_MODEL_ALPHA
	RPI_MODEL_COMPUTE_MODULE
	RPI_MODEL_ZERO
)

// Manufacturer
const (
	RPI_MANUFACTURER_UNKNOWN = iota
	RPI_MANUFACTURER_SONY
	RPI_MANUFACTURER_EGOMAN
	RPI_MANUFACTURER_EMBEST
	RPI_MANUFACTURER_QISDA
)

// PCBRevision
const (
	RPI_PCB_UNKNOWN = iota
	RPI_PCB_REV0
	RPI_PCB_REV1
	RPI_PCB_REV2
)

type Revision struct {
	Model        uint
	Processor    uint
	Memory       uint
	Manufacturer uint
	PCBRevision  uint
	WarrantyBit  bool
}

func GetRevision() *Revision {
	// Create a new revision structure
	r = new(Revision)
	// TODO
	return r
}
