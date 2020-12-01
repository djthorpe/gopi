package ircodec

import (
	"fmt"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

/*

The Sony remote control is based on the Pulse-Width signal coding scheme. The
code exists of 12 bits sent on a 40kHz carrier wave. The code starts with a
header of 2.4ms or 4 times T where T is 600µS. The header is followed by 7
command bits and 5 address bits.

The address and commands exists of logical ones and zeros. A logical one is formed by
a space of 600µS or 1T and a pulse of 1200 µS or 2T. A logical zero is formed by a
space of 600 µS and pulse of 600µS.

The space between 2 transmitted codes when a button is being pressed is 40mS

The bits are transmitted least significant bits first. The total length of a bitstream
is always 45ms.

References:
  http://users.telenet.be/davshomepage/sony.htm
  http://picprojects.org.uk/projects/sirc/sonysirc.pdf
  https://www.sbprojects.net/knowledge/ir/sirc.php

*/

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SonyCodec struct {
	codec    CodecType
	bits     uint
	state    state
	value    uint32
	duration uint32
	length   uint
	repeat   bool
}

type state uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// state
	STATE_EXPECT_HEADER_PULSE state = iota
	STATE_EXPECT_HEADER_SPACE
	STATE_EXPECT_BIT
	STATE_EXPECT_SPACE
	STATE_EXPECT_TRAIL
	STATE_EXPECT_REPEAT
)

const (
	TOLERANCE   = 35    // 35% tolerance on values
	TX_DURATION = 45000 // 45ms between each transmission
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	HEADER_PULSE  = NewMarkSpace(gopi.LIRC_TYPE_PULSE, 2400, TOLERANCE)
	ONEZERO_SPACE = NewMarkSpace(gopi.LIRC_TYPE_SPACE, 575, TOLERANCE)
	ONE_PULSE     = NewMarkSpace(gopi.LIRC_TYPE_PULSE, 1200, TOLERANCE)
	ZERO_PULSE    = NewMarkSpace(gopi.LIRC_TYPE_PULSE, 575, TOLERANCE)
	TRAIL_PULSE   = NewMarkSpace(gopi.LIRC_TYPE_PULSE, 1200, TOLERANCE)
	REPEAT_SPACE  = NewMarkSpace(gopi.LIRC_TYPE_SPACE, TX_DURATION, TOLERANCE)
)

