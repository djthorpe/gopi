/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bus

import (
	"context"
	"sync"
	"time"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Bus struct{}

type bus struct {
	handlers map[gopi.EventNS]map[string][]handlerWithTimeout
	defaults map[gopi.EventNS]gopi.EventHandler

	base.Unit
	sync.Mutex
	sync.WaitGroup
}

type handlerWithTimeout struct {
	fn      gopi.EventHandler
	timeout time.Duration
}

///////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Bus) Name() string { return "gopi.Bus" }

func (config Bus) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(bus)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	return this, nil
}

func (this *bus) Close() error {
	// Wait for handlers to complete
	this.WaitGroup.Wait()

	// Release resources
	this.handlers = nil
	this.defaults = nil

	// Return success
	return this.Unit.Close()
}

///////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Bus

func (this *bus) Emit(evt gopi.Event) {
	// Set NullEvent
	if evt == nil {
		evt = gopi.NullEvent
	}

	// TODO: hold cancel functions and call them with Close() to quickly
	// end handlers by cancelling them

	this.Log.Debug("Emit:", evt)

	// Set name and namespace
	name, ns := evt.Name(), evt.NS()
	if handlers := this.handlersForName(name, ns); len(handlers) > 0 {
		for _, handler := range handlers {
			this.WaitGroup.Add(1)
			go func(handler handlerWithTimeout) {
				ctx, cancel := handler.contextWithCancel()
				defer cancel()
				handler.fn(ctx, evt)
				this.WaitGroup.Done()
			}(handler)
		}
	} else if handler := this.defaultHandlerForNS(ns); handler != nil {
		this.WaitGroup.Add(1)
		go func(fn gopi.EventHandler) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			fn(ctx, evt)
			this.WaitGroup.Done()
		}(handler)
	} else {
		this.Log.Debug("Unhandled event:", evt)
	}
}

func (this *bus) NewHandler(name string, handler gopi.EventHandler) error {
	return this.NewHandlerEx(name, gopi.EVENT_NS_DEFAULT, handler, 0)
}

func (this *bus) NewHandlerEx(name string, ns gopi.EventNS, handler gopi.EventHandler, timeout time.Duration) error {
	// Check incoming parameters
	if name == "" {
		return gopi.ErrBadParameter.WithPrefix("name")
	} else if handler == nil {
		return gopi.ErrBadParameter.WithPrefix("handler")
	} else if timeout < 0 {
		return gopi.ErrBadParameter.WithPrefix("timeout")
	}
	// Set handler
	return this.addHandlerForName(name, ns, handler, timeout)
}

func (this *bus) DefaultHandler(ns gopi.EventNS, handler gopi.EventHandler) error {
	return this.setDefaultHandler(ns, handler)
}

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *bus) setDefaultHandler(ns gopi.EventNS, handler gopi.EventHandler) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Make defaults
	if this.defaults == nil {
		this.defaults = make(map[gopi.EventNS]gopi.EventHandler, 1)
	}
	// Set or delete handler for namespace
	if handler == nil {
		delete(this.defaults, ns)
	} else {
		this.defaults[ns] = handler
	}
	// Return success
	return nil
}

func (this *bus) addHandlerForName(name string, ns gopi.EventNS, handler gopi.EventHandler, timeout time.Duration) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Make handlers map
	if this.handlers == nil {
		this.handlers = make(map[gopi.EventNS]map[string][]handlerWithTimeout, 1)
	}

	// Make namespace map
	if _, exists := this.handlers[ns]; exists == false {
		this.handlers[ns] = make(map[string][]handlerWithTimeout, 1)
	}

	// Make name map
	if _, exists := this.handlers[ns][name]; exists == false {
		this.handlers[ns][name] = make([]handlerWithTimeout, 0, 1)
	}

	// Append a new handler
	this.handlers[ns][name] = append(this.handlers[ns][name], handlerWithTimeout{handler, timeout})

	// Return success
	return nil
}

func (this *bus) handlersForName(name string, ns gopi.EventNS) []handlerWithTimeout {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.handlers == nil {
		return nil
	} else if names, exists := this.handlers[ns]; exists == false {
		return nil
	} else if handlers, exists := names[name]; exists == false {
		return nil
	} else {
		return handlers
	}
}

func (this *bus) defaultHandlerForNS(ns gopi.EventNS) gopi.EventHandler {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if this.defaults == nil {
		return nil
	} else if handler, exists := this.defaults[ns]; exists == false {
		return nil
	} else {
		return handler
	}
}

func (this handlerWithTimeout) contextWithCancel() (context.Context, context.CancelFunc) {
	if this.timeout == 0 {
		return context.WithCancel(context.Background())
	} else {
		return context.WithTimeout(context.Background(), this.timeout)
	}
}
