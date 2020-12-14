package argonone_test

import (
	"testing"
	"time"

	"github.com/djthorpe/gopi/v3/pkg/dev/argonone"
)

////////////////////////////////////////////////////////////////////////////////
// TESTS

func Test_Value_001(t *testing.T) {
	v := argonone.NewValueWithDelta(time.Second)
	if v, c := v.Get(); v != nil {
		t.Error("Expected nil value from v.Get()")
	} else if c != false {
		t.Error("Expected false changed value from v.Get()")
	}
	if v, c := v.Set(1); v != 1 {
		t.Error("Unexpected value from v.Set()", v)
	} else if c != true {
		t.Error("Expected changed value from v.Set()")
	}
	if v, c := v.Set(2); v != 1 {
		t.Error("Unexpected value from v.Set()", v)
	} else if c != false {
		t.Error("Expected changed value from v.Set()")
	}
	time.Sleep(500 * time.Millisecond)
	if v, c := v.Set(2); v != 1 {
		t.Error("Unexpected value from v.Set()", v)
	} else if c != false {
		t.Error("Expected changed value from v.Set()")
	}
	time.Sleep(500 * time.Millisecond)
	if v, c := v.Set(2); v != 2 {
		t.Error("Unexpected value from v.Set()", v)
	} else if c != true {
		t.Error("Expected changed value from v.Set()")
	}
}
