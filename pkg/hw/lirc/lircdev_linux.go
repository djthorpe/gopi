// +build linux

package lirc

import (
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type lircdev struct {
	sync.Mutex

	dev                            *os.File
	features                       linux.LIRCFeature
	send, recv                     bool
	recv_mode, send_mode           gopi.LIRCMode
	recv_dutycycle, send_dutycycle uint32
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewDevice(path string, features linux.LIRCFeature) (*lircdev, error) {
	mode := linux.LIRCMode(0)
	this := new(lircdev)

	if features&linux.LIRC_CAN_REC_MASK > 0 {
		mode |= linux.LIRC_MODE_RCV
		this.recv = true
	}
	if features&linux.LIRC_CAN_SEND_MASK > 0 {
		mode |= linux.LIRC_MODE_SEND
		this.send = true
	}
	if mode == 0 {
		return nil, fmt.Errorf("Device can neither send nor receive")
	}
	if fh, err := linux.LIRCOpenDevice(path, mode); err != nil {
		return nil, err
	} else {
		this.features = features
		this.dev = fh
	}

	// Set modes
	if this.send {
		if mode, err := linux.LIRCSendMode(this.Fd()); err != nil {
			this.dev.Close()
			return nil, err
		} else {
			this.send_mode = mode
		}
	}
	if this.recv {
		if mode, err := linux.LIRCRcvMode(this.Fd()); err != nil {
			this.dev.Close()
			return nil, err
		} else {
			this.recv_mode = mode
		}
	}

	// Return success
	return this, nil
}

func (this *lircdev) Close() error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var result error

	// Close device
	if this.dev != nil {
		result = this.dev.Close()
		this.dev = nil
	}

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *lircdev) String() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	str := "<lircdev"

	if this.dev != nil {
		str += " fd=" + fmt.Sprint(this.dev.Fd())
		str += " name=" + strconv.Quote(this.dev.Name())
	}
	if this.features != 0 {
		str += " features=" + fmt.Sprint(this.features)
	}
	if this.send {
		str += " send_mode=" + fmt.Sprint(this.send_mode)
	}
	if this.recv {
		str += " recv_mode=" + fmt.Sprint(this.recv_mode)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *lircdev) Fd() uintptr {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev != nil {
		return this.dev.Fd()
	} else {
		return 0
	}
}

func (this *lircdev) Name() string {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.dev != nil {
		return this.dev.Name()
	} else {
		return ""
	}
}

////////////////////////////////////////////////////////////////////////////////
// READ FROM DEVICE

func (this *lircdev) ReadEvent(uintptr, gopi.FilePollFlags) (gopi.LIRCEvent, error) {
	var value uint32

	if err := binary.Read(this.dev, binary.LittleEndian, &value); err != nil {
		return nil, err
	} else if err != nil {
		return nil, err
	} else if evt := NewEvent(this.dev.Name(), this.recv_mode, value); evt == nil {
		return nil, nil
	} else {
		return evt, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// WRITE TO DEVICE (IN SEND MODE)

// Send Pulse Mode, values are in milliseconds
func (this *lircdev) PulseSend(values []uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check device can send
	if this.send == false {
		return gopi.ErrOutOfOrder.WithPrefix("PulseSend")
	}
	// Check for odd number of values
	if len(values) == 0 || len(values)%2 == 0 {
		return gopi.ErrBadParameter.WithPrefix("PulseSend")
	}
	// Set send mode
	if this.send_mode != gopi.LIRC_MODE_PULSE {
		if err := this.SetSendMode(gopi.LIRC_MODE_PULSE); err != nil {
			return err
		}
	}
	// Send data
	if err := binary.Write(this.dev, binary.LittleEndian, values); err != nil {
		return err
	}
	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SEND AND RECEIVE MODE

func (this *lircdev) RcvMode() gopi.LIRCMode {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.recv {
		return this.recv_mode
	} else {
		return gopi.LIRC_MODE_NONE
	}
}

func (this *lircdev) SendMode() gopi.LIRCMode {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.send {
		return this.send_mode
	} else {
		return gopi.LIRC_MODE_NONE
	}
}

func (this *lircdev) SetRcvMode(mode gopi.LIRCMode) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.recv == false {
		return gopi.ErrOutOfOrder.WithPrefix("SetRcvMode")
	}

	switch mode {
	case gopi.LIRC_MODE_RAW:
		if this.features&linux.LIRC_CAN_REC_RAW == 0 {
			return gopi.ErrNotImplemented.WithPrefix("SetRcvMode")
		}
	case gopi.LIRC_MODE_PULSE:
		if this.features&linux.LIRC_CAN_REC_RAW == 0 {
			return gopi.ErrNotImplemented.WithPrefix("SetRcvMode")
		}
	case gopi.LIRC_MODE_MODE2:
		if this.features&linux.LIRC_CAN_REC_MODE2 == 0 {
			return gopi.ErrNotImplemented.WithPrefix("SetRcvMode")
		}
	case gopi.LIRC_MODE_LIRCCODE:
		if this.features&linux.LIRC_CAN_REC_LIRCCODE == 0 {
			return gopi.ErrNotImplemented.WithPrefix("SetRcvMode")
		}
	default:
		return gopi.ErrNotImplemented.WithPrefix("SetRcvMode")
	}

	// Set mode
	if err := linux.LIRCSetRcvMode(this.dev.Fd(), mode); err != nil {
		return err
	} else if mode_, err := linux.LIRCRcvMode(this.dev.Fd()); err != nil {
		return err
	} else if mode != mode_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("SetRcvMode")
	} else {
		this.recv_mode = mode
	}

	// Success
	return nil
}

func (this *lircdev) SetSendMode(mode gopi.LIRCMode) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.send == false {
		return gopi.ErrOutOfOrder.WithPrefix("SetSendMode")
	}
	switch mode {
	case gopi.LIRC_MODE_RAW:
		if this.features&linux.LIRC_CAN_SEND_RAW == 0 {
			return gopi.ErrNotImplemented.WithPrefix("SetSendMode")
		}
	case gopi.LIRC_MODE_PULSE:
		if this.features&linux.LIRC_CAN_SEND_PULSE == 0 {
			return gopi.ErrNotImplemented.WithPrefix("SetSendMode")
		}
	case gopi.LIRC_MODE_MODE2:
		if this.features&linux.LIRC_CAN_SEND_MODE2 == 0 {
			return gopi.ErrNotImplemented.WithPrefix("SetSendMode")
		}
	case gopi.LIRC_MODE_LIRCCODE:
		if this.features&linux.LIRC_CAN_SEND_LIRCCODE == 0 {
			return gopi.ErrNotImplemented.WithPrefix("SetSendMode")
		}
	default:
		return gopi.ErrNotImplemented.WithPrefix("SetSendMode")
	}

	// Set mode
	if err := linux.LIRCSetSendMode(this.dev.Fd(), mode); err != nil {
		return err
	} else if mode_, err := linux.LIRCSendMode(this.dev.Fd()); err != nil {
		return err
	} else if mode != mode_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("SetSendMode")
	} else {
		this.send_mode = mode
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// SEND AND RECEIVE DUTY CYCLE

func (this *lircdev) SendDutyCycle() uint32 {
	return this.send_dutycycle
}

func (this *lircdev) RcvDutyCycle() uint32 {
	return this.recv_dutycycle
}

func (this *lircdev) SetSendDutyCycle(value uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.send == false {
		return gopi.ErrOutOfOrder.WithPrefix("SetSendDutyCycle")
	}
	if this.features&linux.LIRC_CAN_SET_SEND_DUTY_CYCLE == 0 {
		return gopi.ErrNotImplemented.WithPrefix("SetSendDutyCycle")
	}
	if value < 1 || value > 99 {
		return gopi.ErrBadParameter.WithPrefix("SetSendDutyCycle")
	}
	if err := linux.LIRCSetSendDutyCycle(this.dev.Fd(), value); err != nil {
		return err
	} else {
		this.send_dutycycle = value
	}

	// Success
	return nil
}

func (this *lircdev) SetRcvDutyCycle(value uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.recv == false {
		return gopi.ErrOutOfOrder.WithPrefix("SetRcvDutyCycle")
	}
	if this.features&linux.LIRC_CAN_SET_REC_DUTY_CYCLE == 0 {
		return gopi.ErrNotImplemented.WithPrefix("SetRcvDutyCycle")
	}
	if value < 1 || value > 99 {
		return gopi.ErrBadParameter.WithPrefix("SetRcvDutyCycle")
	}
	if err := linux.LIRCSetRcvDutyCycle(this.dev.Fd(), value); err != nil {
		return err
	} else {
		this.recv_dutycycle = value
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RECEIVE TIMEOUTS

func (this *lircdev) SetRcvTimeout(micros uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.recv == false {
		return gopi.ErrOutOfOrder.WithPrefix("SetRcvTimeout")
	}
	if this.features&linux.LIRC_CAN_SET_REC_TIMEOUT == 0 {
		return gopi.ErrNotImplemented.WithPrefix("SetRcvTimeout")
	}
	if min, max, err := linux.LIRCMinMaxTimeoutMicros(this.Fd()); err != nil {
		return err
	} else if micros != 0 && (micros < min || micros > max) {
		return gopi.ErrBadParameter.WithPrefix("SetRcvTimeout")
	}
	if err := linux.LIRCSetRcvTimeoutMicros(this.dev.Fd(), micros); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *lircdev) SetRcvTimeoutReports(enable bool) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.recv == false {
		return gopi.ErrOutOfOrder.WithPrefix("SetRcvTimeoutReports")
	}
	if this.features&linux.LIRC_CAN_SET_REC_TIMEOUT == 0 {
		return gopi.ErrNotImplemented.WithPrefix("SetRcvTimeoutReports")
	}
	if err := linux.LIRCSetRcvTimeoutReports(this.dev.Fd(), enable); err != nil {
		return err
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RECEIVE CARRIER

func (this *lircdev) SetRcvCarrierHz(value uint32) error {
	return gopi.ErrNotImplemented
}

func (this *lircdev) SetRcvCarrierRangeHz(min uint32, max uint32) error {
	return gopi.ErrNotImplemented
}

////////////////////////////////////////////////////////////////////////////////
// RECEIVE RESOLUTION

func (this *lircdev) RcvResolutionMicros() (uint32, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.recv == false {
		return 0, gopi.ErrOutOfOrder.WithPrefix("GetRcvResolutionMicros")
	}
	if this.features&linux.LIRC_CAN_GET_REC_RESOLUTION == 0 {
		return 0, gopi.ErrNotImplemented.WithPrefix("GetRcvResolutionMicros")
	}
	if resolution, err := linux.LIRCRcvResolutionMicros(this.dev.Fd()); err != nil {
		return 0, err
	} else {
		return resolution, nil
	}
}
