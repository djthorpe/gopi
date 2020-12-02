package codec

import (
	"context"

	gopi "github.com/djthorpe/gopi/v3"
)

/*
Code for both the NEC32 protocol and the Legacy AppleTV protocols

The NEC IR transmission protocol uses pulse distance encoding
of the message bits. Each pulse burst (mark – RC transmitter ON)
is 562.5µs in length, at a carrier frequency of 38kHz (26.3µs).
Logical bits are transmitted as follows:

  * Logical '0' – a 562.5µs pulse burst followed by a 562.5µs space,
    with a total transmit time of 1.125ms
  * Logical '1' – a 562.5µs pulse burst followed by a 1.6875ms space,
    with a total transmit time of 2.25ms

When a key is pressed on the remote controller, the message transmitted
consists of the following, in order:

 * A 9ms leading pulse burst (16 times the pulse burst length used for
   a logical data bit)
 * A 4.5ms space
 * The 8-bit address for the receiving device
 * The 8-bit logical inverse of the address
 * The 8-bit command
 * The 8-bit logical inverse of the command
 * A final 562.5µs pulse burst to signify the end of message transmission.

REPEAT CODES

If the key on the remote controller is kept depressed, a repeat code will
be issued, typically around 40ms after the pulse burst that signified the
end of the message. A repeat code will continue to be sent out at 108ms
intervals, until the key is finally released. The repeat code consists of
the following, in order:

 * A 9ms leading pulse burst
 * A 2.25ms space
 * A 562.5µs pulse burst to mark the end of the space (and hence end
   of the transmitted repeat code).


Apple Remote IR Code

This document covers the old white apple remote.
Carrier: ~38KHz
Start: 9ms high, 4.5ms low
Pulse width: ~0.58ms (~853Hz)
Uses pulse-distance modulation.

Bit encoding:
  0: 1 pulse-width high, 1 pulse-width low
  1: 1 pulse-width high, 3 pulse-widths low
  4 octets are transmitted, LSB first.
  First two octets in normal transmission are 0x77 0xE1. (Different for pair command, which is 0x07 0xE1) Third octet is command. Fourth octet is remote ID.
  One of these bits is used for the low-battery indication. I haven't yet identified which one.

Example codes (in transmission order):

01110111 11100001 01000000 11101011 MENU
01110111 11100001 10110000 11101011 VOL-
01110111 11100001 11010000 11101011 VOL+
01110111 11100001 00100000 11101011 PLAY
01110111 11100001 11100000 11101011 NEXT
01110111 11100001 00010000 11101011 PREV
01110111 11100001 01101000 11101011 ??? (MENU+VOLUP)
00000111 11100001 11000000 11101011 PAIR (MENU+NEXT)

Normal Command format: 000XXXXP, where XXXX is the command and P is a parity bit(even parity).

Commands (MSB First):

0x01: 000 0001 0 MENU
0x02: 000 0010 0 PLAY
0x03: 000 0011 1 NEXT
0x04: 000 0100 0 PREV
0x05: 000 0101 1 VOL+
0x06: 000 0110 1 VOL-
0x0B: 000 1011 0 ???? (MENU+VOLUP)
0x0C: 000 1100 1 ???? (MENU+VOLDOWN)
Special Commands:

MENU+NEXT : Pair
MENU+PLAY : Increment Remote Code
MENU+PREV : ???

As far as I can tell, there are 8 normal signals that the remote can emit,
and 3 special signals---for a total of 11 signals.

Different remote
01110111 11100001 00100000 01010000 PLAY
00000111 11100001 11000000 01010000 PAIR
00000111 11100001 10100000 01010000 ??? (MENU+PREV)
00000111 11100001 01000000 11010000 CHANGE CODE (MENU+PLAY)
01110111 11100001 10011000 00110000 ??? (MENU+VOLDOWN)

Reference:
  http://techdocs.altium.com/display/FPGA/NEC+Infrared+Transmission+Protocol
  https://gist.github.com/darconeous/4437f79a34e3b6441628

*/

////////////////////////////////////////////////////////////////////////////////
// TYPES

type NEC struct {
	codec  gopi.InputDevice
	length uint
	state  necstate
	bits   []bool
	value  uint32
}

type necstate uint

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	STATE_NEC_HEADER_PULSE necstate = iota
	STATE_NEC_HEADER_SPACE
	STATE_NEC_PULSE
	STATE_NEC_SPACE
	STATE_NEC_REPEAT
	STATE_NEC_TRAILER
	STATE_NEC_KEYPRESS
	STATE_NEC_KEYREPEAT
)

const (
	NEC_TOLERANCE = 25     // 25% tolerance on values
	APPLETV_CODE  = 0x77E1 // The device code used by the AppleTV
)

