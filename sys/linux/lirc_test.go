// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package linux_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v2/sys/linux"
)

func Test_LIRC_000(t *testing.T) {
	t.Log("Test_LIRC_000")
}
func Test_LIRC_001(t *testing.T) {
	for f := linux.LIRC_FEATURE_MIN; f != 0; f <<= 1 {
		str := fmt.Sprint(f)
		if strings.HasPrefix(str, "LIRC_") {
			t.Logf("%08X => %v", uint32(f), str)
		}
	}
}

func Test_LIRC_002(t *testing.T) {
	for bus := uint(0); bus < 10; bus++ {
		if _, err := os.Stat(linux.LIRCDevice(bus)); os.IsNotExist(err) {
			continue
		} else {
			t.Log(linux.LIRCDevice(bus))
		}
	}
}

func Test_LIRC_003(t *testing.T) {
	for bus := uint(0); bus < 10; bus++ {
		if _, err := os.Stat(linux.LIRCDevice(bus)); os.IsNotExist(err) {
			continue
		} else if fh, err := linux.LIRCOpenDevice(bus, linux.LIRC_MODE_RCV|linux.LIRC_MODE_SEND); err != nil {
			t.Error(err)
		} else {
			defer fh.Close()
			if features, err := linux.LIRCFeatures(fh.Fd()); err != nil {
				t.Error(err)
			} else if features&linux.LIRC_CAN_REC_MASK != 0 && features&linux.LIRC_CAN_SEND_MASK != 0 {
				t.Log(bus, "[send and rcv] =>", features)
			} else if features&linux.LIRC_CAN_REC_MASK != 0 {
				t.Log(bus, "[rcv] =>", features)
			} else if features&linux.LIRC_CAN_SEND_MASK != 0 {
				t.Log(bus, "[send] =>", features)
			} else {
				t.Log(bus, "[????] =>", features)
			}
		}
	}
}
