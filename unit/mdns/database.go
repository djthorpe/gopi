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

func (this *Database) RegisterRecord(r gopi.RPCEvent) {
	srv := r.Service()
	if key := keyForService(srv); key == "" {
		// Ignore if invalid key
		return
	} else if r.TTL() == 0 {
		// TODO also delete name if there are no records with the name
		if exists := this.DeleteRecord(key); exists {
			fmt.Println("del r=", key)
		}
	} else if exists := this.SetRecord(key, srv, time.Now().Add(r.TTL())); exists == false {
		// TODO: add name as well
		fmt.Println("add r=", key)
		// TODO: emit record
	} else {
		// TODO: compare old record to new record and emit if any part has changed
		fmt.Println("update r=", key)
	}
}

func (this *Database) RegisterName(r gopi.RPCEvent) {
	srv := r.Service()
	if r.TTL() == 0 {
		this.DeleteName(srv.Name)
	} else if exists := this.SetName(srv.Name, time.Now().Add(r.TTL())); exists == false {
		fmt.Println("add name=", srv.Name)
	}
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
func (this *Database) SetRecord(key string, srv gopi.RPCServiceRecord, expires time.Time) bool {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Get existing record
	record, exists := this.records[key]

	// Record has expired
	if exists && time.Now().After(record.expiry) {
		exists = false
	}

	// Set new expiry
	this.records[key] = Record{srv, expires}

	// Return exists
	return exists
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

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func keyForService(record gopi.RPCServiceRecord) string {
	if record.Name == "" {
		return ""
	} else {
		return Quote(record.Name) + "." + record.Service
	}
}
