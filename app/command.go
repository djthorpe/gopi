/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app

import (
	"fmt"
	"os"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type command struct {
	base
	main gopi.MainCommandFunc
}

////////////////////////////////////////////////////////////////////////////////
// gopi.App implementation for command-line tool

func NewCommandLineTool(main gopi.MainCommandFunc, units ...string) (gopi.App, error) {
	this := new(command)

	// Check parameters
	if main == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("gopi.MainCommandFunc")
	} else if err := this.base.Init(units); err != nil {
		return nil, err
	} else {
		this.main = main
	}

	// Success
	return this, nil
}

func (this *command) Run() int {

	fmt.Fprintf(os.Stderr, "Run() is not imple,ented")
	return -1
}
