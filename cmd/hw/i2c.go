package main

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *app) RunI2C(ctx context.Context) error {
	if this.I2C == nil {
		return fmt.Errorf("No I2C interface enabled")
	} else if devices := this.I2C.Devices(); len(devices) == 0 {
		return fmt.Errorf("No I2C interface enabled")
	}

	bus, err := this.OpenI2C()
	if err != nil {
		return err
	}

	args := this.Args()
	if len(args) == 0 {
		return this.I2CDetectSlave(ctx, bus)
	}

	switch args[0] {
	case "write":
		return this.I2CWrite(ctx, bus, args[1:])
	default:
		return gopi.ErrBadParameter.WithPrefix(args[0])
	}
}

func (this *app) OpenI2C() (gopi.I2CBus, error) {
	// Set bus from argument or default
	bus := gopi.I2CBus(*this.i2cbus)
	if bus == 0 {
		bus = this.I2C.Devices()[0]
	}
	// Check to make sure I2C bus exists
	for _, device := range this.I2C.Devices() {
		if bus == device {
			return bus, nil
		}
	}
	// Bus not found
	return 0, gopi.ErrBadParameter.WithPrefix("-i2c.bus")
}

func (this *app) I2CDetectSlave(ctx context.Context, bus gopi.I2CBus) error {
	for slave := uint8(0x00); slave <= uint8(0x7F); slave++ {
		if slave%0x10 == 0x00 {
			fmt.Printf("%02X ", slave)
		}
		if detected, err := this.I2C.DetectSlave(bus, slave); err != nil {
			return err
		} else if detected {
			fmt.Printf("%02X ", slave)
		} else {
			fmt.Print("-- ")
		}
		if slave%0x10 == 0x0F {
			fmt.Println()
		}
	}

	return nil
}

func (this *app) I2CWrite(ctx context.Context, bus gopi.I2CBus, args []string) error {
	data := [][]byte{}
	for _, arg := range args {
		if bytes, err := hex.DecodeString(arg); err != nil {
			return err
		} else {
			data = append(data, bytes)
		}
	}
	for _, bytes := range data {
		if n, err := this.I2C.Write(bus, bytes); err != nil {
			return err
		} else {
			fmt.Println(hex.EncodeToString(bytes), "=>", n, "bytes written")
		}
	}
	return nil
}
