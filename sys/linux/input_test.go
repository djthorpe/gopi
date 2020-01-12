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

func Test_Input_000(t *testing.T) {
	t.Log("Test_Input_000")
}
func Test_Input_001(t *testing.T) {
	for bus := uint(0); bus < 10; bus++ {
		if _, err := os.Stat(linux.EVDevice(bus)); os.IsNotExist(err) {
			// Ignore
		} else {
			t.Log(bus, "=>", linux.EVDevice(bus))
		}
	}
}

func Test_Input_002(t *testing.T) {
	for bus := uint(0); bus < 10; bus++ {
		if _, err := os.Stat(linux.EVDevice(bus)); os.IsNotExist(err) {
			// Ignore
		} else if dev, err := linux.EVOpenDevice(bus); err != nil {
			t.Error(err)
		} else {
			defer dev.Close()
			if name, err := linux.EVGetName(dev.Fd()); err != nil {
				t.Error(err)
			} else if info, err := linux.EVGetInfo(dev.Fd()); err != nil {
				t.Error(err)
			} else if events, err := linux.EVGetSupportedEventTypes(dev.Fd()); err != nil {
				t.Error(err)
			} else {
				phys, _ := linux.EVGetPhys(dev.Fd())
				uniq, _ := linux.EVGetUniq(dev.Fd())

				t.Log(bus, "=>", linux.EVDevice(bus))
				t.Log("  name", name)
				t.Log("  info", info, "(bus,vendor,product,version)")
				t.Log("  phys", phys)
				t.Log("  uniq", uniq)
				t.Log("  events", events)
			}
		}
	}
}
