// +build dvb

package dvb_test

import (
	"os"
	"testing"
	"time"

	"github.com/djthorpe/gopi/v3/pkg/sys/dvb"
)

func Test_Frontend_000(t *testing.T) {
	devices := dvb.Devices()
	if len(devices) == 0 {
		t.Skip("Skipping test, no devices available")
	}
	for _, device := range devices {
		file, err := device.FEOpen(os.O_RDONLY)
		if err != nil {
			t.Error(err)
		}
		defer file.Close()
		if info, err := dvb.FEGetInfo(file.Fd()); err != nil {
			t.Error(err)
		} else {
			t.Log(device)
			t.Log("  GetInfo =>", info)
		}
	}
}

func Test_Frontend_001(t *testing.T) {
	devices := dvb.Devices()
	if len(devices) == 0 {
		t.Skip("Skipping test, no devices available")
	}
	for _, device := range devices {
		file, err := device.FEOpen(os.O_RDONLY)
		if err != nil {
			t.Error(err)
		}
		defer file.Close()
		if major, minor, err := dvb.FEGetVersion(file.Fd()); err != nil {
			t.Error(err)
		} else {
			t.Log(device, "version=", major, ".", minor)
		}
	}
}

func Test_Frontend_002(t *testing.T) {
	devices := dvb.Devices()
	if len(devices) == 0 {
		t.Skip("Skipping test, no devices available")
	}
	for _, device := range devices {
		file, err := device.FEOpen(os.O_RDONLY)
		if err != nil {
			t.Error(err)
		}
		defer file.Close()
		if sys, err := dvb.FEEnumDeliverySystems(file.Fd()); err != nil {
			t.Error(err)
		} else {
			t.Log(device, "FEEnumDeliverySystems=", sys)
		}
	}
}

func Test_Frontend_003(t *testing.T) {
	devices := dvb.Devices()
	if len(devices) == 0 {
		t.Skip("Skipping test, no devices available")
	}
	// Use first device for tuning, and first channel in tune table
	fh, err := os.Open(FILES[0])
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if params, err := dvb.ReadTuneParamsTable(fh); err != nil {
		t.Fatal(err)
	} else if len(params) == 0 {
		t.Fatal("No tune parameters in", FILES[0])
	} else if dev, err := devices[0].FEOpen(os.O_RDWR); err != nil {
		t.Fatal(err)
	} else {
		defer dev.Close()
		t.Log("Tuning for", params[0])
		if err := dvb.FETune(dev.Fd(), params[0]); err != nil {
			t.Error(err)
		}
	}
}
func Test_Frontend_004(t *testing.T) {
	devices := dvb.Devices()
	if len(devices) == 0 {
		t.Skip("Skipping test, no devices available")
	}
	// Use first device for tuning, and first channel in tune table
	fh, err := os.Open(FILES[0])
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	if params, err := dvb.ReadTuneParamsTable(fh); err != nil {
		t.Fatal(err)
	} else if len(params) == 0 {
		t.Fatal("No tune parameters in", FILES[0])
	} else if dev, err := devices[0].FEOpen(os.O_RDWR); err != nil {
		t.Fatal(err)
	} else {
		defer dev.Close()
		t.Log("Tuning for", params[0])
		if err := dvb.FETune(dev.Fd(), params[0]); err != nil {
			t.Fatal(err)
		}
		ticker := time.NewTicker(time.Millisecond * 100)
		timer := time.NewTimer(5 * time.Second)
		defer ticker.Stop()
		defer timer.Stop()
	FOR_LOOP:
		for {
			select {
			case <-timer.C:
				t.Log("Tune timeout")
				break FOR_LOOP
			case <-ticker.C:
				if status, err := dvb.FEReadStatus(dev.Fd()); err != nil {
					t.Fatal(err)
				} else if status != dvb.FE_NONE {
					t.Log("  status=", status)
					if status&dvb.FE_HAS_LOCK == dvb.FE_HAS_LOCK {
						break FOR_LOOP
					}
				}
			}
		}
	}
}
