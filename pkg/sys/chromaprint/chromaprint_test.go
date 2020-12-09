// +build chromaprint

package chromaprint_test

import (
	"testing"

	chromaprint "github.com/djthorpe/gopi/v3/pkg/sys/chromaprint"
)

////////////////////////////////////////////////////////////////////////////////
// TEST ENUMS

func Test_chromaprint_000(t *testing.T) {
	t.Log("Test_chromaprint_000")
}

func Test_chromaprint_001(t *testing.T) {
	if version := chromaprint.Version(); version != "" {
		t.Log("Version:", version)
	} else {
		t.Error("Expected non-empty version")
	}
}

func Test_chromaprint_002(t *testing.T) {
	if ctx := chromaprint.NewChromaprint(chromaprint.ALGORITHM_DEFAULT); ctx == nil {
		t.Error("Unexpected nil return from NewChromaprint")
	} else {
		ctx.Free()
	}
}

func Test_chromaprint_003(t *testing.T) {
	ctx := chromaprint.NewChromaprint(chromaprint.ALGORITHM_DEFAULT)
	if ctx == nil {
		t.Error("Unexpected nil return from NewChromaprint")
	}
	defer ctx.Free()
	rate := 44100
	ch := 1
	size := rate * 5 * ch // Uint16 samples for 5 seconds
	if err := ctx.Start(44100, 2); err != nil {
		t.Error(err)
	}
	buf := make([]byte, size*2) // Two bytes per sample
	for i := 0; i < 5; i++ {
		t.Log("Feeding 5 seconds of silence...")
		if err := ctx.Feed(buf); err != nil {
			t.Error(err)
		}
	}
	if err := ctx.Finish(); err != nil {
		t.Error(err)
	}
	if fp, err := ctx.GetFingerprint(); err != nil {
		t.Error(err)
	} else {
		t.Log("Fingerprint=", fp)
	}

}
