// +build rpi

/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package gpiorpi

import (
	"fmt"
	"sync"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

type GPIO struct {
}

type gpio struct {
	log gopi.Logger

	base.Unit
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION gopi.Unit

func (GPIO) Name() string { return "gopi/gpio/rpi" }

func (config GPIO) New(log gopi.Logger) (gopi.Unit, error) {
	this := new(gpio)
	if err := this.Unit.Init(log); err != nil {
		return nil, err
	}
	if err := this.Init(config); err != nil {
		return nil, err
	}
	return this, nil
}

////////////////////////////////////////////////////////////////////////////////
// INIT & CLOSE

func (this *gpio) Init(config GPIO) error {
	this.Lock()
	defer this.Unlock()

	// Return success
	return nil
}

func (this *gpio) Close() error {
	this.Lock()
	defer this.Unlock()

	// Return success
	return this.Unit.Close()
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *gpio) String() string {
	if this.Unit.Closed {
		return fmt.Sprintf("<%v>", this.Log.Name())
	} else {
		return fmt.Sprintf("<%v>", this.Log.Name())
	}
}
