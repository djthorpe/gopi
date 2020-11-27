package influxdb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/djthorpe/gopi/v3"
)

// Ref:
// https://docs.influxdata.com/influxdb/v1.8/write_protocols/line_protocol_tutorial/

// QuoteMeasurement returns a line for a measurement
func QuoteMeasurement(m gopi.Measurement) (string, error) {
	str := ""
	// Append name
	if name := m.Name(); name == "" {
		return "", gopi.ErrBadParameter.WithPrefix("name")
	} else {
		str += name
	}

	// Append tags
	if tags, err := QuoteFields(m.Tags()); err != nil {
		return "", err
	} else if tags != "" {
		str += "," + tags
	}

	// Append metrics
	if metrics, err := QuoteFields(m.Metrics()); err != nil {
		return "", err
	} else if metrics == "" {
		return "", gopi.ErrBadParameter.WithPrefix("metrics")
	} else {
		str += " " + metrics
	}

	// Append timestamp
	if ts := m.Time(); ts.IsZero() == false {
		str += " " + fmt.Sprint(ts.UnixNano())
	}

	// Return line
	return str, nil
}

func QuoteFields(fields []gopi.Field) (string, error) {
	str := ""
	for _, f := range fields {
		// Skip tags which have nil values
		if f.IsNil() {
			continue
		}
		// Add comma to split the tags
		if str != "" {
			str += ","
		}
		// Add name=value
		if value, err := QuoteField(f); err != nil {
			return "", err
		} else {
			str += value
		}
	}
	return str, nil
}

func QuoteField(f gopi.Field) (string, error) {
	if name := QuoteFieldName(f); name == "" {
		return "", gopi.ErrBadParameter
	} else if value := QuoteFieldValue(f); value == "" {
		return "", gopi.ErrBadParameter
	} else {
		return name + "=" + value, nil
	}
}

// QuoteFieldName returns quoted version of the name, returning
// empty string if invalid
func QuoteFieldName(f gopi.Field) string {
	name := strings.TrimSpace(f.Name())
	if name == "" || IsReservedName(name) {
		return ""
	}
	// TODO: Equals, space, backslash and comma are quoted
	return name
}

// QuoteFieldName returns quoted version of the value, returning
// empty string if invalid or nil
func QuoteFieldValue(f gopi.Field) string {
	if f.IsNil() {
		return ""
	}
	switch f.Kind() {
	case "nil":
		return ""
	case "string":
		return strconv.Quote(f.Value().(string))
	case "bool", "float32", "float64":
		return fmt.Sprint(f.Value())
	case "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64":
		return fmt.Sprint(f.Value()) + "i"
	case "time.Time":
		return fmt.Sprint(f.Value().(time.Time).UnixNano())
	default:
		return ""
	}
}

func IsReservedName(name string) bool {
	return name == "time" || name == "_field" || name == "_measurement"
}
