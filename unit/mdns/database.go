/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package mdns

import (
	"fmt"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Database represents the current state of mDNS, service names
// and records
type Database struct {
	// Map of names with expiry time
	names map[string]time.Time

	// Map of records
	records map[string]Record

	sync.RWMutex
}

type Record struct {
	srv    gopi.RPCServiceRecord
	expiry time.Time
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *Database) Init() {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	this.names = make(map[string]time.Time)
	this.records = make(map[string]Record)
}

func (this *Database) Close() {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	this.names = nil
	this.records = nil
}

////////////////////////////////////////////////////////////////////////////////
// REGISTER METHODS

func (this *Database) RegisterRecord(r *event) gopi.RPCEvent {
	srv := r.Service()
	if key := keyForService(srv); key == "" {
		// Ignore if invalid key
	} else if r.TTL() == 0 {
		if exists := this.DeleteRecord(key); exists {
			// TODO: if there are no other records with this name, then delete
			// the name
			r.type_ = gopi.RPC_EVENT_SERVICE_REMOVED
			return r
		}
	} else if exists, modified := this.SetRecord(key, srv, time.Now().Add(r.TTL())); exists == false {
		// Add the name
		if exists := this.SetName(srv.Service, time.Now().Add(r.TTL())); exists == false {
			fmt.Println("add name=", srv.Service)
		}
		r.type_ = gopi.RPC_EVENT_SERVICE_ADDED
		return r
	} else if modified {
		r.type_ = gopi.RPC_EVENT_SERVICE_UPDATED
		return r
	}

	// Return nil - no event emitted
	return nil
}

func (this *Database) RegisterName(r *event) gopi.RPCEvent {
	srv := r.Service()
	if r.TTL() == 0 {
		if exists := this.DeleteName(srv.Name); exists {
			fmt.Println("del name=", srv.Name)
		}
	} else if exists := this.SetName(srv.Name, time.Now().Add(r.TTL())); exists == false {
		return r
	}

	// Return nil - no event emitted
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// ADD AND REMOVE NAMES

// SetName and return true if the name previously existed
func (this *Database) SetName(name string, expires time.Time) bool {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Get existing record
	expiry, exists := this.names[name]

	// Record has expired
	if exists && time.Now().After(expiry) {
		exists = false
	}

	// Set new expiry
	this.names[name] = expires

	// Return exists
	return exists
}

// DeleteName and return true if the name previously existed
func (this *Database) DeleteName(name string) bool {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if _, exists := this.names[name]; exists {
		delete(this.names, name)
		return true
	} else {
		return false
	}
}

// ExistsName returns true if name exists and isn't expired
func (this *Database) ExistsName(name string) bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	if expiry, exists := this.names[name]; exists == false {
		return false
	} else if time.Now().After(expiry) {
		return false
	} else {
		return true
	}
}

// Return all unexpired names
func (this *Database) Names() []string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	names := make([]string, 0, len(this.names))
	for name, expiry := range this.names {
		if time.Now().After(expiry) == false {
			names = append(names, name)
		}
	}
	return names
}

////////////////////////////////////////////////////////////////////////////////
// ADD AND REMOVE SERVICE RECORDS

// SetRecord and return true if the record previously existed
func (this *Database) SetRecord(key string, srv gopi.RPCServiceRecord, expires time.Time) (bool, bool) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Get existing record
	record, exists := this.records[key]

	// Record has expired
	if exists && time.Now().After(record.expiry) {
		exists = false
	}

	modified := false
	if exists && serviceEquals(record.srv, srv) == false {
		modified = true
	}

	// Set new record and expiry
	this.records[key] = Record{srv, expires}

	// Return exists
	return exists, modified
}

// DeleteRecord and return true if the name previously existed
func (this *Database) DeleteRecord(key string) bool {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if _, exists := this.records[key]; exists {
		delete(this.records, key)
		return true
	} else {
		return false
	}
}

// Return all unexpired records that match a service
func (this *Database) Records(srv string) []gopi.RPCServiceRecord {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	records := make([]gopi.RPCServiceRecord, 0, len(this.records))
	for _, r := range this.records {
		if srv == "" || srv == r.srv.Service {
			if time.Now().After(r.expiry) == false {
				records = append(records, r.srv)
			}
		}
	}
	return records
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func keyForService(record gopi.RPCServiceRecord) string {
	if record.Name == "" {
		return ""
	} else {
		return Quote(record.Name) + "." + record.Service
	}
}

func serviceEquals(a, b gopi.RPCServiceRecord) bool {
	if a.Name != b.Name {
		return false
	}
	if a.Service != b.Service {
		return false
	}
	if a.Host != b.Host {
		return false
	}
	if a.Port != b.Port {
		return false
	}
	if len(a.Txt) != len(b.Txt) {
		return false
	}
	for i, txt := range a.Txt {
		if txt != b.Txt[i] {
			return false
		}
	}
	if len(a.Addrs) != len(b.Addrs) {
		return false
	}
	addrs := make(map[string]bool, len(a.Addrs))
	for _, addr := range a.Addrs {
		addrs[addr.String()] = true
	}
	for _, addr := range b.Addrs {
		if _, exists := addrs[addr.String()]; exists == false {
			return false
		}
	}
	// Everything matches
	return true
}
