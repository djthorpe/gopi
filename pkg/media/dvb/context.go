// +build dvb

package dvb

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
	ts "github.com/djthorpe/gopi/v3/pkg/media/internal/ts"
)

type Context struct {
	err     error
	nit     uint16
	service map[uint16]*Service
}

type Service struct {
	pid uint16
}

func NewContext(pat *ts.Section, err error) gopi.DVBContext {
	if pat.TableId != ts.PAT {
		return nil
	}

	this := new(Context)
	this.err = err
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
			this.service[key] = NewService(program.Pid)
		}
	}

	// Return success
	return this
}

func NewService(pid uint16) *Service {
	return &Service{pid}
}

func (this *Context) Err() error {
	return this.err
}

func (this *Context) String() string {
	str := "<dvb.context"
	if err := this.Err(); err != nil {
		str += fmt.Sprint(" err=", err)
	}
	if nit := this.nit; nit != 0 {
		str += fmt.Sprintf(" nit_pid=0x%04X", nit)
	}
	for id, service := range this.service {
		str += fmt.Sprintf(" %04X=%v", id, service)
	}
	return str + ">"
}

func (this *Service) String() string {
	str := "<dvb.service"
	str += fmt.Sprintf(" pid=0x%04X", this.pid)
	return str + ">"
}
