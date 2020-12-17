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
	name, usage, syntax string           // The name, usage and syntax information for the command
	args                []string         // The arguments for the command
	fn                  gopi.CommandFunc // The function called
	commands            []*command       // Subcommands
}

///////////////////////////////////////////////////////////////////////////////
// NEW

func NewCommand(name, usage, syntax string, args []string, fn gopi.CommandFunc) *command {
	return &command{name, usage, syntax, args, fn, make([]*command, 0)}
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *command) Name() string {
	return this.name
}

func (this *command) Usage() string {
	return this.usage
}

func (this *command) Syntax() string {
	return this.syntax
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
	if this.syntax != "" {
		str += " syntax=" + strconv.Quote(this.syntax)
	}
	if len(this.args) > 0 {
		str += " args=" + fmt.Sprint(this.args)
	}
	if len(this.commands) > 0 {
		str += " subcommands="
		for i, cmd := range this.commands {
			if i > 0 {
				str += ","
			}
			str += strconv.Quote(cmd.name)
		}
	}
	return str + ">"
}
