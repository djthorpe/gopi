package ts

import (
	"encoding/binary"
	"fmt"
	"io"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type PMTSection struct {
	Header
	ClockPid uint16
	DTable
	STable
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (s *PMTSection) Read(r io.Reader, length int) error {
	if err := s.Header.Read(r); err != nil {
		return err
	} else if err := binary.Read(r, binary.BigEndian, &s.ClockPid); err != nil {
		return err
	} else if err := s.DTable.Read(r); err != nil {
		return err
	} else if err := s.STable.Read(r); err != nil {
		return err
	}

	// Mask clock pid
	s.ClockPid &= 0x1FFF

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s PMTSection) String() string {
	str := "<pmt"
	str += fmt.Sprint(" ", s.Header)
	str += fmt.Sprintf(" pcr_pid=0x%04X", s.ClockPid)
	str += fmt.Sprint(" ", s.DTable)
	str += fmt.Sprint(" ", s.STable)
	return str + ">"
}
