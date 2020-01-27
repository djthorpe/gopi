// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bitmap_test

import (
	"encoding/hex"
	"testing"

	// Frameworks
	rpi "github.com/djthorpe/gopi/v2/sys/rpi"
	bitmap "github.com/djthorpe/gopi/v2/unit/surfaces/bitmap"
)

func init() {
	rpi.DXInit()
}

func Test_Data_000(t *testing.T) {
	t.Log("Test_Data_000")
}

func Test_Data_001(t *testing.T) {
	data := new(bitmap.Data)
	if data.Bytes() != nil {
		t.Error("Unexpected value for Bytes", data.Bytes())
	}
}

func Test_Data_002(t *testing.T) {
	data := new(bitmap.Data)
	if err := data.SetCapacity(1); err != nil {
		t.Error(err)
	}
	if data.Bytes() == nil {
		t.Error("Unexpected value for Bytes", data.Bytes())
	}
	if data.Uint32() == nil {
		t.Error("Unexpected value for Uint32", data.Uint32())
	}
}

func Test_Data_003(t *testing.T) {
	data := new(bitmap.Data)
	if err := data.SetCapacity(1); err != nil {
		t.Error(err)
	}
	if len(data.Bytes()) != 1 {
		t.Error("Unexpected value for cap=1 Bytes", data.Bytes())
	}
	if len(data.Uint32()) != 1 {
		t.Error("Unexpected value for cap=1 Uint32", data.Uint32())
	}
	if err := data.SetCapacity(4); err != nil {
		t.Error(err)
	}
	if len(data.Bytes()) != 4 {
		t.Error("Unexpected value for cap=4 Bytes", data.Bytes())
	}
	if len(data.Uint32()) != 1 {
		t.Error("Unexpected value for cap=4 Uint32", data.Uint32())
	}
	if err := data.SetCapacity(5); err != nil {
		t.Error(err)
	}
	if len(data.Bytes()) != 5 {
		t.Error("Unexpected value for cap=5 Bytes", data.Bytes())
	}
	if len(data.Uint32()) != 2 {
		t.Error("Unexpected value for cap=5 Uint32", data.Uint32())
	}
	if err := data.SetCapacity(7); err != nil {
		t.Error(err)
	}
	if len(data.Bytes()) != 7 {
		t.Error("Unexpected value for cap=7 Bytes", data.Bytes())
	}
	if len(data.Uint32()) != 2 {
		t.Error("Unexpected value for cap=7 Uint32", data.Uint32())
	}
	if err := data.SetCapacity(8); err != nil {
		t.Error(err)
	}
	if len(data.Bytes()) != 8 {
		t.Error("Unexpected value for cap=8 Bytes", data.Bytes())
	}
	if len(data.Uint32()) != 2 {
		t.Error("Unexpected value for cap=8 Uint32", data.Uint32())
	}
	if err := data.SetCapacity(9); err != nil {
		t.Error(err)
	}
	if len(data.Bytes()) != 9 {
		t.Error("Unexpected value for cap=9 Bytes", data.Bytes())
	}
	if len(data.Uint32()) != 3 {
		t.Error("Unexpected value for cap=9 Uint32", data.Uint32())
	}
}

func Test_Data_004(t *testing.T) {
	data := new(bitmap.Data)
	if err := data.SetCapacity(4); err != nil {
		t.Error(err)
	}
	bytes := data.Bytes()
	bytes[0] = 0xFF
	bytes[1] = 0xFE
	bytes[2] = 0xFD
	bytes[3] = 0xFA
	uints := data.Uint32()
	if uints[0] != 0xFAFDFEFF {
		t.Errorf("Unexpected uint32 value %04X with bytes %v", uints[0], hex.EncodeToString(bytes))
	}

	data.SetCapacity(0)
	if data.Bytes() != nil {
		t.Error("Unexpected Bytes() value", data.Bytes())
	}
	if data.Uint32() != nil {
		t.Error("Unexpected Uint32() value", data.Uint32())
	}
}
