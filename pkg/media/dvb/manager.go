// +build dvb

package dvb

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	dvb "github.com/djthorpe/gopi/v3/pkg/sys/dvb"
	"github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Logger
	sync.RWMutex

	frontend map[uint]*os.File
	demux    map[uint]*os.File
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger)

	// Set up file descriptor maps
	this.frontend = make(map[uint]*os.File)
	this.demux = make(map[uint]*os.File)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error

	// Close file descriptors
	for _, fh := range this.demux {
		if err := fh.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}
	for _, fh := range this.frontend {
		if err := fh.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.demux = nil
	this.frontend = nil

	// Return any errors
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) Tuners() []gopi.DVBTuner {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	tuners := []gopi.DVBTuner{}
	for _, device := range dvb.Devices() {
		if tuner, err := NewTuner(device); err != nil {
			this.Debug("Tuners:", err)
		} else {
			tuners = append(tuners, tuner)
		}
	}
	return tuners
}

func (this *Manager) ParseTunerParams(r io.Reader) ([]gopi.DVBTunerParams, error) {
	result := []gopi.DVBTunerParams{}
	params, err := dvb.ReadTuneParamsTable(r)
	if err != nil {
		return nil, err
	}
	for _, param := range params {
		result = append(result, NewParams(param))
	}
	return result, nil
}

func (this *Manager) Tune(ctx context.Context, tuner gopi.DVBTuner, params gopi.DVBTunerParams) error {
	var fd uintptr
	var err error

	if tuner_, ok := tuner.(*Tuner); ok == false || tuner_ == nil {
		return gopi.ErrBadParameter
	} else if params_, ok := params.(*Params); ok == false || params_ == nil {
		return gopi.ErrBadParameter
	} else if fd, err = this.getFrontend(tuner_); err != nil {
		return err
	} else if err := tuner_.Validate(params_); err != nil {
		return err
	} else if err := dvb.FETune(fd, params_.TuneParams); err != nil {
		return err
	}

	// Lock for reading
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Now loop until status is tuned, whilst emitting the
	// current status
	ticker := time.NewTicker(time.Millisecond * 100)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			status, err := dvb.FEReadStatus(fd)
			if err != nil {
				return err
			}
			switch {
			case status&dvb.FE_HAS_LOCK == dvb.FE_HAS_LOCK:
				return nil
			case status == dvb.FE_NONE:
				// Do nothing, no tune status
			default:
				this.Debug("  status=", status)
			}
		}
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Manager) String() string {
	str := "<dvb.manager"
	for _, tuner := range this.Tuners() {
		str += " " + fmt.Sprint(tuner)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// getFrontend returns file descriptor for frontend
func (this *Manager) getFrontend(tuner *Tuner) (uintptr, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if tuner == nil {
		return 0, gopi.ErrBadParameter.WithPrefix("GetFrontend")
	}

	key := tuner.Id()
	if fh, exists := this.frontend[key]; exists {
		return fh.Fd(), nil
	} else if fh, err := tuner.OpenFrontend(); err != nil {
		return 0, err
	} else {
		this.frontend[key] = fh
		return fh.Fd(), nil
	}
}

// getDemux returns file descriptor for demux
func (this *Manager) getDemux(tuner *Tuner) (uintptr, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if tuner == nil {
		return 0, gopi.ErrBadParameter.WithPrefix("GetDemux")
	}

	key := tuner.Id()
	if fh, exists := this.demux[key]; exists {
		return fh.Fd(), nil
	} else if fh, err := tuner.OpenDemux(); err != nil {
		return 0, err
	} else {
		this.demux[key] = fh
		return fh.Fd(), nil
	}
}

// disposeFrontend closes file descriptor for frontend
func (this *Manager) disposeFrontend(tuner *Tuner) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if tuner == nil {
		return gopi.ErrBadParameter.WithPrefix("DisposeFrontend")
	}

	key := tuner.Id()
	fh, exists := this.frontend[key]
	if exists == false {
		return gopi.ErrNotFound
	}
	defer delete(this.frontend, key)
	if err := fh.Close(); err != nil {
		return err
	}

	// Return success
	return nil
}

// disposeDemux closes file descriptor for demux
func (this *Manager) disposeDemux(tuner *Tuner) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if tuner == nil {
		return gopi.ErrBadParameter.WithPrefix("DisposeDemux")
	}

	key := tuner.Id()
	fh, exists := this.demux[key]
	if exists == false {
		return gopi.ErrNotFound
	}
	defer delete(this.demux, key)
	if err := fh.Close(); err != nil {
		return err
	}

	// Return success
	return nil
}
