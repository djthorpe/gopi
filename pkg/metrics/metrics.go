package metrics

import (
	"fmt"
	"strings"
	"sync"
	"text/scanner"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type metrics struct {
	gopi.Unit
	sync.RWMutex
	gopi.Publisher

	m map[string]*measurement
}

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *metrics) New(cfg gopi.Config) error {
	this.m = make(map[string]*measurement)

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - METRICS

// Define a measurement with metric definitions and optional tag fields
func (this *metrics) NewMeasurement(name, metrics string, tags ...gopi.Field) (gopi.Measurement, error) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for duplicate name
	if _, exists := this.m[name]; exists {
		return nil, gopi.ErrDuplicateEntry.WithPrefix(name)
	}

	if measurement, err := NewMeasurement(name, metrics, tags...); err != nil {
		return nil, err
	} else {
		key := measurement.Name()
		this.m[key] = measurement
	}

	// Return success
	return this.m[name], nil
}

// NewFields returns array of fields. Some elements may be set to nil
// where a parse error occured
func (this *metrics) NewFields(values ...string) []gopi.Field {
	fields := make([]gopi.Field, len(values))
	for i, value := range values {
		f, _ := parseField(value)
		fields[i] = f
	}
	return fields
}

// Emit metrics for a named measurement, omitting timestamp
func (this *metrics) Emit(name string, values ...interface{}) error {
	return this.EmitTS(name, time.Time{}, values...)
}

// EmitTS emits metrics for a named measurement, with defined timestamp
func (this *metrics) EmitTS(name string, ts time.Time, values ...interface{}) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Clone measurement
	// TODO
	/*	if m, exists := this.m[name]; exists == false {
			return gopi.ErrBadParameter.WithPrefix(name)
		} else if m2, err := m.Clone(ts, values...); err != nil {
			return err
		} else {
			this.Publisher.Emit(m2)
		}
	*/
	// Return success
	return nil
}

func (this *metrics) Measurements() []gopi.Measurement {
	m := make([]gopi.Measurement, 0, len(this.m))
	for _, v := range this.m {
		m = append(m, v)
	}
	return m
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *metrics) String() string {
	str := "<metrics"
	for k, v := range this.m {
		str += " " + k + "=" + fmt.Sprint(v)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func parseField(src string) (gopi.Field, error) {
	var field gopi.Field
	var s scanner.Scanner
	s.Init(strings.NewReader(src))

	// Start state looking for an identifier
	state := stateIdent

	// Scan tokens
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		text := s.TokenText()
		switch state {
		case stateIdent:
			if field = NewField(text); field == nil {
				return nil, gopi.ErrBadParameter.WithPrefix(text)
			}
			state = stateEquals
		case stateEquals:
			if text != "=" {
				return nil, gopi.ErrBadParameter.WithPrefix(field.Name())
			}
			state = stateValue
		case stateValue:
			if value := parseValue(text); value == nil {
				return nil, gopi.ErrBadParameter.WithPrefix(field.Name())
			} else if err := field.SetValue(value); err != nil {
				return nil, err
			}
			state = stateDone
		case stateDone:
			return nil, gopi.ErrBadParameter.WithPrefix(field.Name())
		default:
			return nil, gopi.ErrInternalAppError
		}
	}

	// Field without default value
	if state == stateEquals || state == stateDone {
		return field, nil
	}

	// Return success
	return nil, gopi.ErrNotImplemented
}

func parseValue(src string) interface{} {
	return nil
}
