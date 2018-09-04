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
	"syscall"
	"unsafe"

	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2C struct {
	Bus uint
}

type I2CFunction uint32

type i2c struct {
	log   gopi.Logger
	bus   uint
	slave uint8
	dev   *os.File
	funcs I2CFunction
	lock  sync.Mutex
}

type i2c_smbus_ioctl_data struct {
	rw      uint8
	command uint8
	size    uint32
	data    uintptr
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	I2C_DEV                   = "/dev/i2c"
	I2C_SLAVE_NONE      uint8 = 0xFF
	I2C_SMBUS_BLOCK_MAX       = 32 /* As specified in SMBus standard */
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
	// i2c functions
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

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new I2C object, returns error if not possible
func (config I2C) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.linux.I2C>Open")

	// create new GPIO driver
	this := new(i2c)

	// Set logging & device
	this.log = log
	this.bus = config.Bus
	this.slave = I2C_SLAVE_NONE

	// Open the device
	if dev, err := i2c_open_device(config.Bus); err != nil {
		return nil, err
	} else {
		this.dev = dev
	}

	// Get functionality
	if funcs, err := this.i2cFuncs(); err != nil {
		this.dev.Close()
		return nil, err
	} else {
		this.funcs = funcs
	}

	// success
	return this, nil
}

