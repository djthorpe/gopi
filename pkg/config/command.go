package config

import (
	"context"
	"fmt"
	"strconv"

	"github.com/djthorpe/gopi/v3"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type command struct {
	name, usage string
	args        []string
	fn          gopi.CommandFunc
}

///////////////////////////////////////////////////////////////////////////////
// NEW

func NewCommand(name, usage string, args []string, fn gopi.CommandFunc) gopi.Command {
	return &command{
		name, usage, args, fn,
	}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *command) Name() string {
	return this.name
}

func (this *command) Usage() string {
	return this.usage
}

func (this *command) Args() []string {
	return this.args
}

func (this *command) Run(ctx context.Context) error {
	if this.fn == nil {
		return gopi.ErrNotImplemented
	} else {
		return this.fn(ctx)
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *command) String() string {
	str := "<command"
	str += " name=" + strconv.Quote(this.name)
	if this.usage != "" {
		str += " usage=" + strconv.Quote(this.usage)
	}
	if len(this.args) > 0 {
		str += " args=" + fmt.Sprint(this.args)
	}
	return str + ">"
}
