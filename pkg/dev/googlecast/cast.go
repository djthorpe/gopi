package googlecast

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Cast struct {
	sync.RWMutex

	id, fn string
	md, rs string
	st     uint
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func NewCastFromRecord(r gopi.ServiceRecord) *Cast {
	this := new(Cast)

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

	return this
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

// Id returns the identifier for a chromecast
func (this *Cast) Id() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
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

// State returns 0 if backdrop, else returns 1
func (this *Cast) State() uint {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return this.st
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Cast) String() string {
	str := "<cast.device"
	str += " id=" + this.Id()
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
// PRIVATE METHODS

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
