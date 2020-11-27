package graph

import (
	"fmt"
	"sync"
)

type counter struct {
	sync.RWMutex
	c    chan struct{}
	v    int
	done bool
}

func (this *counter) Add(i int) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if this.c == nil {
		this.c = make(chan struct{})
	}
	this.v += i
	if this.v == 0 && this.done {
		this.c <- struct{}{}
	}
}

func (this *counter) Sub(i int) {
	this.RWMutex.Lock()
	defer this.RWMutex.Unlock()
	if this.c == nil {
		this.c = make(chan struct{})
	}
	this.v -= i
	if this.v == 0 && this.done {
		this.c <- struct{}{}
	}
}

func (this *counter) Done() <-chan struct{} {
	// Emits when counter is zero
	return this.c
}

func (this *counter) String() string {
	this.RWMutex.RLock()
	defer this.RWMutex.RUnlock()
	return fmt.Sprint(this.v)
}
