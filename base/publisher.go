/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package base

import (
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Publisher struct {
	queues map[uint]*queue
	token  gopi.Channel

	sync.Mutex
}

type queue struct {
	cb map[gopi.Channel]func(value interface{})

	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *Publisher) Init(num uint) *queue {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.queues == nil {
		this.queues = make(map[uint]*queue, 1)
	}
	if q, exists := this.queues[num]; exists {
		return q
	} else if q = NewQueue(); q != nil {
		this.queues[num] = q
		return q
	} else {
		return nil
	}
}

func (this *Publisher) Close() {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.queues != nil {
		for _, queue := range this.queues {
			queue.Close()
		}
		this.queues = nil
	}
}

func (this *Publisher) Emit(num uint, value interface{}) {
	if q, exists := this.queues[num]; exists {
		q.Emit(value)
	}
}

func (this *Publisher) Subscribe(num uint, callback func(value interface{})) gopi.Channel {
	if q := this.Init(num); q != nil && callback != nil {
		return q.Add(callback, this.NextToken())
	} else {
		return 0
	}
}

func (this *Publisher) Unsubscribe(token gopi.Channel) {
	for _, q := range this.queues {
		q.Remove(token)
	}
}

func (this *Publisher) NextToken() gopi.Channel {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()
	this.token = this.token + gopi.Channel(1)
	return this.token
}

////////////////////////////////////////////////////////////////////////////////
// QUEUES

func NewQueue() *queue {
	this := &queue{
		cb: make(map[gopi.Channel]func(value interface{})),
	}
	return this
}

func (this *queue) Add(callback func(value interface{}), token gopi.Channel) gopi.Channel {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if callback != nil {
		if _, exists := this.cb[token]; exists == false {
			this.cb[token] = callback
			return token
		}
	}
	return 0
}

func (this *queue) Remove(token gopi.Channel) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if _, exists := this.cb[token]; exists {
		delete(this.cb, token)
	}
}

func (this *queue) Emit(value interface{}) {
	for _, cb := range this.cb {
		if cb != nil {
			cb(value)
		}
	}
}

func (this *queue) Close() {
	this.cb = nil
}
