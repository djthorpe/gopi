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
	ts      gopi.Field
	metrics []gopi.Field
	tags    []gopi.Field
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	stateIdent = iota
	stateEquals
	stateValue
)

var (
	reMeasurementName = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_\\-]*$")
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func NewMeasurement(name, metrics string, tags ...gopi.Field) (*measurement, error) {
	// Check measurement name
	if reMeasurementName.MatchString(name) == false {
		return nil, gopi.ErrBadParameter.WithPrefix(name)
	}

	// Check tags
	if nilElement(tags) {
		return nil, gopi.ErrBadParameter.WithPrefix(name)
	} else if dup := duplicateName(tags); dup != "" {
		return nil, gopi.ErrDuplicateEntry.WithPrefix(name)
	}

	// Parse metrics and check metrics
	if metrics, err := parseMetrics(metrics); err != nil {
		return nil, err
	} else if len(metrics) == 0 || nilElement(metrics) {
		return nil, gopi.ErrBadParameter.WithPrefix(name)
	} else {
		return &measurement{name, nil, metrics, tags}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - MEASUREMENT

func (this *measurement) Name() string {
	return this.name
}

func (this *measurement) Time() time.Time {
	if this.ts == nil {
		return time.Time{}
	} else {
		return this.ts.Value().(time.Time)
	}
}

func (this *measurement) Tags() []gopi.Field {
	return this.tags
}

func (this *measurement) Metrics() []gopi.Field {
	return this.metrics
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
		case stateIdent, stateEquals:
			if f := NewField(value); f == nil {
				return nil, gopi.ErrBadParameter.WithPrefix(value)
			} else {
				metrics = append(metrics, f)
			}
			state = stateValue
		case stateValue:
			if value == "," {
				state = stateEquals
			} else {
				for _, f := range metrics {
					if f.Kind() == "nil" {
						if err := f.(*field).SetKind(value); err != nil {
							return nil, gopi.ErrBadParameter.WithPrefix(value)
						}
					}
				}
			}
			state = stateIdent
		default:
			return nil, gopi.ErrInternalAppError
		}
	}

	// Check state is as expected
	if state != stateIdent {
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

func nilElement(fields []gopi.Field) bool {
	for _, field := range fields {
		if field == nil {
			return true
		}
	}
	return false
}
