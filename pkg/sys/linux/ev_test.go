// +build linux

package linux_test

import (
	"testing"

	// Frameworks
	"github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

func Test_Event_000(t *testing.T) {
	devices, err := linux.EVDevices()
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		t.Log("dev=", linux.EVDevice(device))
	}
}

func Test_Event_001(t *testing.T) {
	devices, err := linux.EVDevices()
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		if fh, err := linux.EVOpenDevice(device); err != nil {
			t.Error(err)
		} else {
			defer fh.Close()
			if name, err := linux.EVGetName(fh.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Logf("dev=%v name=%q", device, name)
			}
		}
	}
}

func Test_Event_002(t *testing.T) {
	devices, err := linux.EVDevices()
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		if fh, err := linux.EVOpenDevice(device); err != nil {
			t.Error(err)
		} else {
			defer fh.Close()
			if name, err := linux.EVGetPhys(fh.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Logf("dev=%v phys=%q", device, name)
			}
		}
	}
}
func Test_Event_003(t *testing.T) {
	devices, err := linux.EVDevices()
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		if fh, err := linux.EVOpenDevice(device); err != nil {
			t.Error(err)
		} else {
			defer fh.Close()
			if name, err := linux.EVGetUniq(fh.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Logf("dev=%v uniq=%q", device, name)
			}
		}
	}
}
func Test_Event_004(t *testing.T) {
	devices, err := linux.EVDevices()
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		if fh, err := linux.EVOpenDevice(device); err != nil {
			t.Error(err)
		} else {
			defer fh.Close()
			if info, err := linux.EVGetInfo(fh.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Logf("dev=%v info=%v", device, info)
			}
		}
	}
}

func Test_Event_005(t *testing.T) {
	devices, err := linux.EVDevices()
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		if fh, err := linux.EVOpenDevice(device); err != nil {
			t.Error(err)
		} else {
			defer fh.Close()
			if evts, err := linux.EVGetSupportedEventTypes(fh.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Logf("dev=%v supported_events=%v", device, evts)
			}
		}
	}
}

func Test_Event_006(t *testing.T) {
	devices, err := linux.EVDevices()
	if err != nil {
		t.Fatal(err)
	}
	for _, device := range devices {
		if fh, err := linux.EVOpenDevice(device); err != nil {
			t.Error(err)
		} else {
			defer fh.Close()
			if leds, err := linux.EVGetLEDState(fh.Fd()); err != nil {
				t.Error(err)
			} else {
				t.Logf("dev=%v leds=%v", device, leds)
			}
		}
	}
}
