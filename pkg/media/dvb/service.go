// +build dvb

package dvb

import (
	"fmt"
	"time"

	ts "github.com/djthorpe/gopi/v3/pkg/media/internal/ts"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	pid     uint16
	id      uint16
	ts      time.Time
	streams []ts.ESRow
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewService(pid, id uint16) *Service {
	this := new(Service)
	this.id = id
	this.pid = pid
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Service) Id() uint16 {
	return this.id
}

func (this *Service) Pid() uint16 {
	return this.pid
}

func (this *Service) Streams() bool {
	return this.streams == nil
}

func (this *Service) SetStreams(streams []ts.ESRow) {
	this.streams = streams
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Service) String() string {
	str := "<dvb.service"
	str += fmt.Sprintf(" id=0x%04X", this.id)
	str += fmt.Sprintf(" pid=0x%04X", this.pid)
	for _, stream := range this.streams {
		str += fmt.Sprint(" stream=", stream)
	}
	return str + ">"
}
