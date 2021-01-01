package metrics

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/scanner"
	"time"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type measurement struct {
	name    string
	ts      time.Time
	metrics []gopi.Field
	tags    []gopi.Field
	fields  map[string]gopi.Field
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	stateIdent = iota
	stateIdent2
	stateValue
	stateDone
)

var (
	reMeasurementName = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_\\-]*$")
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewMeasurement(name, metrics string, tags ...gopi.Field) (*measurement, error) {
	// Create measurement
	this := new(measurement)

	// Check measurement name
	if reMeasurementName.MatchString(name) == false {
		return nil, gopi.ErrBadParameter.WithPrefix(name)
	} else {
		this.name = name
	}

	// Check tags
	if hasNilElement(tags) {
		return nil, gopi.ErrBadParameter.WithPrefix(name)
	} else if dup := duplicateName(tags); dup != "" {
		return nil, gopi.ErrDuplicateEntry.WithPrefix(name)
	} else {
		this.tags = tags
	}

	// Parse metrics and check metrics
	metrics_, err := parseMetrics(metrics)
	if err != nil {
		return nil, err
	} else if len(metrics_) == 0 || hasNilElement(metrics_) {
		return nil, gopi.ErrBadParameter.WithPrefix(name)
	} else {
		this.metrics = metrics_
	}

	// Map fields
	this.fields = make(map[string]gopi.Field, len(tags)+len(metrics))
	for _, field := range this.metrics {
		key := field.Name()
		if _, exists := this.fields[key]; exists {
			return nil, gopi.ErrDuplicateEntry.WithPrefix(key)
		} else {
			this.fields[key] = field
		}
	}
	for _, field := range this.tags {
		key := field.Name()
		if _, exists := this.fields[key]; exists {
			return nil, gopi.ErrDuplicateEntry.WithPrefix(key)
		} else {
			this.fields[key] = field
		}
	}

	// Return success
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *measurement) Name() string {
	return this.name
}

func (this *measurement) Time() time.Time {
	return this.ts
}

func (this *measurement) Tags() []gopi.Field {
	return this.tags
}

func (this *measurement) Metrics() []gopi.Field {
	return this.metrics
}

func (this *measurement) Get(name string) interface{} {
	if field, exists := this.fields[name]; exists == false {
		return nil
	} else {
		return field.Value()
	}
}
func (this *measurement) Set(name string, value interface{}) error {
	if field, exists := this.fields[name]; exists == false {
		return nil
	} else {
		return field.SetValue(value)
	}
}

func (this *measurement) Clone(ts time.Time, tags []gopi.Field, values ...interface{}) (*measurement, error) {
	// Check correct number of arguments
	if len(values) != len(this.metrics) {
		return nil, gopi.ErrBadParameter.WithPrefix("Clone")
	}

	that := new(measurement)
	that.name = this.name
	that.ts = ts
	that.fields = make(map[string]gopi.Field, len(this.fields))

	// Index new tags and use them instead of defaults
	for _, tag := range tags {
		fmt.Println("TODO: Clone", tag)
	}

	// Clone tags
	that.tags = make([]gopi.Field, len(this.tags))
	for i, value := range this.tags {
		field := value.Copy()
		that.tags[i] = field
		that.fields[field.Name()] = field
	}

	// Clone metrics and set new values
	that.metrics = make([]gopi.Field, len(this.metrics))
	for i, value := range this.metrics {
		field := value.Copy()
		key := field.Name()
		if err := field.SetValue(values[i]); err != nil {
			return nil, fmt.Errorf("Clone: %q: %w", key, err)
		}
		that.metrics[i] = field
		that.fields[key] = field
	}

	// Return success
	return that, nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *measurement) String() string {
	str := "<measurement"
	str += " name=" + strconv.Quote(this.name)
	if ts := this.Time(); ts.IsZero() == false {
		str += " ts=" + ts.Format(time.RFC3339)
	}
	if len(this.tags) > 0 {
		str += " tags=" + fmt.Sprint(this.tags)
	}
	if len(this.metrics) > 0 {
		str += " metrics=" + fmt.Sprint(this.metrics)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func parseMetrics(src string) ([]gopi.Field, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(src))

	// Start state looking for an identifier
	state := stateIdent
	metrics := []gopi.Field{}

	// Scan tokens
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		value := s.TokenText()
		switch state {
		case stateIdent2:
			if value == "," {
				state = stateIdent
				break
			}
			fallthrough
		case stateIdent:
			if f := NewField(value); f == nil {
				return nil, gopi.ErrBadParameter.WithPrefix(value)
			} else {
				metrics = append(metrics, f)
			}
			state = stateValue
		case stateValue:
			if value == "," {
				state = stateIdent
				break
			}
			for _, f := range metrics {
				if f.Kind() == "nil" {
					if err := f.(*field).SetKind(value); err != nil {
						return nil, gopi.ErrBadParameter.WithPrefix(value)
					}
				}
			}
			state = stateIdent2
		default:
			return nil, gopi.ErrInternalAppError
		}
	}

	// Check state is as expected
	if state != stateIdent && state != stateIdent2 {
		return nil, gopi.ErrBadParameter.WithPrefix("metrics")
	}

	// Check for duplicates
	if dup := duplicateName(metrics); dup != "" {
		return nil, gopi.ErrDuplicateEntry.WithPrefix(dup)
	}

	// Return success
	return metrics, nil
}

func duplicateName(fields []gopi.Field) string {
	m := make(map[string]bool, len(fields))
	for _, field := range fields {
		key := field.Name()
		if _, exists := m[key]; exists {
			return key
		} else {
			m[key] = true
		}
	}
	// No duplicate found
	return ""
}

func hasNilElement(fields []gopi.Field) bool {
	for _, field := range fields {
		if field == nil {
			return true
		}
	}
	return false
}
