/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package util_test

import (
	"testing"
	"fmt"
)

import (
    util "github.com/djthorpe/gopi/util"
)

var (
	PARSE_TESTS_GOOD = map[string]float64 {
		"100in": 100.00,
		"100 in": 100.00,
		"   9 in": 9.00,
		"   9in   ": 9.00,
		"   7.5   in   ": 7.50,
		"   110   mm   ": 4.33,
		"   250 cm   ": 98.43,
		" 1x1 in   ": 1.41,
		" 56.1 x 49.9 cm   ": 29.56,
		" 50 x 100 mm   ": 4.40,
	}
	PARSE_TESTS_BAD = []string{
		"0.0", "", "    ","string 56in", "56", "-100in",
	}
)

func Test000_Parse_Good(t *testing.T) {
	for k,v := range PARSE_TESTS_GOOD {
		value, err := util.ParseLengthString(k)
		if err != nil {
			t.Error("String=",k," returned error ",err)
			continue
		}
		got := fmt.Sprintf("%.2f",value)
		expected := fmt.Sprintf("%.2f",v)
		if got != expected {
			t.Error("String=",k," expected ",expected," but got ",got)
		}
	}
}

func Test001_Parse_Bad(t *testing.T) {
	for _,k := range PARSE_TESTS_BAD {
		_, err := util.ParseLengthString(k)
		if err == nil {
			t.Error("String=",k," expected error, but returned nil")
		}
	}
}
