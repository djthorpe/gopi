// +build pulse

package pulse_test

import (
	"math"
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
	// Both x86 and ARM are Little Endian
	spec := pulse.NewSampleSpec(pulse.PA_SAMPLE_FLOAT32LE, 44100, 1)
	if handle, err := pulse.PulseNewSimple("", t.Name(), pulse.PA_STREAM_PLAYBACK, "", "Sine Wave", spec, nil, nil); err != nil {
		t.Error(err)
	} else {
		defer handle.Free()
		buf := make([]float32, spec.Rate()*uint32(spec.Channels())) // One second of samples
		freq := float64(1000)                                       // Frequency
		// Create buffer
		for i := 0; i < len(buf); i++ {
			sample := math.Sin(2 * math.Pi * freq * float64(i) / float64(spec.Rate()))
			buf[i] = float32(sample)
		}
		// Iterate twice to get two seconds of tone
		for i := 0; i < 2; i++ {
			if err := handle.WriteFloat32(buf); err != nil {
				t.Error(err)
				break
			}
		}
		// Flush
		if err := handle.Flush(); err != nil {
			t.Error(err)
		}
	}
}

func Test_Simple_002(t *testing.T) {
	// Both x86 and ARM are Little Endian
	spec := pulse.NewSampleSpec(pulse.PA_SAMPLE_FLOAT32LE, 44100, 1)
	if handle, err := pulse.PulseNewSimple("", t.Name(), pulse.PA_STREAM_RECORD, "", "Record", spec, nil, nil); err != nil {
		t.Error(err)
	} else {
		defer handle.Free()
		buf := make([]float32, spec.Rate()*uint32(spec.Channels())*2) // Two seconds of samples
		t.Log("Recording..")
		if err := handle.ReadFloat32(buf); err != nil {
			t.Error(err)
		}
		t.Log(buf)
	}
}