var (
	NEC_HEADER_PULSE        = NewMarkSpace(gopi.LIRC_TYPE_PULSE, 9000, NEC_TOLERANCE) // 9ms
	NEC_HEADER_SPACE        = NewMarkSpace(gopi.LIRC_TYPE_SPACE, 4500, NEC_TOLERANCE) // 4.5ms
	NEV_HEADER_REPEAT_SPACE = NewMarkSpace(gopi.LIRC_TYPE_SPACE, 2250, NEC_TOLERANCE) // 2.25ms

	NEC_BIT_PULSE  = NewMarkSpace(gopi.LIRC_TYPE_PULSE, 563, NEC_TOLERANCE)  // 650ns
	NEC_ONE_SPACE  = NewMarkSpace(gopi.LIRC_TYPE_SPACE, 1688, NEC_TOLERANCE) // 1.6ms
	NEC_ZERO_SPACE = NewMarkSpace(gopi.LIRC_TYPE_SPACE, 563, NEC_TOLERANCE)  // 500ns
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewNEC(codec gopi.InputDevice) *NEC {
	this := new(NEC)
	if length := bitLengthForCodec(codec); length == 0 {
		return nil
	} else {
		this.codec = codec
		this.length = length
	}

	// Return success
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func bitLengthForCodec(codec gopi.InputDevice) uint {
	switch codec {
	case gopi.INPUT_DEVICE_NEC_32:
		return 32
	case gopi.INPUT_DEVICE_APPLETV_32:
		return 32
	default:
		return 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *NEC) Run(ctx context.Context, publisher gopi.Publisher) error {
	// Subscribe to LIRCEvent messages
	ch := publisher.Subscribe()
	defer publisher.Unsubscribe(ch)

	// Process LIRCEvent messages
	for {
		select {
		case evt := <-ch:
			if lircevent, ok := evt.(gopi.LIRCEvent); ok {
				this.Process(lircevent, publisher)
			}
		case <-ctx.Done():
			return nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *NEC) Process(evt gopi.LIRCEvent, publisher gopi.Publisher) {
	// Set new state
	this.state = nextState(evt, this.state, func(bit bool) bool {
		this.bits = append(this.bits, bit)
		return len(this.bits) == int(this.length)
	})

	// Determine any action from the new state
	action := gopi.INPUT_EVENT_NONE
	switch this.state {
	case STATE_NEC_KEYPRESS:
		this.value = valueFromBits(this.bits)
		action = gopi.INPUT_EVENT_KEYPRESS
		this.state = STATE_NEC_HEADER_PULSE
	case STATE_NEC_KEYREPEAT:
		action = gopi.INPUT_EVENT_KEYREPEAT
		this.state = STATE_NEC_HEADER_PULSE
	}

	// Perform the action
	if action != gopi.INPUT_EVENT_NONE && this.value > 0 {
		if code := codeForCodec(this.codec, this.value); code != 0 {
			publisher.Emit(&CodecEvent{action, this.codec, uint32(code)}, true)
		}
	}

	// Reset state
	if this.state == STATE_NEC_HEADER_PULSE {
		this.bits = nil
	}
}

func valueFromBits(bits []bool) uint32 {
	value := uint32(0)
	for i := 0; i < len(bits); i++ {
		value <<= 1
		if bits[i] {
			value |= 1
		}
	}
	return value
}

func codeForCodec(codec gopi.InputDevice, value uint32) uint32 {
	switch codec {
	case gopi.INPUT_DEVICE_APPLETV_32:
		if addr := value & 0xFFFF0000 >> 16; addr != APPLETV_CODE {
			return 0
		}
		return (value & 0xFF00 >> 8)
	case gopi.INPUT_DEVICE_NEC_32:
		if value == 0 {
			return 0
		} else if addr := value & 0xFFFF0000 >> 16; addr == APPLETV_CODE {
			return 0
		}
		// Check to make sure scancode and reverse of scancode match
		scancode1 := value & 0xFF00FF00 >> 8
		scancode2 := value & 0x00FF00FF
		if scancode1 != scancode2^0xFF00FF {
			return 0
		}
		return (scancode2 & 0xFF) | (scancode2&0xFF0000)>>8
	default:
		return 0
	}
}

func nextState(evt gopi.LIRCEvent, state necstate, fn func(bit bool) bool) necstate {
	switch state {
	case STATE_NEC_HEADER_PULSE:
		if NEC_HEADER_PULSE.Matches(evt) {
			return STATE_NEC_HEADER_SPACE
		}
	case STATE_NEC_HEADER_SPACE:
		if NEC_HEADER_SPACE.Matches(evt) {
			return STATE_NEC_PULSE
		} else if NEV_HEADER_REPEAT_SPACE.Matches(evt) {
			return STATE_NEC_REPEAT
		}
	case STATE_NEC_REPEAT:
		if NEC_BIT_PULSE.Matches(evt) {
			return STATE_NEC_KEYREPEAT
		}
	case STATE_NEC_PULSE:
		if NEC_BIT_PULSE.Matches(evt) {
			return STATE_NEC_SPACE
		}
	case STATE_NEC_SPACE:
		if NEC_ONE_SPACE.Matches(evt) {
			if fn(true) {
				return STATE_NEC_TRAILER
			} else {
				return STATE_NEC_PULSE
			}
		} else if NEC_ZERO_SPACE.Matches(evt) {
			if fn(false) {
				return STATE_NEC_TRAILER
			} else {
				return STATE_NEC_PULSE
			}
		}
	case STATE_NEC_TRAILER:
		if NEC_BIT_PULSE.Matches(evt) {
			// Key has been pressed
			return STATE_NEC_KEYPRESS
		}
	}

	// By default, reset to start state
	return STATE_NEC_HEADER_PULSE
}
