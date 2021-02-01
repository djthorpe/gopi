package ts

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type PMTSection struct {
	Header
	ClockPid uint16 // 13 bits
	Length   uint16 // 12 bits
	DTable
	ESTable
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (s *PMTSection) Read(r io.Reader, length int) error {
	if err := s.Header.Read(r); err != nil {
		return err
	} else if err := binary.Read(r, binary.BigEndian, &s.ClockPid); err != nil {
		return err
	} else if err := binary.Read(r, binary.BigEndian, &s.Length); err != nil {
		return err
	} else {
		s.ClockPid &= 0x1FFF // 13 bits
		s.Length &= 0x0FFF   // 12 bits
	}

	/* Section and Last Section should always be zero */
	if s.Header.Section != 0 || s.Header.LastSection != 0 {
		return gopi.ErrUnexpectedResponse.WithPrefix("PMT")
	}

	/* Read descriptors */
	if err := s.DTable.Read(r, s.Length); err != nil {
		return err
	}

	/* Read elementary streams */
	if err := s.ESTable.Read(r, uint16(length)-(s.Length+4)); err != nil {
		return err
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s PMTSection) String() string {
	str := "<pmt"
	str += fmt.Sprint(" ", s.Header)
	str += fmt.Sprintf(" pcr_pid=0x%04X", s.ClockPid)
	str += fmt.Sprintf(" d=%v", s.DTable)
	str += fmt.Sprintf(" es=%v", s.ESTable)
	return str + ">"
}
