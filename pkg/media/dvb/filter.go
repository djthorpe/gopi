// +build dvb

package dvb

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ts "github.com/djthorpe/gopi/v3/pkg/media/internal/ts"
	dvb "github.com/djthorpe/gopi/v3/pkg/sys/dvb"
	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Filter struct {
	sync.RWMutex
	dev *os.File
}

type SectionFilter struct {
	Filter
	*dvb.DMXSectionFilter
}

type StreamFilter struct {
	Filter
	*dvb.DMXStreamFilter
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewSectionFilter(tuner *Tuner, pid uint16, tids ...ts.TableType) (*SectionFilter, error) {
	this := new(SectionFilter)

	// Check incoming parameters
	if tuner == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("NewSectionFilter")
	}
	if len(tids) == 0 {
		return nil, gopi.ErrBadParameter.WithPrefix("NewSectionFilter")
	}

	// Open device
	if dev, err := tuner.DMXOpen(); err != nil {
		return nil, err
	} else {
		this.dev = dev
	}

	// Create filter with 0ms timeout (no timeout)
	this.DMXSectionFilter = dvb.NewSectionFilter(pid, 0, dvb.DMX_NONE)
	for i, tid := range tids {
		this.DMXSectionFilter.Set(i, uint8(tid), 0xFF, 0x00)
	}

	// Set filter
	if err := dvb.DMXSetSectionFilter(this.dev.Fd(), this.DMXSectionFilter); err != nil {
		this.dev.Close()
		return nil, err
	}

	// Return success
	return this, nil
}

func NewStreamFilter(tuner *Tuner, pid uint16, in dvb.DMXInput, out dvb.DMXOutput, stream dvb.DMXStreamType) (*StreamFilter, error) {
	this := new(StreamFilter)

	// Check incoming parameters
	if tuner == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("NewStreamFilter")
	}

	// Open device
	if dev, err := tuner.DMXOpen(); err != nil {
		return nil, err
	} else {
		this.dev = dev
	}

	// Create filter
	this.DMXStreamFilter = dvb.NewStreamFilter(pid, in, out, stream, dvb.DMX_NONE)

	// Set filter
	if err := dvb.DMXSetStreamFilter(this.dev.Fd(), this.DMXStreamFilter); err != nil {
		this.dev.Close()
		return nil, err
	}

	// Return success
	return this, nil
}

func (this *Filter) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error
	if this.dev != nil {
		if err := this.dev.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Release resources
	this.dev = nil

	// Return success
	return result
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Filter) Fd() uintptr {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if this.dev == nil {
		return 0
	} else {
		return this.dev.Fd()
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Filter) Start() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.dev == nil {
		return gopi.ErrOutOfOrder.WithPrefix("Start")
	}

	return dvb.DMXStart(this.dev.Fd())
}

func (this *Filter) Stop() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.dev == nil {
		return gopi.ErrOutOfOrder.WithPrefix("Stop")
	}

	return dvb.DMXStop(this.dev.Fd())
}

func (this *Filter) AddPid(pid uint16) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.dev == nil {
		return gopi.ErrOutOfOrder.WithPrefix("AddPid")
	}

	return dvb.DMXAddPid(this.dev.Fd(), pid)
}

func (this *Filter) AddPids(pids []uint16) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.dev == nil {
		return gopi.ErrOutOfOrder.WithPrefix("AddPids")
	}

	var result error
	for _, pid := range pids {
		if err := dvb.DMXAddPid(this.dev.Fd(), pid); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Success
	return result
}

func (this *Filter) SetBufferSize(size uint32) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.dev == nil {
		return gopi.ErrOutOfOrder.WithPrefix("SetBufferSize")
	}

	return dvb.DMXSetBufferSize(this.dev.Fd(), size)
}

func (this *Filter) RemovePid(pid uint16) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if this.dev == nil {
		return gopi.ErrOutOfOrder.WithPrefix("RemovePid")
	}

	return dvb.DMXRemovePid(this.dev.Fd(), pid)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *SectionFilter) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<dvb.sectionfilter"
	if this.dev != nil {
		str += " dev=" + strconv.Quote(this.dev.Name())
		str += " filter=" + fmt.Sprint(this.DMXSectionFilter)
	}
	return str + ">"
}

func (this *StreamFilter) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<dvb.streamfilter"
	if this.dev != nil {
		str += " dev=" + strconv.Quote(this.dev.Name())
		str += " filter=" + fmt.Sprint(this.DMXStreamFilter)
	}
	return str + ">"
}
