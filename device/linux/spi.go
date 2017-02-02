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
	util "github.com/djthorpe/gopi/util"
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
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type SPI struct {
	Bus     uint
	Channel uint
}

type spiDriver struct {
	log     *util.LoggerDevice // logger
	dev     *os.File
	bus     uint
	channel uint
	lock    sync.Mutex
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
func (config SPI) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug("<linux.SPI>Open")

	var err error

	// create new GPIO driver
	this := new(spiDriver)
	this.bus = config.Bus
	this.channel = config.Channel

	// Set logging & device
	this.log = log

	// Open the device
	this.dev, err = os.OpenFile(fmt.Sprintf("%v%v.%v", SPI_DEV, this.bus, this.channel), os.O_RDWR, 0)
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
	mode, _ := this.GetMode()
	speed, _ := this.GetMaxSpeedHz()
	bits, _ := this.GetBitsPerWord()
	return fmt.Sprintf("<linux.SPI>{ bus=%v channel=%v mode=%v max_speed=%vHz bits_per_word=%v }", this.bus, this.channel, mode, speed, bits)
}

////////////////////////////////////////////////////////////////////////////////
// READ AND WRITE MODE, SPEED AND BITS PER WORD

func (this *spiDriver) GetMode() (hw.SPIMode, error) {
	var mode uint8
	err := spiIoctl(this.dev.Fd(), SPI_IOC_RD_MODE, unsafe.Pointer(&mode))
	if err != 0 {
		return hw.SPI_MODE_NONE, err
	}
	return hw.SPIMode(mode), nil
}

func (this *spiDriver) GetMaxSpeedHz() (uint32, error) {
	var speed uint32
	err := spiIoctl(this.dev.Fd(), SPI_IOC_RD_MAX_SPEED_HZ, unsafe.Pointer(&speed))
	if err != 0 {
		return 0, err
	}
	return speed, nil
}

func (this *spiDriver) GetBitsPerWord() (uint8, error) {
	var bits uint8
	err := spiIoctl(this.dev.Fd(), SPI_IOC_RD_BITS_PER_WORD, unsafe.Pointer(&bits))
	if err != 0 {
		return 0, err
	}
	return bits, nil
}

func (this *spiDriver) SetMode(mode hw.SPIMode) error {
	err := spiIoctl(this.dev.Fd(), SPI_IOC_WR_MODE, unsafe.Pointer(&mode))
	if err != 0 {
		return err
	}
	return nil
}

func (this *spiDriver) SetMaxSpeedHz(speed uint32) error {
	err := spiIoctl(this.dev.Fd(), SPI_IOC_WR_MAX_SPEED_HZ, unsafe.Pointer(&speed))
	if err != 0 {
		return err
	}
	return nil
}

func (this *spiDriver) SetBitsPerWord(bits uint8) error {
	err := spiIoctl(this.dev.Fd(), SPI_IOC_WR_BITS_PER_WORD, unsafe.Pointer(&bits))
	if err != 0 {
		return err
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Call ioctl
func spiIoctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
