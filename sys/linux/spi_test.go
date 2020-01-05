// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package linux_test

import (
	"os"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/sys/linux"
)

func Test_SPI_000(t *testing.T) {
	for bus := uint(0); bus <= 2; bus++ {
		for slave := uint(0); slave <= 2; slave++ {
			dev := linux.SPIDevice(bus, slave)
			if _, err := os.Stat(dev); os.IsNotExist(err) {
				// Ignore
			} else {
				t.Log("DEV", dev)
			}
		}
	}
}

func Test_SPI_001(t *testing.T) {
	for bus := uint(0); bus <= 2; bus++ {
		for slave := uint(0); slave <= 2; slave++ {
			dev := linux.SPIDevice(bus, slave)
			if _, err := os.Stat(dev); os.IsNotExist(err) {
				// Ignore
			} else if spi, err := linux.SPIOpenDevice(bus, slave); err != nil {
				t.Error(err)
			} else if mode, err := linux.SPIMode(spi.Fd()); err != nil {
				t.Error(err)
			} else if speedhz, err := linux.SPISpeedHz(spi.Fd()); err != nil {
				t.Error(err)
			} else if bits, err := linux.SPIBitsPerWord(spi.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Log("DEV", dev)
				t.Log("  MODE", mode)
				t.Log("  SPEED", speedhz)
				t.Log("  BITS/WORD", bits)
			}
		}
	}
}

func Test_SPI_002(t *testing.T) {
	for bus := uint(0); bus <= 2; bus++ {
		for slave := uint(0); slave <= 2; slave++ {
			dev := linux.SPIDevice(bus, slave)
			if _, err := os.Stat(dev); os.IsNotExist(err) {
				// Ignore
			} else if spi, err := linux.SPIOpenDevice(bus, slave); err != nil {
				t.Error(err)
			} else if mode, err := linux.SPIMode(spi.Fd()); err != nil {
				t.Error(err)
			} else if speedhz, err := linux.SPISpeedHz(spi.Fd()); err != nil {
				t.Error(err)
			} else if bits, err := linux.SPIBitsPerWord(spi.Fd()); err != nil {
				t.Error(err)
			} else if err := linux.SPISetMode(spi.Fd(), mode); err != nil {
				t.Error(err)
			} else if err := linux.SPISetSpeedHz(spi.Fd(), speedhz); err != nil {
				t.Error(err)
			} else if err := linux.SPISetBitsPerWord(spi.Fd(), bits); err != nil {
				t.Error(err)
			} else {
				t.Log("DEV", dev)
				t.Log("  MODE", mode)
				t.Log("  SPEED", speedhz)
				t.Log("  BITS/WORD", bits)
				spi.Close()
			}
		}
	}
}
