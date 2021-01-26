// +build dvb

package dvb_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/dvb"
)

func Test_Demux_000(t *testing.T) {
	filter := dvb.NewSectionFilter(0, 0, 0)
	t.Log(filter)
}

func Test_Demux_001(t *testing.T) {
	filter := dvb.NewStreamFilter(0, 0, 0, 0, 0)
	t.Log(filter)
}

func Test_Demux_002(t *testing.T) {
	devices := dvb.Devices()
	if len(devices) == 0 {
		t.Skip("Skipping test, no devices available")
	}
	for _, device := range devices {
		dev, err := device.DMXOpen()
		if err != nil {
			t.Error(err)
		}
		defer dev.Close()
		filter := dvb.NewStreamFilter(0, 0, 0, 0, 0)
		if err := dvb.DMXSetStreamFilter(dev.Fd(), filter); err != nil {
			t.Error(err)
		} else if err := dvb.DMXStart(dev.Fd()); err != nil {
			t.Error(err)
		} else if pids, err := dvb.DMXGetStreamPids(dev.Fd()); err != nil {
			t.Error(err)
		} else if err := dvb.DMXStop(dev.Fd()); err != nil {
			t.Error(err)
		} else {
			t.Log(dev.Name())
			for k, v := range pids {
				t.Logf("  pid 0x%04X => %v", v, k)
			}
		}
	}
}
