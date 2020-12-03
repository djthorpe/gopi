package gpio

import (
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

type event struct {
	name string
	gopi.GPIOPin
	gopi.GPIOEdge
}

////////////////////////////////////////////////////////////////////////////////
// NEW

func NewEvent(name string, pin gopi.GPIOPin, edge gopi.GPIOEdge) gopi.GPIOEvent {
	return &event{name, pin, edge}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *event) Name() string {
	return this.name
}

func (this *event) Pin() gopi.GPIOPin {
	return this.GPIOPin
}

func (this *event) Edge() gopi.GPIOEdge {
	return this.GPIOEdge
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<event.gpio"
	if this.GPIOPin != gopi.GPIO_PIN_NONE {
		str += " pin=" + fmt.Sprint(this.GPIOPin)
	}
	if this.GPIOEdge != gopi.GPIO_EDGE_NONE {
		str += " edge=" + fmt.Sprint(this.GPIOEdge)
	}
	return str + ">"
}
