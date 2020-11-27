// +build linux

package linux

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	I2C_DEV             = "/dev/i2c"
	I2C_SMBUS_BLOCK_MAX = 32 /* As specified in SMBus standard */
)

const (
	// i2c ioctl commands
	I2C_RETRIES     = 0x0701 /* number of times a device address should be polled when not acknowledging */
	I2C_TIMEOUT     = 0x0702 /* set timeout in units of 10 ms */
	I2C_SLAVE       = 0x0703 /* Use this slave address */
	I2C_SLAVE_FORCE = 0x0706 /* Use this slave address, even if it is already in use by a driver! */
	I2C_TENBIT      = 0x0704 /* 0 for 7 bit addrs, != 0 for 10 bit */
	I2C_FUNCS       = 0x0705 /* Get the adapter functionality mask */
	I2C_RDWR        = 0x0707 /* Combined R/W transfer (one STOP only) */
	I2C_PEC         = 0x0708 /* != 0 to use PEC with SMBus */
	I2C_SMBUS       = 0x0720 /* SMBus transfer */
)

const (
	// i2c_smbus_xfer read or write markers
	I2C_SMBUS_READ  uint8 = 0x01
	I2C_SMBUS_WRITE uint8 = 0x00
)

const (
	// SMBus transaction types
	I2C_SMBUS_QUICK            uint32 = 0
	I2C_SMBUS_BYTE             uint32 = 1
	I2C_SMBUS_BYTE_DATA        uint32 = 2
	I2C_SMBUS_WORD_DATA        uint32 = 3
	I2C_SMBUS_PROC_CALL        uint32 = 4
	I2C_SMBUS_BLOCK_DATA       uint32 = 5
	I2C_SMBUS_I2C_BLOCK_BROKEN uint32 = 6
	I2C_SMBUS_BLOCK_PROC_CALL  uint32 = 7
	I2C_SMBUS_I2C_BLOCK_DATA   uint32 = 8
)

