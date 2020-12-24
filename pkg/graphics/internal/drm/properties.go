// +build drm

package drm

import (
	"fmt"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	drm "github.com/djthorpe/gopi/v3/pkg/sys/drm"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Properties struct {
	sync.RWMutex

	fd    uintptr
	id    uint32
	props map[string]*property
	dirty bool
}

type property struct {
	id    uint32
	value uint64
	dirty bool
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *Properties) New(fd uintptr, id uint32) error {
	if fd == 0 || id == 0 {
		return gopi.ErrInternalAppError.WithPrefix("NewProperties")
	} else {
		this.fd = fd
		this.id = id
	}

	props := drm.GetAnyProperties(fd, id)
	if props == nil {
		return gopi.ErrInternalAppError.WithPrefix("NewProperties")
	}
	defer props.Free()
	if this.props = this.init(props); this.props == nil {
		return gopi.ErrInternalAppError.WithPrefix("NewProperties")
	}

	// Return success
	return nil
}

func (this *Properties) Dispose() error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	this.props = nil
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

// Dirty returns true is properties have been changed
func (this *Properties) Dirty() bool {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.dirty
}

// Clean removes any dirty flags
func (this *Properties) Clean() {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	for name := range this.props {
		this.props[name].dirty = false
	}
	this.dirty = false
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// GetProperty returns the current property value, and a boolean
// indicating if that property exists
func (this *Properties) GetProperty(name string) (uint64, bool) {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	property, exists := this.props[name]
	if exists {
		return property.value, true
	} else {
		return 0, false
	}
}

// SetProperty sets a property value. Returns false if that
// property does not exist
func (this *Properties) SetProperty(name string, value uint64) bool {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if property, exists := this.props[name]; exists {
		if property.value != value {
			this.props[name].value = value
			this.props[name].dirty = true
			this.dirty = true
		}
		return true
	} else {
		return false
	}
}

// SetProperties sets property values in bulk. Returns an
// error if any property is not supported
func (this *Properties) SetProperties(tuples map[string]uint64) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var result error
	for name, value := range tuples {
		if property, exists := this.props[name]; exists {
			if property.value != value {
				this.props[name].value = value
				this.props[name].dirty = true
				this.dirty = true
			}
		} else {
			return gopi.ErrBadParameter.WithPrefix("SetProperties: ", name)
		}
	}

	// Return any errors
	return result
}

// GetDirtyProperties returns properties which have been changed,
// keyed by their property ID rather than their name
func (this *Properties) GetDirtyProperties() map[uint32]uint64 {
	// Check for case where no values have changed
	if this.dirty == false {
		return nil
	}
	// Iterate through values
	result := make(map[uint32]uint64, len(this.props))
	for _, v := range this.props {
		if v.dirty == false {
			continue
		}
		result[v.id] = v.value
	}
	// Return result
	return result
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Properties) String() string {
	str := "<properties"
	for k, v := range this.props {
		str += fmt.Sprintf(" %v=%v", k, v)
	}
	return str + ">"
}

func (this *property) String() string {
	if this.dirty {
		return fmt.Sprintf("<%d:%d>[dirty]", this.id, this.value)
	} else {
		return fmt.Sprintf("<%d:%d>", this.id, this.value)
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Properties) init(props *drm.Properties) map[string]*property {
	keys, values := props.Keys(), props.Values()
	result := make(map[string]*property, len(keys))
	for i, key := range keys {
		value := values[i]
		prop := drm.NewProperty(this.fd, key)
		if prop == nil {
			continue
		}
		defer prop.Free()
		name := prop.Name()
		result[name] = &property{key, value, false}
	}
	return result
}
