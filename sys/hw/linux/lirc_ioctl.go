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
	"github.com/djthorpe/gopi"
)

/*
	#include <sys/ioctl.h>
	#include <linux/types.h>
	static int _LIRC_GET_FEATURES() { return _IOR('i', 0x00000000, __u32); }
	static int _LIRC_GET_SEND_MODE() { return _IOR('i', 0x00000001, __u32); }
	static int _LIRC_GET_REC_MODE() { return _IOR('i', 0x00000002, __u32); }
	static int _LIRC_GET_REC_RESOLUTION() { return _IOR('i', 0x00000007, __u32); }
	static int _LIRC_GET_MIN_TIMEOUT() { return _IOR('i', 0x00000008, __u32); }
 	static int _LIRC_GET_MAX_TIMEOUT() { return _IOR('i', 0x00000009, __u32); }
	// code length in bits, currently only for LIRC_MODE_LIRCCODE
	static int _LIRC_GET_LENGTH() { return _IOR('i', 0x0000000F, __u32); }
	static int _LIRC_SET_SEND_MODE() { return _IOW('i', 0x00000011, __u32); }
	static int _LIRC_SET_REC_MODE() { return _IOW('i', 0x00000012, __u32); }
	// Note: these can reset the according pulse_width
	static int _LIRC_SET_SEND_CARRIER() { return _IOW('i', 0x00000013, __u32); }
	static int _LIRC_SET_REC_CARRIER() { return _IOW('i', 0x00000014, __u32); }
	static int _LIRC_SET_SEND_DUTY_CYCLE() { return _IOW('i', 0x00000015, __u32); }
	static int _LIRC_SET_REC_DUTY_CYCLE() { return _IOW('i', 0x00000016, __u32); }
	static int _LIRC_SET_TRANSMITTER_MASK() { return _IOW('i', 0x00000017, __u32); }
	// When a timeout != 0 is set the driver will send a LIRC_MODE2_TIMEOUT data packet
	// otherwise LIRC_MODE2_TIMEOUT is never sent, timeout is disabled by default
	static int _LIRC_SET_REC_TIMEOUT() { return _IOW('i', 0x00000018, __u32); }
	// 1 enables, 0 disables timeout reports in MODE2
	static int _LIRC_SET_REC_TIMEOUT_REPORTS() { return _IOW('i', 0x00000019, __u32); }
	// if enabled from the next key press on the driver will send LIRC_MODE2_FREQUENCY packets
	static int _LIRC_SET_MEASURE_CARRIER_MODE() { return _IOW('i', 0x0000001d, __u32); }
	// to set a range use LIRC_SET_REC_CARRIER_RANGE with the lower bound first and later
	// LIRC_SET_REC_CARRIER with the upper bound
	static int _LIRC_SET_REC_CARRIER_RANGE() { return _IOW('i', 0x0000001f, __u32); }
	static int _LIRC_SET_WIDEBAND_RECEIVER() { return _IOW('i', 0x00000023, __u32); }
*/
import "C"

import (
	"os"
	"syscall"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	LIRC_GET_FEATURES             = uintptr(C._LIRC_GET_FEATURES())
	LIRC_GET_SEND_MODE            = uintptr(C._LIRC_GET_SEND_MODE())
	LIRC_GET_REC_MODE             = uintptr(C._LIRC_GET_REC_MODE())
	LIRC_GET_REC_RESOLUTION       = uintptr(C._LIRC_GET_REC_RESOLUTION())
	LIRC_GET_MIN_TIMEOUT          = uintptr(C._LIRC_GET_MIN_TIMEOUT())
	LIRC_GET_MAX_TIMEOUT          = uintptr(C._LIRC_GET_MAX_TIMEOUT())
	LIRC_GET_LENGTH               = uintptr(C._LIRC_GET_LENGTH())
	LIRC_SET_SEND_MODE            = uintptr(C._LIRC_SET_SEND_MODE())
	LIRC_SET_REC_MODE             = uintptr(C._LIRC_SET_REC_MODE())
	LIRC_SET_SEND_CARRIER         = uintptr(C._LIRC_SET_SEND_CARRIER())
	LIRC_SET_REC_CARRIER          = uintptr(C._LIRC_SET_REC_CARRIER())
	LIRC_SET_SEND_DUTY_CYCLE      = uintptr(C._LIRC_SET_SEND_DUTY_CYCLE())
	LIRC_SET_REC_DUTY_CYCLE       = uintptr(C._LIRC_SET_REC_DUTY_CYCLE())
	LIRC_SET_TRANSMITTER_MASK     = uintptr(C._LIRC_SET_TRANSMITTER_MASK())
	LIRC_SET_REC_TIMEOUT          = uintptr(C._LIRC_SET_REC_TIMEOUT())
	LIRC_SET_REC_TIMEOUT_REPORTS  = uintptr(C._LIRC_SET_REC_TIMEOUT_REPORTS())
	LIRC_SET_MEASURE_CARRIER_MODE = uintptr(C._LIRC_SET_MEASURE_CARRIER_MODE())
	LIRC_SET_REC_CARRIER_RANGE    = uintptr(C._LIRC_SET_REC_CARRIER_RANGE())
	LIRC_SET_WIDEBAND_RECEIVER    = uintptr(C._LIRC_SET_WIDEBAND_RECEIVER())
)

