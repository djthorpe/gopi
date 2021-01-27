package ts

import (
	"encoding/binary"
	"fmt"
	"io"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type PATSection struct {
	Header
	Programs []PATProgram
}

type PATProgram struct {
	Program uint16
	Pid     uint16
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (p *PATSection) Read(r io.Reader, length int) error {
	if err := p.Header.Read(r); err != nil {
		return err
	}

	// Read entries until length is 0
	for length > 0 {
		var row PATProgram
		if err := binary.Read(r, binary.BigEndian, &row); err != nil {
			return err
		} else {
			row.Pid &= 0x1FFF // Use 13 bits for Pid
			p.Programs = append(p.Programs, row)
		}
		length -= 4
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (p PATSection) String() string {
	str := "<pat"
	str += fmt.Sprint(" ", p.Header)
	for _, row := range p.Programs {
		str += fmt.Sprint(" ", row)
	}
	return str + ">"
}

func (r PATProgram) String() string {
	return fmt.Sprintf("<id=0x%04X pid=0x%04X>", r.Program, r.Pid)
}
