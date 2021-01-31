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
		// Create context
		if ctx, err := this.CreateContext(tuner_, section); err != nil {
			this.Debug("Tune: ", err)
		} else {
			cb(ctx)
		}
		// Remove section filter sometime in the future
		go func() {
			if err := this.RemoveSectionFilter(tuner_, filter); err != nil {
				this.Debug("Tune: ", err)
			}
		}()
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
		if section, err := filter.Read(); err != nil {
			this.Debug("SectionFilter: ", err)
			cb(nil)
		} else {
			cb(section)
		}
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

func (this *Manager) CreateContext(tuner *Tuner, pat *ts.Section) (*Context, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if tuner == nil || pat == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("CreateContext")
	}

	// Create context
	ctx := NewContext(pat)
	if ctx == nil {
		return nil, gopi.ErrInternalAppError.WithPrefix("CreateContext")
	} else {
		this.context[tuner] = ctx
	}

	// Return context
	return ctx, nil
}

func (this *Manager) ScanPMT(tuner *Tuner, pid uint16) error {
	if filter, err := tuner.NewSectionFilter(pid, ts.PMT, dvb.DMX_ONESHOT); err != nil {
		return err
	} else if err := this.StartSectionFilter(filter, func(section *ts.Section) {
		// Oneshot filter
		this.Debug("PMT: ", section)
		// Remove section filter sometime in the future
		go func() {
			if err := this.RemoveSectionFilter(tuner, filter); err != nil {
				this.Debug("Tune: ", err)
			}
		}()
	}); err != nil {
		return err
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
		if service := ctx.NextServiceScan(); service != nil {
			this.Debug("ScanPMT for tuner:", tuner.Id(), " service: ", service)
			if err := this.ScanPMT(tuner, service.Pid()); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Return any errors
	return result
}

/*

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

*/
