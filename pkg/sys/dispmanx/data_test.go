package dispmanx_test

import (
	"testing"

	dx "github.com/djthorpe/gopi/v3/pkg/sys/dispmanx"
)

func Test_Data_001(t *testing.T) {
	data := dx.NewData(10)
	if data == nil {
		t.Fatal("NewData returned nil")
	}
	defer data.Dispose()
	// Capacity is on uint32 boundaries
	if cap := data.Cap(); cap != 12 {
		t.Error("Unexpected capacity", cap)
	}
	// Clear data to 0xFF
	data.SetUint8(0, data.Cap(), 0xFF)

	// Check data
	for i, val := range data.Byte(0) {
		if val != 0xFF {
			t.Errorf("Unexpected value at i=%d val=0x%02X", i, val)
		}
	}

	// Set data using offset
	for i := uint32(0); i < data.Cap(); i++ {
		data.SetUint8(uintptr(i), 1, uint8(i))
	}

	// Check data
	for i, val := range data.Byte(0) {
		if val != uint8(i) {
			t.Errorf("Unexpected value at i=%d val=0x%02X", i, val)
		}
	}

	// Set data in twos
	for i := uint32(0); i < data.Cap(); i += 2 {
		data.SetUint8(uintptr(i), 2, uint8(i))
	}

	// Check data
	for i, val := range data.Byte(0) {
		if val != uint8(i)>>1<<1 {
			t.Errorf("Unexpected value at i=%d val=0x%02X (expected 0x%02X)", i, val, uint8(i)>>1<<1)
		}
	}
}

func Test_Data_002(t *testing.T) {
	data := dx.NewData(10)

	// Clear data to 0x1234
	data.SetUint16(0, data.Cap()>>1, 0x1234)

	// Check data
	for i, val := range data.Uint16(0) {
		if val != 0x1234 {
			t.Errorf("Unexpected value at i=%d val=0x%04X", i, val)
		}
	}

	// Set data using offset
	for i := uint32(0); i < data.Cap()>>1; i++ {
		data.SetUint16(uintptr(i)<<1, 1, uint16(i))
	}

	// Check data
	for i, val := range data.Uint16(0) {
		if val != uint16(i) {
			t.Errorf("Unexpected value at i=%d val=0x%04X", i, val)
		}
	}

	// Set data in twos
	for i := uint32(0); i < data.Cap()>>1; i += 2 {
		data.SetUint16(uintptr(i<<1), 2, uint16(0xFFF0+i))
	}

	// Check data
	for i, val := range data.Uint16(0) {
		if val != uint16(0xFFF0+i)>>1<<1 {
			t.Errorf("Unexpected value at i=%d val=0x%02X (expected 0x%02X)", i, val, uint16(i)>>1<<1)
		}
	}
}

func Test_Data_003(t *testing.T) {
	data := dx.NewData(24) // 24 bytes equals 6 pixels

	// Clear data to 0x12345678
	data.SetUint32(0, data.Cap()>>2, 0x12345678)

	// Check data
	for i, val := range data.Uint32(0) {
		if val != 0x12345678 {
			t.Errorf("Unexpected value at i=%d val=0x%08X", i, val)
		}
	}

	// Set data using offset
	for i := uint32(0); i < data.Cap()>>2; i++ {
		data.SetUint32(uintptr(i)<<2, 1, uint32(i))
	}

	// Check data
	for i, val := range data.Uint32(0) {
		if val != uint32(i) {
			t.Errorf("Unexpected value at i=%d val=0x%08X", i, val)
		}
	}
}

func Test_Data_004(t *testing.T) {
	data := dx.NewData(2) // 2 bytes is 16 bits

	// Clear data to 0xFFFF
	data.SetBit(0, 16, true)
	t.Log("set all first 16 bits to true =>", data)
	data.SetBit(0, 15, false)
	t.Log("set first 15 bits to false =>", data)
	data.SetBit(0, 10, true)
	t.Log("set first 10 bits to true =>", data)
}
