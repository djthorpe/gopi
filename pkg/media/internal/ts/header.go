package ts

import (
	"encoding/binary"
	"fmt"
	"io"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Header struct {
	NetworkId   uint16
	Version     uint8
	Section     uint8
	LastSection uint8
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewHeader(r io.Reader) (*Header, error) {
	this := new(Header)

	if err := binary.Read(r, binary.LittleEndian, this); err != nil {
		return nil, err
	} else {
		// TODO: Current
		//Current = this.Version&0x01 != 0x00
		this.Version = (this.Version >> 1) & 0x1F
	}

	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Header) String() string {
	str := "<header"
	str += fmt.Sprintf(" network_id=0x%04X", this.NetworkId)
	str += fmt.Sprintf(" version=0x%02X", this.Version)
	str += fmt.Sprintf(" section=0x%02X", this.Section)
	str += fmt.Sprintf(" last_section=0x%02X", this.LastSection)
	return str + ">"
}
