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
	"strings"
	"encoding/hex"
	"fmt"
	"os"
	"sync"
	"syscall"
	"unsafe"

	// Frameworks
	"github.com/djthorpe/gopi"
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
	Bus uint

	// Slave number
	Slave uint

	// Transfer delay between blocks, in microseconds
	Delay uint16
}

type spi struct {
	log           gopi.Logger  // logger
	dev           *os.File     // device
	bus           uint         // bus number
	slave         uint         // slave number
	mode          gopi.SPIMode // mode
	speed_hz      uint32       // maximum speed in hertz
	bits_per_word uint8        // bits per word
	delay_usec    uint16       // Transfer delay
	lock          sync.Mutex   // mutex lock
}

type spi_message struct {
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

// Create new SPI object or returns error
func (config SPI) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<sys.hw.linux.SPI>Open{ bus=%v slave=%v delay=%vus }", config.Bus, config.Slave, config.Delay)

	// create new GPIO driver
	this := new(spi)
	this.bus = config.Bus
	this.slave = config.Slave
	this.delay_usec = config.Delay
	this.log = log

	// Open the device
	if dev, err := os.OpenFile(fmt.Sprintf("%v%v.%v", SPI_DEV, this.bus, this.slave), os.O_RDWR, 0); err != nil {
		return nil, err
	} else {
		this.dev = dev
	}

	// Get current mode, speed and bits per word
	if mode, err := this.getMode(); err != nil {
		return nil, err
	} else {
		this.mode = mode
	}
	if speed_hz, err := this.getMaxSpeedHz(); err != nil {
		return nil, err
	} else {
		this.speed_hz = speed_hz
	}
	if bits_per_word, err := this.getBitsPerWord(); err != nil {
		return nil, err
	} else {
		this.bits_per_word = bits_per_word
	}

	// success
	return this, nil
}

// Close SPI connection
func (this *spi) Close() error {
	this.log.Debug("<sys.hw.linux.SPI>Close")

	err := this.dev.Close()
	this.dev = nil
	return err
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *spi) String() string {
	return fmt.Sprintf("<sys.hw.linux.SPI>{ bus=%v slave=%v mode=%v delay=%vus max_speed=%vHz bits_per_word=%v }", this.bus, this.slave, this.mode, this.delay_usec, this.speed_hz, this.bits_per_word)
}

////////////////////////////////////////////////////////////////////////////////
// GET AND SET PARAMETERS

func (this *spi) Mode() gopi.SPIMode {
	return this.mode
}

func (this *spi) MaxSpeedHz() uint32 {
	return this.speed_hz
}

func (this *spi) BitsPerWord() uint8 {
	return this.bits_per_word
}

func (this *spi) SetMode(mode gopi.SPIMode) error {
	this.log.Debug2("<sys.hw.linux.SPI.SetMode>{ mode=%v }", mode)
	if err := this.spi_ioctl(this.dev.Fd(), SPI_IOC_WR_MODE, unsafe.Pointer(&mode)); err != 0 {
		return os.NewSyscallError("SetMode", err)
	} else if mode, err := this.getMode(); err != nil {
		return err
	} else {
		this.mode = mode
		return nil
	}
}

func (this *spi) SetMaxSpeedHz(speed uint32) error {
	this.log.Debug2("<sys.hw.linux.SPI.SetMaxSpeedHz>{ speed=%v }", speed)
	if err := this.spi_ioctl(this.dev.Fd(), SPI_IOC_WR_MAX_SPEED_HZ, unsafe.Pointer(&speed)); err != 0 {
		return os.NewSyscallError("SetMaxSpeedHz", err)
	} else if speed_hz, err := this.getMaxSpeedHz(); err != nil {
		return err
	} else {
		this.speed_hz = speed_hz
		return nil
	}
}

