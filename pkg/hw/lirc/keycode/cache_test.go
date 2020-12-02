package keycode_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/hw/lirc/keycode"
	"github.com/djthorpe/gopi/v3/pkg/tool"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type CacheApp struct {
	gopi.Unit
	*keycode.Cache
}

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Cache_001(t *testing.T) {
	tool.Test(t, nil, new(CacheApp), func(app *CacheApp) {
		if app.Cache == nil {
			t.Error("nil Cache unit")
		} else {
			t.Log(app.Cache)
		}
	})
}

func Test_Cache_002(t *testing.T) {
	tool.Test(t, nil, new(CacheApp), func(app *CacheApp) {
		for k := gopi.KEYCODE_NONE; k <= gopi.KEYCODE_MAX; k++ {
			name := app.Cache.KeycodeName(k)
			if k2 := app.Cache.LookupKeycode(name); k2 != k {
				t.Error("Unexpected name->keycode", name, k2, " (expected", k, ")")
			}
		}
	})
}

func Test_Cache_003(t *testing.T) {
	tool.Test(t, nil, new(CacheApp), func(app *CacheApp) {
		for d := gopi.INPUT_DEVICE_MIN; d != gopi.INPUT_DEVICE_MAX; d <<= 1 {
			name := app.Cache.DeviceName(d)
			if d2 := app.Cache.LookupDevice(name); d2 != d {
				t.Error("Unexpected name->device", name, d2, " (expected", d, ")")
			}
		}
	})
}

func Test_Cache_004(t *testing.T) {
	tool.Test(t, nil, new(CacheApp), func(app *CacheApp) {
		good := []string{
			"",
			"# This is a test",
			"KEYCODE_Q RC5_14 0x01",
			"KEYCODE_Q SONY_20 0x0101 # With a comment",
			"KEYCODE_Q RC5_14 0xFF",
			"KEYCODE_Q SONY_12 0xFFEE # With a comment",
		}
		for i, line := range good {
			if entry, err := app.Cache.DecodeLine(line); err != nil {
				t.Errorf("Test %d: 1: Unexpected error returned for %q: %v", i, line, err)
			} else if line_ := app.Cache.EncodeLine(entry); line_ != line {
				t.Errorf("Test %d: 2: Unexpected difference for %q (expected %q)", i, line, line_)
			} else if entry_, err := app.Cache.DecodeLine(line_); err != nil {
				t.Errorf("Test %d: 3: Unexpected error returned for %q: %v", i, line_, err)
			} else if entry_.Equals(entry) == false {
				t.Errorf("Test %d: 4: Unexpected difference for %q", i, line_)
			} else {
				t.Logf("Test %d: %q => %v", i, line, entry)
			}
		}
	})
}
