package argonone

import (
	"time"
)

// Value represents a value which changes over time once the final value has been
// established, like a Schmitt trigger
type Value struct {
	// The current value and proposed new value
	current, new interface{}
	// How long the new value needs to persist for
	delta time.Duration
	// When the new value last changed
	change time.Time
}

// Construct a value which remains constant for at least d
// time duration.
func NewValueWithDelta(d time.Duration) *Value {
	return &Value{nil, nil, d, time.Time{}}
}

// Get the current value
func (v *Value) Get() (interface{}, bool) {
	changed := false
	if v.change.IsZero() {
		return v.current, changed
	}
	if time.Since(v.change) >= v.delta {
		changed = true
		v.current = v.new
	}
	return v.current, changed
}

// Set a new value, and return the current value
func (v *Value) Set(new interface{}) (interface{}, bool) {
	if v.change.IsZero() {
		v.new = new
		v.current = new
		v.change = time.Now()
		return new, true
	}
	if new != v.new {
		v.new = new
		v.change = time.Now()
	}
	return v.Get()
}
