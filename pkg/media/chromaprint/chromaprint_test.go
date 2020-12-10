// +build chromaprint

package chromaprint_test

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"

	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

const (
	SAMPLE_FILE = "../../../etc/media/int16_44100_2ch_audio.raw"
)

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Chromaprint_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Manager == nil {
			t.Error("manager is nil")
		} else {
			t.Log(app.Manager)
		}
	})
}

func Test_Chromaprint_002(t *testing.T) {
	// Make fingerprint from five seconds of silence
	tool.Test(t, nil, new(App), func(app *App) {
		rate := 44100 // Samples per second
		ch := 1       // Channels
		duration := 5 // Seconds
		stream, err := app.Manager.NewStream(rate, ch)
		if err != nil {
			t.Error(err)
		} else if r := stream.Rate(); r != rate {
			t.Error("Unexpected rate value", r)
		} else if c := stream.Channels(); c != ch {
			t.Error("Unexpected channels value", c)
		} else {
			t.Log(stream)
		}
		buf := make([]int16, rate*duration*ch)
		for i := 0; i < duration; i++ {
			if err := stream.Write(buf); err != nil {
				t.Error(err)
			}
		}
		if fp, err := stream.GetFingerprint(); err != nil {
			t.Error(err)
		} else {
			fmt.Printf("fp=%v dur=%v\n", fp, stream.Duration())
		}
	})
}

func Test_Chromaprint_003(t *testing.T) {
	// Make fingerprint from five seconds of silence
	tool.Test(t, nil, new(App), func(app *App) {
		rate := 44100 // Samples per second
		ch := 2       // Channels
		stream, err := app.Manager.NewStream(rate, ch)
		if err != nil {
			t.Fatal(err)
		}
		fh, err := os.Open(SAMPLE_FILE)
		if err != nil {
			t.Fatal(err)
		}
		defer fh.Close()
		buf := make([]int16, rate*ch) // One second buffer
		for {
			if err := binary.Read(fh, binary.LittleEndian, buf); err == io.EOF {
				break
			} else if err != nil {
				t.Fatal(err)
			} else if err := stream.Write(buf); err != nil {
				t.Error(err)
			}
		}
		if fp, err := stream.GetFingerprint(); err != nil {
			t.Error(err)
		} else {
			t.Log("stream=", stream)
			t.Log("fp=", fp)
		}
	})
}