////////////////////////////////////////////////////////////////////////////////
// IOCTL CALLS

func (this *lirc) getFeatures() (lirc_feature, error) {
	var features lirc_feature
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_GET_FEATURES, unsafe.Pointer(&features)); err != 0 {
		return 0, os.NewSyscallError("getFeatures", err)
	}
	return features, nil
}

func (this *lirc) getSendMode() (gopi.LIRCMode, error) {
	var value gopi.LIRCMode
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_GET_SEND_MODE, unsafe.Pointer(&value)); err != 0 {
		return 0, os.NewSyscallError("getSendMode", err)
	}
	return value, nil
}

func (this *lirc) setSendMode(value gopi.LIRCMode) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_SEND_MODE, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setSendMode", err)
	}
	return nil
}

func (this *lirc) getRcvMode() (gopi.LIRCMode, error) {
	var value gopi.LIRCMode
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_GET_REC_MODE, unsafe.Pointer(&value)); err != 0 {
		return 0, os.NewSyscallError("getRcvMode", err)
	}
	return value, nil
}

func (this *lirc) setRcvMode(value gopi.LIRCMode) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_REC_MODE, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setRcvMode", err)
	}
	return nil
}

func (this *lirc) getRcvResolutionMicros() (uint32, error) {
	var value uint32
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_GET_REC_RESOLUTION, unsafe.Pointer(&value)); err != 0 {
		return 0, os.NewSyscallError("getRcvResolutionMicros", err)
	}
	return value, nil
}

func (this *lirc) getMinMaxTimeoutMicros() (uint32, uint32, error) {
	var min, max uint32
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_GET_MIN_TIMEOUT, unsafe.Pointer(&min)); err != 0 {
		return 0, 0, os.NewSyscallError("getMinMaxTimeoutMicros", err)
	}
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_GET_MAX_TIMEOUT, unsafe.Pointer(&max)); err != 0 {
		return 0, 0, os.NewSyscallError("getMinMaxTimeoutMicros", err)
	}
	return min, max, nil
}

func (this *lirc) setRcvTimeoutMicros(value uint32) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_REC_TIMEOUT, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setRcvTimeoutMicros", err)
	}
	return nil
}

func (this *lirc) getLength() (uint32, error) {
	var value uint32
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_GET_LENGTH, unsafe.Pointer(&value)); err != 0 {
		return 0, os.NewSyscallError("getLength", err)
	}
	return value, nil
}

func (this *lirc) setRcvCarrierHz(value uint32) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_REC_CARRIER, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setRcvCarrierHz", err)
	}
	return nil
}

func (this *lirc) setRcvCarrierRangeHz(value uint32) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_REC_CARRIER_RANGE, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setRcvCarrierRangeHz", err)
	}
	return nil
}

func (this *lirc) setSendCarrierHz(value uint32) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_SEND_CARRIER, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setSendCarrierHz", err)
	}
	return nil
}

func (this *lirc) setSendDutyCycle(value uint32) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_SEND_DUTY_CYCLE, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setSendDutyCycle", err)
	}
	return nil
}

func (this *lirc) setRcvDutyCycle(value uint32) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_REC_DUTY_CYCLE, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setRcvDutyCycle", err)
	}
	return nil
}

func (this *lirc) setTransmitterMask(value uint32) error {
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_TRANSMITTER_MASK, unsafe.Pointer(&value)); err != 0 {
		return os.NewSyscallError("setTransmitterMask", err)
	}
	return nil
}

func (this *lirc) setRcvTimeoutReports(value bool) error {
	value2 := bool2uint32(value)
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_REC_TIMEOUT_REPORTS, unsafe.Pointer(&value2)); err != 0 {
		return os.NewSyscallError("setRcvTimeoutReports", err)
	}
	return nil
}

func (this *lirc) setMeasureCarrierMode(value bool) error {
	value2 := bool2uint32(value)
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_MEASURE_CARRIER_MODE, unsafe.Pointer(&value2)); err != 0 {
		return os.NewSyscallError("setMeasureCarrierMode", err)
	}
	return nil
}

func (this *lirc) setWidebandReceiver(value bool) error {
	value2 := bool2uint32(value)
	if err := this.lirc_ioctl(this.dev.Fd(), LIRC_SET_WIDEBAND_RECEIVER, unsafe.Pointer(&value2)); err != 0 {
		return os.NewSyscallError("setWidebandReceiver", err)
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Call ioctl
func (this *lirc) lirc_ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}

// Convert false -> 0 and true -> 1
func bool2uint32(value bool) uint32 {
	if value {
		return 1
	} else {
		return 0
	}
}
