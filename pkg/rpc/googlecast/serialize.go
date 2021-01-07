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
	state := this.GetState()
	if state == nil {
		return nil
	} else {
		return fromProtoCast(state.GetCast())
	}
}

func (this *event) Volume() (float32, bool) {
	state := this.GetState()
	if state == nil {
		return 0, false
	} else if volume := state.GetVolume(); volume == nil {
		return 0, false
	} else {
		return volume.GetLevel(), volume.GetMuted()
	}
}

func (this *event) App() gopi.CastApp {
	state := this.GetState()
	if state == nil {
		return nil
	} else {
		return fromProtoApp(state.GetApp())
	}
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
// APPLICATION INTERFACE

type app struct {
	*App
}

func fromProtoApp(pb *App) gopi.CastApp {
	if pb == nil {
		return nil
	} else {
		return &app{pb}
	}
}

func (this *app) Id() string {
	return this.GetId()
}

func (this *app) Name() string {
	return this.GetName()
}

func (this *app) Status() string {
	return this.GetStatus()
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

func toProtoVolume(level float32, muted bool) *Volume {
	return &Volume{
		Level: level,
		Muted: muted,
	}
}

func toProtoApp(app gopi.CastApp) *App {
	if app == nil || app.(gopi.CastApp) == nil {
		return nil
	}
	return &App{
		Id:     app.Id(),
		Name:   app.Name(),
		Status: app.Status(),
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
		State: &CastState{
			Cast:   toProtoCast(evt.Cast()),
			Volume: toProtoVolume(evt.Volume()),
			App:    toProtoApp(evt.App()),
		},
		Changed: CastEvent_Flag(evt.Flags()),
	}
}

func fromProtoEvent(evt *CastEvent) gopi.CastEvent {
	if evt == nil || evt.State == nil || evt.State.Cast == nil {
		return nil
	} else {
		return &event{evt}
	}
}

func toProtoNull() *CastEvent {
	return &CastEvent{}
}
