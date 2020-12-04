package metrics

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"strings"
	"sync"
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

const (
	retrycount      = 5
	retrydurationms = 10
)

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

func (this *metrics) Field(name string, value ...interface{}) gopi.Field {
	if len(value) > 1 {
		return nil
	}
	field := NewField(name)
	if field == nil {
		return nil
	} else if len(value) == 1 {
		if err := field.SetValue(value[0]); err != nil {
			return nil
		}
	}
	return field
}

// Emit metrics for a named measurement, omitting timestamp
func (this *metrics) Emit(name string, values ...interface{}) error {
	return this.EmitTS(name, time.Time{}, values...)
}

// EmitTS emits metrics for a named measurement, with defined timestamp
// will retry if the channel is temporarily full
func (this *metrics) EmitTS(name string, ts time.Time, values ...interface{}) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Clone measurement
	if m, exists := this.m[name]; exists == false {
		return gopi.ErrBadParameter.WithPrefix("Emit", name)
	} else if m, err := m.Clone(ts, values...); err != nil {
		return err
	} else {
		for i := 0; i < retrycount; i++ {
			if err := this.Publisher.Emit(m, false); errors.Is(err, gopi.ErrChannelFull) {
				time.Sleep(time.Millisecond * retrydurationms)
			} else if err != nil {
				return err
			} else {
				return nil
			}
		}
		return gopi.ErrChannelFull.WithPrefix("Emit", name)
	}
}

func (this *metrics) Measurements() []gopi.Measurement {
	m := make([]gopi.Measurement, 0, len(this.m))
	for _, v := range this.m {
		m = append(m, v)
	}
	return m
}

func (this *metrics) HostTag() gopi.Field {
	host, _ := os.Hostname()
	return NewField("host", host)
}

func (this *metrics) UserTag() gopi.Field {
	if user, _ := user.Current(); user != nil {
		return NewField("user", user.Username)
	} else {
		return nil
	}
}

func (this *metrics) EnvTag(name string) gopi.Field {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil
	} else {
		return NewField(name, os.Getenv(name))
	}
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
/*
func parseField(src string) (gopi.Field, error) {
	var field gopi.Field
	var s scanner.Scanner
	s.Init(strings.NewReader(src))

	// Start state looking for an identifier
	state := stateIdent
	value := ""

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
			value += text
		default:
			return nil, gopi.ErrInternalAppError
		}
	}

	// Return if we only have the identifier
	if state == stateEquals {
		return field, nil
	} else if state != stateValue || value == "" {
		return nil, gopi.ErrBadParameter.WithPrefix(field.Name())
	}

	if value := parseValue(value); value == nil {
		return nil, gopi.ErrBadParameter.WithPrefix(field.Name())
	} else if err := field.SetValue(value); err != nil {
		return nil, err
	}

	// Return success
	return field, nil
}

func parseValue(src string) interface{} {
	// bool
	if v, err := strconv.ParseBool(src); err == nil {
		return v
	}
	// time in RFC3339 format
	if v, err := time.Parse(time.RFC3339, src); err == nil {
		return v
	}
	// If ends with 'i' then parse integer
	if strings.HasSuffix(src, "i") {
		src_ := strings.TrimSuffix(src, "i")
		if v, err := strconv.ParseInt(src_, 0, 64); err == nil {
			return v
		}
	}
	// Parse uint
	if v, err := strconv.ParseUint(src, 0, 64); err == nil {
		return v
	}
	// Parse int
	if v, err := strconv.ParseInt(src, 0, 64); err == nil {
		return v
	}
	// Parse float
	if v, err := strconv.ParseFloat(src, 64); err == nil {
		return v
	}
	// Parse quoted string
	if v, err := strconv.Unquote(src); err == nil {
		return v
	}
	// Unable to interpret
	return nil
}
*/
