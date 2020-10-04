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

type spi struct {
	sync.Mutex
	Device

	delay_usec    uint16
	dev           *os.File
	mode          gopi.SPIMode
	speed_hz      uint32
	bits_per_word uint8
}

////////////////////////////////////////////////////////////////////////////////
// Globals

const (
	MaxBus = 9 // Maximum bus number
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Devices) Enumerate() []Device {
	devices := []Device{}
	for bus := uint(0); bus <= MaxBus; bus++ {
		if _, err := os.Stat(linux.SPIDevice(bus, 0)); os.IsNotExist(err) == false {
			devices = append(devices, Device{bus, 0})
		}
		if _, err := os.Stat(linux.SPIDevice(bus, 1)); os.IsNotExist(err) == false {
			devices = append(devices, Device{bus, 1})
		}
	}
	return devices
}

// Open SPI device
func (this *Devices) Open(dev Device, delay uint16) (gopi.SPI, error) {
	// Check for already opened
	if this.Get(dev) != nil {
		return nil, gopi.ErrBadParameter
	}

	// Open the device
	spi := new(spi)
	if fd, err := linux.SPIOpenDevice(dev.Bus, dev.Slave); err != nil {
		return nil, err
	} else {
		spi.dev = fd
		spi.delay_usec = delay
	}

	this.RWMutex.Lock()
	this.devices[dev] = spi
	this.RWMutex.Unlock()

	// Get current mode, speed and bits per word
	if mode, err := linux.SPIMode(spi.dev.Fd()); err != nil {
		return nil, err
	} else {
		spi.mode = mode
	}
	if speed_hz, err := linux.SPISpeedHz(spi.dev.Fd()); err != nil {
		return nil, err
	} else {
		spi.speed_hz = speed_hz
	}
	if bits_per_word, err := linux.SPIBitsPerWord(spi.dev.Fd()); err != nil {
		return nil, err
	} else {
		spi.bits_per_word = bits_per_word
	}

	// Return success
	return spi, nil
}

func (this *Devices) Close(dev Device) error {
	// Get the device
	spi, _ := this.Get(dev).(*spi)
	if spi == nil {
		return gopi.ErrBadParameter
	} else {
		this.Delete(dev)
	}

	// Close the filehandle
	spi.Mutex.Lock()
	defer spi.Mutex.Unlock()
	if spi.dev != nil {
		if err := spi.dev.Close(); err != nil {
			return err
		} else {
			spi.dev = nil
		}
	}

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *spi) String() string {
	str := "<spi"
	if this.dev != nil {
		if this.mode != gopi.SPI_MODE_NONE {
			str += " mode=" + fmt.Sprint(this.mode)
		}
		if this.speed_hz != 0 {
			str += " max_speed=" + fmt.Sprint(this.speed_hz) + "Hz"
		}
		if this.delay_usec != 0 {
			str += " delay=" + fmt.Sprint(uint16(this.delay_usec)) + "us"
		}
		if this.bits_per_word != 0 {
			str += " bits_per_word=" + fmt.Sprint(this.bits_per_word)
		}
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.SPI

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

func (this *spi) SetMaxSpeedHz(speed uint32) error {
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

func (this *spi) SetBitsPerWord(bits uint8) error {
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

func (this *spi) Transfer(send []byte) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(send) == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("send")
	} else if receive, err := linux.SPITransfer(this.dev.Fd(), send, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return nil, err
	} else {
		return receive, nil
	}
}

func (this *spi) Read(len uint32) ([]byte, error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("len")
	} else if receive, err := linux.SPIRead(this.dev.Fd(), len, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return nil, err
	} else {
		return receive, nil
	}
}

func (this *spi) Write(send []byte) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if len(send) == 0 {
		return gopi.ErrBadParameter.WithPrefix("send")
	} else if err := linux.SPIWrite(this.dev.Fd(), send, this.speed_hz, this.delay_usec, this.bits_per_word); err != nil {
		return err
	} else {
		return nil
	}
}
