package influxdb

import (
	"fmt"
	"strings"

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
	if tags := m.Tags(); len(tags) > 0 {
		for _, tag := range tags {
			str += ","
			if value, err := QuoteField(tag); err != nil {
				return "", err
			} else {
				str += value
			}
		}
	}
	// Append metrics
	if metrics := m.Metrics(); len(metrics) > 0 {
		str += " "
		for i, metric := range metrics {
			if i > 0 {
				str += ","
			}
			if value, err := QuoteField(metric); err != nil {
				return "", err
			} else {
				str += value
			}
		}
	}
	// Append timestamp
	if ts := m.Time(); ts.IsZero() == false {
		str += " " + fmt.Sprint(ts.UnixNano())
	}
	// Return line
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
// empty string if invalid
func QuoteFieldValue(f gopi.Field) string {

}

func IsReservedName(name string) bool {
	return name == "time" || name == "_field" || name == "_measurement"
}
