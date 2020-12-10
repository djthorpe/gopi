// +build pulse

package pulse_test

import (
	"encoding/binary"
	"io"
	"math"
	"os"
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/pulse"
)

const (
	SAMPLE_FILE = "../../../etc/media/int16_44100_2ch_audio.raw"
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
	spec := pulse.NewSampleSpec(pulse.PA_SAMPLE_S16LE, 44100, 2)
	handle, err := pulse.PulseNewSimple("", t.Name(), pulse.PA_STREAM_PLAYBACK, "", "Sample", spec, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer handle.Free()
	fh, err := os.Open(SAMPLE_FILE)
	if err != nil {
		t.Fatal(err)
	}
	defer fh.Close()
	buf := make([]int16, 44100*2)
	for {
		if err := binary.Read(fh, binary.LittleEndian, buf); err == io.EOF {
			break
		} else if err != nil {
			t.Error(err)
		} else if err := handle.WriteInt16(buf); err != nil {
			t.Error(err)
		}
	}
	// Flush
	if err := handle.Flush(); err != nil {
		t.Error(err)
	}
}

func Test_Simple_003(t *testing.T) {
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