var (
	timestamp = time.Now()
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewSony(codec CodecType) *SonyCodec {
	this := new(SonyCodec)
	if bits := bitsForCodecType(codec); bits == 0 {
		return nil
	} else {
		this.codec = codec
		this.bits = bits
	}

	// Reset state
	this.Reset()

	// Return success
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *SonyCodec) Type() CodecType {
	return this.codec
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *SonyCodec) Reset() {
	this.state = STATE_EXPECT_HEADER_PULSE
	this.value = 0
	this.length = 0
	this.duration = 0
	this.repeat = false
}

func (this *SonyCodec) Emit(value uint32, repeat bool) {
	if scancode, device, err := codeForCodec(this.codec, value); err != nil {
		fmt.Println("SONY EMIT", scancode, device, repeat)
	}
}

func (this *SonyCodec) Process(evt gopi.LIRCEvent) {
	// check mode and set duration
	if evt.Mode() != gopi.LIRC_MODE_MODE2 {
		return
	}
	duration := evt.Value().(uint32)

	switch this.state {
	case STATE_EXPECT_HEADER_PULSE:
		if HEADER_PULSE.Matches(evt) {
			this.state = STATE_EXPECT_SPACE
			this.duration += duration
		} else {
			this.Reset()
		}
	case STATE_EXPECT_SPACE:
		if ONEZERO_SPACE.Matches(evt) {
			this.value <<= 1
			this.state = STATE_EXPECT_BIT
			this.duration += duration
		} else {
			REPEAT_SPACE.Set(TX_DURATION-this.duration, TOLERANCE)
			if REPEAT_SPACE.Matches(evt) && this.length == this.bits {
				this.Emit(this.value, this.repeat)
				this.value = 0
				this.length = 0
				this.duration = 0
				this.repeat = true
				this.state = STATE_EXPECT_HEADER_PULSE
			} else {
				this.Reset()
			}
		}
	case STATE_EXPECT_BIT:
		if ONE_PULSE.Matches(evt) {
			this.value |= 1
			this.length += 1
			this.state = STATE_EXPECT_SPACE
			this.duration += duration
		} else if ZERO_PULSE.Matches(evt) {
			this.value |= 0
			this.length += 1
			this.state = STATE_EXPECT_SPACE
			this.duration += duration
		} else {
			this.Reset()
		}
	default:
		this.Reset()
	}
}

/*
////////////////////////////////////////////////////////////////////////////////
// SENDING

func (this *codec) Send(device uint32, scancode uint32, repeats uint) error {
	this.log.Debug2("<remotes.codec.sony>Send{ codec_type=%v device=0x%08X scancode=0x%08X repeats=%v }", this.codec_type, device, scancode, repeats)

	// Array of pulses
	pulses := make([]uint32, 0, 100)

	// Make bits and pulses
	if bits, err := bitsForCodec(this.codec_type, device, scancode); err != nil {
		return err
	} else {
		for j := uint(0); j < (repeats + 1); j++ {
			length := HEADER_PULSE.Value
			pulses = append(pulses, HEADER_PULSE.Value)

			// Send the bits
			for i := 0; i < len(bits); i++ {
				pulses = append(pulses, ONEZERO_SPACE.Value)
				length += ONEZERO_SPACE.Value
				if bits[i] {
					pulses = append(pulses, ONE_PULSE.Value)
					length += ONEZERO_SPACE.Value
				} else {
					pulses = append(pulses, ZERO_PULSE.Value)
					length += ZERO_PULSE.Value
				}
			}

			// If repeats then send trail
			if repeats > 0 && j < repeats {
				pulses = append(pulses, TX_DURATION-length)
			}
		}
	}

	// Perform the sending
	return this.lirc.PulseSend(pulses)
}
*/

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func bitsForCodecType(codec CodecType) uint {
	switch codec {
	case CODEC_SONY_12:
		return 12
	case CODEC_SONY_15:
		return 15
	case CODEC_SONY_20:
		return 20
	default:
		return 0
	}
}

func codeForCodec(codec CodecType, value uint32) (uint32, uint32, error) {
	switch codec {
	case CODEC_SONY_12:
		// 7 scancode bits and 5 device bits
		return (value & 0x0FE0) >> 5, (value & 0x001F), nil
	case CODEC_SONY_15:
		// 15 bit codes are similar, with 7 command bits and 8 device bits
		return (value & 0x7F00) >> 8, (value & 0xFF), nil
	case CODEC_SONY_20:
		// 20 bit codes have 7 command bits and 13 device bits
		return (value & 0xFE000) >> 13, (value & 0x1FFF), nil
	default:
		return 0, 0, gopi.ErrBadParameter
	}
}

func bitsForCodec(codec CodecType, device uint32, scancode uint32) ([]bool, error) {
	bits := make([]bool, 0, bitsForCodecType(codec))
	switch codec {
	case CODEC_SONY_12:
		// 7 scancode bits and 5 device bits
		bits = bitsAppend(bits, scancode, 7)
		bits = bitsAppend(bits, device, 5)
	case CODEC_SONY_15:
		// 7 scancode bits and 8 device bits
		bits = bitsAppend(bits, scancode, 7)
		bits = bitsAppend(bits, device, 8)
	case CODEC_SONY_20:
		// 7 scancode bits and 13 device bits
		bits = bitsAppend(bits, scancode, 7)
		bits = bitsAppend(bits, device, 13)
	default:
		return nil, gopi.ErrBadParameter
	}
	return bits, nil
}

func bitsAppend(array []bool, value uint32, length uint) []bool {
	mask := uint32(1) << (length - 1)
	for i := uint(0); i < length; i++ {
		array = append(array, value&mask != 0)
		mask >>= 1
	}
	return array
}

func (s state) String() string {
	switch s {
	case STATE_EXPECT_HEADER_PULSE:
		return "STATE_EXPECT_HEADER_PULSE"
	case STATE_EXPECT_HEADER_SPACE:
		return "STATE_EXPECT_HEADER_SPACE"
	case STATE_EXPECT_BIT:
		return "STATE_EXPECT_BIT"
	case STATE_EXPECT_SPACE:
		return "STATE_EXPECT_SPACE"
	case STATE_EXPECT_TRAIL:
		return "STATE_EXPECT_TRAIL"
	case STATE_EXPECT_REPEAT:
		return "STATE_EXPECT_REPEAT"
	default:
		return "[?? Invalid state]"
	}
}
