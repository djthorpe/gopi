package googlecast

import (
	"fmt"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	ptypes "github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/duration"
)

/////////////////////////////////////////////////////////////////////
// EVENT INTERFACE

type event struct {
	*CastEvent
}

func (this *event) Name() string {
	return "cast.event"
}

func (this *event) Cast() gopi.Cast {
	return fromProtoCast(this.GetCast())
}

func (this *event) Flags() gopi.CastFlag {
	return gopi.CastFlag(this.Changed)
}

func (this *event) String() string {
	str := "<cast.event"
	if cast := this.Cast(); cast != nil {
		str += " cast=" + fmt.Sprint(cast)
	}
	if flags := this.Flags(); flags != 0 {
		str += " changed=" + fmt.Sprint(flags)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// CAST DEVICE INTERFACE

type cast struct {
	*Cast
}

func (this *cast) Id() string {
	return this.Cast.GetId()
}

func (this *cast) Name() string {
	return this.Cast.GetName()
}

func (this *cast) Model() string {
	return this.Cast.GetModel()
}

func (this *cast) Service() string {
	return this.Cast.GetService()
}

func (this *cast) State() uint {
	st := this.Cast.GetState()
	if st == Cast_NONE || st == Cast_UNKNOWN {
		return 0
	}
	if st&Cast_IDLE != 0 || st&Cast_BACKDROP != 0 {
		return 0
	}
	return 1
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func fromProtoCast(pb *Cast) gopi.Cast {
	if pb == nil {
		return nil
	}
	return &cast{pb}
}

func toProtoCast(cast gopi.Cast) *Cast {
	if cast == nil {
		return nil
	}
	return &Cast{
		Id:      cast.Id(),
		Name:    cast.Name(),
		Model:   cast.Model(),
		Service: cast.Service(),
		State:   toProtoState(cast),
	}
}

func toProtoState(cast gopi.Cast) Cast_CastState {
	if cast.State() == 0 {
		return Cast_IDLE
	} else {
		return Cast_ACTIVE
	}
}

func toProtoDuration(value time.Duration) *duration.Duration {
	return ptypes.DurationProto(value)
}

func toProtoEvent(evt gopi.CastEvent) *CastEvent {
	if evt == nil || evt.Cast() == nil {
		return nil
	}
	return &CastEvent{
		Cast:    toProtoCast(evt.Cast()),
		Changed: CastEvent_Flag(evt.Flags()),
	}
}

func fromProtoEvent(evt *CastEvent) gopi.CastEvent {
	if evt == nil || evt.Cast == nil {
		return nil
	} else {
		return &event{evt}
	}
}

func toProtoNull() *CastEvent {
	return &CastEvent{}
}
