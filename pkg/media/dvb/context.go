// +build dvb

package dvb

import (
	"fmt"

	ts "github.com/djthorpe/gopi/v3/pkg/media/internal/ts"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Context struct {
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
		key := program.Program
		if key == 0 {
			if program.Pid != 0 {
				this.nit = program.Pid
			}
		} else {
			this.service[key] = NewService(key, program.Pid)
		}
	}

	// Return success
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Context) NextServiceScan() *Service {
	// Returns a service which hasn't been populated with PMT information (streams)
	// yet or nil if all services have been scanned
	for _, service := range this.service {
		if service.pmt == false {
			service.pmt = true
			return service
		}
	}
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
