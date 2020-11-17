package graph

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	gopi "github.com/djthorpe/gopi/v3"
	multierror "github.com/hashicorp/go-multierror"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type graph struct {
	sync.RWMutex
	sync.WaitGroup

	units map[reflect.Type]reflect.Value
	objs  []reflect.Value
	logfn func(...interface{})
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
	this.logfn = fn
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
	return result
}

// Call Define for each unit object
func (this *graph) Define(cfg gopi.Config) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var result error
	seen := make(map[reflect.Type]bool, len(this.units))
	for _, obj := range this.objs {
		if err := this.do("Define", obj, []reflect.Value{reflect.ValueOf(cfg)}, seen); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

// Call New for each unit object
func (this *graph) New(cfg gopi.Config) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var result error
	seen := make(map[reflect.Type]bool, len(this.units))
	for _, obj := range this.objs {
		if err := this.do("New", obj, []reflect.Value{reflect.ValueOf(cfg)}, seen); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

// Call Dispose for each unit object
func (this *graph) Dispose() error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	var result error
	seen := make(map[reflect.Type]bool, len(this.units))
	for _, obj := range this.objs {
		if err := this.do("Dispose", obj, []reflect.Value{}, seen); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result
}

// Call Run for each unit object
func (this *graph) Run(ctx context.Context) error {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()

	seen := make(map[reflect.Type]bool, len(this.units))
	cancels := []context.CancelFunc{}
	errs := make(chan error)

	// Collect errors
	var result error
	go func() {
		for err := range errs {
			if err != nil && errors.Is(err, context.Canceled) == false {
				result = multierror.Append(result, err)
			}
			this.WaitGroup.Done()
		}
	}()

	// Send cancels on context end
	go func() {
		// Wait until the context is done
		<-ctx.Done()

		// Call cancels
		for _, cancel := range cancels {
			cancel()
		}
	}()

	// Call run functions
	for _, obj := range this.objs {
		cancels = append(cancels, this.run(obj, errs, seen)...)
	}

	// Wait for Run() functions to complete
	this.WaitGroup.Wait()

	// Close err channel
	close(errs)

	// Return the context cancel reason
	if result == nil {
		return ctx.Err()
	} else {
		return multierror.Append(result, ctx.Err())
	}
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

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *graph) graph(unit reflect.Value) error {
	// Check incoming parameter
	if isUnitType(unit.Type()) == false {
		return gopi.ErrBadParameter.WithPrefix(unit.Type().String())
	}

	// For each field, initialise either by mapping an interface to
	// a registered unit type or directly
	return forEachField(unit, func(f reflect.StructField, i int) error {
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

func (this *graph) do(fn string, unit reflect.Value, args []reflect.Value, seen map[reflect.Type]bool) error {
	// Check incoming parameter
	if isUnitType(unit.Type()) == false {
		return gopi.ErrBadParameter.WithPrefix(unit.Type().String())
	}

	if this.logfn != nil {
		this.logfn(fn, "=>", unit.Type())
	}

	// For each field, call function
	if err := forEachField(unit, func(f reflect.StructField, i int) error {
		if t := this.unitTypeForField(f); t == nil {
			return nil
		} else if err := this.do(fn, this.units[t], args, seen); err != nil {
			return err
		} else {
			seen[t] = true
		}

		// Return success
		return nil
	}); err != nil {
		return err
	}

	// Call the function and return the error
	return callFn(fn, unit, args)
}

func (this *graph) run(unit reflect.Value, errs chan<- error, seen map[reflect.Type]bool) []context.CancelFunc {
	cancels := []context.CancelFunc{}

	if this.logfn != nil {
		this.logfn("Run", "=>", unit.Type())
	}

	// Recurse into run
	forEachField(unit, func(f reflect.StructField, i int) error {
		if _, exists := seen[f.Type]; exists {
			return nil
		}
		if isUnitType(f.Type) == false {
			return nil
		}
		seen[f.Type] = true
		cancels = append(cancels, this.run(this.units[f.Type], errs, seen)...)
		return nil
	})

	// Now call Run in a goroutine, which passes error back to channel
	ctx, cancel := context.WithCancel(context.Background())
	this.WaitGroup.Add(1)
	go func() {
		errs <- callFn("Run", unit, []reflect.Value{reflect.ValueOf(ctx)})
	}()

	return append(cancels, cancel)
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
