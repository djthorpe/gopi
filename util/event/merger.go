package event

import (
	"reflect"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi"
)

type Merger struct {
	sync.Mutex
	Publisher

	// the publishers we are watching
	publishers map[gopi.Publisher]<-chan gopi.Event

	// all the channels which are being merged
	in []<-chan gopi.Event

	// signal changes
	change chan struct{}
	done   chan struct{}
}

func (this *Merger) Merge(publisher gopi.Publisher) {
	this.Lock()
	defer this.Unlock()

	if this.in == nil {
		this.in = make([]<-chan gopi.Event, 0, 1)
	}
	if this.publishers == nil {
		this.publishers = make(map[gopi.Publisher]<-chan gopi.Event, 1)
	}
	// Sanity check publisher
	if publisher == nil {
		return
	}
	// If no change then start go routine
	if this.change == nil {
		this.change = make(chan struct{})
		this.done = make(chan struct{})
		go this.mergeInBackground()
	}
	// We cannot merge twice, so ignore
	if _, exists := this.publishers[publisher]; exists {
		return
	}
	// Perform subscriptiom
	channel := publisher.Subscribe()
	this.publishers[publisher] = channel
	this.in = append(this.in, channel)
}

func (this *Merger) Unmerge(publisher gopi.Publisher) {
	this.Lock()
	defer this.Unlock()

	// Sanity check
	if this.in == nil || this.publishers == nil || publisher == nil {
		return
	}
	if channel, exists := this.publishers[publisher]; exists == false {
		return
	} else {
		publisher.Unsubscribe(channel)
		delete(this.publishers, publisher)
	}
}

func (this *Merger) Close() {
	// Close publisher
	this.Publisher.Close()
	for publisher := range this.publishers {
		this.Unmerge(publisher)
	}
	// Close change channel
	if this.change != nil {
		close(this.change)
		// Wait for done signal
		<-this.done
	}
	// Empty data structures
	this.publishers = nil
	this.in = nil

}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *Merger) cases() []reflect.SelectCase {
	cases := make([]reflect.SelectCase, 1, len(this.in)+1)
	// Add the change channel
	cases[0] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(this.change),
	}
	// Add the remaining channels - ignoring nil channels
	// which have been closed
	if this.in != nil {
		for i := range this.in {
			if this.in[i] != nil {
				cases = append(cases, reflect.SelectCase{
					Dir:  reflect.SelectRecv,
					Chan: reflect.ValueOf(this.in[i]),
				})
			}
		}
	}
	// return all cases
	return cases
}

func (this *Merger) mergeInBackground() {
	// Continue loop until chanhe channel is closed
	cases := this.cases()
FOR_LOOP:
	for {
		// Deal with zero cases condition
		if len(cases) == 0 {
			break FOR_LOOP
		}
		// select cases
		i, v, ok := reflect.Select(cases)
		if i == 0 && ok == false {
			// We need to reload the cases. If zero then end
			break FOR_LOOP
		} else if i == 0 {
			// Reload cases
			cases = this.cases()
		} else if ok {
			this.Emit(v.Interface().(gopi.Event))
		} else if i > 0 {
			// Set channel to nil to remove from cases
			this.in[i-1] = nil
			// Rebuild cases
			cases = this.cases()
		}
	}
	// Indicate the background thread is done
	this.done <- gopi.DONE
}
