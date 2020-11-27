package influxdb_test

import (
	"testing"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/db/influxdb"
	"github.com/djthorpe/gopi/v3/pkg/metrics"
)

func Test_LineProtocol_000(t *testing.T) {
	reserved := []string{
		"time",
		"_field",
		"_measurement",
	}
	for _, word := range reserved {
		if influxdb.IsReservedName(word) {
			t.Logf("Reserved word %q", word)
		} else {
			t.Errorf("Expected reserved word %q", word)
		}
	}
}

func Test_LineProtocol_001(t *testing.T) {
	tests := []struct{ in, out string }{
		{"time", ""},
		{"t", "t"},
	}
	for i, test := range tests {
		if f := metrics.NewField(test.in); f == nil {
			t.Errorf("Unexpected nil")
		} else if out := influxdb.QuoteFieldName(f); out != test.out {
			t.Errorf("Test %v: Unexpected output %q for input %q", i, out, test.in)
		} else {
			t.Logf("Test %v: %q => %q", i, test.in, test.out)
		}
	}
}

func Test_LineProtocol_002(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		out   string
	}{
		{"nil", nil, ""},
		{"uint8", uint8(8), "8i"},
		{"int8", int8(-8), "-8i"},
		{"uint32", uint32(32), "32i"},
		{"int32", int32(-32), "-32i"},
		{"bool", false, "false"},
		{"bool", true, "true"},
		{"float32", float32(3.14), "3.14"},
		{"float64", float64(-3), "-3"},
		{"string", "test", "\"test\""},
		{"string", "test test", "\"test test\""},
		{"string", "test\"test", "\"test\\\"test\""},
		{"string", "test\\test", "\"test\\\\test\""},
	}
	for i, test := range tests {
		if f := metrics.NewField(test.name, test.value); f == nil {
			t.Errorf("Unexpected nil")
		} else if out := influxdb.QuoteFieldValue(f); out != test.out {
			if test.name == "string" {
				t.Errorf("Test %v: Unexpected output %q for input %q", i, out, test.value.(string))
			} else {
				t.Errorf("Test %v: Unexpected output %q for input %v", i, out, test.value)
			}
		} else {
			t.Logf("Test %v: %v => %q", i, test.value, test.out)
		}
	}
}

func Test_LineProtocol_003(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
		out   string
	}{
		{"nil", nil, ""},
		{"uint8", uint8(8), "uint8=8i"},
		{"int8", int8(-8), "uint8=8i,int8=-8i"},
		{"nil", nil, "uint8=8i,int8=-8i"},
		{"bool", true, "uint8=8i,int8=-8i,bool=true"},
		{"string", "test test", "uint8=8i,int8=-8i,bool=true,string=\"test test\""},
	}
	fields := []gopi.Field{}
	for i, test := range tests {
		if f := metrics.NewField(test.name, test.value); f == nil {
			t.Errorf("Unexpected nil")
		} else {
			fields = append(fields, f)
		}
		if out, err := influxdb.QuoteFields(fields); err != nil {
			t.Error(err)
		} else if out != test.out {
			t.Errorf("Test %v: Unexpected output %q for input %q", i, out, test.value.(string))
		} else {
			t.Logf("Test %v: %v", i, out)
		}
	}
}

func Test_LineProtocol_004(t *testing.T) {
	if m, err := metrics.NewMeasurement("test", "metric bool", metrics.NewField("tag", "tag")); err != nil {
		t.Error(err)
	} else if in, err := m.Clone(time.Time{}, true); err != nil {
		t.Error(err)
	} else if out, err := influxdb.QuoteMeasurement(in); err != nil {
		t.Error(err)
	} else if out != "test,tag=\"tag\" metric=true" {
		t.Log(out)
	}
}
