package ircodec

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
)

/*
The Philips RC5 IR transmission protocol uses Manchester encoding of the message bits. Each
pulse burst (mark – RC transmitter ON) is 889us in length, at a carrier frequency of 36kHz (27.7us).
Logical bits are transmitted as follows:

  * Logical '0' – an 889us pulse burst followed by an 889us space, with a total
    transmit time of 1.778ms
  * Logical '1' – an 889us space followed by an 889us pulse burst, with a total
    transmit time of 1.778ms

When a key is pressed on the remote controller, the message frame transmitted consists of
the following 14 bits, in order:

  * Two Start bits (S1 and S2), both logical '1'.
  * A Toggle bit (T). This bit is inverted each time a key is released and pressed again.
  * The 5-bit address for the receiving device
  * The 6-bit command.

The address and command bits are each sent most significant bit first. The Toggle bit (T) is
used by the receiver to distinguish weather the key has been pressed repeatedly, or weather
it is being held depressed. As long as the key on the remote controller is kept depressed,
the message frame will be repeated every 114ms. The Toggle bit will retain the same logic
level during all of these repeated message frames. It is up to the receiver software to interpret
this auto-repeat feature of the protocol.

Reference:
  https://techdocs.altium.com//display/FPGA/Philips+RC5+Infrared+Transmission+Protocol

*/

////////////////////////////////////////////////////////////////////////////////
// TYPES

type RC5Codec struct {
	codec  CodecType
	length uint
	state  rc5state
	bits   []bool
	toggle bool
}

type rc5state uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	RC5_EXPECT_FIRST_PULSE rc5state = iota
	RC5_EXPECT_PULSE
	RC5_EXPECT_SPACE
)

const (
	RC5_TOLERANCE = 35 // 35% tolerance on values
)

var (
	RC5_LONG_PULSE  = NewMarkSpace(gopi.LIRC_TYPE_PULSE, 1778, TOLERANCE)
	RC5_LONG_SPACE  = NewMarkSpace(gopi.LIRC_TYPE_SPACE, 1778, TOLERANCE)
	RC5_SHORT_PULSE = NewMarkSpace(gopi.LIRC_TYPE_PULSE, 889, TOLERANCE)
	RC5_SHORT_SPACE = NewMarkSpace(gopi.LIRC_TYPE_SPACE, 889, TOLERANCE)
	RC5_TIMEOUT     = NewMarkSpace(gopi.LIRC_TYPE_TIMEOUT, 9000, TOLERANCE)
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewRC5(codec CodecType) *RC5Codec {
	this := new(RC5Codec)
	if codec != CODEC_RC5_14 {
		return nil
	} else {
		this.codec = codec
		this.length = 14
	}

	// Reset state
	this.Reset()

	// Return success
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *RC5Codec) Type() CodecType {
	return this.codec
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *RC5Codec) Reset() {
	this.state = RC5_EXPECT_FIRST_PULSE
	this.bits = nil
}

func (this *RC5Codec) Process(evt gopi.LIRCEvent) {
	switch this.state {
	case RC5_EXPECT_FIRST_PULSE:
		if RC5_SHORT_PULSE.Matches(evt) {
			this.rc5Eject(false, true) // 01
			this.state = RC5_EXPECT_SPACE
		} else if RC5_LONG_PULSE.Matches(evt) {
			this.rc5Eject(false, true, true) // 011
			this.state = RC5_EXPECT_SPACE
		} else {
			this.Reset()
		}
	case RC5_EXPECT_PULSE:
		if RC5_LONG_PULSE.Matches(evt) {
			this.rc5Eject(true, true) // 11
			this.state = RC5_EXPECT_SPACE
		} else if RC5_SHORT_PULSE.Matches(evt) {
			this.rc5Eject(true) // 1
			this.state = RC5_EXPECT_SPACE
		} else {
			this.Reset()

		}
	case RC5_EXPECT_SPACE:
		if RC5_LONG_SPACE.Matches(evt) {
			this.rc5Eject(false, false) // 00
			this.state = RC5_EXPECT_PULSE
			break
		} else if RC5_SHORT_SPACE.Matches(evt) {
			this.rc5Eject(false) // 0
			this.state = RC5_EXPECT_PULSE
			break
		} else if RC5_TIMEOUT.GreaterThan(evt) {
			this.rc5Eject(false) // 0
		}
		fallthrough
	default:
		this.Reset()
	}
}

func (this *RC5Codec) rc5Eject(bits ...bool) {
	// Append bits
	this.bits = append(this.bits, bits...)
	if len(this.bits) < int(this.length)*2 {
		return
	}

	// Manchester decoding
	value := uint16(0)
	for i, j := 0, 1; i < int(this.length*2); i, j = i+2, j+2 {
		x, y := this.bits[i], this.bits[j]
		switch {
		case x && !y:
			value = value << 1
		case !x && y:
			value = value<<1 | 1
		}
	}

	// Stop bits should be 3
	if stop, toggle, addr, cmd := rc5Decode(value); stop == 0x03 {
		if toggle == this.toggle {
			fmt.Printf("%v addr=%v cmd=%v\n", gopi.INPUT_EVENT_KEYREPEAT, addr, cmd)
		} else {
			fmt.Printf("%v addr=%v cmd=%v\n", gopi.INPUT_EVENT_KEYPRESS, addr, cmd)
			this.toggle = toggle
		}
	}

	// Reset
	this.bits = nil
}

func rc5Decode(value uint16) (uint8, bool, uint8, uint8) {
	stop := uint8(value & 0x3000 >> 12)
	toggle := uint8(value & 0x0800 >> 11)
	addr := uint8(value & 0x07C0 >> 6)
	cmd := uint8(value & 0x003F >> 0)
	return stop, toggle != 0, addr, cmd
}
