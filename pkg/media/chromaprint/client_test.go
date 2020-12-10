// +build chromaprint

package chromaprint_test

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	chromaprint "github.com/djthorpe/gopi/v3/pkg/media/chromaprint"
	tool "github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	gopi.Unit
	chromaprint.Manager
	*chromaprint.Client
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Client_001(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if app.Client == nil {
			t.Error("client is nil")
		} else {
			t.Log(app.Client)
		}
	})
}

func Test_Client_002(t *testing.T) {
	tool.Test(t, nil, new(App), func(app *App) {
		if matches, err := app.Client.Lookup("AQAAT0mUaEkSRZEGAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", 5*time.Second, chromaprint.META_TRACK); err != nil {
			t.Error(err)
		} else if len(matches) != 0 {
			t.Error("Unexpected matches")
		}
	})
}

func Test_Client_003(t *testing.T) {
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
			} else {
				fmt.Println("buf=", buf)
			}
		}
		if fp, err := stream.GetFingerprint(); err != nil {
			t.Error(err)
		} else if matches, err := app.Client.Lookup(fp, stream.Duration(), chromaprint.META_ALL); err != nil {
			t.Error(err)
		} else {
			t.Log("stream=", stream)
			t.Log("matches=", matches)
		}
	})
}
