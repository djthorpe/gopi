package googlecast

import (
	"sync"
	"time"

	multierror "github.com/hashicorp/go-multierror"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type promise struct {
	expires time.Time
	fn      func() error
}

type promises struct {
	sync.Mutex
	timeout   time.Duration
	callbacks map[int]promise
}

////////////////////////////////////////////////////////////////////////////////
// LIFECYCLE

func (this *promises) InitWithTimeout(timeout time.Duration) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	this.timeout = timeout
	this.callbacks = make(map[int]promise)
}

func (this *promises) Dispose() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	this.callbacks = nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *promises) Set(id int, fn func() error) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Set promise with timeout
	this.callbacks[id] = promise{time.Now().Add(this.timeout), fn}

	// Expire old promises
	this.expire()
}

func (this *promises) Call(id int) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	var result error

	// Callback
	if promise, exists := this.callbacks[id]; exists {
		if promise.expires.After(time.Now()) {
			if err := promise.fn(); err != nil {
				result = multierror.Append(err)
			}
		}
		delete(this.callbacks, id)
	}

	// Expire old promises
	this.expire()

	return result
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *promises) expire() {
	for id, promise := range this.callbacks {
		if promise.expires.After(time.Now()) == false {
			delete(this.callbacks, id)
		}
	}
}
