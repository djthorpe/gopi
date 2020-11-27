package event

import gopi "github.com/djthorpe/gopi/v3"

type nullevent struct{}

func NewNullEvent() gopi.Event {
	return new(nullevent)
}

func (*nullevent) Name() string {
	return "null"
}

func (*nullevent) String() string {
	return "<nullevent>"
}
