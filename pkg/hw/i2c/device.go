package i2c

import (
	"fmt"
	"os"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	bus   gopi.I2CBus
	fh    *os.File
	funcs I2CFunction
	slave uint8
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewDevice(bus gopi.I2CBus, fh *os.File, funcs I2CFunction) *device {
	return &device{bus, fh, funcs, 0}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *device) Fd() uintptr {
	if this.fh != nil {
		return this.fh.Fd()
	} else {
		return 0
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	str := "<device"
	str += " bus=" + fmt.Sprint(this.bus)
	if this.slave != 0 {
		str += fmt.Sprintf(" slave=0x%02X", this.slave)
	}
	if this.funcs != 0 {
		str += " funcs=" + fmt.Sprint(this.funcs)
	}
	return str + ">"
}
