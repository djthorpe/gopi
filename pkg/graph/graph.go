package graph

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type graph struct {
	sync.RWMutex
	sync.WaitGroup

	units   map[reflect.Type]reflect.Value
	objs    []reflect.Value
	Logfn   func(...interface{})
	errs    chan error
	cancels []context.CancelFunc
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	Global = NewGraph(nil)
	iface  = make(map[reflect.Type]reflect.Type)
	stubs  = make(map[string]reflect.Type)
)

/////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

// Construct empty graph
func NewGraph(fn func(...interface{})) *graph {
	this := new(graph)
	this.units = make(map[reflect.Type]reflect.Value)
	this.Logfn = fn
	return this
}

/////////////////////////////////////////////////////////////////////
// REGISTRATION FUNCTIONS

func RegisterUnit(t, i reflect.Type) {
	if err := registerUnit(t, i); err != nil {
		panic(err)
	}
}

func RegisterServiceStub(s string, t reflect.Type) {
	if err := registerServiceStub(s, t); err != nil {
		panic(err)
	}
}

func registerUnit(t, i reflect.Type) error {
	if t == nil || i == nil {
		return gopi.ErrBadParameter.WithPrefix("RegisterUnit")
	}
	for i.Kind() == reflect.Ptr {
		i = i.Elem()
	}
	if i.Kind() != reflect.Interface {
		return gopi.ErrBadParameter.WithPrefix(i, "Not an interface")
	}
	if isUnitType(t) == false {
		return gopi.ErrBadParameter.WithPrefix(t, "Not a gopi.Unit")
	}
	if t.Implements(i) == false {
		return fmt.Errorf("%v does not implement interface %v", t, i)
	}
	if _, exists := iface[i]; exists {
		return gopi.ErrDuplicateEntry.WithPrefix(i)
	} else {
		iface[i] = t
	}

	// Return success
	return nil
}

func registerServiceStub(s string, t reflect.Type) error {
	// Check that type implements the stub interface
	if isServiceStubType(t) == false {
		return gopi.ErrNotImplemented.WithPrefix(s)
	} else if _, exists := stubs[s]; exists {
		return gopi.ErrDuplicateEntry.WithPrefix(s)
	} else {
		stubs[s] = t
	}

	// Return success
	return nil
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func Create(objs ...interface{}) (*graph, error) {
	if err := Global.Create(objs...); err != nil {
		return nil, err
	} else {
		return Global, nil
	}
}

func (this *graph) Create(objs ...interface{}) error {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()

	// Check for parameters. Create cannot be called twice
	if len(objs) == 0 {
		return gopi.ErrBadParameter
	} else if len(this.objs) != 0 {
		return gopi.ErrOutOfOrder
	}

	var result error
	for _, obj := range objs {
		obj_ := reflect.ValueOf(obj)
		if err := this.graph(obj_); err != nil {
			result = multierror.Append(result, err)
		} else {
			this.objs = append(this.objs, obj_)
		}
	}
	return unwrap(result)
}

// Call Define for each unit object
func (this *graph) Define(cfg gopi.Config) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var result error
	seen := make(map[reflect.Type]bool, len(this.units))
	for _, obj := range this.objs {
		if err := this.do("Define", obj, []reflect.Value{reflect.ValueOf(cfg)}, seen, 0); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return unwrap(result)
}

// Call New for each unit object
func (this *graph) New(cfg gopi.Config) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	seen := make(map[reflect.Type]bool, len(this.units))
	for _, obj := range this.objs {
		if err := this.do("New", obj, []reflect.Value{reflect.ValueOf(cfg)}, seen, 0); err != nil {
			return err
		}
	}

	// Return success
	return nil
}

// Call Dispose for each unit object. At the moment, the order of
// the Dispose is not considered.
func (this *graph) Dispose() error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var result error
	seen := make(map[reflect.Type]bool, len(this.units))
	for _, obj := range this.objs {
		if err := this.do("Dispose", obj, []reflect.Value{}, seen, 0); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return unwrap(result)
}

// Call Run for each unit object and wait for all to complete
func (this *graph) Run(ctx context.Context, done bool) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	// Set up channel for receiving errors from Run invocations
	this.errs = make(chan error)

	// Collect errors
	var result error
	go func() {
		for err := range this.errs {
			if err != nil && errors.Is(err, context.Canceled) == false {
				result = multierror.Append(result, err)
				this.cancelWithError(unwrap(result))
			}
		}
	}()

	// Make new context with cancel
	ctx2, cancel := context.WithCancel(ctx)
	this.cancels = append(this.cancels, cancel)

	// Call Run functions
	seen := make(map[reflect.Type]bool, len(this.units))
	for _, obj := range this.objs {
		if err := this.do("Run", obj, nil, seen, 0); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Wait for ctx.Done, then send cancels
	<-ctx2.Done()
	this.cancelWithError(nil)

	// Wait for all Run functions to complete, then finish
	// collecting errors
	this.WaitGroup.Wait()
	close(this.errs)

	// Return the result
	return unwrap(result)
}

