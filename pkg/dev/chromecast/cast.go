package chromecast

import (
	"context"
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
	if port := r.Port(); port == 0 {
		return nil
	} else {
		this.port = port
	}
	if ips := r.Addrs(); len(ips) == 0 {
		return nil
	} else {
		this.ips = ips
	}

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

func (this *Manager) Run(ctx context.Context) error {
	// Subscribe to events
	ch := this.Publisher.Subscribe()
	defer this.Publisher.Unsubscribe(ch)

	// Loop handling messages until done
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case evt := <-ch:
			fmt.Println(evt)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Cast) String() string {
	str := "<cast.device"
	str += fmt.Sprintf(" id=%q", this.id)
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *Cast) Id() string {
	return this.id
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
