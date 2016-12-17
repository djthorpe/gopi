/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package util_test

import (
	"fmt"
	"testing"
)

import (
	util "github.com/djthorpe/gopi/util"
)

type PPIInput struct {
	W, H  uint
	Value string
}

var (
	PARSE_TESTS_GOOD_000 = map[string]float64{
		"100in":              100.00,
		"100 in":             100.00,
		"   9 in":            9.00,
		"   9in   ":          9.00,
		"   7.5   in   ":     7.50,
		"   110   mm   ":     4.33,
		"   250 cm   ":       98.43,
		" 1x1 in   ":         1.41,
		" 56.1 x 49.9 cm   ": 29.56,
		" 50 x 100 mm   ":    4.40,
	}
	PARSE_TESTS_BAD_001 = []string{
		"0.0", "", "    ", "string 56in", "56", "-100in",
	}
	PARSE_TESTS_GOOD_002 = map[uint]PPIInput{
		101: {0, 0, "101"},
		100: {800, 500, "8in"},
		72:  {800, 500, "72"},
		0:   {0, 0, ""},
	}
)

func Test000_Parse_Good(t *testing.T) {
	for k, v := range PARSE_TESTS_GOOD_000 {
		value, err := util.ParseLengthString(k)
		if err != nil {
			t.Error("String=", k, " returned error ", err)
			continue
		}
		got := fmt.Sprintf("%.2f", value)
		expected := fmt.Sprintf("%.2f", v)
		if got != expected {
			t.Error("String=", k, " expected ", expected, " but got ", got)
		}
	}
}

func Test001_Parse_Bad(t *testing.T) {
	for _, k := range PARSE_TESTS_BAD_001 {
		_, err := util.ParseLengthString(k)
		if err == nil {
			t.Error("String=", k, " expected error, but returned nil")
		}
	}
}

func Test002_Parse_Good(t *testing.T) {
	for v, k := range PARSE_TESTS_GOOD_002 {
		ppi, err := util.PixelsPerInch(k.W, k.H, k.Value)
		if err != nil {
			t.Error("Input=", k, "returned error:", err)
			continue
		}
		if ppi != v {
			t.Error("Input=", k, "expected:", v, "but got:", ppi)
		}
	}
}
