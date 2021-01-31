// +build dvb

package dvb

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	ts "github.com/djthorpe/gopi/v3/pkg/media/internal/ts"
	dvb "github.com/djthorpe/gopi/v3/pkg/sys/dvb"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Manager struct {
	gopi.Unit
	gopi.Logger
	gopi.FilePoll
	sync.RWMutex

	context map[*Tuner]*Context
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	// deltaStats defines interval that frontend measurements are taken
	deltaStats = 4 * time.Second
	// deltaState defines interval that internal state is updated from tuners
	deltaState = time.Millisecond * 100
)

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger, this.FilePoll)

	// Contexts store current state against open tuners
	this.context = make(map[*Tuner]*Context)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Dispose of any open tuners
	var result error
	for tuner := range this.context {
		if err := tuner.Dispose(this.FilePoll); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.context = nil

	// Return any errors
	return result
}

func (this *Manager) Run(ctx context.Context) error {
	state := time.NewTicker(deltaState)
	defer state.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-state.C:
			if err := this.updateState(); err != nil {
				this.Print("UpdateState:", err)
			}
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Manager) Tuners() []gopi.DVBTuner {
	// Append any tuners which aren't tuned
	devices := dvb.Devices()
	tuners := make([]gopi.DVBTuner, 0, len(devices))
	for _, device := range devices {
		if tuner := this.getTunerForId(device.Adapter); tuner != nil {
			tuners = append(tuners, tuner)
		} else if tuner, err := NewTuner(device); err != nil {
			this.Debug("Tuners:", err)
		} else {
			tuners = append(tuners, tuner)
		}
	}

	// Return tuners
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

func (this *Manager) Tune(ctx context.Context, tuner gopi.DVBTuner, params gopi.DVBTunerParams, cb gopi.DVBTuneCallack) error {

	// Check parameters
	params_, ok := params.(*Params)
	if ok == false || params_ == nil || cb == nil {
		return gopi.ErrBadParameter.WithPrefix("Tune")
	}

	// Dispose of context and tuner if already open
	var result error
	if tuner_ := this.getTunerForId(tuner.Id()); tuner_ != nil {
		if err := tuner_.Dispose(this.FilePoll); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Validate parameters, open frontend and tune, then ask for PAT
	if tuner_, ok := tuner.(*Tuner); ok == false || tuner_ == nil {
		return gopi.ErrBadParameter.WithPrefix("Tune")
	} else if err := tuner_.Validate(params_); err != nil {
		return err
	} else if err := tuner_.OpenFrontend(); err != nil {
		return err
	} else if err := tuner_.Tune(ctx, params_); err != nil {
		return err
	} else if filter, err := tuner_.NewSectionFilter(0, ts.PAT, dvb.DMX_ONESHOT); err != nil {
		return err
	} else if err := this.StartSectionFilter(filter, func(section *ts.Section) {
		// Oneshot filter
		this.Debug("PAT: ", section)
		if err := this.RemoveSectionFilter(tuner_, filter); err != nil {
			this.Debug("Tune: ", err)
		}

		/*
			this.RWMutex.Lock()
			defer this.RWMutex.Unlock()
			// Create context, callback
			this.context[tuner_] = NewContext(section)
			cb(this.context[tuner_])*/
	}); err != nil {
		return err
	}

	// Return any errors
	return result
}

// Close will remove any filters close device
func (this *Manager) Close(tuner gopi.DVBTuner) error {
	tuner_ := this.getTunerForId(tuner.Id())
	if tuner_ == nil {
		return gopi.ErrNotFound.WithPrefix("Close")
	}

	// Dispose of tuner (stop watching, closing filters, etc)
	if err := tuner_.Dispose(this.FilePoll); err != nil {
		return err
	}

	// Delete context
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	delete(this.context, tuner_)

	// Return success
	return nil
}

func (this *Manager) StartSectionFilter(filter *SectionFilter, cb func(*ts.Section)) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Filewatch for sections, and read them as they appear
	if err := this.FilePoll.Watch(filter.Fd(), gopi.FILEPOLL_FLAG_READ, func(uintptr, gopi.FilePollFlags) {
		go func() {
			if section, err := filter.Read(); err != nil {
				this.Debug("SectionFilter: ", err)
				cb(nil)
			} else {
				cb(section)
			}
		}()
	}); err != nil {
		filter.Dispose()
		return err
	}

	// Start filtering
	if err := filter.Start(); err != nil {
		this.FilePoll.Unwatch(filter.Fd())
		filter.Dispose()
		return err
	}

	// Return success
	return nil
}

func (this *Manager) RemoveSectionFilter(tuner *Tuner, filter *SectionFilter) error {
	var result error

	// Stop filtering and remove from tuner
	if err := filter.Stop(); err != nil {
		result = multierror.Append(result, err)
	}
	if err := tuner.RemoveSectionFilter(this.FilePoll, filter); err != nil {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
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

func (this *Manager) getTunerForId(key uint) *Tuner {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	for tuner := range this.context {
		if tuner.Id() == key {
			return tuner
		}
	}
	return nil
}

// updateState
func (this *Manager) updateState() error {
	var result error
	for tuner, ctx := range this.context {
		fmt.Println("TODO:", tuner, ctx)
	}
	/*
			if service := ctx.NextServiceScan(); service != nil {
				this.Debug("ScanPMT for tuner:", tuner.Id(), " service: ", service)
				if err := this.ScanPMT(tuner, service.Pid()); err != nil {
					result = multierror.Append(result, err)
				}
			}
		}
	*/
	// Return any errors
	return result
}

/*
// removeSectionFilter
func removeSectionFilter(arr []*SectionFilter, filter *SectionFilter) []*SectionFilter {
	for i, elem := range arr {
		if elem == filter {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	// Filter not in array
	return arr
}


// emitStats
func (this *Manager) emitStats() error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var result error
	for id, fh := range this.frontend {
		if status, err := dvb.FEReadStatus(fh.Fd()); err != nil {
			result = multierror.Append(result, err)
		} else if status == dvb.FE_NONE {
			continue
		} else {
			this.Print(id, "->", status)
			if stats, err := dvb.FEGetPropStats(fh.Fd(),
				dvb.DTV_STAT_SIGNAL_STRENGTH,
				dvb.DTV_STAT_CNR,
			); err != nil {
				result = multierror.Append(result, err)
			} else {
				this.Print(id, "->", stats)
			}
		}
	}

	// Return any errors
	return result
}

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
		return nil
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
		return nil
	}
	defer delete(this.demux, key)
	if err := fh.Close(); err != nil {
		return err
	}

	// Return success
	return nil
}

// disposeContext removes state
func (this *Manager) disposeContext(tuner *Tuner) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if tuner == nil {
		return gopi.ErrBadParameter.WithPrefix("DisposeDemux")
	} else if _, exists := this.context[tuner]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("DisposeDemux")
	} else {
		delete(this.context, tuner)
	}

	// Return success
	return nil
}
*/
