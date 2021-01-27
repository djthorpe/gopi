package ts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

// Ref: https://www.etsi.org/deliver/etsi_ts/101200_101299/101211/01.11.01_60/ts_101211v011101p.pdf
// Ref: https://www.etsi.org/deliver/etsi_en/300400_300499/300468/01.03.01_60/en_300468v010301p.pdf

////////////////////////////////////////////////////////////////////////////////
// TYPES

type TableType uint8

type SectionHeader struct {
	TableId TableType
	Length  uint16
}

type Section struct {
	SectionHeader
	*NITSection
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	PAT       TableType = 0x00
	CAT       TableType = 0x01
	PMT       TableType = 0x02
	NIT       TableType = 0x40
	NIT_OTHER TableType = 0x41
	SDT       TableType = 0x42
	SDT_OTHER TableType = 0x46
	BAT       TableType = 0x4A
	EIT       TableType = 0x4E
	EIT_OTHER TableType = 0x4F
	TDT       TableType = 0x70
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewSection(r io.Reader) (*Section, error) {
	this := new(Section)

	// Read table type and length
	if err := binary.Read(r, binary.LittleEndian, &this.SectionHeader); err != nil {
		return nil, err
	} else {
		this.Length = this.Length & 0x0FFF
	}

	// Read buffer of data
	data := make([]byte, this.Length)
	if _, err := r.Read(data); err != nil {
		return nil, err
	}

	// Parse data
	switch this.TableId {
	case NIT, NIT_OTHER:
		if nit, err := NewNITSection(bytes.NewReader(data)); err != nil {
			return nil, err
		} else {
			this.NITSection = nit
		}
	}

	// Success
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Section) String() string {
	str := "<dvb.section"
	str += " table_id=" + fmt.Sprint(this.TableId)
	switch this.TableId {
	case NIT, NIT_OTHER:
		str += " " + fmt.Sprint(this.NITSection)
	default:
		str += " length=" + fmt.Sprint(this.Length)
	}
	return str + ">"
}

func (f TableType) String() string {
	switch f {
	case PAT:
		return "PAT"
	case CAT:
		return "CAT"
	case PMT:
		return "PMT"
	case NIT:
		return "NIT"
	case NIT_OTHER:
		return "NIT_OTHER"
	case SDT:
		return "SDT"
	case SDT_OTHER:
		return "SDT_OTHER"
	case BAT:
		return "BAT"
	case EIT:
		return "EIT"
	case EIT_OTHER:
		return "EIT_OTHER"
	case TDT:
		return "TDT"
	default:
		return fmt.Sprintf("0x%02X", uint8(f))
	}
}