const (
	// i2c functions
	I2C_FUNC_NONE                   I2CFunction = 0x00000000
	I2C_FUNC_I2C                    I2CFunction = 0x00000001
	I2C_FUNC_10BIT_ADDR             I2CFunction = 0x00000002
	I2C_FUNC_PROTOCOL_MANGLING      I2CFunction = 0x00000004 /* I2C_M_IGNORE_NAK etc. */
	I2C_FUNC_SMBUS_PEC              I2CFunction = 0x00000008
	I2C_FUNC_NOSTART                I2CFunction = 0x00000010 /* I2C_M_NOSTART */
	I2C_FUNC_SMBUS_BLOCK_PROC_CALL  I2CFunction = 0x00008000 /* SMBus 2.0 */
	I2C_FUNC_SMBUS_QUICK            I2CFunction = 0x00010000
	I2C_FUNC_SMBUS_READ_BYTE        I2CFunction = 0x00020000
	I2C_FUNC_SMBUS_WRITE_BYTE       I2CFunction = 0x00040000
	I2C_FUNC_SMBUS_READ_BYTE_DATA   I2CFunction = 0x00080000
	I2C_FUNC_SMBUS_WRITE_BYTE_DATA  I2CFunction = 0x00100000
	I2C_FUNC_SMBUS_READ_WORD_DATA   I2CFunction = 0x00200000
	I2C_FUNC_SMBUS_WRITE_WORD_DATA  I2CFunction = 0x00400000
	I2C_FUNC_SMBUS_PROC_CALL        I2CFunction = 0x00800000
	I2C_FUNC_SMBUS_READ_BLOCK_DATA  I2CFunction = 0x01000000
	I2C_FUNC_SMBUS_WRITE_BLOCK_DATA I2CFunction = 0x02000000
	I2C_FUNC_SMBUS_READ_I2C_BLOCK   I2CFunction = 0x04000000 /* I2C-like block xfer  */
	I2C_FUNC_SMBUS_WRITE_I2C_BLOCK  I2CFunction = 0x08000000 /* w/ 1-byte reg. addr. */
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2CFunction uint32

type i2c_smbus_ioctl_data struct {
	rw      uint8
	command uint8
	size    uint32
	data    uintptr
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func I2CDevice(bus uint) string {
	return fmt.Sprintf("%v-%v", I2C_DEV, bus)
}

func I2COpenDevice(bus uint) (*os.File, error) {
	if file, err := os.OpenFile(I2CDevice(bus), os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func I2CFunctions(fd uintptr) (I2CFunction, error) {
	var funcs I2CFunction
	if err := i2c_ioctl(fd, I2C_FUNCS, uintptr(unsafe.Pointer(&funcs))); err != nil {
		return funcs, err
	} else {
		return funcs, nil
	}
}

func I2CSetSlave(fd uintptr, slave uint8) error {
	return i2c_ioctl(fd, I2C_SLAVE, uintptr(slave))
}

func I2CDetectSlave(fd uintptr, slave uint8, funcs I2CFunction) (bool, error) {
	if err := I2CSetSlave(fd, slave); err != nil {
		return false, err
	} else if funcs&I2C_FUNC_SMBUS_QUICK != 0 {
		if err := i2c_smbus_write_quick(fd, 0); err == nil {
			return true, nil
		} else {
			return false, nil
		}
	} else if funcs&I2C_FUNC_SMBUS_READ_BYTE != 0 {
		if _, err := i2c_smbus_read_byte(fd); err == nil {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, gopi.ErrNotImplemented.WithPrefix("I2CDetectSlave")
	}
}

func I2CWriteQuick(fd uintptr, value uint8, funcs I2CFunction) error {
	if funcs&I2C_FUNC_SMBUS_QUICK == 0 {
		return gopi.ErrNotImplemented.WithPrefix("I2CWriteQuick")
	} else {
		return i2c_smbus_write_quick(fd, value)
	}
}

func I2CReadUint8(fd uintptr, reg uint8, funcs I2CFunction) (uint8, error) {
	if funcs&I2C_FUNC_SMBUS_READ_BYTE_DATA == 0 {
		return 0, gopi.ErrNotImplemented.WithPrefix("I2CReadUint8")
	} else {
		return i2c_smbus_read_byte_data(fd, reg)
	}
}

func I2CReadInt8(fd uintptr, reg uint8, funcs I2CFunction) (int8, error) {
	value, err := I2CReadUint8(fd, reg, funcs)
	return int8(value), err
}

func I2CReadUint16(fd uintptr, reg uint8, funcs I2CFunction) (uint16, error) {
	if funcs&I2C_FUNC_SMBUS_READ_WORD_DATA == 0 {
		return 0, gopi.ErrNotImplemented.WithPrefix("I2CReadUint16")
	} else {
		return i2c_smbus_read_word_data(fd, reg)
	}
}

func I2CReadInt16(fd uintptr, reg uint8, funcs I2CFunction) (int16, error) {
	value, err := I2CReadUint16(fd, reg, funcs)
	return int16(value), err
}

func I2CReadBlock(fd uintptr, reg, length uint8, funcs I2CFunction) ([]byte, error) {
	if funcs&I2C_FUNC_SMBUS_READ_I2C_BLOCK == 0 {
		return nil, gopi.ErrNotImplemented.WithPrefix("I2CReadBlock")
	} else {
		return i2c_smbus_read_i2c_block_data(fd, reg, length)
	}
}

func I2CWriteUint8(fd uintptr, reg, value uint8, funcs I2CFunction) error {
	if funcs&I2C_FUNC_SMBUS_WRITE_BYTE_DATA == 0 {
		return gopi.ErrNotImplemented.WithPrefix("I2CWriteUint8")
	} else {
		return i2c_smbus_write_byte_data(fd, reg, value)
	}
}

func I2CWriteInt8(fd uintptr, reg uint8, value int8, funcs I2CFunction) error {
	return I2CWriteUint8(fd, reg, uint8(value), funcs)
}

func I2CWriteUint16(fd uintptr, reg uint8, value uint16, funcs I2CFunction) error {
	if funcs&I2C_FUNC_SMBUS_WRITE_WORD_DATA == 0 {
		return gopi.ErrNotImplemented.WithPrefix("I2CWriteUint16")
	} else {
		return i2c_smbus_write_word_data(fd, reg, value)
	}
}

func I2CWriteInt16(fd uintptr, reg uint8, value int16, funcs I2CFunction) error {
	return I2CWriteUint16(fd, reg, uint16(value), funcs)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func i2c_ioctl(fd, cmd, arg uintptr) error {
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0); err != 0 {
		return err
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// SMBUS PRIVATE METHODS

func i2c_smbus_access(fd uintptr, rw uint8, command uint8, size uint32, data uintptr) error {
	args := &i2c_smbus_ioctl_data{
		rw:      rw,
		command: command,
		size:    size,
		data:    data,
	}
	return i2c_ioctl(fd, I2C_SMBUS, uintptr(unsafe.Pointer(args)))
}

func i2c_smbus_write_quick(fd uintptr, value uint8) error {
	return i2c_smbus_access(fd, value, uint8(0), I2C_SMBUS_QUICK, 0)
}

func i2c_smbus_read_byte(fd uintptr) (uint8, error) {
	var data uint8
	if err := i2c_smbus_access(fd, I2C_SMBUS_READ, 0, I2C_SMBUS_BYTE, uintptr(unsafe.Pointer(&data))); err != nil {
		return uint8(0), err
	}
	return data, nil
}

func i2c_smbus_write_byte(fd uintptr, value uint8) error {
	return i2c_smbus_access(fd, I2C_SMBUS_WRITE, value, I2C_SMBUS_BYTE, 0)
}

func i2c_smbus_read_byte_data(fd uintptr, command uint8) (uint8, error) {
	var data uint8
	if err := i2c_smbus_access(fd, I2C_SMBUS_READ, command, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&data))); err != nil {
		return uint8(0), err
	}
	return data, nil
}

func i2c_smbus_write_byte_data(fd uintptr, command, value uint8) error {
	return i2c_smbus_access(fd, I2C_SMBUS_WRITE, command, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&value)))
}

func i2c_smbus_read_word_data(fd uintptr, command uint8) (uint16, error) {
	var data uint16
	if err := i2c_smbus_access(fd, I2C_SMBUS_READ, command, I2C_SMBUS_WORD_DATA, uintptr(unsafe.Pointer(&data))); err != nil {
		return uint16(0), err
	}
	return data, nil
}

func i2c_smbus_write_word_data(fd uintptr, command uint8, value uint16) error {
	return i2c_smbus_access(fd, I2C_SMBUS_WRITE, command, I2C_SMBUS_WORD_DATA, uintptr(unsafe.Pointer(&value)))
}

func i2c_smbus_process_call(fd uintptr, command uint8, value uint16) (uint16, error) {
	if err := i2c_smbus_access(fd, I2C_SMBUS_WRITE, command, I2C_SMBUS_PROC_CALL, uintptr(unsafe.Pointer(&value))); err != nil {
		return value, err
	}
	return value, nil
}

func i2c_smbus_read_block_data(fd uintptr, command uint8) ([]byte, error) {
	var data [I2C_SMBUS_BLOCK_MAX + 2]byte
	if err := i2c_smbus_access(fd, I2C_SMBUS_READ, command, I2C_SMBUS_BLOCK_DATA, uintptr(unsafe.Pointer(&data))); err != nil {
		return nil, err
	}
	block := make([]byte, data[0])
	for i := uint8(0); i < data[0]; i++ {
		block[i] = data[i+2]
	}
	return block, nil
}

func i2c_smbus_read_i2c_block_data(fd uintptr, command uint8, length uint8) ([]byte, error) {
	var data [I2C_SMBUS_BLOCK_MAX + 2]byte

	size := I2C_SMBUS_I2C_BLOCK_DATA
	data[0] = length
	if length > I2C_SMBUS_BLOCK_MAX {
		length = I2C_SMBUS_BLOCK_MAX
		size = I2C_SMBUS_I2C_BLOCK_BROKEN
	}
	if err := i2c_smbus_access(fd, I2C_SMBUS_READ, command, size, uintptr(unsafe.Pointer(&data))); err != nil {
		return nil, err
	}
	block := make([]byte, data[0])
	for i := uint8(0); i < data[0]; i++ {
		block[i] = data[i+1]
	}
	return block, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f I2CFunction) String() string {
	str := ""
	if f == 0 {
		return f.FlagString()
	}
	flag := I2CFunction(I2C_FUNC_I2C)
	for {
		if f&flag != I2CFunction(0) {
			str += flag.FlagString() + ","
		}
		flag = flag << 1
		if flag > I2C_FUNC_SMBUS_WRITE_I2C_BLOCK {
			break
		}
	}
	return strings.TrimSuffix(str, ",")
}

func (f I2CFunction) FlagString() string {
	switch f {
	case I2C_FUNC_NONE:
		return "I2C_FUNC_NONE"
	case I2C_FUNC_I2C:
		return "I2C_FUNC_I2C"
	case I2C_FUNC_10BIT_ADDR:
		return "I2C_FUNC_10BIT_ADDR"
	case I2C_FUNC_PROTOCOL_MANGLING:
		return "I2C_FUNC_PROTOCOL_MANGLING"
	case I2C_FUNC_SMBUS_PEC:
		return "I2C_FUNC_SMBUS_PEC"
	case I2C_FUNC_NOSTART:
		return "I2C_FUNC_NOSTART"
	case I2C_FUNC_SMBUS_BLOCK_PROC_CALL:
		return "I2C_FUNC_SMBUS_BLOCK_PROC_CALL"
	case I2C_FUNC_SMBUS_QUICK:
		return "I2C_FUNC_SMBUS_QUICK"
	case I2C_FUNC_SMBUS_READ_BYTE:
		return "I2C_FUNC_SMBUS_READ_BYTE"
	case I2C_FUNC_SMBUS_WRITE_BYTE:
		return "I2C_FUNC_SMBUS_WRITE_BYTE"
	case I2C_FUNC_SMBUS_READ_BYTE_DATA:
		return "I2C_FUNC_SMBUS_READ_BYTE_DATA"
	case I2C_FUNC_SMBUS_WRITE_BYTE_DATA:
		return "I2C_FUNC_SMBUS_WRITE_BYTE_DATA"
	case I2C_FUNC_SMBUS_READ_WORD_DATA:
		return "I2C_FUNC_SMBUS_READ_WORD_DATA"
	case I2C_FUNC_SMBUS_WRITE_WORD_DATA:
		return "I2C_FUNC_SMBUS_WRITE_WORD_DATA"
	case I2C_FUNC_SMBUS_PROC_CALL:
		return "I2C_FUNC_SMBUS_PROC_CALL"
	case I2C_FUNC_SMBUS_READ_BLOCK_DATA:
		return "I2C_FUNC_SMBUS_READ_BLOCK_DATA"
	case I2C_FUNC_SMBUS_WRITE_BLOCK_DATA:
		return "I2C_FUNC_SMBUS_WRITE_BLOCK_DATA"
	case I2C_FUNC_SMBUS_READ_I2C_BLOCK:
		return "I2C_FUNC_SMBUS_READ_I2C_BLOCK"
	case I2C_FUNC_SMBUS_WRITE_I2C_BLOCK:
		return "I2C_FUNC_SMBUS_WRITE_I2C_BLOCK"
	default:
		return "[?? Unknown I2CFunction value]"
	}
}
