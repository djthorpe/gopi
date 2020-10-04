package event

import (
	"fmt"

	"github.com/djthorpe/gopi/v3"
	"github.com/djthorpe/gopi/v3/pkg/log"
)

type Publisher struct {
	gopi.Unit
	*log.Log
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Publisher) String() string {
	str := "<publisher"
	if this == nil {
		str += " nil"
	} else {
		str += " log=" + fmt.Sprint(this.Log)
	}
	return str + ">"
}
