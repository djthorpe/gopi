// +build pulse

package pulse_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/pulse"
)

func Test_Simple_000(t *testing.T) {
	ss := pulse.NewSampleSpec(pulse.PA_SAMPLE_U8, 44100, 2)
	if handle, err := pulse.PulseNewSimple("", t.Name(), pulse.PA_STREAM_PLAYBACK, "", "Music", ss, nil, nil); err != nil {
		t.Error(err)
	} else {
		handle.Free()
	}
}

func Test_Simple_001(t *testing.T) {
	ss := pulse.NewSampleSpec(pulse.PA_SAMPLE_U8, 44100, 2)
	if handle, err := pulse.PulseNewSimple("", t.Name(), pulse.PA_STREAM_PLAYBACK, "", "Music", ss, nil, nil); err != nil {
		t.Error(err)
	} else {
		defer handle.Free()
		buf := make([]byte, 100)
		if err := handle.Write(buf); err != nil {
			t.Error(err)
		}
		if err := handle.Flush(); err != nil {
			t.Error(err)
		}
	}
}
