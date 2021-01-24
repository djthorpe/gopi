// +build dvb

package dvb_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/dvb"
)

func Test_Device_000(t *testing.T) {
	devices := dvb.Devices()
	if devices == nil {
		t.Skip("Skipping test, no devices available")
	}
	for _, device := range devices {
		t.Log(device)
	}
}
