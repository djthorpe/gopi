package ts

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tag uint8

type DTable struct {
	Length uint16
	Rows   []DRow
}

type STable struct {
	Length uint16
	Rows   []SRow
}

type DRow struct {
	header struct {
		Tag
		Length uint8
	}
	data []byte
}

type SRow struct {
	header struct {
		StreamId  uint16
		NetworkId uint16
	}
	DTable
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	network_name                Tag = 0x40
	service_list                Tag = 0x41
	stuffing                    Tag = 0x42
	satellite_delivery_system   Tag = 0x43
	cable_delivery_system       Tag = 0x44
	bouquet_name                Tag = 0x47
	service                     Tag = 0x48
	country_availability        Tag = 0x49
	linkage                     Tag = 0x4A
	nvod_reference              Tag = 0x4B
	time_shifted_service        Tag = 0x4C
	short_event                 Tag = 0x4D
	extended_event              Tag = 0x4E
	time_shifted_event          Tag = 0x4F
	component                   Tag = 0x50
	mosaic                      Tag = 0x51
	stream_identifier           Tag = 0x52
	ca_identifier               Tag = 0x53
	content                     Tag = 0x54
	parental_rating             Tag = 0x55
	teletext                    Tag = 0x56
	telephone                   Tag = 0x57
	local_time_offset           Tag = 0x58
	subtitling                  Tag = 0x59
	terrestrial_delivery_system Tag = 0x5A
	multilingual_network_name   Tag = 0x5B
	multilingual_bouquet_name   Tag = 0x5C
	multilingual_service_name   Tag = 0x5D
	multilingual_component      Tag = 0x5E
	private_data_specifier      Tag = 0x5F
	service_move                Tag = 0x60
	short_smoothing_buffer      Tag = 0x61
	frequency_list              Tag = 0x62
	partial_transport_stream    Tag = 0x63
	data_broadcast              Tag = 0x64
	ca_system                   Tag = 0x65
	data_broadcast_id           Tag = 0x66
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (t *DTable) Read(r io.Reader) error {
	// Read length of table that follows in bytes
	if err := binary.Read(r, binary.BigEndian, &t.Length); err != nil {
		return err
	} else {
		// Top four bits are reserved for future use, length is 12 bits
		t.Length &= 0x0FFF
	}

	// Read rows until length is zero
	for length := t.Length; length > 0; {
		var row DRow
		if err := binary.Read(r, binary.LittleEndian, &row.header); err != nil {
			return err
		}
		row.data = make([]byte, int(row.header.Length))
		if _, err := r.Read(row.data); err != nil {
			return err
		}
		// Append row, decrement by 2 bytes and length of data
		t.Rows = append(t.Rows, row)
		length -= 2 + uint16(row.header.Length)
	}

	// Return success
	return nil
}

func (t *STable) Read(r io.Reader) error {
	// Read length of descriptor that follow in bytes
	if err := binary.Read(r, binary.BigEndian, &t.Length); err != nil {
		return err
	} else {
		// Top four bits are reserved for future use, length is 12 bits
		t.Length &= 0x0FFF
	}

	// Read rows until length is zero
	for length := t.Length; length > 0; {
		var row SRow
		if err := binary.Read(r, binary.LittleEndian, &row.header); err != nil {
			return err
		} else if err := row.DTable.Read(r); err != nil {
			return err
		}
		// Append row, decrement by 3 x uint16 and length of descriptor table
		t.Rows = append(t.Rows, row)
		length -= 6 + uint16(row.DTable.Length)
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (t DTable) String() string {
	str := "<dvb.descriptors"
	for _, row := range t.Rows {
		str += fmt.Sprint(" ", row)
	}
	return str + ">"
}

func (t STable) String() string {
	str := "<dvb.streams"
	for _, row := range t.Rows {
		str += fmt.Sprint(" ", row)
	}
	return str + ">"
}

func (r DRow) String() string {
	switch r.header.Tag {
	case bouquet_name, network_name:
		return fmt.Sprintf("<%v=%q>", r.header.Tag, string(r.data))
	default:
		return fmt.Sprintf("<%v=%v>", r.header.Tag, hex.EncodeToString(r.data))
	}
}

func (r SRow) String() string {
	return fmt.Sprintf("<stream_id=0x%04X network_id=0x%04X %v>", r.header.StreamId, r.header.NetworkId, r.DTable)
}

func (t Tag) String() string {
	switch t {
	case network_name:
		return "network_name"
	case service_list:
		return "service_list"
	case stuffing:
		return "stuffing"
	case satellite_delivery_system:
		return "satellite_delivery_system"
	case cable_delivery_system:
		return "cable_delivery_system"
	case bouquet_name:
		return "bouquet_name"
	case service:
		return "service"
	case country_availability:
		return "country_availability"
	case linkage:
		return "linkage"
	case nvod_reference:
		return "nvod_reference"
	case time_shifted_service:
		return "time_shifted_service"
	case short_event:
		return "short_event"
	case extended_event:
		return "extended_event"
	case time_shifted_event:
		return "time_shifted_event"
	case component:
		return "component"
	case mosaic:
		return "mosaic"
	case stream_identifier:
		return "stream_identifier"
	case ca_identifier:
		return "ca_identifier"
	case content:
		return "content"
	case parental_rating:
		return "parental_rating"
	case teletext:
		return "teletext"
	case telephone:
		return "telephone"
	case local_time_offset:
		return "local_time_offset"
	case subtitling:
		return "subtitling"
	case terrestrial_delivery_system:
		return "terrestrial_delivery_system"
	case multilingual_network_name:
		return "multilingual_network_name"
	case multilingual_bouquet_name:
		return "multilingual_bouquet_name"
	case multilingual_service_name:
		return "multilingual_service_name"
	case multilingual_component:
		return "multilingual_component"
	case private_data_specifier:
		return "private_data_specifier"
	case service_move:
		return "service_move"
	case short_smoothing_buffer:
		return "short_smoothing_buffer"
	case frequency_list:
		return "frequency_list"
	case partial_transport_stream:
		return "partial_transport_stream"
	case data_broadcast:
		return "data_broadcast"
	case ca_system:
		return "ca_system"
	case data_broadcast_id:
		return "data_broadcast_id"
	default:
		return fmt.Sprint("0x%02X", uint8(t))
	}

}
