package graph

import (
	"reflect"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	unitType   = reflect.TypeOf((*gopi.Unit)(nil)).Elem()
	stubType   = reflect.TypeOf((*gopi.ServiceStub)(nil)).Elem()
	loggerType = reflect.TypeOf((*gopi.Logger)(nil)).Elem()
)

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// forEachField calls a function for each field of a struct ptr
// and returns all errors, or immediately with a single error if
// immediate is set
func forEachField(unit reflect.Value, immediate bool, fn func(reflect.StructField, int) error) error {
	var result error
	if isStructPtr(unit.Type()) {
		t := unit.Elem().Type()
		for i := 0; i < t.NumField(); i++ {
			if err := fn(t.Field(i), i); err != nil {
				if immediate {
					return err
				} else {
					result = multierror.Append(result, err)
				}
			}
		}
	}
	return result
}

// isStructPtr returns true if the type is a pointer to a struct
func isStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

// isUnitType returns true if a struct ptr contains a gopi.Unit
// type
func isUnitType(t reflect.Type) bool {
	if isStructPtr(t) == false {
		return false
	}
	t = t.Elem()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous && equalsType(f.Type, unitType) {
			return true
		}
	}
	return false
}

// isLoggerType returns true if a struct ptr is a gopi.Logger
func isLoggerType(t reflect.Type) bool {
	if isStructPtr(t) == false {
		return false
	} else if t.Implements(loggerType) == false {
		return false
	}
	// Return success
	return true
}

// isServiceStubType returns true if a type implements a stubservice
func isServiceStubType(t reflect.Type) bool {
	if isStructPtr(t) == false {
		return false
	} else if t.Implements(stubType) == false {
		return false
	}
	// Return success
	return true
}

// equalsType returns true if two types are equivalent
func equalsType(a, b reflect.Type) bool {
	return a.Name() == b.Name() && a.PkgPath() == b.PkgPath()
}

// callFn will call a function on a struct and pass arguments
// but expects the first returned argument to be an error
func callFn(name string, unit reflect.Value, args []reflect.Value) error {
	if fn := unit.MethodByName(name); fn.IsValid() == false {
		return nil
	} else if ret := fn.Call(args); len(ret) != 1 {
		return nil
	} else if err, ok := ret[0].Interface().(error); ok {
		return err
	} else if err == nil {
		return nil
	} else {
		return gopi.ErrBadParameter.WithPrefix(name)
	}
}
