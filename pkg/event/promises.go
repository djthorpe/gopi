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
type PromiseFinally func(interface{}, error)

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

func (this *promise) Finally(fn func(interface{}, error), wait bool) {
	var wg sync.WaitGroup
	if fn != nil {
		this.finally = fn
	}
	wg.Add(1)
	go func() {
		// Run the chain of promises
		this.run()

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
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

// Run the promise in foreground and call finally when completed
func (this *promise) run() {
	for {
		select {
		case <-this.parent.Done():
			if this.finally != nil {
				this.finally(this.value, this.parent.Err())
			}
			return
		default:
			if len(this.chain) == 0 {
				if this.finally != nil {
					this.finally(this.value, nil)
				}
				return
			} else if out, err := this.chain[0](this.parent, this.value); err != nil {
				if this.finally != nil {
					this.finally(out, err)
				}
				return
			} else {
				this.chain = this.chain[1:]
				this.value = out
			}
		}
	}
}
