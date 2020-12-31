package googlecast

import gopi "github.com/djthorpe/gopi/v3"

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
