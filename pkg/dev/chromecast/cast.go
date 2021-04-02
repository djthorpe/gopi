package chromecast

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

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
	ips    []net.IP
	port   uint16
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCastFromRecord(r gopi.ServiceRecord) *Cast {
	this := new(Cast)

	// Set addr and port
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

func (this *Cast) Disconnect() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Return success
	return nil
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
	return this.rs
}

// State returns 0 if backdrop (no app running), else returns 1
func (this *Cast) State() uint {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.st
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

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Cast) updateFrom(other *Cast) {
	this.RWMutex.Lock()
	other.RWMutex.RLock()
	defer this.RWMutex.Unlock()
	defer other.RWMutex.RUnlock()

	this.id = other.id
	this.fn = other.fn
	this.md = other.md
	this.rs = other.rs
	this.st = other.st
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
