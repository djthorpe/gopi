// +build linux

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2017
  All Rights Reserved

  Documentation http://djthorpe.github.io/gopi/
  For Licensing and Usage information, please see LICENSE.md
*/

package linux

import (
	"fmt"
	"os"
	"strings"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type LIRC struct {
	// Device path
	Device string

	// Filepoller
	FilePoll FilePollInterface
}

type lirc struct {
	// device
	dev *os.File

	// logger
	log gopi.Logger

	// poll
	filepoll FilePollInterface

	// mutex lock
	lock sync.Mutex

	// features
	features lirc_feature

	// modes
	rcv_mode, send_mode gopi.LIRCMode
}

type lirc_feature uint32

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// LIRC_DEV is the default device path
	LIRC_DEV = "/dev/lirc0"
	// LIRC_CARRIER_FREQUENCY is the default carrier frequency
	LIRC_CARRIER_FREQUENCY = 38000
	// LIRC_DUTY_CYCLE is the default duty cycle
	LIRC_DUTY_CYCLE = 50
)

const (
	LIRC_MODE2SEND uint32 = 0
	LIRC_MODE2REC  uint32 = 16
)

const (
	LIRC_CAN_SEND_RAW                 lirc_feature = lirc_feature(gopi.LIRC_MODE_RAW) << LIRC_MODE2SEND
	LIRC_CAN_SEND_PULSE               lirc_feature = lirc_feature(gopi.LIRC_MODE_PULSE) << LIRC_MODE2SEND
	LIRC_CAN_SEND_MODE2               lirc_feature = lirc_feature(gopi.LIRC_MODE_MODE2) << LIRC_MODE2SEND
	LIRC_CAN_SEND_LIRCCODE            lirc_feature = lirc_feature(gopi.LIRC_MODE_LIRCCODE) << LIRC_MODE2SEND
	LIRC_CAN_SEND_MASK                lirc_feature = 0x0000003F
	LIRC_CAN_SET_SEND_CARRIER         lirc_feature = 0x00000100
	LIRC_CAN_SET_SEND_DUTY_CYCLE      lirc_feature = 0x00000200
	LIRC_CAN_SET_TRANSMITTER_MASK     lirc_feature = 0x00000400
	LIRC_CAN_REC_RAW                  lirc_feature = lirc_feature(gopi.LIRC_MODE_RAW) << LIRC_MODE2REC
	LIRC_CAN_REC_PULSE                lirc_feature = lirc_feature(gopi.LIRC_MODE_PULSE) << LIRC_MODE2REC
	LIRC_CAN_REC_MODE2                lirc_feature = lirc_feature(gopi.LIRC_MODE_MODE2) << LIRC_MODE2REC
	LIRC_CAN_REC_LIRCCODE             lirc_feature = lirc_feature(gopi.LIRC_MODE_LIRCCODE) << LIRC_MODE2REC
	LIRC_CAN_REC_MASK                 lirc_feature = LIRC_CAN_SEND_MASK << LIRC_MODE2REC
	LIRC_CAN_SET_REC_CARRIER          lirc_feature = LIRC_CAN_SET_SEND_CARRIER << LIRC_MODE2REC
	LIRC_CAN_SET_REC_DUTY_CYCLE       lirc_feature = LIRC_CAN_SET_SEND_DUTY_CYCLE << LIRC_MODE2REC
	LIRC_CAN_SET_REC_DUTY_CYCLE_RANGE lirc_feature = 0x40000000
	LIRC_CAN_SET_REC_CARRIER_RANGE    lirc_feature = 0x80000000
	LIRC_CAN_GET_REC_RESOLUTION       lirc_feature = 0x20000000
	LIRC_CAN_SET_REC_TIMEOUT          lirc_feature = 0x10000000
	LIRC_CAN_SET_REC_FILTER           lirc_feature = 0x08000000
	/*
		LIRC_CAN_MEASURE_CARRIER          lirc_feature = 0x02000000
		LIRC_CAN_USE_WIDEBAND_RECEIVER    lirc_feature = 0x04000000
	*/
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open creates a new LIRC object, returns error if not possible
func (config LIRC) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.linux.LIRC.Open>{ device=%v }", config.Device)
	var err error

	if config.Device == "" {
		config.Device = LIRC_DEV
	}

	// create new driver
	this := new(lirc)
	this.log = log

	// File Poll module is required or else returns ErrBadParameter
	if config.FilePoll != nil {
		this.filepoll = config.FilePoll
	} else {
		return nil, gopi.ErrBadParameter
	}

	// Open the device
	this.dev, err = os.OpenFile(config.Device, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	// Get features
	if this.features, err = this.getFeatures(); err != nil {
		this.dev.Close()
		return nil, err
	}

	// Get modes
	if this.rcv_mode, err = this.getRcvMode(); err != nil {
		this.dev.Close()
		return nil, err
	}

	// return driver
	return this, nil
}

// Close connection
func (this *lirc) Close() error {
	this.log.Debug("<sys.hw.linux.LIRC.Close>{ }")

	if err := this.dev.Close(); err != nil {
		return err
	} else {
		this.dev = nil
	}

	// Blank out
	this.filepoll = nil

	return nil
}

////////////////////////////////////////////////////////////////////////////////
// GET AND SET PROPERTIES

func (this *lirc) RcvMode() gopi.LIRCMode {
	return this.rcv_mode
}

func (this *lirc) SendMode() gopi.LIRCMode {
	return this.send_mode
}

func (this *lirc) SetRcvMode(m gopi.LIRCMode) error {
	// Check to make sure feature is supported
	switch m {
	case gopi.LIRC_MODE_RAW:
		if this.features&LIRC_CAN_REC_RAW == 0 {
			return gopi.ErrNotImplemented
		}
	case gopi.LIRC_MODE_PULSE:
		if this.features&LIRC_CAN_REC_RAW == 0 {
			return gopi.ErrNotImplemented
		}
	case gopi.LIRC_MODE_MODE2:
		if this.features&LIRC_CAN_REC_MODE2 == 0 {
			return gopi.ErrNotImplemented
		}
	case gopi.LIRC_MODE_LIRCCODE:
		if this.features&LIRC_CAN_REC_LIRCCODE == 0 {
			return gopi.ErrNotImplemented
		}
	default:
		return gopi.ErrNotImplemented
	}
	if err := this.setRcvMode(m); err != nil {
		return err
	}
	this.rcv_mode = m
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *lirc) String() string {
	features := make([]string, 0)
	for i := uint(0); i < 32; i++ {
		mask := (lirc_feature(1) << i)
		if this.features&mask != 0 {
			features = append(features, fmt.Sprint(mask))
		}
	}
	return fmt.Sprintf("<sys.hw.linux.LIRC>{ features=%v rcv_mode=%v send_mode=%v }", strings.Join(features, ","), this.rcv_mode, this.send_mode)
}

func (f lirc_feature) String() string {
	switch f {
	case LIRC_CAN_SEND_RAW:
		return "LIRC_CAN_SEND_RAW"
	case LIRC_CAN_SEND_PULSE:
		return "LIRC_CAN_SEND_PULSE"
	case LIRC_CAN_SEND_MODE2:
		return "LIRC_CAN_SEND_MODE2"
	case LIRC_CAN_SEND_LIRCCODE:
		return "LIRC_CAN_SEND_LIRCCODE"
	case LIRC_CAN_SEND_MASK:
		return "LIRC_CAN_SEND_MASK"
	case LIRC_CAN_SET_SEND_CARRIER:
		return "LIRC_CAN_SET_SEND_CARRIER"
	case LIRC_CAN_SET_SEND_DUTY_CYCLE:
		return "LIRC_CAN_SET_SEND_DUTY_CYCLE"
	case LIRC_CAN_SET_TRANSMITTER_MASK:
		return "LIRC_CAN_SET_TRANSMITTER_MASK"
	case LIRC_CAN_REC_RAW:
		return "LIRC_CAN_REC_RAW"
	case LIRC_CAN_REC_PULSE:
		return "LIRC_CAN_REC_PULSE"
	case LIRC_CAN_REC_MODE2:
		return "LIRC_CAN_REC_MODE2"
	case LIRC_CAN_REC_LIRCCODE:
		return "LIRC_CAN_REC_LIRCCODE"
	case LIRC_CAN_REC_MASK:
		return "LIRC_CAN_REC_MASK"
	case LIRC_CAN_SET_REC_CARRIER:
		return "LIRC_CAN_SET_REC_CARRIER"
	case LIRC_CAN_SET_REC_DUTY_CYCLE:
		return "LIRC_CAN_SET_REC_DUTY_CYCLE"
	case LIRC_CAN_SET_REC_DUTY_CYCLE_RANGE:
		return "LIRC_CAN_SET_REC_DUTY_CYCLE_RANGE"
	case LIRC_CAN_SET_REC_CARRIER_RANGE:
		return "LIRC_CAN_SET_REC_CARRIER_RANGE"
	case LIRC_CAN_GET_REC_RESOLUTION:
		return "LIRC_CAN_GET_REC_RESOLUTION"
	case LIRC_CAN_SET_REC_TIMEOUT:
		return "LIRC_CAN_SET_REC_TIMEOUT"
	case LIRC_CAN_SET_REC_FILTER:
		return "LIRC_CAN_SET_REC_FILTER"
	/*
		case LIRC_CAN_MEASURE_CARRIER:
			return "LIRC_CAN_MEASURE_CARRIER"
		case LIRC_CAN_USE_WIDEBAND_RECEIVER:
			return "LIRC_CAN_USE_WIDEBAND_RECEIVER"
	*/
	default:
		return "[?? Invalid lirc_feature value]"
	}
}
