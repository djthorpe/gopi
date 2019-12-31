/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app

import (
	// Frameworks
	"errors"
	"flag"
	"fmt"
	"os"

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

	// Name of command
	name := os.Args[0]

	// Check parameters
	if main == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("gopi.MainCommandFunc")
	} else if err := this.base.Init(name, units); err != nil {
		return nil, err
	} else {
		this.main = main
	}

	// Success
	return this, nil
}

func (this *command) Run() int {
	if returnValue := this.base.Run(); returnValue != 0 {
		return returnValue
	} else if err := this.main(this, this.Flags().Args()); errors.Is(err, gopi.ErrHelp) || errors.Is(err, flag.ErrHelp) {
		this.flags.Usage(os.Stderr)
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, this.flags.Name()+":", err)
		return -1
	}

	// Success
	return 0
}
