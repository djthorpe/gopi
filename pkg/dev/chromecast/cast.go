package chromecast

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	// Modules
	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Cast struct {
	sync.RWMutex

	id, fn string
	md, rs string
	st     uint
	host   string
	ips    []net.IP
	port   uint16

	vol *Volume
	app *App
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCastFromRecord(r gopi.ServiceRecord) *Cast {
	this := new(Cast)

	// Set addr and port
	this.host = r.Host()
	this.port = r.Port()
	this.ips = r.Addrs()

	// Set properties
	tuples := txtToMap(r.Txt())
	if id, exists := tuples["id"]; exists && id != "" {
		this.id = id
	} else {
		return nil
	}
	if fn, exists := tuples["fn"]; exists && fn != "" {
		this.fn = fn
	} else {
		this.fn = this.id
	}
	if md, exists := tuples["md"]; exists {
		this.md = md
	}
	if rs, exists := tuples["rs"]; exists {
		this.rs = rs
	}
	if st, exists := tuples["st"]; exists {
		if st, err := strconv.ParseUint(st, 0, 64); err == nil {
			this.st = uint(st)
		}
	}

	// Return success
	return this
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Cast) String() string {
	str := "<cast.device"

	str += fmt.Sprintf(" id=%q", this.Id())
	if name := this.Name(); name != "" {
		str += " name=" + strconv.Quote(name)
	}
	if model := this.Model(); model != "" {
		str += " model=" + strconv.Quote(model)
	}
	if service := this.Service(); service != "" {
		str += " service=" + strconv.Quote(service)
	}
	str += " state=" + fmt.Sprint(this.State())

	if vol := this.volume(); vol != nil {
		str += fmt.Sprint(" vol=", vol)
	}
	if app := this.App(); app != nil {
		str += fmt.Sprint(" app=", app)
	}

	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Cast) Id() string {
	return this.id
}

// Name returns the readable name for a chromecast
func (this *Cast) Name() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.fn
}

// Model returns the reported model information
func (this *Cast) Model() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.md
}

// Service returns the currently running service
func (this *Cast) Service() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if this.app != nil && this.app.DisplayName != "" {
		return this.app.DisplayName
	} else {
		return this.rs
	}
}

// State returns 0 if backdrop (no app running), else returns 1
func (this *Cast) State() uint {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.st
}

// Return volume or nil if volume is not known
func (this *Cast) volume() *Volume {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if this.vol == nil {
		return nil
	} else {
		// Make a copy of volume
		vol := *this.vol
		return &vol
	}
}

// Return volume or (0,false) if volume is not known
func (this *Cast) Volume() (float32, bool) {
	if vol := this.volume(); vol == nil {
		return 0, false
	} else if vol.Level == 0.0 {
		return 0, true
	} else {
		return vol.Level, vol.Muted
	}
}

func (this *Cast) App() *App {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if this.app == nil {
		return nil
	} else {
		// Make a copy of app
		app := *this.app
		return &app
	}
}

////////////////////////////////////////////////////////////////////////////////
// METHODS

func (this *Cast) Equals(other *Cast) gopi.CastFlag {
	flags := gopi.CAST_FLAG_NONE
	if other == nil {
		return flags
	}
	// Any change to Name, Model or Id
	if this.Id() != other.Id() || this.Name() != other.Name() || this.Model() != other.Model() {
		flags |= gopi.CAST_FLAG_NAME
	}
	// Any change to service or state
	if this.Service() != other.Service() || this.State() != other.State() {
		flags |= gopi.CAST_FLAG_APP
	}
	// Return changed flags
	return flags
}

func (this *Cast) ConnectWithTimeout(ch gopi.Publisher, timeout time.Duration) (*Conn, error) {
	// Use hostname to connect
	addr := fmt.Sprintf("%v:%v", this.host, this.port)

	// Update state
	this.vol = nil
	this.app = nil

	// Perform the connection
	return NewConnWithTimeout(this.id, addr, timeout, ch)
}

func (this *Cast) Disconnect(conn *Conn) error {
	// Update state
	this.vol = nil
	this.app = nil

	// Close connection
	return conn.Close()
}

func (this *Cast) UpdateState(state *State) gopi.CastFlag {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Changes in volume
	flags := gopi.CAST_FLAG_NONE
	if this.vol == nil || state.volume.Equals(*this.vol) == false {
		this.vol = &state.volume
		flags |= gopi.CAST_FLAG_VOLUME
	}

	// Changes in app
	if this.app == nil && len(state.apps) > 0 {
		this.app = &state.apps[0]
		flags |= gopi.CAST_FLAG_APP
	} else if len(state.apps) == 0 && this.app != nil {
		this.app = nil
		flags |= gopi.CAST_FLAG_APP
	} else if this.app != nil && len(state.apps) > 0 && this.app.Equals(state.apps[0]) == false {
		this.app = &state.apps[0]
		flags |= gopi.CAST_FLAG_APP
	}

	// Return any changed state
	return flags
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Cast) updateFrom(other *Cast) {
	this.RWMutex.Lock()
	other.RWMutex.RLock()
	defer this.RWMutex.Unlock()
	defer other.RWMutex.RUnlock()

	// This seems dumb
	this.id = other.id
	this.fn = other.fn
	this.md = other.md
	this.rs = other.rs
	this.st = other.st
	this.host = other.host
	this.ips = other.ips
	this.port = other.port
}

func txtToMap(txt []string) map[string]string {
	result := make(map[string]string, len(txt))
	for _, r := range txt {
		if kv := strings.SplitN(r, "=", 2); len(kv) == 2 {
			result[kv[0]] = kv[1]
		} else if len(kv) == 1 {
			result[kv[0]] = ""
		}
	}
	return result
}