func (this *spi) SetBitsPerWord(bits uint8) error {
	this.log.Debug2("<sys.hw.linux.SPI.SetBitsPerWord>{ bits=%v }", bits)
	if err := this.spi_ioctl(this.dev.Fd(), SPI_IOC_WR_BITS_PER_WORD, unsafe.Pointer(&bits)); err != 0 {
		return os.NewSyscallError("SetBitsPerWord", err)
	} else if bits_per_word, err := this.getBitsPerWord(); err != nil {
		return err
	} else {
		this.bits_per_word = bits_per_word
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// TRANSFER

func (this *spi) Transfer(send []byte) ([]byte, error) {
	this.log.Debug2("<sys.hw.linux.SPI.Transfer>{ send=%v }", strings.ToUpper(hex.EncodeToString(send))
	buffer_size := len(send)
	if buffer_size == 0 {
		return []byte{}, nil
	}
	recv := make([]byte, buffer_size)
	message := spi_message{
		tx_buf:        uint64(uintptr(unsafe.Pointer(&send[0]))),
		rx_buf:        uint64(uintptr(unsafe.Pointer(&recv[0]))),
		len:           uint32(buffer_size),
		speed_hz:      this.speed_hz,
		delay_usecs:   this.delay_usec,
		bits_per_word: this.bits_per_word,
	}

	if err := this.spi_ioctl(this.dev.Fd(), uintptr(C._SPI_IOC_MESSAGE(C.int(1))), unsafe.Pointer(&message)); err != 0 {
		return nil, os.NewSyscallError("Transfer", err)
	} else {
		return recv, nil
	}
}

func (this *spi) Read(buffer_size uint32) ([]byte, error) {
	this.log.Debug2("<sys.hw.linux.SPI.Read>{ buffer_size=%v }", buffer_size)
	if buffer_size == 0 {
		return []byte{}, nil
	}
	recv := make([]byte, buffer_size)
	message := spi_message{
		tx_buf:        0,
		rx_buf:        uint64(uintptr(unsafe.Pointer(&recv[0]))),
		len:           buffer_size,
		speed_hz:      this.speed_hz,
		delay_usecs:   this.delay_usec,
		bits_per_word: this.bits_per_word,
	}
	if err := this.spi_ioctl(this.dev.Fd(), uintptr(C._SPI_IOC_MESSAGE(C.int(1))), unsafe.Pointer(&message)); err != 0 {
		return nil, os.NewSyscallError("Read", err)
	} else {
		return recv, nil
	}
}

func (this *spi) Write(send []byte) error {
	this.log.Debug2("<sys.hw.linux.SPI.Write>{ send=%v }", strings.ToUpper(hex.EncodeToString(send))
	buffer_size := len(send)
	if buffer_size == 0 {
		return nil
	}
	message := spi_message{
		tx_buf:        uint64(uintptr(unsafe.Pointer(&send[0]))),
		rx_buf:        0,
		len:           uint32(buffer_size),
		speed_hz:      this.speed_hz,
		delay_usecs:   this.delay_usec,
		bits_per_word: this.bits_per_word,
	}
	if err := this.spi_ioctl(this.dev.Fd(), uintptr(C._SPI_IOC_MESSAGE(C.int(1))), unsafe.Pointer(&message)); err != 0 {
		return os.NewSyscallError("Write", err)
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *spi) getMode() (gopi.SPIMode, error) {
	var mode uint8
	if err := this.spi_ioctl(this.dev.Fd(), SPI_IOC_RD_MODE, unsafe.Pointer(&mode)); err != 0 {
		return gopi.SPI_MODE_NONE, os.NewSyscallError("spi_ioctl", err)
	} else {
		return gopi.SPIMode(mode), nil
	}
}

func (this *spi) getMaxSpeedHz() (uint32, error) {
	var speed_hz uint32
	if err := this.spi_ioctl(this.dev.Fd(), SPI_IOC_RD_MAX_SPEED_HZ, unsafe.Pointer(&speed_hz)); err != 0 {
		return 0, os.NewSyscallError("spi_ioctl", err)
	} else {
		return speed_hz, nil
	}
}

func (this *spi) getBitsPerWord() (uint8, error) {
	var bits_per_word uint8
	if err := this.spi_ioctl(this.dev.Fd(), SPI_IOC_RD_BITS_PER_WORD, unsafe.Pointer(&bits_per_word)); err != 0 {
		return 0, os.NewSyscallError("spi_ioctl", err)
	} else {
		return bits_per_word, nil
	}
}

// Call ioctl
func (this *spi) spi_ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	this.lock.Lock()
	defer this.lock.Unlock()
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
