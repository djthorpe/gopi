package keycode

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type keycodedb struct {
	sync.RWMutex

	path    string
	name    string
	dirty   bool
	mapping []*mapentry
	lookup  map[uint64]*mapentry
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewDatabase(path, name string, cache *Cache) (*keycodedb, error) {
	this := new(keycodedb)
	this.name = name
	this.path = path

	if mapping, err := cache.DecodeFile(path); os.IsNotExist(err) {
		// Set dirty flag to write database
		this.dirty = true
	} else if err != nil {
		return nil, err
	} else {
		this.mapping = mapping
		this.lookup = make(map[uint64]*mapentry, len(this.mapping))
	}

	// Create lookup table
	for i, entry := range this.mapping {
		if key := keyForEntry(entry.Device, entry.Code); key != 0 {
			entry.Index = i
			this.lookup[key] = entry
		}
	}

	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *keycodedb) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<keycode.db"
	if this.path != "" {
		str += " filename=" + strconv.Quote(filepath.Base(this.path))
	}
	if this.name != "" {
		str += " name=" + strconv.Quote(this.name)
	}
	for _, v := range this.lookup {
		str += fmt.Sprintf(" { %v,%v }=>%v", v.Device, scancodeString(v.Code), v.Key)
	}
	if this.dirty {
		str += " modified=" + fmt.Sprint(this.dirty)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *keycodedb) Modified() bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.dirty
}

func (this *keycodedb) Path() string {
	return this.path
}

func (this *keycodedb) Name() string {
	return this.name
}

func (this *keycodedb) Write(cache *Cache) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Clear dirty flag regardless of error
	this.dirty = false

	// Write file
	if err := cache.EncodeFile(this.path, this.mapping); err != nil {
		return err
	}

	// Return success
	return nil
}

func (this *keycodedb) Lookup(device gopi.InputDeviceType, code uint32) gopi.KeyCode {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if hash := keyForEntry(device, code); hash == 0 {
		return gopi.KEYCODE_NONE
	} else if entry, exists := this.lookup[hash]; exists == false {
		return gopi.KEYCODE_NONE
	} else {
		return entry.Key
	}
}

func (this *keycodedb) Set(device gopi.InputDeviceType, code uint32, key gopi.KeyCode) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check parameters
	if device == gopi.INPUT_DEVICE_NONE || key == gopi.KEYCODE_NONE {
		return gopi.ErrBadParameter
	}

	// Get existing a mapentry
	hash := keyForEntry(device, code)
	if hash == 0 {
		return gopi.ErrBadParameter
	}
	entry, exists := this.lookup[hash]

	// Create a new map entry or modify the existing one
	if exists == false {
		entry = &mapentry{Device: device, Code: code, Key: key, Index: len(this.mapping)}
		this.mapping = append(this.mapping, entry)
		this.lookup[hash] = entry
	} else {
		entry.Key = key
	}

	// Set dirty flag
	this.dirty = true

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func keyForEntry(device gopi.InputDeviceType, code uint32) uint64 {
	// Creates a lookup key for device/scancode combo or
	// returns zero
	if device == gopi.INPUT_DEVICE_NONE {
		return 0
	}
	return uint64(device)<<32 | uint64(code)
}
