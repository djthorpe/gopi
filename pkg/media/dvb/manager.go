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

	frontend map[uint]*os.File
	demux    map[uint]*os.File
	section  map[uint][]*SectionFilter
	stream   map[uint][]*StreamFilter
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Manager) New(gopi.Config) error {
	this.Require(this.Logger, this.FilePoll)

	// Set up file descriptor maps
	this.frontend = make(map[uint]*os.File)
	this.demux = make(map[uint]*os.File)

	// Set up filter maps
	this.section = make(map[uint][]*SectionFilter)
	this.stream = make(map[uint][]*StreamFilter)

	// Return success
	return nil
}

func (this *Manager) Dispose() error {
	var result error

	// Close filters
	for key, filters := range this.section {
		for _, filter := range filters {
			if err := this.StopSectionFilter(key, filter); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}
	for key, filters := range this.stream {
		for _, filter := range filters {
			if err := this.StopStreamFilter(key, filter); err != nil {
				result = multierror.Append(result, err)
			}
		}
	}

	// Lock for exclusive access
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

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
	this.section = nil
	this.stream = nil

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

func (this *Manager) Tune(ctx context.Context, tuner gopi.DVBTuner, params gopi.DVBTunerParams, cb gopi.DVBTuneCallack) error {
	var fd uintptr
	var err error

	if tuner_, ok := tuner.(*Tuner); ok == false || tuner_ == nil {
		return gopi.ErrBadParameter.WithPrefix("Tune")
	} else if params_, ok := params.(*Params); ok == false || params_ == nil {
		return gopi.ErrBadParameter.WithPrefix("Tune")
	} else if fd, err = this.getFrontend(tuner_); err != nil {
		return err
	} else if err := tuner_.Validate(params_); err != nil {
		return err
	} else if err := this.StopFilters(tuner_); err != nil {
		return err
	} else if err := dvb.FETune(fd, params_.TuneParams); err != nil {
		return err
	}

	// TODO: Create a context which will contain state information

	// Now loop until status is tuned, whilst emitting the current status
	ticker := time.NewTicker(time.Millisecond * 100)
	this.Debug("Tune ", tuner.Name(), " adapter=", tuner.Id())
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
				this.Debug("  ->", status)
				if err := this.ScanPAT(tuner, cb); err != nil {
					this.Close(tuner)
					return err
				}
				return nil
			case status == dvb.FE_NONE:
				// Do nothing, no tune status
			default:
				this.Debug("  ->", status)
			}
		}
	}

	// Return success
	return nil
}

// Close will remove any filters close device
func (this *Manager) Close(tuner gopi.DVBTuner) error {
	// Stop filters
	if tuner_, ok := tuner.(*Tuner); ok == false || tuner_ == nil {
		return gopi.ErrBadParameter.WithPrefix("Close")
	} else if err := this.StopFilters(tuner_); err != nil {
		return err
	}

	// Dispose
	var result error
	if err := this.disposeDemux(tuner.(*Tuner)); err != nil {
		result = multierror.Append(result, err)
	}
	if err := this.disposeFrontend(tuner.(*Tuner)); err != nil {
		result = multierror.Append(result, err)
	}

	// Return any errors
	return result
}

// StartSectionFilter starts scanning for specific pid & tid
func (this *Manager) StartSectionFilter(tuner *Tuner, pid uint16, tid ts.TableType, flags dvb.DMXFlag, cb func(*ts.Section, error)) (*SectionFilter, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check parameters
	if tuner == nil || cb == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("StartSectionFilter")
	}

	// Create filter
	filter, err := NewSectionFilter(tuner, pid, tid, flags)
	if err != nil {
		return nil, err
	}

	// Watch for sections, and read them as they appear
	if err := this.FilePoll.Watch(filter.Fd(), gopi.FILEPOLL_FLAG_READ, func(uintptr, gopi.FilePollFlags) {
		cb(filter.Read())
	}); err != nil {
		filter.Dispose()
		return nil, err
	}

	// Start filtering
	if err := filter.Start(); err != nil {
		this.FilePoll.Unwatch(filter.Fd())
		filter.Dispose()
		return nil, err
	} else {
		key := tuner.Id()
		this.section[key] = append(this.section[key], filter)
	}

	// Return success
	return filter, nil
}

// StopFilters stops all filters for a tuner
func (this *Manager) StopFilters(tuner *Tuner) error {
	// Check parameters
	if tuner == nil {
		return gopi.ErrBadParameter.WithPrefix("StopFilters")
	}

	// Tuner key
	key := tuner.Id()

	// Sensitive section - get filters
	this.RWMutex.RLock()
	sections, _ := this.section[key]
	streams, _ := this.stream[key]
	this.RWMutex.RUnlock()

	// Stop filters
	var result error
	for _, filter := range sections {
		if err := this.StopSectionFilter(key, filter); err != nil {
			result = multierror.Append(result, err)
		}
	}
	for _, filter := range streams {
		if err := this.StopStreamFilter(key, filter); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return any errors
	return result
}

// StopSectionFilter stops section filter
func (this *Manager) StopSectionFilter(key uint, filter *SectionFilter) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check parameters
	if filter == nil {
		return gopi.ErrBadParameter.WithPrefix("StopSectionFilter")
	}
	if _, exists := this.section[key]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("StopSectionFilter")
	}

	// Remove from list of filters
	this.section[key] = removeSectionFilter(this.section[key], filter)

	// Unwatch
	var result error
	if err := this.FilePoll.Unwatch(filter.Fd()); err != nil {
		result = multierror.Append(result, err)
	}
	// Stop filtering
	if err := filter.Stop(); err != nil {
		result = multierror.Append(result, err)
	}
	// Dispose filter
	if err := filter.Dispose(); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

// StopStreamFilter stops stream filter
func (this *Manager) StopStreamFilter(key uint, filter *StreamFilter) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	return gopi.ErrNotImplemented
}

// ScanPAT starts filtering for PAT table section, and then callback when returned
// or an error on timeout
func (this *Manager) ScanPAT(tuner gopi.DVBTuner, cb gopi.DVBTuneCallack) error {
	_, err := this.StartSectionFilter(tuner.(*Tuner), uint16(0x0000), ts.PAT, dvb.DMX_ONESHOT, func(pat *ts.Section, err error) {
		if ctx := NewContext(pat, err); ctx == nil {
			this.Debug("  Context is nil")
			this.Close(tuner)
		} else if err := cb(ctx); err != nil {
			this.Debug("  Error from callback", err)
			this.Close(tuner)
		}
	})
	return err
}

/*
// ScanPMT starts filtering for PMT table section
func (this *Manager) ScanPMT(tuner gopi.DVBTuner) error {
	if _, err := this.StartSectionFilter(tuner.(*Tuner), uint16(0x0000), ts.PAT, 0); err != nil {
		return err
	} else {
		return nil
	}
}

// ScanNIT starts filtering for NIT table section
func (this *Manager) ScanNIT(tuner gopi.DVBTuner) error {
	if _, err := this.StartSectionFilter(tuner.(*Tuner), uint16(0x0010), ts.NIT, 0); err != nil {
		return err
	} else {
		return nil
	}
}
*/

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
