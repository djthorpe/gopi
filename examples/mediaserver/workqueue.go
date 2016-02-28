package main

import (
	"runtime"
	"sync"
)

type WorkHandler func(interface{})

type WorkQueue struct {
	Handler   WorkHandler
	Workers   int
	push      chan interface{}
	pop       chan struct{}
	suspend   chan bool
	suspended bool
	stop      chan struct{}
	stopped   bool
	buffer    []interface{}
	count     int
	wg        sync.WaitGroup
}

func NewWorkQueue(workers int, handler WorkHandler) *WorkQueue {
	q := new(WorkQueue)
	q.Handler = handler
	q.Workers = workers
	q.push = make(chan interface{})
	q.pop = make(chan struct{})
	q.suspend = make(chan bool)
	q.stop = make(chan struct{})

	go q.run()
	runtime.SetFinalizer(q, (*WorkQueue).Stop)
	return q
}

func (this *WorkQueue) Push(val interface{}) {
	this.push <- val
}

func (this *WorkQueue) Stop() {
	this.stop <- struct{}{}
	runtime.SetFinalizer(this, nil)
}

func (this *WorkQueue) Wait() {
	this.wg.Wait()
}

func (this *WorkQueue) Len() (_, _ int) {
	return this.count, len(this.buffer)
}

func (this *WorkQueue) run() {
	defer func() {
		this.wg.Add(-len(this.buffer))
		this.buffer = nil
	}()

	for {
		select {
		case val := <-this.push:
			this.buffer = append(this.buffer, val)
			this.wg.Add(1)
		case <-this.pop:
			this.count--
		case suspend := <-this.suspend:
			if suspend != this.suspended {
				if suspend {
					this.wg.Add(1)
				} else {
					this.wg.Done()
				}
				this.suspended = suspend
			}
		case <-this.stop:
			this.stopped = true
		}

		if this.stopped && this.count == 0 {
			return
		}

		for (this.count < this.Workers || this.Workers == 0) && len(this.buffer) > 0 && !(this.suspended || this.stopped) {
			val := this.buffer[0]
			this.buffer = this.buffer[1:]
			this.count++
			go func() {
				defer func() {
					this.pop <- struct{}{}
					this.wg.Done()
				}()
				this.Handler(val)
			}()
		}

	}

}
