package metrics

import (
	"fmt"
	"testing"
	"time"
)

func Test_Field_000(t *testing.T) {
	if f := NewField("test"); f != nil {
		t.Log(f)
	} else {
		t.Error("NewField returns nil value")
	}
	if f := NewField(""); f != nil {
		t.Error("NewField should return nil value")
	}
	if f := NewField("_"); f != nil {
		t.Error("NewField should return nil value")
	}
	if f := NewField("_a"); f != nil {
		t.Error("NewField should return nil value")
	}
	if f := NewField("0a"); f != nil {
		t.Error("NewField should return nil value")
	}
	if f := NewField("-a"); f != nil {
		t.Error("NewField should return nil value")
	}
	if f := NewField("a0"); f == nil {
		t.Error("NewField returns nil value")
	}
	if f := NewField("a-"); f == nil {
		t.Error("NewField returns nil value")
	}
	if f := NewField("a_"); f == nil {
		t.Error("NewField returns nil value")
	}
}

func Test_Field_001(t *testing.T) {
	if f := NewField("test", nil, nil); f != nil {
		t.Error("NewField should return nil value")
	}
	if f := NewField("test", nil); f == nil {
		t.Error("NewField should not return nil value")
	}
}

func Test_Field_002(t *testing.T) {
	tests := []struct {
		v interface{}
		k string
	}{
		{true, "bool"},
		{false, "bool"},
		{int8(0), "int8"},
		{int16(0), "int16"},
		{int32(0), "int32"},
		{int64(0), "int64"},
		{uint8(0), "uint8"},
		{uint16(0), "uint16"},
		{uint32(0), "uint32"},
		{uint64(0), "uint64"},
		{time.Time{}, "time.Time"},
		{"string", "string"},
		{"", "string"},
		{float32(3.14), "float32"},
		{float64(31.4), "float64"},
	}
	for i, test := range tests {
		c := fmt.Sprintf("%v_%v", t.Name(), i)
		if f := NewField(c, test.v); f == nil {
			t.Error("NewField should not return nil value")
		} else if kind := f.Kind(); kind != test.k {
			t.Error("NewField kind is incorrect", kind, "!=", test.k)
		}
	}
}
