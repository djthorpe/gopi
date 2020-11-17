package metrics

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"time"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type field struct {
	name  string
	value interface{}
	kind
	sync.RWMutex
}

type kind int

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	// Field name can't include numeric, minus or underscore as prefix
	reFieldName = regexp.MustCompile("^[A-Za-z][A-Za-z0-9_\\-]*$")
)

var (
	mapKind = make(map[string]kind, 10)
)

const (
	kNone kind = iota
	kString
	kBool
	kUint8
	kUint16
	kUint32
	kUint64
	kInt8
	kInt16
	kInt32
	kInt64
	kFloat32
	kFloat64
	kTime
	kMax
)

////////////////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

func NewField(name string, value ...interface{}) gopi.Field {
	// Check incoming parameters
	if name == "" || len(value) > 1 {
		return nil
	} else if reFieldName.MatchString(name) == false {
		return nil
	}

	// Create a new field
	this := new(field)
	this.name = name

	// Set field value
	if len(value) == 1 {
		if err := this.SetValue(value[0]); err != nil {
			return nil
		}
	}

	// Return field
	return this
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *field) Name() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.name
}

func (this *field) Kind() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	return this.kind.String()
}

func (this *field) IsNil() bool {
	return this.kind == kNone || this.value == nil
}

func (this *field) Value() interface{} {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Returns zero-value for type if nil
	if this.IsNil() {
		return this.kind.ZeroValue()
	} else {
		return this.value
	}
}

func (this *field) SetValue(v interface{}) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	if v == nil {
		this.value = nil
		this.kind = kNone
	} else {
		switch v.(type) {
		case uint8:
			this.value = v
			this.kind = kUint8
		case uint16:
			this.value = v
			this.kind = kUint16
		case uint32:
			this.value = v
			this.kind = kUint32
		case uint64:
			this.value = v
			this.kind = kUint64
		case int8:
			this.value = v
			this.kind = kInt8
		case int16:
			this.value = v
			this.kind = kInt16
		case int32:
			this.value = v
			this.kind = kInt32
		case int64:
			this.value = v
			this.kind = kInt64
		case string:
			this.value = v
			this.kind = kString
		case bool:
			this.value = v
			this.kind = kBool
		case float32:
			this.value = v
			this.kind = kFloat32
		case float64:
			this.value = v
			this.kind = kFloat64
		case time.Time:
			this.value = v
			this.kind = kTime
		default:
			return gopi.ErrBadParameter.WithPrefix(this.name)
		}
	}

	// Return success
	return nil
}

func (this *field) SetKind(k string) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	var once sync.Once
	// Initialize map
	once.Do(func() {
		for v := kNone; v < kMax; v++ {
			key := fmt.Sprint(v)
			mapKind[key] = v
		}
	})
	// Reset the value to nil
	if kind, exists := mapKind[k]; exists == false {
		return gopi.ErrBadParameter.WithPrefix(k)
	} else {
		this.value = nil
		this.kind = kind
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *field) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	str := "<field"
	str += " name=" + strconv.Quote(this.name)
	if this.kind != kNone {
		str += " type=" + fmt.Sprint(this.kind)
		if this.IsNil() == false {
			if this.kind == kString {
				str += " default=" + strconv.Quote(this.Value().(string))
			} else {
				str += " default=" + fmt.Sprint(this.Value())
			}
		}
	}
	return str + ">"
}

func (k kind) String() string {
	switch k {
	case kNone:
		return "nil"
	case kString:
		return "string"
	case kBool:
		return "bool"
	case kUint8:
		return "uint8"
	case kUint16:
		return "uint16"
	case kUint32:
		return "uint32"
	case kUint64:
		return "uint64"
	case kInt8:
		return "int8"
	case kInt16:
		return "int16"
	case kInt32:
		return "int32"
	case kInt64:
		return "int64"
	case kFloat32:
		return "float32"
	case kFloat64:
		return "float64"
	case kTime:
		return "time.Time"
	default:
		return "[?? Invalid kind]"
	}
}

func (k kind) ZeroValue() interface{} {
	switch k {
	case kString:
		return ""
	case kBool:
		return false
	case kUint8:
		return uint8(0)
	case kUint16:
		return uint16(0)
	case kUint32:
		return uint32(0)
	case kUint64:
		return uint64(0)
	case kInt8:
		return int8(0)
	case kInt16:
		return int16(0)
	case kInt32:
		return int32(0)
	case kInt64:
		return int64(0)
	case kFloat32:
		return float32(0)
	case kFloat64:
		return float64(0)
	case kTime:
		return time.Time{}
	default:
		return nil
	}
}
