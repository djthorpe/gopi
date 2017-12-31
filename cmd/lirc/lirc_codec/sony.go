package lirc_codec

import (
	"github.com/djthorpe/gopi"
)

const (
	SONY_TOLERANCE           = 35
	SONY_HEADER_PULSE uint32 = 2499
	SONY_HEADER_SPACE uint32 = 526
	SONY_ONE_PULSE    uint32 = 1277
	SONY_ZERO_PULSE   uint32 = 676
	SONY_SPACE_BIT    uint32 = 523
	SONY_PULSE_TRAIL  uint32 = 676
	SONY_SPACE_REPEAT uint32 = 24534
	SONY_SPACE_TRAIL1 uint32 = 19200 // probably 'repeat'
	SONY_SPACE_TRAIL2 uint32 = 27500 // probably 'end of message'
)

type SonyDecoder struct {
	state                              uint
	header_pulse_min, header_pulse_max uint32
	header_space_min, header_space_max uint32
	one_pulse_min, one_pulse_max       uint32
	zero_pulse_min, zero_pulse_max     uint32
	space_bit_min, space_bit_max       uint32
	space_repeat_min, space_repeat_max uint32
	log                                gopi.Logger
	value                              uint64
	bits                               uint
}

func NewSonyDecoder(log gopi.Logger) *SonyDecoder {
	this := new(SonyDecoder)
	this.log = log
	this.header_pulse_min, this.header_pulse_max = lengthWithTolerance(SONY_HEADER_PULSE, SONY_TOLERANCE)
	this.header_space_min, this.header_space_max = lengthWithTolerance(SONY_HEADER_SPACE, SONY_TOLERANCE)
	this.one_pulse_min, this.one_pulse_max = lengthWithTolerance(SONY_ONE_PULSE, SONY_TOLERANCE)
	this.zero_pulse_min, this.zero_pulse_max = lengthWithTolerance(SONY_ZERO_PULSE, SONY_TOLERANCE)
	this.space_bit_min, this.space_bit_max = lengthWithTolerance(SONY_SPACE_BIT, SONY_TOLERANCE)
	this.space_repeat_min, this.space_repeat_max = lengthWithTolerance(SONY_SPACE_REPEAT, SONY_TOLERANCE)
	return this
}

func (this *SonyDecoder) Receive(e gopi.LIRCEvent) {

	// Reset value
	if this.state == 0 {
		this.value = 0
		this.bits = 0
	}

	switch this.state {
	case 0:
		if isPulse(e, this.header_pulse_min, this.header_pulse_max) == true {
			this.state = 1
		} else {
			this.state = 0
		}
	case 1:
		if isSpace(e, this.header_space_min, this.header_space_max) == true {
			this.state = 2
		} else {
			this.state = 0
		}
	case 2:
		if isPulse(e, this.one_pulse_min, this.one_pulse_max) {
			this.value = this.value | 1
			this.state = 3
		} else if isPulse(e, this.zero_pulse_min, this.zero_pulse_max) {
			this.value = this.value | 0
			this.state = 3
		} else {
			this.state = 0
		}
	case 3:
		if isSpace(e, this.space_bit_min, this.space_bit_max) {
			this.value = this.value << 1
			this.bits = this.bits + 1
			if this.bits == 11 {
				this.log.Debug("EJECT=%X", this.value)
				this.state = 0
			} else {
				this.state = 2
			}
		} else if isSpace(e, this.space_repeat_min, this.space_repeat_max) {
			if this.bits == 11 {
				this.log.Debug("EJECT=%X", this.value)
			}
			this.state = 0
		} else {
			this.state = 0
		}
	default:
		this.log.Debug("State=%v - reseting to state 0", this.state)
		this.state = 0
	}
}

func isPulse(e gopi.LIRCEvent, min, max uint32) bool {
	if e.Type() != gopi.LIRC_TYPE_PULSE {
		return false
	}
	return isLength(e, min, max)
}

func isSpace(e gopi.LIRCEvent, min, max uint32) bool {
	if e.Type() != gopi.LIRC_TYPE_SPACE {
		return false
	}
	return isLength(e, min, max)
}

func isLength(e gopi.LIRCEvent, min, max uint32) bool {
	v := e.Value()
	return v >= min && v <= max
}

func lengthWithTolerance(length uint32, tolerance float32) (uint32, uint32) {
	return length - uint32(float32(length)*tolerance/100.0), length + uint32(float32(length)*tolerance/100.0)
}
