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
	"testing"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type debug struct {
	main gopi.MainTestFunc
	t    *testing.T
	args []string
	base.App
}

////////////////////////////////////////////////////////////////////////////////
// gopi.App implementation for debug tool

func NewTestTool(t *testing.T, main gopi.MainTestFunc, args []string, units ...string) (gopi.App, error) {
	this := new(debug)

	// Name of test
	name := t.Name()

	// Check parameters
	if main == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("main")
	} else if err := this.App.Init(name, units); err != nil {
		return nil, err
	} else {
		this.main = main
		this.args = args
		this.t = t
	}

	// Success
	return this, nil
}

func (this *debug) Run() int {
	if returnValue := this.App.Start(this.args); returnValue != 0 {
		return returnValue
	}

	// Defer closing of instances to exit
	defer func() {
		if err := this.App.Close(); err != nil {
			fmt.Fprintln(os.Stderr, this.App.Flags().Name()+":", err)
		}
	}()

	// Run main function - doesn't return any errors since they are
	// handled by the testing package
	this.main(this, this.t)

	// Always return 0
	return 0
}
