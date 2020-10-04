// +build linux

package spi

import (
	"fmt"
	"os"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Spi struct {
	gopi.Unit
	sync.Mutex

	bus, slave, delay_usec *uint
	dev                    *os.File
	mode                   gopi.SPIMode
	speed_hz               uint32
	bits_per_word          uint8
}

////////////////////////////////////////////////////////////////////////////////

func (this *Spi) Define(cfg gopi.Config) error {
	this.bus = cfg.FlagUint("spi.bus", 0, "SPI Bus")
	this.slave = cfg.FlagUint("spi.slave", 0, "SPI Slave")
	this.delay_usec = cfg.FlagUint("spi.delay", 0, "SPI Transfer delay in microseconds")
	return nil
}

func (this *Spi) New(_ gopi.Config) error {
	// Open the device
	if dev, err := linux.SPIOpenDevice(*this.bus, *this.slave); err != nil {
		return err
	} else {
		this.dev = dev
	}

	// Get current mode, speed and bits per word
	if mode, err := linux.SPIMode(this.dev.Fd()); err != nil {
		return err
	} else {
		this.mode = mode
	}
	if speed_hz, err := linux.SPISpeedHz(this.dev.Fd()); err != nil {
		return err
	} else {
		this.speed_hz = speed_hz
	}
	if bits_per_word, err := linux.SPIBitsPerWord(this.dev.Fd()); err != nil {
		return err
	} else {
		this.bits_per_word = bits_per_word
	}

	// Success
	return nil
}

func (this *Spi) Dispose() error {
	if this.dev != nil {
		if err := this.dev.Close(); err != nil {
			return err
		}
	}

	// Release resources
	this.dev = nil

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Spi) String() string {
	str := "<spi"
	if this.bus != nil {
		str += " bus=" + fmt.Sprint(*this.bus)
	}
	if this.slave != nil {
		str += " slave=" + fmt.Sprint(*this.slave)
	}
	if this.mode != gopi.SPI_MODE_NONE {
		str += " mode=" + fmt.Sprint(this.mode)
	}
	if this.speed_hz != 0 {
		str += " max_speed=" + fmt.Sprint(this.speed_hz) + "Hz"
	}
	if this.delay_usec != nil {
		str += " delay=" + fmt.Sprint(uint16(*this.delay_usec)) + "us"
	}
	if this.bits_per_word != 0 {
		str += " bits_per_word=" + fmt.Sprint(this.bits_per_word)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.SPI

func (this *Spi) SetMode(mode gopi.SPIMode) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetMode(this.dev.Fd(), mode); err != nil {
		return err
	} else if mode_, err := linux.SPIMode(this.dev.Fd()); err != nil {
		return err
	} else if mode != mode_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("mode")
	} else {
		this.mode = mode
		return nil
	}
}

func (this *Spi) SetMaxSpeedHz(speed uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetSpeedHz(this.dev.Fd(), speed); err != nil {
		return err
	} else if speed_, err := linux.SPISpeedHz(this.dev.Fd()); err != nil {
		return err
	} else if speed != speed_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("speed")
	} else {
		this.speed_hz = speed
		return nil
	}
}

func (this *Spi) SetBitsPerWord(bits uint8) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetBitsPerWord(this.dev.Fd(), bits); err != nil {
		return err
	} else if bits_, err := linux.SPIBitsPerWord(this.dev.Fd()); err != nil {
		return err
	} else if bits != bits_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("bits")
	} else {
		this.bits_per_word = bits
		return nil
	}
}

func (this *Spi) Transfer(send []byte) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(send) == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("send")
	} else if receive, err := linux.SPITransfer(this.dev.Fd(), send, this.speed_hz, uint16(*this.delay_usec), this.bits_per_word); err != nil {
		return nil, err
	} else {
		return receive, nil
	}
}

func (this *Spi) Read(len uint32) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("len")
	} else if receive, err := linux.SPIRead(this.dev.Fd(), len, this.speed_hz, uint16(*this.delay_usec), this.bits_per_word); err != nil {
		return nil, err
	} else {
		return receive, nil
	}
}

func (this *Spi) Write(send []byte) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(send) == 0 {
		return gopi.ErrBadParameter.WithPrefix("send")
	} else if err := linux.SPIWrite(this.dev.Fd(), send, this.speed_hz, uint16(*this.delay_usec), this.bits_per_word); err != nil {
		return err
	} else {
		return nil
	}
}
