// +build linux

package spi

import (
	"fmt"
	"os"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	sync.Mutex

	delay_usec    uint16
	fd            *os.File
	mode          gopi.SPIMode
	speed_hz      uint32
	bits_per_word uint8
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewDevice(bus gopi.SPIBus, delay_usec uint16) (*device, error) {
	dev := new(device)
	if fd, err := linux.SPIOpenDevice(bus.Bus, bus.Slave); err != nil {
		return nil, err
	} else {
		dev.fd = fd
		dev.delay_usec = delay_usec
	}
	if mode, err := linux.SPIMode(dev.fd.Fd()); err != nil {
		defer dev.fd.Close()
		return nil, err
	} else {
		dev.mode = mode
	}
	if speed_hz, err := linux.SPISpeedHz(dev.fd.Fd()); err != nil {
		defer dev.fd.Close()
		return nil, err
	} else {
		dev.speed_hz = speed_hz
	}
	if bits_per_word, err := linux.SPIBitsPerWord(dev.fd.Fd()); err != nil {
		defer dev.fd.Close()
		return nil, err
	} else {
		dev.bits_per_word = bits_per_word
	}

	return dev, nil
}

func (dev *device) Close() error {
	dev.Mutex.Lock()
	defer dev.Mutex.Unlock()

	var result error
	if dev.fd != nil {
		result = dev.fd.Close()
	}
	dev.fd = nil

	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (dev *device) String() string {
	str := "<spi"
	if dev.mode != gopi.SPI_MODE_NONE {
		str += " mode=" + fmt.Sprint(dev.mode)
	}
	if dev.speed_hz != 0 {
		str += " max_speed=" + fmt.Sprint(dev.speed_hz) + "Hz"
	}
	if dev.delay_usec != 0 {
		str += " delay=" + fmt.Sprint(uint16(dev.delay_usec)) + "us"
	}
	if dev.bits_per_word != 0 {
		str += " bits_per_word=" + fmt.Sprint(dev.bits_per_word)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *device) Mode() gopi.SPIMode {
	return this.mode
}

func (this *device) MaxSpeedHz() uint32 {
	return this.speed_hz
}

func (this *device) BitsPerWord() uint8 {
	return this.bits_per_word
}

func (this *device) SetMode(mode gopi.SPIMode) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetMode(this.fd.Fd(), mode); err != nil {
		return err
	} else if mode_, err := linux.SPIMode(this.fd.Fd()); err != nil {
		return err
	} else if mode != mode_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("SetMode")
	} else {
		this.mode = mode
		return nil
	}
}

func (this *device) SetMaxSpeedHz(speed uint32) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetSpeedHz(this.fd.Fd(), speed); err != nil {
		return err
	} else if speed_, err := linux.SPISpeedHz(this.fd.Fd()); err != nil {
		return err
	} else if speed != speed_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("SetMaxSpeedHz")
	} else {
		this.speed_hz = speed
		return nil
	}
}

func (this *device) SetBitsPerWord(bits_per_word uint8) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if err := linux.SPISetBitsPerWord(this.fd.Fd(), bits_per_word); err != nil {
		return err
	} else if bits_per_word_, err := linux.SPIBitsPerWord(this.fd.Fd()); err != nil {
		return err
	} else if bits_per_word != bits_per_word_ {
		return gopi.ErrUnexpectedResponse.WithPrefix("SetBitsPerWord")
	} else {
		this.bits_per_word = bits_per_word
		return nil
	}
}

func (this *device) Transfer(data []byte) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(data) == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("Transfer")
	} else if receive, err := linux.SPITransfer(this.fd.Fd(), data, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return nil, err
	} else {
		return receive, nil
	}
}

func (this *device) Read(data []byte) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(data) == 0 {
		return gopi.ErrBadParameter.WithPrefix("Read")
	} else if err := linux.SPIRead(this.fd.Fd(), data, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *device) Write(data []byte) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(data) == 0 {
		return gopi.ErrBadParameter.WithPrefix("Write")
	} else if err := linux.SPIWrite(this.fd.Fd(), data, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return err
	} else {
		return nil
	}
}