func (this *graph) cancelWithError(err error) {
	if err != nil && errors.Is(err, context.Canceled) == false {
		if logger := this.GetLogger(); logger != nil {
			logger.Debug("Cancelling with error: ", err)
		}
	}
	for _, cancel := range this.cancels {
		cancel()
	}
}

func unwrap(err error) error {
	if err != nil {
		if err_, ok := err.(*multierror.Error); ok {
			if len(err_.Errors) == 0 {
				return nil
			} else if len(err_.Errors) == 1 {
				return err_.Errors[0]
			}
		}
	}
	return err
}

func NewServiceStub(s string) gopi.ServiceStub {
	return Global.NewServiceStub(s)
}

func (this *graph) NewServiceStub(s string) gopi.ServiceStub {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if t, exists := stubs[s]; exists == false {
		return nil
	} else if stub := reflect.New(t.Elem()); stub.IsValid() == false {
		return nil
	} else {
		return stub.Interface().(gopi.ServiceStub)
	}
}

// GetLogger returns a logger object if used, or nil
func (this *graph) GetLogger() gopi.Logger {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	if t, exists := iface[loggerType]; exists == false {
		return nil
	} else if unit, exists := this.units[t]; exists == false {
		return nil
	} else if isLoggerType(t) == false {
		return nil
	} else {
		return unit.Interface().(gopi.Logger)
	}
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *graph) graph(unit reflect.Value) error {
	// Check incoming parameter
	if isUnitType(unit.Type()) == false {
		return gopi.ErrBadParameter.WithPrefix(unit.Type().String())
	}

	// For each field, initialise either by mapping an interface to
	// a registered unit type or directly
	return forEachField(unit, false, func(f reflect.StructField, i int) error {
		t := this.unitTypeForField(f)
		if t == nil {
			return nil
		}

		// Create a unit
		if _, exists := this.units[t]; exists == false {
			this.units[t] = reflect.New(t.Elem())
			if err := this.graph(this.units[t]); err != nil {
				return err
			}
		}

		// Set field to unit
		field := unit.Elem().Field(i)
		field.Set(this.units[t])

		// Return success
		return nil
	})
}

func (this *graph) isAppObject(unit reflect.Value) bool {
	for _, obj := range this.objs {
		if obj == unit {
			return true
		}
	}
	return false
}

func (this *graph) unitTypeForField(f reflect.StructField) reflect.Type {
	if f.Type.Kind() == reflect.Interface {
		if _, exists := iface[f.Type]; exists {
			return iface[f.Type]
		}
	} else if isUnitType(f.Type) {
		return f.Type
	}
	// Not found
	return nil
}

func (this *graph) do(fn string, unit reflect.Value, args []reflect.Value, seen map[reflect.Type]bool, indent int) error {
	// Check incoming parameter
	if isUnitType(unit.Type()) == false {
		return gopi.ErrBadParameter.WithPrefix(unit.Type().String())
	}

	var result error
	if fn == "Dispose" {
		if this.Logfn != nil {
			this.Logfn(strings.Repeat(" ", indent*2), fn, "=>", unit.Type())
		}
		if err := callFn(fn, unit, args); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// For each field, call function
	if err := forEachField(unit, fn == "New", func(f reflect.StructField, i int) error {
		if t := this.unitTypeForField(f); t == nil {
			return nil
		} else if _, exists := seen[t]; exists {
			return nil
		} else if err := this.do(fn, this.units[t], args, seen, indent+1); err != nil {
			seen[t] = true
			return fmt.Errorf("%w (in %v)", err, t)
		} else {
			seen[t] = true
		}
		return nil
	}); err != nil {
		result = multierror.Append(result, err)
	}

	if fn == "Run" {
		if this.Logfn != nil {
			this.Logfn(strings.Repeat(" ", indent*2), fn, "=>", unit.Type())
		}
		this.callRun(unit)
	} else if fn != "Dispose" {
		if this.Logfn != nil {
			this.Logfn(strings.Repeat(" ", indent*2), fn, "=>", unit.Type())
		}
		if err := callFn(fn, unit, args); err != nil {
			result = multierror.Append(result, err)
		}
	}

	return result
}

func (this *graph) callRun(unit reflect.Value) {
	this.WaitGroup.Add(1)
	ctx, cancel := context.WithCancel(context.Background())
	this.cancels = append(this.cancels, cancel)
	go func() {
		err := callFn("Run", unit, []reflect.Value{reflect.ValueOf(ctx)})
		if this.isAppObject(unit) {
			// Run ends when any application Run function ends
			this.cancelWithError(err)
		}
		if this.Logfn != nil {
			this.Logfn("Run", "<=", unit.Type())
		}
		this.errs <- err
		this.WaitGroup.Done()
		<-ctx.Done()
	}()
}

/////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *graph) String() string {
	str := "<graph"
	for k, v := range this.objs {
		str += fmt.Sprint(" ", k, "=>", v)
	}
	return str + ">"
}
