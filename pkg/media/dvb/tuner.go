// +build dvb

package dvb

import (
	"fmt"
	"os"

	gopi "github.com/djthorpe/gopi/v3"
	dvb "github.com/djthorpe/gopi/v3/pkg/sys/dvb"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tuner struct {
	dvb.Device
	dvb.FEInfo

	version string
	sys     []dvb.FEDeliverySystem
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewTuner(d dvb.Device) (*Tuner, error) {
	this := new(Tuner)

	// Read capabilities from tuner
	fh, err := d.FEOpen(os.O_RDONLY)
	if err != nil {
		return nil, err
	} else {
		this.Device = d
	}
	defer fh.Close()

	if info, err := dvb.FEGetInfo(fh.Fd()); err != nil {
		return nil, err
	} else {
		this.FEInfo = info
	}
	if major, minor, err := dvb.FEGetVersion(fh.Fd()); err != nil {
		return nil, err
	} else {
		this.version = fmt.Sprint(minor, ".", major)
	}
	if sys, err := dvb.FEEnumDeliverySystems(fh.Fd()); err != nil {
		return nil, err
	} else {
		this.sys = sys
	}

	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Tuner) Id() uint {
	return this.Device.Adapter
}

func (this *Tuner) Name() string {
	return this.FEInfo.Name()
}

func (this *Tuner) Version() string {
	return this.version
}

func (this *Tuner) Sys() []dvb.FEDeliverySystem {
	return this.sys
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Tuner) String() string {
	str := "<dvb.tuner"
	str += fmt.Sprint(" adapter=", this.Device.Adapter)
	if name := this.Name(); name != "" {
		str += fmt.Sprintf(" name=%q", name)
	}
	if version := this.Version(); version != "" {
		str += fmt.Sprintf(" version=%q", version)
	}
	if sys := this.Sys(); len(sys) > 0 {
		str += fmt.Sprint(" sys=", sys)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// OpenFrontend returns a file descriptor for frontend
func (this *Tuner) OpenFrontend() (*os.File, error) {
	return this.Device.FEOpen(os.O_RDWR)
}

// OpenDemux returns a file descriptor for demux
func (this *Tuner) OpenDemux() (*os.File, error) {
	return this.Device.DMXOpen()
}

// Validate determines if parameters are supported by the tuner
func (this *Tuner) Validate(params *Params) error {
	if this.hasDeliverySystem(params) == false {
		return gopi.ErrBadParameter.WithPrefix("DeliverySystem")
	}

	// TODO: Validate more params here
	fmt.Println("TODO: Validate: More validation here")

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Tuner) hasDeliverySystem(params *Params) bool {
	for _, supported := range this.sys {
		if params.DeliverySystem == supported {
			return true
		}
	}
	return false
}
