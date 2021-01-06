package googlecast

import (
	"fmt"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type event struct {
	cast   *Cast
	app    *App
	volume *Volume
	reqId  int
	flags  gopi.CastFlag
}

type state struct {
	key    string
	req    int
	err    error
	dbg    string
	close  bool
	values []interface{}
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewEvent(device *Cast, app *App, volume *Volume, flags gopi.CastFlag, reqId int) gopi.CastEvent {
	return &event{device, app, volume, reqId, flags}
}

func NewState(key string, req int, values ...interface{}) state {
	return state{key, req, nil, "", false, values}
}

func NewError(key string, err error) state {
	return state{key, 0, err, "", false, nil}
}

func Close(key string) state {
	return state{key, 0, nil, "", true, nil}
}

func Debug(key string, format string, a ...interface{}) state {
	return state{key, 0, nil, fmt.Sprintf(format, a...), false, nil}
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

func (this *event) App() gopi.CastApp {
	return this.app
}

func (this *event) Volume() (float32, bool) {
	if this.volume != nil {
		return this.volume.Level, this.volume.Muted
	} else {
		return 0, false
	}
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
	if this.reqId != 0 {
		str += " reqId=" + fmt.Sprint(this.reqId)
	}
	if this.cast != nil {
		str += " device=" + fmt.Sprint(this.cast)
	}
	return str + ">"
}
