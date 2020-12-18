// +build linux

package spi

import (
	"os"

	gopi "github.com/djthorpe/gopi/v3"
	linux "github.com/djthorpe/gopi/v3/pkg/sys/linux"
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *spi) Devices() []gopi.SPIBus {
	devices := []gopi.SPIBus{}
	for bus := uint(0); bus <= maxBus; bus++ {
		if _, err := os.Stat(linux.SPIDevice(bus, 0)); os.IsNotExist(err) == false {
			devices = append(devices, gopi.SPIBus{bus, 0})
		}
		if _, err := os.Stat(linux.SPIDevice(bus, 1)); os.IsNotExist(err) == false {
			devices = append(devices, gopi.SPIBus{bus, 1})
		}
	}
	return devices
}

func (this *spi) Mode(bus gopi.SPIBus) gopi.SPIMode {
	if device, err := this.Open(bus); err != nil {
		return gopi.SPI_MODE_NONE
	} else {
		return device.Mode()
	}
}

func (this *spi) SetMode(bus gopi.SPIBus, mode gopi.SPIMode) error {
	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return device.SetMode(mode)
	}
}

func (this *spi) MaxSpeedHz(bus gopi.SPIBus) uint32 {
	if device, err := this.Open(bus); err != nil {
		return 0
	} else {
		return device.MaxSpeedHz()
	}
}

func (this *spi) SetMaxSpeedHz(bus gopi.SPIBus, speed_hz uint32) error {
	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return device.SetMaxSpeedHz(speed_hz)
	}
}

func (this *spi) BitsPerWord(bus gopi.SPIBus) uint8 {
	if device, err := this.Open(bus); err != nil {
		return 0
	} else {
		return device.BitsPerWord()
	}
}

func (this *spi) SetBitsPerWord(bus gopi.SPIBus, bits_per_word uint8) error {
	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return device.SetBitsPerWord(bits_per_word)
	}
}

func (this *spi) Transfer(bus gopi.SPIBus, data []byte) ([]byte, error) {
	if device, err := this.Open(bus); err != nil {
		return nil, err
	} else {
		return device.Transfer(data)
	}
}

func (this *spi) Read(bus gopi.SPIBus, data []byte) error {
	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return device.Read(data)
	}
}

func (this *spi) Write(bus gopi.SPIBus, data []byte) error {
	if device, err := this.Open(bus); err != nil {
		return err
	} else {
		return device.Write(data)
	}
}
