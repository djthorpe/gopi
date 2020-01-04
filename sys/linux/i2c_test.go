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

func Test_I2C_000(t *testing.T) {
	for i := uint(0); i <= 2; i++ {
		dev := linux.I2CDevice(i)
		if _, err := os.Stat(dev); os.IsNotExist(err) {
			// Ignore
		} else {
			t.Log("DEV", dev)
		}
	}
}

func Test_I2C_001(t *testing.T) {
	for i := uint(0); i <= 2; i++ {
		dev := linux.I2CDevice(i)
		if _, err := os.Stat(dev); os.IsNotExist(err) {
			// Ignore
		} else if dev, err := linux.I2COpenDevice(i); err != nil {
			t.Error(err)
		} else {
			defer dev.Close()
			if funcs, err := linux.I2CFunctions(dev.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Log("functions=", funcs)
			}
		}
	}
}

func Test_I2C_002(t *testing.T) {
	for i := uint(0); i <= 2; i++ {
		dev := linux.I2CDevice(i)
		if _, err := os.Stat(dev); os.IsNotExist(err) {
			// Ignore
		} else if dev, err := linux.I2COpenDevice(i); err != nil {
			t.Error(err)
		} else {
			defer dev.Close()
			if funcs, err := linux.I2CFunctions(dev.Fd()); err != nil {
				t.Error(err)
			} else {
				for slave := byte(0); slave <= 0x7F; slave++ {
					if detected, err := linux.I2CDetectSlave(dev.Fd(), slave, funcs); err != nil {
						t.Error(err)
					} else {
						t.Logf("%02X => %v", slave, detected)
					}
				}
			}
		}
	}
}
