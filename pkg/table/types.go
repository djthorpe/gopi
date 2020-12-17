package table

import (
	"reflect"
	"strconv"
	"time"
)

/////////////////////////////////////////////////////////////////////
// TYPES

// Kind represents supported kinds of data
type Kind uint

// Types represents possible kinds for a value
type Types map[Kind]bool

/////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	String Kind = iota
	Float
	Int
	Uint
	Bool
	Time
	Duration
	Nil
	kindMin = String
	kindMax = Nil
)

/////////////////////////////////////////////////////////////////////
// DATE FORMATS

var (
	dateFmts = []string{
		time.RFC3339,
		time.RFC3339Nano,
		time.ANSIC,
		time.StampMicro,
		time.StampMilli,
		time.Stamp,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
		"2006/01/02",
		"02-01-2006",
		"02/01/2006",
	}
)

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (t Types) Del(k Kind) {
	delete(t, k)
}

func (t Types) Set(k Kind) {
	t[k] = true
}

func (t Types) Get(k Kind) bool {
	_, exists := t[k]
	return exists
}

func (t Types) Parse(v interface{}) {
	// Initialise Types
	if len(t) == 0 {
		for k := kindMin; k <= kindMax; k++ {
			t.Set(k)
		}
	}
	// Parse from possibles
	rv := reflect.ValueOf(v)
	for k := range t {
		if isKind(rv, k) == false {
			t.Del(k)
		}
	}
}

// Kind returns the most likely kind of data for this value
func (t Types) Kind() Kind {
	for k := kindMax; k > kindMin; k-- {
		if t.Get(k) {
			return k
		}
	}
	return kindMin
}

func (t Types) Kinds() []Kind {
	result := []Kind{}
	for k := kindMax; k > kindMin; k-- {
		if t.Get(k) {
			result = append(result, k)
		}
	}
	return append(result, kindMin)
}

func (k Kind) String() string {
	switch k {
	case String:
		return "string"
	case Nil:
		return "nil"
	case Uint:
		return "uint"
	case Int:
		return "int"
	case Float:
		return "float"
	case Bool:
		return "bool"
	case Time:
		return "time"
	case Duration:
		return "duration"
	}
	return ""
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func alignmentForKind(k Kind) Alignment {
	switch k {
	case Uint, Int, Float, Bool, Time, Duration:
		return Right
	default:
		return Left
	}
}

func isKind(rv reflect.Value, k Kind) bool {
	switch k {
	case Nil:
		if isNil(rv) {
			return true
		} else if rv.Kind() == reflect.String {
			return rv.IsZero()
		} else {
			return false
		}
	case String:
		return true
	case Uint:
		return isUint(rv)
	case Int:
		return isInt(rv)
	case Float:
		return isFloat(rv)
	case Bool:
		return isBool(rv)
	case Time:
		return isTime(rv)
	case Duration:
		return isDuration(rv)
	}
	return false
}

func isUint(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Uintptr:
		return true
	case reflect.String:
		if rv.IsZero() {
			return true
		} else if _, err := strconv.ParseUint(rv.String(), 0, 64); err == nil {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

func isInt(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.String:
		if rv.IsZero() {
			return true
		} else if _, err := strconv.ParseInt(rv.String(), 0, 64); err == nil {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

func isFloat(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	case reflect.String:
		if rv.IsZero() {
			return true
		} else if _, err := strconv.ParseFloat(rv.String(), 64); err == nil {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

func isBool(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Bool:
		return true
	case reflect.String:
		if rv.IsZero() {
			return true
		} else if _, err := strconv.ParseBool(rv.String()); err == nil {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}
func isTime(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Struct:
		return rv.Type() == reflect.TypeOf(time.Time{})
	case reflect.String:
		str := rv.String()
		for _, dateFmt := range dateFmts {
			if _, err := time.Parse(dateFmt, str); err == nil {
				return true
			}
		}
	}
	return false
}

func isDuration(rv reflect.Value) bool {
	switch rv.Kind() {
	case reflect.Struct:
		return rv.Type() == reflect.TypeOf(time.Duration(0))
	case reflect.String:
		if rv.IsZero() {
			return true
		} else if _, err := time.ParseDuration(rv.String()); err == nil {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}
