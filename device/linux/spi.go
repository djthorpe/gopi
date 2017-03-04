/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package linux /* import "github.com/djthorpe/gopi/device/linux" */

import (
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	hw "github.com/djthorpe/gopi/hw"
)

////////////////////////////////////////////////////////////////////////////////
// CGO INTERFACE

/*
#include <sys/ioctl.h>
#include <linux/spi/spidev.h>
static int _SPI_IOC_RD_MODE() { return SPI_IOC_RD_MODE; }
static int _SPI_IOC_WR_MODE() { return SPI_IOC_WR_MODE; }
static int _SPI_IOC_RD_LSB_FIRST() { return SPI_IOC_RD_LSB_FIRST; }
static int _SPI_IOC_WR_LSB_FIRST() { return SPI_IOC_WR_LSB_FIRST; }
static int _SPI_IOC_RD_BITS_PER_WORD() { return SPI_IOC_RD_BITS_PER_WORD; }
static int _SPI_IOC_WR_BITS_PER_WORD() { return SPI_IOC_WR_BITS_PER_WORD; }
static int _SPI_IOC_RD_MAX_SPEED_HZ() { return SPI_IOC_RD_MAX_SPEED_HZ; }
static int _SPI_IOC_WR_MAX_SPEED_HZ() { return SPI_IOC_WR_MAX_SPEED_HZ; }
static int _SPI_IOC_RD_MODE32() { return SPI_IOC_RD_MODE32; }
static int _SPI_IOC_WR_MODE32() { return SPI_IOC_WR_MODE32; }
static int _SPI_IOC_MESSAGE(int n) { return SPI_IOC_MESSAGE(n); }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SPI struct {
	// Bus number
	Bus     uint

	// Slave number
	Slave uint

	// Transfer delay between blocks, in microseconds
	Delay   uint16
}

type spiDriver struct {
	// logger
	log gopi.Logger

	// device
	dev *os.File

	// bus number
	bus uint

	// slave number
	slave uint

	// mode
	mode hw.SPIMode

	// maximum speed in hertz
	speed_hz uint32

	// bits per word
	bits_per_word uint8

	// Transfer delay
	delay_usec uint16

	// mutex lock
	lock sync.Mutex
}

type spiMessage struct {
	tx_buf        uint64
	rx_buf        uint64
	len           uint32
	speed_hz      uint32
	delay_usecs   uint16
	bits_per_word uint8
	cs_change     uint8
	tx_nbits      uint8
	rx_nbits      uint8
	pad           uint16
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	SPI_DEV       = "/dev/spidev"
	SPI_IOC_MAGIC = 107
)

////////////////////////////////////////////////////////////////////////////////
// VARIABLES

var (
	SPI_IOC_RD_MODE          = uintptr(C._SPI_IOC_RD_MODE())
	SPI_IOC_WR_MODE          = uintptr(C._SPI_IOC_WR_MODE())
	SPI_IOC_RD_LSB_FIRST     = uintptr(C._SPI_IOC_RD_LSB_FIRST())
	SPI_IOC_WR_LSB_FIRST     = uintptr(C._SPI_IOC_WR_LSB_FIRST())
	SPI_IOC_RD_BITS_PER_WORD = uintptr(C._SPI_IOC_RD_BITS_PER_WORD())
	SPI_IOC_WR_BITS_PER_WORD = uintptr(C._SPI_IOC_WR_BITS_PER_WORD())
	SPI_IOC_RD_MAX_SPEED_HZ  = uintptr(C._SPI_IOC_RD_MAX_SPEED_HZ())
	SPI_IOC_WR_MAX_SPEED_HZ  = uintptr(C._SPI_IOC_WR_MAX_SPEED_HZ())
	SPI_IOC_RD_MODE32        = uintptr(C._SPI_IOC_RD_MODE32())
	SPI_IOC_WR_MODE32        = uintptr(C._SPI_IOC_WR_MODE32())
)

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Create new SPI object, returns error if not possible
func (config SPI) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<linux.SPI>Open")

	var err error

	// create new GPIO driver
	this := new(spiDriver)
	this.bus = config.Bus
	this.slave = config.Slave
	this.delay_usec = config.Delay

	// Set logging & device
	this.log = log

	// Open the device
	this.dev, err = os.OpenFile(fmt.Sprintf("%v%v.%v", SPI_DEV, this.bus, this.slave), os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	// Get current mode, speed and bits per word
	this.mode, err = this.getMode()
	if err != nil {
		return nil, err
	}
	this.speed_hz, err = this.getMaxSpeedHz()
	if err != nil {
		return nil, err
	}
	this.bits_per_word, err = this.getBitsPerWord()
	if err != nil {
		return nil, err
	}

	// success
	return this, nil
}

// Close SPI connection
func (this *spiDriver) Close() error {
	this.log.Debug("<linux.SPI>Close")

	err := this.dev.Close()
	this.dev = nil
	return err
}

// Strinfigy SPI driver
func (this *spiDriver) String() string {
	return fmt.Sprintf("<linux.SPI>{ bus=%v slave=%v mode=%v delay=%vus max_speed=%vHz bits_per_word=%v }", this.bus, this.slave, this.mode, this.delay_usec, this.speed_hz, this.bits_per_word)
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE MODE, SPEED AND BITS PER WORD

func (this *spiDriver) GetMode() hw.SPIMode {
	return this.mode
}

func (this *spiDriver) GetMaxSpeedHz() uint32 {
	return this.speed_hz
}

func (this *spiDriver) GetBitsPerWord() uint8 {
	return this.bits_per_word
}

func (this *spiDriver) SetMode(mode hw.SPIMode) error {
	err := this.ioctl(this.dev.Fd(), SPI_IOC_WR_MODE, unsafe.Pointer(&mode))
	if err != 0 {
		return err
	}
	var err2 error
	this.mode, err2 = this.getMode()
	return err2
}

func (this *spiDriver) SetMaxSpeedHz(speed uint32) error {
	err := this.ioctl(this.dev.Fd(), SPI_IOC_WR_MAX_SPEED_HZ, unsafe.Pointer(&speed))
	if err != 0 {
		return err
	}
	var err2 error
	this.speed_hz, err2 = this.getMaxSpeedHz()
	return err2
}

func (this *spiDriver) SetBitsPerWord(bits uint8) error {
	err := this.ioctl(this.dev.Fd(), SPI_IOC_WR_BITS_PER_WORD, unsafe.Pointer(&bits))
	if err != 0 {
		return err
	}
	var err2 error
	this.bits_per_word, err2 = this.getBitsPerWord()
	return err2
}

////////////////////////////////////////////////////////////////////////////////
// TRANSFER

func (this *spiDriver) Transfer(send []byte) ([]byte,error) {
	buffer_size := len(send)
	if buffer_size == 0 {
		return []byte{ },nil
	}
	recv := make([]byte,buffer_size)
	message := spiMessage{
		tx_buf: uint64(uintptr(unsafe.Pointer(&send[0]))),
		rx_buf: uint64(uintptr(unsafe.Pointer(&recv[0]))),
		len: uint32(buffer_size),
		speed_hz: this.speed_hz,
		delay_usecs: this.delay_usec,
		bits_per_word: this.bits_per_word,
	}

	err := this.ioctl(this.dev.Fd(),uintptr(C._SPI_IOC_MESSAGE(C.int(1))),unsafe.Pointer(&message))
	if err != 0 {
		return nil, err
	}
	return recv, nil
}

func (this *spiDriver) Read(buffer_size uint32) ([]byte,error) {
	if buffer_size == 0 {
		return []byte{ },nil
	}
	recv := make([]byte,buffer_size)
	message := spiMessage{
		tx_buf: 0,
		rx_buf: uint64(uintptr(unsafe.Pointer(&recv[0]))),
		len: buffer_size,
		speed_hz: this.speed_hz,
		delay_usecs: this.delay_usec,
		bits_per_word: this.bits_per_word,
	}
	err := this.ioctl(this.dev.Fd(),uintptr(C._SPI_IOC_MESSAGE(C.int(1))),unsafe.Pointer(&message))
	if err != 0 {
		return nil, err
	}
	return recv, nil
}

func (this *spiDriver) Write(send []byte) error {
	buffer_size := len(send)
	if buffer_size == 0 {
		return nil
	}
	message := spiMessage{
		tx_buf: uint64(uintptr(unsafe.Pointer(&send[0]))),
		rx_buf: 0,
		len: uint32(buffer_size),
		speed_hz: this.speed_hz,
		delay_usecs: this.delay_usec,
		bits_per_word: this.bits_per_word,
	}
	err := this.ioctl(this.dev.Fd(),uintptr(C._SPI_IOC_MESSAGE(C.int(1))),unsafe.Pointer(&message))
	if err != 0 {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *spiDriver) getMode() (hw.SPIMode, error) {
	var mode uint8

	err := this.ioctl(this.dev.Fd(), SPI_IOC_RD_MODE, unsafe.Pointer(&mode))
	if err != 0 {
		return hw.SPI_MODE_NONE, err
	}
	return hw.SPIMode(mode), nil
}

func (this *spiDriver) getMaxSpeedHz() (uint32, error) {
	var speed_hz uint32

	err := this.ioctl(this.dev.Fd(), SPI_IOC_RD_MAX_SPEED_HZ, unsafe.Pointer(&speed_hz))
	if err != 0 {
		return 0, err
	}
	return speed_hz, nil
}

func (this *spiDriver) getBitsPerWord() (uint8, error) {
	var bits_per_word uint8

	err := this.ioctl(this.dev.Fd(), SPI_IOC_RD_BITS_PER_WORD, unsafe.Pointer(&bits_per_word))
	if err != 0 {
		return 0, err
	}
	return bits_per_word, nil
}

// Call ioctl
func (this *spiDriver) ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
