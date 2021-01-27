package ts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/djthorpe/gopi/v3"
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
	PATSection
	PMTSection
	NITSection
	crc uint32
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

const (
	SECTION_BUFFER_SIZE = 4096
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewSection(r io.Reader, data []byte) (*Section, error) {
	this := new(Section)

	// Read table type and length
	if err := binary.Read(r, binary.BigEndian, &this.SectionHeader); err != nil {
		return nil, err
	} else {
		this.Length &= 0x0FFF
	}

	// Read buffer
	if n, err := r.Read(data); err != nil {
		return nil, err
	} else if n != int(this.Length) {
		return nil, gopi.ErrUnexpectedResponse.WithPrefix("NewSection")
	}

	// Read section data
	r2 := bytes.NewReader(data)
	switch this.TableId {
	case PAT:
		if err := this.PATSection.Read(r2, int(this.Length)-4); err != nil {
			return nil, err
		}
	case PMT:
		if err := this.PMTSection.Read(r2, int(this.Length)-4); err != nil {
			return nil, err
		}
	case NIT, NIT_OTHER:
		if err := this.NITSection.Read(r2, int(this.Length)-4); err != nil {
			return nil, err
		}
	}

	// Read CRC
	if err := binary.Read(r2, binary.BigEndian, &this.crc); err != nil {
		return nil, err
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
	case PAT:
		str += " " + fmt.Sprint(this.PATSection)
	case PMT:
		str += " " + fmt.Sprint(this.PMTSection)
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
