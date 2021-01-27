package ts

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Table struct {
	Rows []Row
}

type RowHeader struct {
	Tag    uint8
	Length uint8
}

type Row struct {
	RowHeader
	data []byte
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewTable(r io.Reader) (*Table, error) {
	this := new(Table)

	// Read length of descriptor that follow in bytes
	var length uint16
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	} else {
		// Top four bits are reserved for future use, length is 12 bits
		length &= 0x0FFF
	}

	// Read rows until length is zero
	for length > 0 {
		fmt.Println("remaining length=", length)
		var row Row
		if err := binary.Read(r, binary.LittleEndian, &row.RowHeader); err != nil {
			return nil, err
		}
		fmt.Printf("  tag=%02X length=%v\n", row.RowHeader.Tag, row.RowHeader.Length)
		row.data = make([]byte, int(row.Length))
		if _, err := r.Read(row.data); err != nil {
			return nil, err
		}
		if row.Tag == 0x40 { // Network name
			fmt.Printf("  data=%q\n", string(row.data))
		}

		// Append row
		this.Rows = append(this.Rows, row)

		// Decrement by 2 bytes and length of data
		length -= 2
		length -= uint16(row.Length)
	}

	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Table) String() string {
	str := ""
	for _, row := range this.Rows {
		str += fmt.Sprintf("<0x%02X=%v> ", row.Tag, hex.EncodeToString(row.data))
	}
	return strings.TrimSuffix(str, " ")
}
