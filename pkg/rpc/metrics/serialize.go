package metrics

import (
	"fmt"
	"strconv"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	ptypes "github.com/golang/protobuf/ptypes"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type measurement struct {
	pb *Measurement
}

type field struct {
	pb *Field
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func fromProtoMeasurement(pb *Measurement) gopi.Measurement {
	if pb == nil {
		return nil
	} else {
		return &measurement{pb}
	}
}

func toProtoNull() *Measurement {
	return &Measurement{}
}

func toProtoMeasurement(m gopi.Measurement) *Measurement {
	if m == nil {
		return nil
	}

	return &Measurement{
		Name:    m.Name(),
		Tags:    toProtoFields(m.Tags()),
		Metrics: toProtoFields(m.Metrics()),
		Ts:      toProtoTimestamp(m.Time()),
	}
}

func toProtoFields(fields []gopi.Field) []*Field {
	result := make([]*Field, len(fields))
	for i, field := range fields {
		result[i] = &Field{
			Name: field.Name(),
			Kind: field.Kind(),
		}
		if field.IsNil() == false {
			result[i].Value = toProtoFieldValue(field)
		}
	}
	return result
}

func toProtoTimestamp(ts time.Time) *timestamp.Timestamp {
	if ts.IsZero() {
		return nil
	} else if proto, err := ptypes.TimestampProto(ts); err == nil {
		return proto
	} else {
		return nil
	}
}

func toProtoFieldValue(field gopi.Field) isField_Value {
	v := field.Value()
	switch value := v.(type) {
	case string:
		return &Field_Str{value}
	case uint8:
		return &Field_Uint{uint64(value)}
	case uint16:
		return &Field_Uint{uint64(value)}
	case uint32:
		return &Field_Uint{uint64(value)}
	case uint64:
		return &Field_Uint{uint64(value)}
	case int8:
		return &Field_Int{int64(value)}
	case int16:
		return &Field_Int{int64(value)}
	case int32:
		return &Field_Int{int64(value)}
	case int64:
		return &Field_Int{int64(value)}
	case float32:
		return &Field_Float{float64(value)}
	case float64:
		return &Field_Float{float64(value)}
	case bool:
		return &Field_Bool{value}
	case time.Time:
		return &Field_Time{toProtoTimestamp(value)}
	default:
		return nil
	}
}

/////////////////////////////////////////////////////////////////////
// MEASUREMENT

func (this *measurement) Name() string {
	return this.pb.GetName()
}

func (this *measurement) Time() time.Time {
	if ts, err := ptypes.Timestamp(this.pb.GetTs()); err != nil {
		return time.Time{}
	} else {
		return ts
	}
}

func (this *measurement) Tags() []gopi.Field {
	tags := this.pb.GetTags()
	result := make([]gopi.Field, len(tags))
	for i, tag := range tags {
		result[i] = &field{tag}
	}
	return result
}

func (this *measurement) Metrics() []gopi.Field {
	metrics := this.pb.GetMetrics()
	result := make([]gopi.Field, len(metrics))
	for i, metric := range metrics {
		result[i] = &field{metric}
	}
	return result
}

func (this *measurement) Get(name string) interface{} {
	// TODO
	return nil
}

func (this *measurement) Set(string, interface{}) error {
	return gopi.ErrNotImplemented
}

func (this *measurement) String() string {
	str := "<measurement"
	str += " name=" + strconv.Quote(this.Name())
	if ts := this.Time(); ts.IsZero() == false {
		str += " ts=" + ts.Format(time.RFC3339)
	}
	if tags := this.Tags(); len(tags) > 0 {
		str += " tags=" + fmt.Sprint(tags)
	}
	if metrics := this.Metrics(); len(metrics) > 0 {
		str += " metrics=" + fmt.Sprint(metrics)
	}
	return str + ">"
}

/////////////////////////////////////////////////////////////////////
// FIELD

func (this *field) Name() string {
	return this.pb.GetName()
}

func (this *field) Kind() string {
	return this.pb.GetKind()
}

func (this *field) IsNil() bool {
	return this.pb.Value == nil
}

func (this *field) Value() interface{} {
	kind := this.pb.GetKind()
	switch kind {
	case "string":
		return this.pb.GetStr()
	case "uint8":
		return uint8(this.pb.GetUint())
	case "uint16":
		return uint16(this.pb.GetUint())
	case "uint32":
		return uint32(this.pb.GetUint())
	case "uint64":
		return uint64(this.pb.GetUint())
	case "int8":
		return int8(this.pb.GetInt())
	case "int16":
		return int16(this.pb.GetInt())
	case "int32":
		return int32(this.pb.GetInt())
	case "int64":
		return int64(this.pb.GetInt())
	case "float32":
		return float32(this.pb.GetFloat())
	case "float64":
		return float64(this.pb.GetFloat())
	case "bool":
		return this.pb.GetBool()
	case "time.Time":
		if ts, err := ptypes.Timestamp(this.pb.GetTime()); err != nil {
			return time.Time{}
		} else {
			return ts
		}
	default:
		return nil
	}
}

func (this *field) SetValue(interface{}) error {
	return gopi.ErrNotImplemented
}

func (this *field) Copy() gopi.Field {
	return nil
}

func (this *field) String() string {
	str := "<field"
	str += " name=" + strconv.Quote(this.Name())
	if kind := this.Kind(); kind != "" {
		str += " type=" + fmt.Sprint(kind)
		if this.IsNil() == false {
			if kind == "string" {
				str += " default=" + strconv.Quote(this.Value().(string))
			} else {
				str += " default=" + fmt.Sprint(this.Value())
			}
		}
	}
	return str + ">"
}
