package waveshare

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	"golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

type EPD struct {
	gopi.Unit
	gopi.GPIO
	gopi.SPI

	bus    gopi.SPIBus
	w, h   *uint
	rotate *int
}

const (
	EPD_PIN_RESET = gopi.GPIOPin(17)
	EPD_PIN_CS    = gopi.GPIOPin(8)
	EPD_PIN_DC    = gopi.GPIOPin(25)
	EPD_PIN_BUSY  = gopi.GPIOPin(24)

	EPD_SPI_SPEED = 4000000
	EPD_SPI_BUS   = 0
	EPD_SPI_SLAVE = 0
	EPD_SPI_MODE  = gopi.SPI_MODE_0
)

const (
	EPD_CMD_SLEEP_MODE      = 0x10
	EPD_CMD_SWRESET         = 0x12
	EPD_CMD_RAM_WRITE_BLACK = 0x24
	EPD_CMD_RAM_WRITE_RED   = 0x26
	EPD_CMD_RAM_XADDRESS    = 0x4E
	EPD_CMD_RAM_YADDRESS    = 0x4F
	EPD_CMD_UPDATE          = 0x20
)

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *EPD) Define(cfg gopi.Config) {
	this.w = cfg.FlagUint("epd.width", 880, "Width of display")
	this.h = cfg.FlagUint("epd.height", 528, "Height of display")
	this.rotate = cfg.FlagInt("epd.rotate", 0, "Image rotation in degrees")
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
	this.bus = gopi.SPIBus{EPD_SPI_BUS, EPD_SPI_SLAVE}

	// Initialise the interfaces
	if err := this.init(); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *EPD) Dispose() error {
	// Put EPD into sleep mode
	if err := this.Sleep(); err != nil {
		return err
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *EPD) Size() gopi.Size {
	return gopi.Size{float32(*this.w), float32(*this.h)}
}

func (this *EPD) Clear(ctx context.Context) error {
	width := *this.w
	height := *this.h
	stride := width >> 3 // bytes per row

	// Set RAM X,Y Addresses to zero
	this.send(EPD_CMD_RAM_XADDRESS, []byte{0x00, 0x00})
	this.send(EPD_CMD_RAM_YADDRESS, []byte{0x00, 0x00})

	// Send data
	buf := make([]byte, stride*height)
	for i := range buf {
		buf[i] = 0xFF
	}
	this.send(EPD_CMD_RAM_WRITE_BLACK, buf)
	this.send(EPD_CMD_RAM_WRITE_RED, buf)

	// Load LUT from MCU(0x32)
	this.send(0x22, []byte{0xF7})
	this.send(EPD_CMD_UPDATE, nil)
	time.Sleep(10 * time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := this.waitUntilIdle(ctx); err != nil {
		return err
	}

	return nil
}

// DrawMono assumes image is already coded as black and white
// pixels and that the image is the correct size
func (this *EPD) DrawMono(ctx context.Context, img image.Image) error {
	width := *this.w
	height := *this.h
	stride := width >> 3 // bytes per row

	// Set RAM X,Y Addresses to zero
	this.send(EPD_CMD_RAM_XADDRESS, []byte{0x00, 0x00})
	this.send(EPD_CMD_RAM_YADDRESS, []byte{0x00, 0x00})

	// Construct bit-per-pixel image
	buf := make([]byte, stride*height)
	for y := uint(0); y < height; y++ {
		for x := uint(0); x < stride; x++ {
			data := uint8(0)
			for bit := uint(0); bit < 8; bit++ {
				data <<= 1
				c := img.At(int(x*8+bit), int(y))
				if y, ok := c.(color.Gray); ok && y.Y != 0 {
					data |= 1
				}
			}
			buf[x+y*stride] = data
		}
	}
	// Write blacks
	this.send(EPD_CMD_RAM_WRITE_BLACK, buf)

	// Load LUT from MCU(0x32)
	this.send(0x22, []byte{0xF7})
	this.send(EPD_CMD_UPDATE, nil)
	time.Sleep(10 * time.Millisecond)

	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := this.waitUntilIdle(ctx); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *EPD) Draw(ctx context.Context, src image.Image) error {
	return this.DrawSized(ctx, 1.0, src)
}

func (this *EPD) DrawSized(ctx context.Context, scale float64, src image.Image) error {
	bounds := image.Rectangle{image.ZP, image.Pt(int(*this.w), int(*this.h))}

	// Create image for the framesize
	scaled := image.NewRGBA(bounds)

	// Determine the best way to fit the image into the frame. We prefer full
	// height images
	rsrc := float64(src.Bounds().Dx()) / float64(src.Bounds().Dy())
	rdst := float64(scaled.Bounds().Dx()) / float64(scaled.Bounds().Dy())
	if rdst > rsrc {
		scale = scale * float64(scaled.Bounds().Dy()) / float64(src.Bounds().Dy())
	} else {
		scale = scale * float64(scaled.Bounds().Dx()) / float64(src.Bounds().Dx())
	}
	transform := NewAffineTransform().Scale(scale, scale)
	draw.ApproxBiLinear.Transform(scaled, f64.Aff3(transform), src, src.Bounds(), draw.Over, nil)

	// Shift image into the middle of the frame
	xdiff := float64(scaled.Bounds().Dx()) - float64(src.Bounds().Dx())*scale
	ydiff := float64(scaled.Bounds().Dy()) - float64(src.Bounds().Dy())*scale
	transform = NewAffineTransform().Translate(xdiff/2, ydiff/2)
	shifted := image.NewRGBA(scaled.Bounds())
	draw.ApproxBiLinear.Transform(shifted, f64.Aff3(transform), scaled, scaled.Bounds(), draw.Over, nil)
	scaled = shifted

	// Perform rotations
	if theta, cx, cy := getRotation(*this.rotate, scaled); theta != 0 {
		transform := NewAffineTransform().Rotate(theta, cx, cy)
		rotated := image.NewRGBA(bounds)
		draw.ApproxBiLinear.Transform(rotated, f64.Aff3(transform), scaled, scaled.Bounds(), draw.Over, nil)
		scaled = rotated
	}

	// Convert to BW using dithering
	dst := image.NewPaletted(scaled.Bounds(), []color.Color{
		color.Gray{Y: 255},
		color.Gray{Y: 0},
	})
	draw.FloydSteinberg.Draw(dst, dst.Bounds(), scaled, image.ZP)

	return this.DrawMono(ctx, dst)
}

func (this *EPD) Sleep() error {
	// Sleep mode 1. To wake up, hardware reset is required
	return this.send(EPD_CMD_SLEEP_MODE, []byte{0x01})
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *EPD) init() error {
	// GPIO Init
	this.GPIO.SetPinMode(EPD_PIN_RESET, gopi.GPIO_OUTPUT)
	this.GPIO.SetPinMode(EPD_PIN_DC, gopi.GPIO_OUTPUT)
	this.GPIO.SetPinMode(EPD_PIN_CS, gopi.GPIO_OUTPUT)
	this.GPIO.SetPinMode(EPD_PIN_BUSY, gopi.GPIO_INPUT)

	// SPI Init
	this.SPI.SetMode(this.bus, EPD_SPI_MODE)
	this.SPI.SetMaxSpeedHz(this.bus, EPD_SPI_SPEED)

	// Toggle reset pin and wait until idle
	this.reset()
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Software reset
	this.send(EPD_CMD_SWRESET, nil)
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Auto Write Red RAM
	this.send(0x46, []byte{0xF7})
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Auto Write B/W RAM
	this.send(0x47, []byte{0xF7})
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

	// Load Temperature and Waveform setting
	this.send(0x22, []byte{0xB1})
	this.send(0x20, nil)
	if err := this.waitUntilIdleTimeout(time.Second); err != nil {
		return err
	}

	// Set XY Counters to zero
	this.send(EPD_CMD_RAM_XADDRESS, []byte{0x00, 0x00})
	this.send(EPD_CMD_RAM_YADDRESS, []byte{0x00, 0x00})

	// Return sucess
	return nil
}

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

	this.GPIO.WritePin(EPD_PIN_DC, gopi.GPIO_HIGH)
	this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_LOW)
	defer this.GPIO.WritePin(EPD_PIN_CS, gopi.GPIO_HIGH)
	for _, b := range data {
		if err := this.SPI.Write(this.bus, []byte{b}); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

// Toggle the hardware reset pin
func (this *EPD) reset() {
	this.GPIO.WritePin(EPD_PIN_RESET, gopi.GPIO_HIGH)
	time.Sleep(200 * time.Millisecond)
	this.GPIO.WritePin(EPD_PIN_RESET, gopi.GPIO_LOW)
	time.Sleep(2 * time.Millisecond)
	this.GPIO.WritePin(EPD_PIN_RESET, gopi.GPIO_HIGH)
	time.Sleep(200 * time.Millisecond)
}

func getRotation(deg int, src image.Image) (float64, float64, float64) {
	theta := math.Pi * float64(deg) / 180
	cx := float64(src.Bounds().Max.X+src.Bounds().Min.X) / 2.0
	cy := float64(src.Bounds().Max.Y+src.Bounds().Min.Y) / 2.0
	return theta, cx, cy
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *EPD) String() string {
	str := "<epd"
	str += " size=" + fmt.Sprint(this.Size())
	str += " gpio=" + fmt.Sprint(this.GPIO)
	str += " spi=" + fmt.Sprint(this.SPI)
	return str + ">"
}
