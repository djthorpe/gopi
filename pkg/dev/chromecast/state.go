package chromecast

import (
	"fmt"
)

type State struct {
	key     string
	req     int
	volume  Volume
	apps    []App
	media   []Media
	payload []byte
	err     error
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewAppState(key string, req int, payload []byte, volume Volume, app ...App) *State {
	return &State{key, req, volume, app, nil, payload, nil}
}

func NewMediaState(key string, req int, payload []byte, media ...Media) *State {
	return &State{key, req, Volume{}, nil, media, payload, nil}
}

func NewPayloadState(key string, req int, payload []byte) *State {
	return &State{key, req, Volume{}, nil, nil, payload, nil}
}

func NewErrorState(key string, req int, payload []byte, err error) *State {
	return &State{key, req, Volume{}, nil, nil, payload, err}
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *State) Name() string {
	return "cast.state"
}

func (this *State) Payload() []byte {
	return this.payload
}

func (this *State) Err() error {
	return this.err
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *State) String() string {
	str := "<cast.state"
	str += fmt.Sprintf(" key=%q", this.key)
	str += fmt.Sprintf(" req=%v", this.req)
	if this.volume.Level > 0 || this.volume.Muted {
		str += fmt.Sprintf(" vol=%v", this.volume)
	}
	if this.apps != nil {
		str += fmt.Sprintf(" apps=%v", this.apps)
	}
	if this.media != nil {
		str += fmt.Sprintf(" media=%v", this.media)
	}
	if this.payload != nil {
		str += fmt.Sprintf(" payload=%q", string(this.payload))
	}
	return str + ">"
}
