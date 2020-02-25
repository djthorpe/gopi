/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package app

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	// Frameworks
	gopi "github.com/djthorpe/gopi/v2"
	base "github.com/djthorpe/gopi/v2/base"
)

////////////////////////////////////////////////////////////////////////////////
// INTERFACES

type command struct {
	main     gopi.MainCommandFunc
	handlers []gopi.EventHandler
	base.App
}

////////////////////////////////////////////////////////////////////////////////
// gopi.App implementation for command-line tool

func NewCommandLineTool(main gopi.MainCommandFunc, handlers []gopi.EventHandler, units ...string) (gopi.App, error) {
	this := new(command)

	// Name of command
	name := filepath.Base(os.Args[0])

	// If there are any handlers, then append "bus" onto required units
	if len(handlers) > 0 {
		units = append(units, "bus")
	}

	// Check parameters
	if main == nil {
		return nil, gopi.ErrBadParameter.WithPrefix("main")
	} else if err := this.App.Init(name, units); err != nil {
		return nil, err
	} else {
		this.main = main
		this.handlers = handlers
	}

	// Success
	return this, nil
}

func (this *command) Run() int {
	if err := this.App.Start(this, os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) == false {
			fmt.Fprintln(os.Stderr, this.App.Flags().Name()+":", err)
			return -1
		} else {
			return 0
		}
	}

	// Defer closing of instances to exit
	defer func() {
		name := this.App.Flags().Name()
		if err := this.App.Close(); err != nil {
			fmt.Fprintln(os.Stderr, name+": Exit Error:", err)
		}
	}()

	// Set up handlers
	if len(this.handlers) > 0 {
		for _, handler := range this.handlers {
			if handler.Name != "" && handler.Handler != nil {
				this.Log().Debug(this.App.Flags().Name()+":", "Set up event handler for", strconv.Quote(handler.Name), "in namespace", handler.EventNS)
				if err := this.Bus().NewHandler(handler); err != nil {
					fmt.Fprintln(os.Stderr, this.App.Flags().Name()+":", err)
					return -1
				}
			}
		}
	}

	// Run main function
	if err := this.main(this, this.Flags().Args()); errors.Is(err, gopi.ErrHelp) || errors.Is(err, flag.ErrHelp) {
		this.App.Flags().Usage(os.Stderr)
		return 0
	} else if err != nil {
		fmt.Fprintln(os.Stderr, this.App.Flags().Name()+":", err)
		return -1
	}

	// Success
	return 0
}