// Close I2C connection
func (this *i2c) Close() error {
	this.log.Debug("<sys.hw.linux.I2C>Close")

	err := this.dev.Close()
	this.dev = nil
	this.slave = I2C_SLAVE_NONE
	return err
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

// Strinfigy I2C object
func (this *i2c) String() string {
	// Gather I2C functions
	flag := I2CFunction(I2C_FUNC_I2C)
	funcs := ""
	for {
		if this.funcs&flag != I2CFunction(0) {
			funcs = funcs + flag.String() + ","
		}
		flag = flag << 1
		if flag > I2C_FUNC_SMBUS_WRITE_I2C_BLOCK {
			break
		}
	}
	slave := fmt.Sprintf("%02X", this.slave)
	if this.slave == I2C_SLAVE_NONE {
		slave = "I2C_SLAVE_NONE"
	}
	return fmt.Sprintf("<sys.hw.linux.I2C>{ bus=%v slave=%v funcs={ %v } }", this.bus, slave, strings.TrimSuffix(funcs, ","))
}

// Stringify I2CFuncs
func (f I2CFunction) String() string {
	switch f {
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

////////////////////////////////////////////////////////////////////////////////
// I2C INTERFACE - SLAVE ADDRESS

// SetSlave sets the slacve address or returns an error if the slave
// address is not found or unsupported
func (this *i2c) SetSlave(slave uint8) error {
	this.log.Debug2("<sys.hw.linux.I2C.SetSlave>{ slave=%v }", slave)
	if this.slave == slave {
		return nil
	} else if slave == I2C_SLAVE_NONE {
		return gopi.ErrBadParameter
	} else if err := i2c_ioctl(this.dev.Fd(), I2C_SLAVE, uintptr(slave)); err != nil {
		return err
	} else {
		this.slave = slave
		return nil
	}
}

// GetSlave returns current slave address, or returns I2C_SLAVE_NONE if no slave
// address has not yet been set
func (this *i2c) GetSlave() uint8 {
	this.log.Debug2("<sys.hw.linux.I2C.GetSlave>{ slave=%v }", this.slave)
	return this.slave
}

// DetectSlave checks to see if there is a device on a certain slave address
func (this *i2c) DetectSlave(slave uint8) (bool, error) {
	this.log.Debug2("<sys.hw.linux.I2C.DetectSlave>{ slave=%v }", slave)

	// Store old slave address and set this one
	old_slave := this.slave
	if slave != old_slave {
		if err := i2c_ioctl(this.dev.Fd(), I2C_SLAVE, uintptr(slave)); err != nil {
			return false, err
		}
	}

	var detect bool
	if this.funcs&I2C_FUNC_SMBUS_QUICK != 0 {
		err := this.i2c_smbus_write_quick(0)
		if err == nil {
			detect = true
		} else {
			detect = false
		}
	} else if this.funcs&I2C_FUNC_SMBUS_READ_BYTE != 0 {
		_, err := this.i2c_smbus_read_byte()
		if err == nil {
			detect = true
		} else {
			detect = false
		}
	} else {
		return false, gopi.ErrNotImplemented
	}

	// Restore slave address
	if old_slave != I2C_SLAVE_NONE {
		if err := i2c_ioctl(this.dev.Fd(), I2C_SLAVE, uintptr(old_slave)); err != nil {
			return false, err
		}
	}
	return detect, nil
}

////////////////////////////////////////////////////////////////////////////////
// QUICK METHOD

func (this *i2c) WriteQuick(value uint8) error {
	this.log.Debug2("<sys.hw.linux.I2C.WriteQuick>{ value=%v }", value)
	if this.slave == I2C_SLAVE_NONE {
		return gopi.ErrBadParameter
	}
	if this.funcs&I2C_FUNC_SMBUS_QUICK == 0 {
		return gopi.ErrNotImplemented
	}
	return this.i2c_smbus_write_quick(value)
}

////////////////////////////////////////////////////////////////////////////////
// READ METHODS

func (this *i2c) ReadUint8(reg uint8) (uint8, error) {
	this.log.Debug2("<sys.hw.linux.I2C.ReadUint8>{ reg=0x%02X }", reg)
	if this.slave == I2C_SLAVE_NONE {
		return uint8(0), gopi.ErrBadParameter
	}
	if this.funcs&I2C_FUNC_SMBUS_READ_BYTE_DATA == 0 {
		return uint8(0), gopi.ErrNotImplemented
	}
	return this.i2c_smbus_read_byte_data(reg)
}

func (this *i2c) ReadInt8(reg uint8) (int8, error) {
	v, e := this.ReadUint8(reg)
	return int8(v), e
}

func (this *i2c) ReadUint16(reg uint8) (uint16, error) {
	this.log.Debug2("<sys.hw.linux.I2C.ReadUint16>{ reg=0x%02X }", reg)
	if this.slave == I2C_SLAVE_NONE {
		return uint16(0), gopi.ErrBadParameter
	}
	if this.funcs&I2C_FUNC_SMBUS_READ_WORD_DATA == 0 {
		return uint16(0), gopi.ErrNotImplemented
	}
	return this.i2c_smbus_read_word_data(reg)
}

func (this *i2c) ReadInt16(reg uint8) (int16, error) {
	v, e := this.ReadUint16(reg)
	return int16(v), e
}

func (this *i2c) ReadBlock(reg, length uint8) ([]byte, error) {
	this.log.Debug2("<sys.hw.linux.I2C.ReadUint16>{ reg=0x%02X length=%v }", reg, length)
	if this.slave == I2C_SLAVE_NONE {
		return nil, gopi.ErrBadParameter
	}
	if this.funcs&I2C_FUNC_SMBUS_READ_I2C_BLOCK == 0 {
		return nil, gopi.ErrNotImplemented
	}
	return this.i2c_smbus_read_i2c_block_data(reg, length)
}

////////////////////////////////////////////////////////////////////////////////
// WRITE METHODS

func (this *i2c) WriteUint8(reg, value uint8) error {
	this.log.Debug2("<sys.hw.linux.I2C.WriteUint8>{ reg=0x%02X value=%v }", reg, value)
	if this.slave == I2C_SLAVE_NONE {
		return gopi.ErrBadParameter
	}
	if this.funcs&I2C_FUNC_SMBUS_WRITE_BYTE_DATA == 0 {
		return gopi.ErrNotImplemented
	}
	return this.i2c_smbus_write_byte_data(reg, value)
}

func (this *i2c) WriteInt8(reg uint8, value int8) error {
	return this.WriteUint8(reg, uint8(value))
}

func (this *i2c) WriteUint16(reg uint8, value uint16) error {
	this.log.Debug2("<sys.hw.linux.I2C.WriteUint16>{ reg=0x%02X value=%v }", reg, value)
	if this.slave == I2C_SLAVE_NONE {
		return gopi.ErrBadParameter
	}
	if this.funcs&I2C_FUNC_SMBUS_WRITE_WORD_DATA == 0 {
		return gopi.ErrNotImplemented
	}
	return this.i2c_smbus_write_word_data(reg, value)
}

func (this *i2c) WriteInt16(reg uint8, value int16) error {
	return this.WriteUint16(reg, uint16(value))
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *i2c) i2cFuncs() (I2CFunction, error) {
	var funcs I2CFunction
	this.lock.Lock()
	defer this.lock.Unlock()
	if err := i2c_ioctl(this.dev.Fd(), I2C_FUNCS, uintptr(unsafe.Pointer(&funcs))); err != nil {
		return funcs, err
	}
	return funcs, nil
}

func i2c_open_device(bus uint) (*os.File, error) {
	if file, err := os.OpenFile(fmt.Sprintf("%v-%v", I2C_DEV, bus), os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func i2c_ioctl(fd, cmd, arg uintptr) error {
	if _, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0); err != 0 {
		return err
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// SMBUS PRIVATE METHODS

func (this *i2c) i2c_smbus_access(rw uint8, command uint8, size uint32, data uintptr) error {
	args := &i2c_smbus_ioctl_data{
		rw:      rw,
		command: command,
		size:    size,
		data:    data,
	}
	this.lock.Lock()
	defer this.lock.Unlock()
	return i2c_ioctl(this.dev.Fd(), I2C_SMBUS, uintptr(unsafe.Pointer(args)))
}

func (this *i2c) i2c_smbus_write_quick(value uint8) error {
	return this.i2c_smbus_access(value, uint8(0), I2C_SMBUS_QUICK, 0)
}

func (this *i2c) i2c_smbus_read_byte() (uint8, error) {
	var data uint8
	if err := this.i2c_smbus_access(I2C_SMBUS_READ, 0, I2C_SMBUS_BYTE, uintptr(unsafe.Pointer(&data))); err != nil {
		return uint8(0), err
	}
	return data, nil
}

func (this *i2c) i2c_smbus_write_byte(value uint8) error {
	if err := this.i2c_smbus_access(I2C_SMBUS_WRITE, value, I2C_SMBUS_BYTE, 0); err != nil {
		return err
	}
	return nil
}

func (this *i2c) i2c_smbus_read_byte_data(command uint8) (uint8, error) {
	var data uint8
	if err := this.i2c_smbus_access(I2C_SMBUS_READ, command, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&data))); err != nil {
		return uint8(0), err
	}
	return data, nil
}

func (this *i2c) i2c_smbus_write_byte_data(command, value uint8) error {
	if err := this.i2c_smbus_access(I2C_SMBUS_WRITE, command, I2C_SMBUS_BYTE_DATA, uintptr(unsafe.Pointer(&value))); err != nil {
		return err
	}
	return nil
}

func (this *i2c) i2c_smbus_read_word_data(command uint8) (uint16, error) {
	var data uint16
	if err := this.i2c_smbus_access(I2C_SMBUS_READ, command, I2C_SMBUS_WORD_DATA, uintptr(unsafe.Pointer(&data))); err != nil {
		return uint16(0), err
	}
	return data, nil
}

func (this *i2c) i2c_smbus_write_word_data(command uint8, value uint16) error {
	if err := this.i2c_smbus_access(I2C_SMBUS_WRITE, command, I2C_SMBUS_WORD_DATA, uintptr(unsafe.Pointer(&value))); err != nil {
		return err
	}
	return nil
}

func (this *i2c) i2c_smbus_process_call(command uint8, value uint16) (uint16, error) {
	if err := this.i2c_smbus_access(I2C_SMBUS_WRITE, command, I2C_SMBUS_PROC_CALL, uintptr(unsafe.Pointer(&value))); err != nil {
		return value, err
	}
	return value, nil
}

func (this *i2c) i2c_smbus_read_block_data(command uint8) ([]byte, error) {
	var data [I2C_SMBUS_BLOCK_MAX + 2]byte
	if err := this.i2c_smbus_access(I2C_SMBUS_READ, command, I2C_SMBUS_BLOCK_DATA, uintptr(unsafe.Pointer(&data))); err != nil {
		return nil, err
	}
	block := make([]byte, data[0])
	for i := uint8(0); i < data[0]; i++ {
		block[i] = data[i+2]
	}
	return block, nil
}

func (this *i2c) i2c_smbus_read_i2c_block_data(command uint8, length uint8) ([]byte, error) {
	var data [I2C_SMBUS_BLOCK_MAX + 2]byte

	size := I2C_SMBUS_I2C_BLOCK_DATA
	data[0] = length
	if length > I2C_SMBUS_BLOCK_MAX {
		length = I2C_SMBUS_BLOCK_MAX
		size = I2C_SMBUS_I2C_BLOCK_BROKEN
	}
	if err := this.i2c_smbus_access(I2C_SMBUS_READ, command, size, uintptr(unsafe.Pointer(&data))); err != nil {
		return nil, err
	}
	block := make([]byte, data[0])
	for i := uint8(0); i < data[0]; i++ {
		block[i] = data[i+1]
	}
	return block, nil
}
