/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
	"strings"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type I2C struct {
	Device gopi.HardwareDriver
	Bus    uint
}

type I2CDriver struct {
	log   *util.LoggerDevice // logger
	bus   uint
	slave uint8
	dev   *os.File
	funcs I2CFunction
}

type I2CFunction uint32

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	I2C_DEV           = "/dev/i2c"
	I2C_SLAVE_NONE uint8 = 0xFF
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

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new I2C object, returns error if not possible
func (config I2C) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<linux.I2C>Open")

	var err error

	// create new GPIO driver
	this := new(I2CDriver)

	// Set logging & device
	this.log = log
	this.bus = config.Bus
	this.slave = I2C_SLAVE_NONE

	// Open the /dev/mem and provide offset & size for accessing memory
	this.dev, err = i2cOpenDevice(config.Bus)
	if err != nil {
		return nil, err
	}

	// Get functionality
	if this.funcs, err = this.i2cFuncs(); err != nil {
		this.dev.Close()
		return nil, err
	}

	// success
	return this, nil
}

// Close I2C connection
func (this *I2CDriver) Close() error {
	this.log.Debug("<linux.I2C>Close")

	err := this.dev.Close()
	this.dev = nil
	this.slave = I2C_SLAVE_NONE
	return err
}

// Strinfigy I2C object
func (this *I2CDriver) String() string {
	// Gather I2C functions
	flag := I2CFunction(I2C_FUNC_I2C)
	funcs := ""
	for {
		if this.funcs & flag != I2CFunction(0) {
			funcs = funcs + flag.String() + ","
		}
		flag = flag << 1
		if flag > I2C_FUNC_SMBUS_WRITE_I2C_BLOCK {
			break
		}
	}
	slave := fmt.Sprintf("%02X",this.slave)
	if this.slave == I2C_SLAVE_NONE {
		slave = "I2C_SLAVE_NONE"
	}
	return fmt.Sprintf("<linux.I2C>{ bus=%v slave=%v funcs={ %v } }", this.bus, slave, strings.TrimSuffix(funcs,","))
}

// Stringify I2CFuncs
func (f I2CFunction) String() string {
	switch(f) {
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
// SLAVE ADDRESS

func (this *I2CDriver) SetSlave(slave uint8) error {
	if this.slave == slave {
		return nil
	}
	if err := i2cIoctl(this.dev.Fd(), I2C_SLAVE, uintptr(slave)); err != nil {
		return err
	}
	this.slave = slave
	return nil
}

// Returns current slave address, or returns I2C_SLAVE_NONE if no slave
// address has not yet been set
func (this *I2CDriver) GetSlave() uint8 {
	return this.slave
}

func (this *I2CDriver) DetectSlave(slave uint8) (bool,error) {
	// Store old slave address and set this one
	old_slave := this.slave
	if slave != old_slave {
		if err := i2cIoctl(this.dev.Fd(), I2C_SLAVE, uintptr(slave)); err != nil {
			return false,err
		}
	}

	var detect bool
	if this.funcs & I2C_FUNC_SMBUS_QUICK != 0 {
		res = i2c_smbus_write_quick(this.dev.Fd(),I2C_SMBUS_WRITE)
		detect := (res >= 0)
	} else if this.funcs & I2C_FUNC_SMBUS_READ_BYTE != 0 {
		res = i2c_smbus_read_byte(this.dev.Fd());
		detect := (res >= 0)
	} else {
		return false,errors.New("I2C bus does not support detection commands")
	}

	// Restore slave address
	if old_slave != I2C_SLAVE_NONE {
		if err := i2cIoctl(this.dev.Fd(), I2C_SLAVE, uintptr(old_slave)); err != nil {
			return false,err
		}
	}
	return detect,nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func i2cOpenDevice(bus uint) (*os.File, error) {
	var file *os.File
	var err error

	if file, err = os.OpenFile(fmt.Sprintf("%v-%v", I2C_DEV, bus), os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	}
	return file, nil
}

func (this *I2CDriver) i2cFuncs() (I2CFunction, error) {
	var funcs I2CFunction
	if err := i2cIoctl(this.dev.Fd(), I2C_FUNCS, uintptr(unsafe.Pointer(&funcs))); err != nil {
		return funcs, err
	}
	return funcs, nil
}

func (this *I2CDriver) i2c_smbus_access(read_write uint8,command uint8,size int,data i2c_smbus_data) error {
	struct i2c_smbus_ioctl_data args;
	__s32 err;

	args.read_write = read_write;
	args.command = command;
	args.size = size;
	args.data = data;

	err = ioctl(this.dev.Fd(),I2C_SMBUS, &args);
	if (err == -1)
		err = -errno;
	return err;
}

func i2cIoctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
