/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package bus

import (
	// Frameworks

	"sync"

	gopi "github.com/djthorpe/gopi/v2"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type Bus struct{}

type bus struct {
	handlers map[string][]gopi.EventHandler

	gopi.UnitBase
	sync.Mutex
}

///////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (Bus) Name() string { return "gopi.Bus" }

func (config Bus) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(bus)
	if err := this.UnitBase.Init(log); err != nil {
		return nil, err
	} else {
		this.handlers = make(map[string][]gopi.EventHandler)
	}
	return this, nil
}

func (this *bus) Close() error {
	// Release resources
	this.handlers = nil

	// Return success
	return this.UnitBase.Close()
}

///////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Bus

func (this *bus) Emit(evt gopi.Event) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	if handlers, exists := this.handlers[evt.Name()]; exists {
		for _, handler := range handlers {
			go handler(evt)
		}
	} else {
		// Unhandled event
	}
}

func (this *bus) NewHandler(name string, handler gopi.EventHandler) error {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	// Check incoming parameters and allocate array for handlers
	if name == "" {
		return gopi.ErrBadParameter.WithPrefix("name")
	} else if handler == nil {
		return gopi.ErrBadParameter.WithPrefix("handler")
	} else if _, exists := this.handlers[name]; exists == false {
		this.handlers[name] = make([]gopi.EventHandler, 0, 1)
	}

	// Append handler
	if handlers, exists := this.handlers[name]; exists == false {
		return gopi.ErrInternalAppError
	} else {
		this.handlers[name] = append(handlers, handler)
	}

	// Success
	return nil
}
