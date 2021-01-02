//+build mmal

package mmal_test

import (
	"testing"

	"github.com/djthorpe/gopi/v3/pkg/sys/mmal"
)

func Test_Types_001(t *testing.T) {
	a := mmal.NewRational(100)
	b := mmal.NewRational(200)
	if a.Add(b).Float32() != 300.0 {
		t.Error("Unexpected result")
	}
	if b.Subtract(a).Float32() != 100.0 {
		t.Error("Unexpected result")
	}
	if b.Divide(a).Float32() != 2.0 {
		t.Error("Unexpected result")
	}
	if b.Multiply(a).Float32() != 100*200 {
		t.Error("Unexpected result")
	}
}

func Test_Types_002(t *testing.T) {
	a := mmal.NewRect(0, 0, 100, 200)
	if x, y := a.Origin(); x != 0 || y != 0 {
		t.Error("Unexpected result")
	} else if w, h := a.Size(); w != 100 || h != 200 {
		t.Error("Unexpected result")
	} else {
		t.Log(a)
	}
}
