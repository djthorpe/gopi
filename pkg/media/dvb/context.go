// +build dvb

package dvb

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	ts "github.com/djthorpe/gopi/v3/pkg/media/internal/ts"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Context struct {
	sync.RWMutex

	nit     uint16
	service map[uint16]*Service
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewContext(pat *ts.Section) *Context {
	if pat.TableId != ts.PAT {
		return nil
	}

	this := new(Context)
	this.nit = uint16(0x0010)
	this.service = make(map[uint16]*Service, len(pat.PATSection.Programs))

	// Iterate through programs, program 0 is the NIT PID
	for _, program := range pat.PATSection.Programs {
		key := program.Pid
		if key == 0 {
			if program.Pid != 0 {
				this.nit = program.Pid
			}
		} else {
			this.service[key] = NewService(key, program.Program)
		}
	}

	// Return success
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Context) NextServiceScan() *Service {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Returns a service which hasn't been populated with PMT information (streams)
	// yet or nil if all services have been scanned
	for _, service := range this.service {
		if service.Streams() == false {
			return service
		}
	}
	return nil
}

func (this *Context) SetPMT(pid uint16, section *ts.Section) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if service, exists := this.service[pid]; exists == false {
		return gopi.ErrNotFound.WithPrefix("SetPMT")
	} else if section.TableId != ts.PMT {
		return gopi.ErrInternalAppError.WithPrefix("SetPMT")
	} else {
		service.SetStreams(section.PMTSection.ESTable.Rows)
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Context) String() string {
	str := "<dvb.context"
	if nit := this.nit; nit != 0 {
		str += fmt.Sprintf(" nit_pid=0x%04X", nit)
	}
	for _, service := range this.service {
		str += fmt.Sprintf(" %v", service)
	}
	return str + ">"
}
