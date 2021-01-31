// +build dvb

package dvb

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	id  uint16
	pid uint16
	pmt bool
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewService(id, pid uint16) *Service {
	return &Service{id, pid, false}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Service) Id() uint16 {
	return this.id
}

func (this *Service) Pid() uint16 {
	return this.pid
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Service) String() string {
	str := "<dvb.service"
	str += fmt.Sprintf(" id=0x%04X", this.id)
	str += fmt.Sprintf(" pid=0x%04X", this.pid)
	return str + ">"
}
