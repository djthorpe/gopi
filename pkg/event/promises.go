package event

import (
	"context"
	"sync"

	"github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Promises struct {
	gopi.Unit
	sync.RWMutex
}

type promise struct {
	parent  context.Context
	chain   []PromiseFunc
	value   interface{}
	finally PromiseFinally
}

// A promise function executes functions sequentially and calls error
// function if an error occurs anywhere in the chain of functions
type PromiseFunc func(context.Context, interface{}) (interface{}, error)

// A Promise error is the last in the chain and acts on any error
type PromiseFinally func(interface{}, error) error

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Promises) String() string {
	str := "<promises"
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Promises) Do(ctx context.Context, fn func(context.Context, interface{}) (interface{}, error), v interface{}) gopi.Promise {
	if ctx == nil {
		ctx = context.Background()
	}
	if fn == nil {
		return nil
	}

	// Return a promise context
	return &promise{
		parent:  ctx,
		chain:   []PromiseFunc{fn},
		value:   v,
		finally: nil,
	}
}

func (this *promise) Then(fn func(context.Context, interface{}) (interface{}, error)) gopi.Promise {
	if fn == nil {
		return nil
	} else {
		this.chain = append(this.chain, fn)
	}
	return this
}

func (this *promise) Finally(fn func(interface{}, error) error, wait bool) error {
	var wg sync.WaitGroup
	var err error
	if fn != nil {
		this.finally = fn
	}
	wg.Add(1)
	go func() {
		// Run the chain of promises
		err = this.run()

		// Release resources
		this.parent = nil
		this.chain = nil
		this.value = nil
		this.finally = nil

		// Indicate done
		wg.Done()
	}()

	// If wait flag is set, then wait until done
	if wait {
		wg.Wait()
		return err
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Run the promise in foreground and call finally when completed
func (this *promise) run() error {
	for {
		select {
		case <-this.parent.Done():
			return this.finally(this.value, this.parent.Err())
		default:
			if len(this.chain) == 0 {
				return this.finally(this.value, nil)
			} else if out, err := this.chain[0](this.parent, this.value); err != nil {
				return this.finally(out, err)
			} else {
				this.chain = this.chain[1:]
				this.value = out
			}
		}
	}
}
