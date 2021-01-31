// +build dvb

package dvb

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	ts "github.com/djthorpe/gopi/v3/pkg/media/internal/ts"
	dvb "github.com/djthorpe/gopi/v3/pkg/sys/dvb"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Tuner struct {
	sync.RWMutex
	dvb.Device
	dvb.FEInfo

	version string
	sys     []dvb.FEDeliverySystem
	dev     *os.File
	section map[*SectionFilter]bool
	stream  map[*StreamFilter]bool
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

	// Add filter maps
	this.section = make(map[*SectionFilter]bool)
	this.stream = make(map[*StreamFilter]bool)

	// Return success
	return this, nil
}

func (this *Tuner) Dispose(fp gopi.FilePoll) error {
	var result error

	for filter := range this.section {
		if err := this.RemoveSectionFilter(fp, filter); err != nil {
			result = multierror.Append(result, err)
		}
	}
	for filter := range this.stream {
		if err := this.RemoveStreamFilter(fp, filter); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Close device
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if this.dev != nil {
		if err := this.dev.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.dev = nil
	this.section = nil
	this.stream = nil

	// Return any errors
	return result
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
// PUBLIC METHODS

func (this *Tuner) NewSectionFilter(pid uint16, tid ts.TableType, flags dvb.DMXFlag) (*SectionFilter, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Create filter
	filter, err := NewSectionFilter(this, pid, tid, flags)
	if err != nil {
		return nil, err
	}

	// Add filter to list of filters
	this.section[filter] = true

	// Return filter
	return filter, nil
}

func (this *Tuner) NewStreamFilter(pid uint16, in dvb.DMXInput, out dvb.DMXOutput, stream dvb.DMXStreamType) (*StreamFilter, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Create filter
	filter, err := NewStreamFilter(this, pid, in, out, stream)
	if err != nil {
		return nil, err
	}

	// Add filter to list of filters
	this.stream[filter] = true

	// Return filter
	return filter, nil
}

func (this *Tuner) RemoveSectionFilter(fp gopi.FilePoll, filter *SectionFilter) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check parameters
	if _, exists := this.section[filter]; exists == false {
		return gopi.ErrNotFound.WithPrefix("RemoveSectionFilter")
	} else {
		delete(this.section, filter)
	}

	// Unwatch, dispose
	var result error
	if err := fp.Unwatch(filter.Fd()); err != nil {
		result = multierror.Append(result, err)
	}
	if err := filter.Dispose(); err != nil {
		result = multierror.Append(result, err)
	}

	// Return errors
	return result
}

func (this *Tuner) RemoveStreamFilter(fp gopi.FilePoll, filter *StreamFilter) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check parameters
	if _, exists := this.stream[filter]; exists == false {
		return gopi.ErrNotFound.WithPrefix("RemoveStreamFilter")
	} else {
		delete(this.stream, filter)
	}

	// Unwatch, dispose
	var result error
	if err := fp.Unwatch(filter.Fd()); err != nil {
		result = multierror.Append(result, err)
	}
	if err := filter.Dispose(); err != nil {
		result = multierror.Append(result, err)
	}

	// Return errors
	return result
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
func (this *Tuner) OpenFrontend() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.dev != nil {
		return nil
	} else if dev, err := this.Device.FEOpen(os.O_RDWR); err != nil {
		return err
	} else {
		this.dev = dev
	}

	// Return success
	return nil
}

// Tune applies parameters to the tuner
func (this *Tuner) Tune(ctx context.Context, params *Params) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.dev == nil {
		return gopi.ErrOutOfOrder.WithPrefix("Tune")
	} else if err := dvb.FETune(this.dev.Fd(), params.TuneParams); err != nil {
		return err
	}

	// Now loop until status is tuned or timeout
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
FOR_LOOP:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			status, err := dvb.FEReadStatus(this.dev.Fd())
			if err != nil {
				return err
			}
			switch {
			case status&dvb.FE_HAS_LOCK == dvb.FE_HAS_LOCK:
				break FOR_LOOP
			case status == dvb.FE_NONE:
				// Do nothing, no tune status
			}
		}
	}

	// Return success
	return nil
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
