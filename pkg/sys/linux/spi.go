// +build linux

package linux

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"

	"github.com/djthorpe/gopi/v3"
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
// PUBLIC METHODS

func SPIDevice(bus, slave uint) string {
	return fmt.Sprintf("%v%v.%v", SPI_DEV, bus, slave)
}

func SPIOpenDevice(bus, slave uint) (*os.File, error) {
	if file, err := os.OpenFile(SPIDevice(bus, slave), os.O_RDWR|os.O_SYNC, 0); err != nil {
		return nil, err
	} else {
		return file, nil
	}
}

func SPIMode(fd uintptr) (gopi.SPIMode, error) {
	var mode uint8
	if err := spi_ioctl(fd, SPI_IOC_RD_MODE, unsafe.Pointer(&mode)); err != 0 {
		return 0, os.NewSyscallError("spi_ioctl", err)
	} else {
		return gopi.SPIMode(mode), nil
	}
}

func SPISpeedHz(fd uintptr) (uint32, error) {
	var speed_hz uint32
	if err := spi_ioctl(fd, SPI_IOC_RD_MAX_SPEED_HZ, unsafe.Pointer(&speed_hz)); err != 0 {
		return 0, os.NewSyscallError("spi_ioctl", err)
	} else {
		return speed_hz, nil
	}
}

func SPIBitsPerWord(fd uintptr) (uint8, error) {
	var bits_per_word uint8
	if err := spi_ioctl(fd, SPI_IOC_RD_BITS_PER_WORD, unsafe.Pointer(&bits_per_word)); err != 0 {
		return 0, os.NewSyscallError("spi_ioctl", err)
	} else {
		return bits_per_word, nil
	}
}

func SPISetMode(fd uintptr, mode gopi.SPIMode) error {
	if err := spi_ioctl(fd, SPI_IOC_WR_MODE, unsafe.Pointer(&mode)); err != 0 {
		return os.NewSyscallError("spi_ioctl", err)
	} else {
		return nil
	}
}

func SPISetSpeedHz(fd uintptr, speed uint32) error {
	if err := spi_ioctl(fd, SPI_IOC_WR_MAX_SPEED_HZ, unsafe.Pointer(&speed)); err != 0 {
		return os.NewSyscallError("spi_ioctl", err)
	} else {
		return nil
	}
}

func SPISetBitsPerWord(fd uintptr, bits uint8) error {
	if err := spi_ioctl(fd, SPI_IOC_WR_BITS_PER_WORD, unsafe.Pointer(&bits)); err != 0 {
		return os.NewSyscallError("spi_ioctl", err)
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// TRANSFER, READ AND WRITE

func SPITransfer(fd uintptr, send []byte, speed uint32, delay uint16, bits uint8) ([]byte, error) {
	buffer_size := len(send)
	if buffer_size == 0 {
		return []byte{}, nil
	}
	recv := make([]byte, buffer_size)
	message := spi_message{
		tx_buf:        uint64(uintptr(unsafe.Pointer(&send[0]))),
		rx_buf:        uint64(uintptr(unsafe.Pointer(&recv[0]))),
		len:           uint32(buffer_size),
		speed_hz:      speed,
		delay_usecs:   delay,
		bits_per_word: bits,
	}

	if err := spi_ioctl(fd, uintptr(C._SPI_IOC_MESSAGE(C.int(1))), unsafe.Pointer(&message)); err != 0 {
		return nil, os.NewSyscallError("spi_ioctl", err)
	} else {
		return recv, nil
	}
}

func SPIRead(fd uintptr, recv []byte, speed uint32, delay uint16, bits uint8) error {
	if len(recv) == 0 {
		return nil
	}
	message := spi_message{
		tx_buf:        0,
		rx_buf:        uint64(uintptr(unsafe.Pointer(&recv[0]))),
		len:           uint32(len(recv)),
		speed_hz:      speed,
		delay_usecs:   delay,
		bits_per_word: bits,
	}
	if err := spi_ioctl(fd, uintptr(C._SPI_IOC_MESSAGE(C.int(1))), unsafe.Pointer(&message)); err != 0 {
		return os.NewSyscallError("spi_ioctl", err)
	} else {
		return nil
	}
}

func SPIWrite(fd uintptr, send []byte, speed uint32, delay uint16, bits uint8) error {
	buffer_size := len(send)
	if buffer_size == 0 {
		return nil
	}
	message := spi_message{
		tx_buf:        uint64(uintptr(unsafe.Pointer(&send[0]))),
		rx_buf:        0,
		len:           uint32(buffer_size),
		speed_hz:      speed,
		delay_usecs:   delay,
		bits_per_word: bits,
	}
	if err := spi_ioctl(fd, uintptr(C._SPI_IOC_MESSAGE(C.int(1))), unsafe.Pointer(&message)); err != 0 {
		return os.NewSyscallError("spi_ioctl", err)
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func spi_ioctl(fd uintptr, name uintptr, data unsafe.Pointer) syscall.Errno {
	_, _, err := syscall.RawSyscall(syscall.SYS_IOCTL, fd, name, uintptr(data))
	return err
}
