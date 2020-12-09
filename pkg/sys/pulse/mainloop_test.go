// +build pulse

package pulse_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/pulse"
)

func Test_Mainloop_000(t *testing.T) {
	if loop := pulse.NewMainloop(); loop == nil {
		t.Error("Unexpected nil value")
	} else {
		loop.Free()
	}
}

func Test_Mainloop_001(t *testing.T) {
	if loop := pulse.NewMainloop(); loop == nil {
		t.Error("Unexpected nil value")
	} else if err := loop.Start(); err != nil {
		t.Error(err)
	} else {
		loop.Stop()
		loop.Free()
	}
}
