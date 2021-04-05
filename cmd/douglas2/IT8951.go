package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	// Modules
	"github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type IT8951 struct {
	gopi.Unit
	gopi.Logger
	gopi.SPI
	gopi.GPIO

	bus gopi.SPIBus
}

type Command uint16

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	SPI_BUS   = 0
	SPI_SLAVE = 0
	SPI_SPEED = 24000000
	SPI_MODE  = gopi.SPI_MODE_0
)

const (
	PIN_CS    = gopi.GPIOPin(8)
	PIN_HRDY  = gopi.GPIOPin(24)
	PIN_RESET = gopi.GPIOPin(17)
	VCOM      = 1500 // e.g. -1.53 = 1530 = 0x5FA
)

const (
	CMD_SYS_RUN      Command = 0x0001
	CMD_STANDBY      Command = 0x0002
	CMD_SLEEP        Command = 0x0003
	CMD_REG_RD       Command = 0x0010
	CMD_REG_WR       Command = 0x0011
	CMD_MEM_BST_RD_T Command = 0x0012
	CMD_MEM_BST_RD_S Command = 0x0013
	CMD_MEM_BST_WR   Command = 0x0014
	CMD_MEM_BST_END  Command = 0x0015
	CMD_LD_IMG       Command = 0x0020
	CMD_LD_IMG_AREA  Command = 0x0021
	CMD_LD_IMG_END   Command = 0x0022
	CMD_DPY_AREA     Command = 0x0034
	CMD_GET_DEV_INFO Command = 0x0302
	CMD_DPY_BUF_AREA Command = 0x0037
	CMD_VCOM         Command = 0x0039
)

var (
	PREAMBLE_READ    = []byte{0x10, 0x00}
	PREAMBLE_WRITE   = []byte{0x00, 0x00}
	PREAMBLE_COMMAND = []byte{0x60, 0x00}
)

/////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *IT8951) Define(cfg gopi.Config) error {
	return nil
}

func (this *IT8951) New(cfg gopi.Config) error {
	this.Require(this.Logger, this.SPI, this.GPIO)

	// Check that GPIO is initialized
	if this.NumberOfPhysicalPins() == 0 {
		return gopi.ErrInternalAppError.WithPrefix("Missing GPIO interface")
	}

	// SPI Init
	this.bus = gopi.SPIBus{SPI_BUS, SPI_SLAVE}
	this.SPI.SetMode(this.bus, SPI_MODE)
	this.SPI.SetMaxSpeedHz(this.bus, SPI_SPEED)

	// GPIO Init
	this.GPIO.SetPinMode(PIN_HRDY, gopi.GPIO_INPUT)
	this.GPIO.SetPullMode(PIN_HRDY, gopi.GPIO_PULL_DOWN)
	this.GPIO.SetPinMode(PIN_RESET, gopi.GPIO_OUTPUT)
	this.GPIO.WritePin(PIN_RESET, gopi.GPIO_HIGH)

	// Init cycle should take less than 4 secs
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*4)
	defer cancel()

	// Reset
	if err := this.Reset(timeout); err != nil {
		return gopi.ErrNotFound.WithPrefix("IT8951")
	}

	// Get VCOM
	if vcom, err := this.GetVCOM(timeout); err != nil {
		return gopi.ErrNotFound.WithPrefix("IT8951")
	} else {
		this.Printf("VCOM=0x%04X", vcom)
	}

	/*
		// Get Device Info
		if err := this.GetDeviceInfo(timeout); errors.Is(err, context.DeadlineExceeded) {
			return gopi.ErrNotFound.WithPrefix("IT8951")
		} else if err != nil {
			return err
		}
	*/

	// Return success
	return nil
}

func (this *IT8951) Run(ctx context.Context) error {
	return nil
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *IT8951) Reset(ctx context.Context) error {
	this.Printf("Reset")
	this.GPIO.WritePin(PIN_RESET, gopi.GPIO_LOW)
	time.Sleep(200 * time.Millisecond)
	this.GPIO.WritePin(PIN_RESET, gopi.GPIO_HIGH)
	return this.WaitForReady(ctx)
}

func (this *IT8951) WaitForReady(ctx context.Context) error {
	this.Printf("WaitForReady")
	this.Print(this.GPIO.ReadPin(PIN_HRDY))
	for this.GPIO.ReadPin(PIN_HRDY) == gopi.GPIO_LOW {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			this.Print(this.GPIO.ReadPin(PIN_HRDY))
			time.Sleep(10 * time.Millisecond)
		}
	}
	this.Print(this.GPIO.ReadPin(PIN_HRDY))
	return nil
}

func (this *IT8951) WriteCommand(ctx context.Context, cmd Command) error {
	out := append(PREAMBLE_COMMAND, byte(cmd>>8), byte(cmd))
	this.Printf("WriteCommand 0x%v", hex.EncodeToString(out))
	if err := this.WaitForReady(ctx); err != nil {
		return err
	} else if err := this.SPI.Write(this.bus, out); err != nil {
		return err
	}
	return nil
}

func (this *IT8951) WriteData(ctx context.Context, data []byte) error {
	out := append(PREAMBLE_WRITE, data...)
	this.Printf("WriteData 0x%v", hex.EncodeToString(out))
	if err := this.WaitForReady(ctx); err != nil {
		return err
	} else if err := this.SPI.Write(this.bus, out); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *IT8951) ReadData(ctx context.Context, size uint32) ([]byte, error) {
	out := append(PREAMBLE_READ, make([]byte, size+2)...)
	this.Printf("ReadData  0x%v", hex.EncodeToString(out))
	if err := this.WaitForReady(ctx); err != nil {
		return nil, err
	} else if in, err := this.SPI.Transfer(this.bus, out); err != nil {
		return nil, err
	} else {
		return in[2:], nil
	}
}

func (this *IT8951) WriteUint16(ctx context.Context, word uint16) error {
	this.Printf("WriteUint16 0x%04X", word)
	return this.WriteData(ctx, []byte{byte(word >> 8), byte(word)})
}

func (this *IT8951) ReadUint16(ctx context.Context) (uint16, error) {
	this.Printf("ReadUint16")
	if in, err := this.ReadData(ctx, 2); err != nil {
		return 0, err
	} else {
		return uint16(in[0])<<8 | uint16(in[1]), nil
	}
}

func (this *IT8951) GetDeviceInfo(ctx context.Context) error {
	if err := this.WriteCommand(ctx, CMD_GET_DEV_INFO); err != nil {
		return err
	}
	if data, err := this.ReadData(ctx, 20); err != nil {
		return err
	} else {
		this.Print("GetDeviceInfo 0x", hex.EncodeToString(data))
	}
	return nil
}

func (this *IT8951) GetVCOM(ctx context.Context) (uint16, error) {
	this.Printf("GetVCOM")
	if err := this.WriteCommand(ctx, CMD_VCOM); err != nil {
		return 0, err
	} else if vcom, err := this.ReadUint16(ctx); err != nil {
		return 0, err
	} else {
		return vcom, nil
	}
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *IT8951) String() string {
	str := "<IT8951"
	str += fmt.Sprint(" gpio=", this.GPIO)
	str += fmt.Sprint(" spi=", this.SPI)
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func SPIUint16(v uint16) []byte {
	return []byte{byte(v >> 8), byte(v)}
}

func SPIByte(v uint8) []byte {
	return []byte{v}
}
