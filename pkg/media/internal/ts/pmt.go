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
	ClockPid uint16 // 13 bits
	Length   uint16 // 10 bits
	D        []byte
	Streams  []ESRow
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
		s.ClockPid &= 0x1FFF
		s.Length &= 0x03FF
	}

	/*
		// Read descriptors
		s.D = make([]byte, s.Length)
		if err := binary.Read(r, binary.BigEndian, &s.D); err != nil {
			return err
		}

		// Read elementary streams
		for length := uint16(length) - (4 + s.Length); length > 0; {
			var row ESRow
			if err := binary.Read(r, binary.BigEndian, &row.header); err != nil {
				return err
			} else {
				row.header.Pid &= 0x1FFF    // 13 bits
				row.header.Length &= 0x03FF // 10 bits
			}
			if err := row.DTable.Read(r, row.header.Length); err != nil {
				return err
			} else {
				s.Streams = append(s.Streams, row)
				length -= 5 + row.header.Length
			}
		}
	*/

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s PMTSection) String() string {
	str := "<pmt"
	str += fmt.Sprint(" ", s.Header)
	str += fmt.Sprintf(" pcr_pid=0x%04X", s.ClockPid)
	str += fmt.Sprintf(" length=%v", s.Length)
	return str + ">"
}
