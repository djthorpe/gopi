// +build dvb

package dvb

import (
	"fmt"
	"os"

	dvb "github.com/djthorpe/gopi/v3/pkg/sys/dvb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tuner struct {
	dvb.Device
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewTuner(d dvb.Device) *Tuner {
	this := new(Tuner)
	this.Device = d
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Tuner) Id() uint {
	return this.Device.Adapter
}

func (this *Tuner) Name() string {
	return "TODO"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Tuner) OpenFrontend() (*os.File, error) {
	return this.Device.FEOpen()
}

func (this *Tuner) Validate(*Params) error {
	fmt.Println("TODO: Validate")
	// TODO
	return nil
}
