package chromecast

import "github.com/djthorpe/gopi/v3"

/////////////////////////////////////////////////////////////////////
// CAST EVENT

type event struct {
	pb *CastEvent
}

func toProtoEvent(evt gopi.CastEvent) *CastEvent {
	return &CastEvent{
		Cast:  toProtoCast(evt.Cast()),
		Flags: toProtoCastFlag(evt.Flags()),
	}
}

func toProtoNull() *CastEvent {
	return &CastEvent{}
}

func fromProtoEvent(pb *CastEvent) gopi.CastEvent {
	if pb == nil {
		return nil
	} else {
		return &event{pb}
	}
}

func (this *event) Name() string {
	return this.pb.Cast.Name
}

func (this *event) Flags() gopi.CastFlag {
	return fromProtoCastFlag(this.pb.Flags)
}

func (this *event) Cast() gopi.Cast {
	return fromProtoCast(this.pb.Cast)
}

func (this *event) String() string {
	str := "<cast.event"
	str += " " + this.pb.String()
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// CAST FLAG

func toProtoCastFlag(flag gopi.CastFlag) CastEvent_Flag {
	return CastEvent_Flag(flag)
}

func fromProtoCastFlag(flag CastEvent_Flag) gopi.CastFlag {
	return gopi.CastFlag(flag)
}

/////////////////////////////////////////////////////////////////////
// CAST STATE

func toProtoState(cast gopi.Cast) Cast_CastState {
	if cast.State() == 0 {
		return Cast_IDLE
	} else {
		return Cast_ACTIVE
	}
}

/////////////////////////////////////////////////////////////////////
// CAST

type cast struct {
	pb *Cast
}

func toProtoCast(cast gopi.Cast) *Cast {
	return &Cast{
		Id:      cast.Id(),
		Name:    cast.Name(),
		Model:   cast.Model(),
		Service: cast.Service(),
		State:   toProtoState(cast),
	}
}

func fromProtoCast(pb *Cast) gopi.Cast {
	if pb == nil {
		return nil
	} else {
		return &cast{pb}
	}
}

func (this *cast) Id() string {
	return this.pb.Id
}

func (this *cast) Name() string {
	return this.pb.Name
}

func (this *cast) Model() string {
	return this.pb.Model
}

func (this *cast) Service() string {
	return this.pb.Service
}

func (this *cast) State() uint {
	if this.pb.State&Cast_ACTIVE != 0 {
		return 1
	} else {
		return 0
	}
}

func (this *cast) String() string {
	str := "<cast.device"
	str += " " + this.pb.String()
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// CAST LIST

func toProtoCastList(casts []gopi.Cast) []*Cast {
	result := make([]*Cast, 0, len(casts))
	for _, cast := range casts {
		result = append(result, toProtoCast(cast))
	}
	return result
}

func fromProtoCastList(casts []*Cast) []gopi.Cast {
	result := make([]gopi.Cast, 0, len(casts))
	for _, cast := range casts {
		result = append(result, fromProtoCast(cast))
	}
	return result
}
