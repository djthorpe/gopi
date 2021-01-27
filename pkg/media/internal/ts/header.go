package ts

import (
	"encoding/binary"
	"fmt"
	"io"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Header struct {
	Id          uint16
	Version     uint8
	Section     uint8
	LastSection uint8
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (h *Header) Read(r io.Reader) error {
	if err := binary.Read(r, binary.BigEndian, h); err != nil {
		return err
	} else {
		// TODO: Current
		//Current = this.Version&0x01 != 0x00
		h.Version = (h.Version >> 1) & 0x1F
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (h Header) String() string {
	str := "<header"
	str += fmt.Sprintf(" id=0x%04X", h.Id)
	str += fmt.Sprintf(" version=0x%02X", h.Version)
	str += fmt.Sprintf(" section=0x%02X", h.Section)
	str += fmt.Sprintf(" last_section=0x%02X", h.LastSection)
	return str + ">"
}
