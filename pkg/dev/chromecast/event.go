package chromecast

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	cast  *Cast
	flags gopi.CastFlag
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCastEvent(cast *Cast, flags gopi.CastFlag) gopi.CastEvent {
	return &event{cast, flags}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *event) Name() string {
	if this.cast != nil {
		return this.cast.Name()
	} else {
		return ""
	}
}

func (this *event) Cast() gopi.Cast {
	return this.cast
}

func (this *event) Flags() gopi.CastFlag {
	return this.flags
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *event) String() string {
	str := "<cast.event"
	if this.flags != gopi.CAST_FLAG_NONE {
		str += " flags=" + fmt.Sprint(this.flags)
	}
	if this.cast != nil {
		str += " cast=" + fmt.Sprint(this.cast)
	}
	return str + ">"
}
