package waveshare

import (
	"context"
	"fmt"
	"image"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/gpio/broadcom"
)

type EPD struct {
	gopi.Unit
	gopi.GPIO
	gopi.SPI

	bus  gopi.SPIBus
	w, h *uint
}

const (
	EPD_PIN_RESET = gopi.GPIOPin(17)
	EPD_PIN_CS    = gopi.GPIOPin(8)
	EPD_PIN_DC    = gopi.GPIOPin(25)
	EPD_PIN_BUSY  = gopi.GPIOPin(24)

	EPD_SPI_SPEED = 10000000
	EPD_SPI_BUS   = 0
	EPD_SPI_SLAVE = 0
	EPD_SPI_MODE  = gopi.SPI_MODE_0
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *EPD) Define(cfg gopi.Config) {
	this.w = cfg.FlagUint("epd.width", 880, "Width of display")
	this.h = cfg.FlagUint("epd.height", 528, "Height of display")
}

func (this *EPD) New(gopi.Config) error {
	if this.GPIO == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing GPIO interface")
	} else if this.GPIO.NumberOfPhysicalPins() == 0 {
		return gopi.ErrInternalAppError.WithPrefix("Missing GPIO interface")
	}
	if this.SPI == nil {
		return gopi.ErrInternalAppError.WithPrefix("Missing SPI interface")
	}

	// Set SPI bus
	this.bus = gopi.SPIBus{0, 0}

	if err := this.Init(); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *EPD) Dispose() error {
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *EPD) Init() error {
	// GPIO Init
	this.GPIO.SetPinMode(EPD_PIN_RESET, gopi.GPIO_OUTPUT)
	this.GPIO.SetPinMode(EPD_PIN_DC, gopi.GPIO_OUTPUT)
	this.GPIO.SetPinMode(EPD_PIN_CS, gopi.GPIO_OUTPUT)
	this.GPIO.SetPinMode(EPD_PIN_BUSY, gopi.GPIO_INPUT)

	// SPI Init
	this.SPI.SetMode(this.bus, gopi.SPI_MODE_0)
	this.SPI.SetMaxSpeedHz(this.bus, 4000000)
	// TODO: Endian, Polarity

	// Toggle reset pin and wait until idle
	this.reset()
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Software reset
	this.send(0x12, nil)
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Auto Write Red RAM
	this.send(0x46, []byte{0xf7})
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Auto Write B/W RAM
	this.send(0x47, []byte{0xf7})
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Soft start setting
	this.send(0x0C, []byte{0xAE, 0xC7, 0xC3, 0xC0, 0x40})

	// Set MUX as 527
	this.send(0x01, []byte{0xAF, 0x02, 0x01})

	// Data entry mode
	this.send(0x11, []byte{0x01})

	// RAM x address start at 0
	this.send(0x44, []byte{0x00, 0x00, 0x6F, 0x03})
	this.send(0x45, []byte{0xFF, 0x03, 0x00, 0x00})

	// VBD, LUT1, for white
	this.send(0x3C, []byte{0x05})
	this.send(0x18, []byte{0x80})

	// Load Temperature and waveform setting
	this.send(0x22, []byte{0xB1})
	this.send(0x20, nil)
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Set RAM x address count to 0
	this.send(0x4E, []byte{0x00, 0x00})
	this.send(0x4F, []byte{0x00, 0x00})

	// Return sucess
	return nil
}

func (this *EPD) Clear(ctx context.Context) error {
	width := *this.w >> 3 // Divide by eight
	height := *this.h

	// Set RAM x address count to 0
	this.send(0x4F, []byte{0x00, 0x00})

	// Send data
	buf := make([]byte, width*height)
	for i := range buf {
		buf[i] = 0xFF
	}
	this.send(0x24, buf)
	this.send(0x26, buf)

	// Load LUT from MCU(0x32)
	this.send(0x22, []byte{0xF7})
	this.send(0x20, nil)
	time.Sleep(10 * time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := this.waitUntilIdle(ctx); err != nil {
		return err
	}

	return nil
}

func (this *EPD) Display(ctx context.Context, img image.Image) error {
	width := *this.w >> 3 // Divide by eight
	height := *this.h

	// Set RAM x address count to 0
	this.send(0x4F, []byte{0x00, 0x00})

	buf := make([]byte, width*height)
	for y := uint(0); y < height; y++ {
		for x := uint(0); x < width; x++ {
			//r, g, b, _ := img.At(int(x), int(y)).RGBA()
			r, _, _, _ := img.At(int(x), int(y)).RGBA()
			buf[x+y*width] = uint8(r) //uint8((uint32(r>>24) + uint32(g>>24) + uint32(b>>24)) / 3)
		}
	}
	this.send(0x24, buf)
	for i := range buf {
		buf[i] = 0xFF
	}
	this.send(0x26, buf)

	// Load LUT from MCU(0x32)
	this.send(0x22, []byte{0xF7})
	this.send(0x20, nil)
	time.Sleep(10 * time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := this.waitUntilIdle(ctx); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *EPD) Sleep() {
	this.sendCommand(0x10)
	this.sendData(0x01)
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// waitUntilIdle waits until busy pin goes low
func (this *EPD) waitUntilIdle(ctx context.Context) error {
	ticker := time.NewTimer(time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if this.GPIO.ReadPin(EPD_PIN_BUSY) == gopi.GPIO_LOW {
				time.Sleep(200 * time.Millisecond)
				return nil
			}
			ticker.Reset(10 * time.Millisecond)
		}
	}
}

// waitUntilIdleTimeout waits until busy pin goes low with timeout
func (this *EPD) waitUntilIdleTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return this.waitUntilIdle(ctx)
}

// send command and then data
func (this *EPD) send(reg uint8, data []byte) error {
	this.GPIO.WritePin(EPD_PIN_DC, gopi.GPIO_LOW)
	this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_LOW)
	if err := this.SPI.Write(this.bus, []byte{reg}); err != nil {
		return err
	}
	this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_HIGH)

	for _, b := range data {
		this.GPIO.WritePin(EPD_PIN_DC, gopi.GPIO_HIGH)
		this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_LOW)
		if err := this.SPI.Write(this.bus, []byte{b}); err != nil {
			return err
		}
		this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_HIGH)
	}

	// Return sucess
	return nil
}

// sendCommand sends a command byte
func (this *EPD) sendCommand(reg uint8) {
	this.GPIO.WritePin(EPD_PIN_DC, gopi.GPIO_LOW)
	this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_LOW)
	this.SPI.Write(this.bus, []byte{reg})
	this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_HIGH)
}

// sendData sends a data byte
func (this *EPD) sendData(data uint8) {
	this.GPIO.WritePin(EPD_PIN_DC, gopi.GPIO_HIGH)
	this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_LOW)
	this.SPI.Write(this.bus, []byte{data})
	this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_HIGH)
}

// reset toggles the reset pin
func (this *EPD) reset() {
	this.GPIO.WritePin(EPD_PIN_RESET, gopi.GPIO_HIGH)
	time.Sleep(200 * time.Millisecond)
	this.GPIO.WritePin(EPD_PIN_RESET, gopi.GPIO_LOW)
	time.Sleep(2 * time.Millisecond)
	this.GPIO.WritePin(EPD_PIN_RESET, gopi.GPIO_HIGH)
	time.Sleep(200 * time.Millisecond)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *EPD) String() string {
	str := "<epd"
	str += " gpio=" + fmt.Sprint(this.GPIO)
	str += " spi=" + fmt.Sprint(this.SPI)
	return str + ">"
}
