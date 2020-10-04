package tool

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
	sync.WaitGroup
	objs  []reflect.Value
	units map[reflect.Type]reflect.Value
}

/////////////////////////////////////////////////////////////////////
// GLOBALS

var (
	unitType = reflect.TypeOf((*gopi.Unit)(nil)).Elem()
)

/////////////////////////////////////////////////////////////////////
// CONSTRUCTOR

// Construct unit objects which are shared
func NewGraph(objs ...interface{}) (*graph, error) {
	this := new(graph)
	this.units = make(map[reflect.Type]reflect.Value)

	// Iterate through the objects, creating units
	var result error
	for _, obj := range objs {
		this.objs = append(this.objs, reflect.ValueOf(obj))
		if err := this.graph(reflect.ValueOf(obj)); err != nil {
			result = multierror.Append(result, err)
		}
	}

	// Return success
	return this, result
}

// Call Define for each unit object
func (this *graph) Define(cfg gopi.Config) error {
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
func (this *graph) Dispose(cfg gopi.Config) error {
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
	seen := make(map[reflect.Type]bool, len(this.units))
	cancels := []context.CancelFunc{}
	errs := make(chan error)

	// Collect errors
	go func() {
		for err := range errs {
			if err != nil && errors.Is(err, context.Canceled) == false {
				fmt.Println("Err=", err)
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
	return ctx.Err()
}

/////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *graph) graph(unit reflect.Value) error {
	// Check incoming parameter
	if isUnitType(unit.Type()) == false {
		return gopi.ErrBadParameter.WithPrefix(unit.Type().String())
	}

	// For each field, initialise
	return forEachField(unit, func(f reflect.StructField, i int) error {
		if isUnitType(f.Type) == false {
			return nil
		}
		// Create a Unit
		if _, exists := this.units[f.Type]; exists == false {
			this.units[f.Type] = reflect.New(f.Type.Elem())
			if err := this.graph(this.units[f.Type]); err != nil {
				return err
			}
		}
		// Set field to unit
		field := unit.Elem().Field(i)
		field.Set(this.units[f.Type])

		// Return success
		return nil
	})
}

func (this *graph) do(fn string, unit reflect.Value, args []reflect.Value, seen map[reflect.Type]bool) error {
	// Check incoming parameter
	if isUnitType(unit.Type()) == false {
		return gopi.ErrBadParameter.WithPrefix(unit.Type().String())
	}

	// For each field, call function
	if err := forEachField(unit, func(f reflect.StructField, i int) error {
		if _, exists := seen[f.Type]; exists {
			return nil
		}
		if isUnitType(f.Type) == false {
			return nil
		}
		if err := this.do(fn, this.units[f.Type], args, seen); err != nil {
			return err
		} else {
			seen[f.Type] = true
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
