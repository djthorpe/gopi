package ircodec

import (
	"math"

	"github.com/djthorpe/gopi/v3"
)

/////////////////////////////////////////////////////////////////////
// Mark Space implementation

// A mark or space value
type MarkSpace struct {
	Type            gopi.LIRCType
	Value, Min, Max uint32
}

func NewMarkSpace(t gopi.LIRCType, value, tolerance uint32) *MarkSpace {
	this := &MarkSpace{Type: t}
	this.Set(value, tolerance)
	return this
}

func (m *MarkSpace) Set(value, tolerance uint32) {
	delta := float64(value) * float64(tolerance) / 100.0
	m.Min = uint32(math.Max(0, float64(value)-delta))
	m.Max = uint32(float64(value) + delta)
	m.Value = value
}

func (m *MarkSpace) Matches(evt gopi.LIRCEvent) bool {
	if evt.Type() != m.Type {
		return false
	}
	if evt.Mode() != gopi.LIRC_MODE_MODE2 {
		return false
	}
	if m.Min > evt.Value().(uint32) {
		return false
	}
	if m.Max < evt.Value().(uint32) {
		return false
	}
	return true
}

func (m *MarkSpace) GreaterThan(evt gopi.LIRCEvent) bool {
	if m.Type != evt.Type() {
		return false
	}
	if evt.Mode() != gopi.LIRC_MODE_MODE2 {
		return false
	}
	if evt.Value().(uint32) < m.Min {
		return false
	}
	return true
}

func (m *MarkSpace) LessThan(evt gopi.LIRCEvent) bool {
	if m.Type != evt.Type() {
		return false
	}
	if evt.Mode() != gopi.LIRC_MODE_MODE2 {
		return false
	}
	if evt.Value().(uint32) > m.Max {
		return false
	}
	return true
}
