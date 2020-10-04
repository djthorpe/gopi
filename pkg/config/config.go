package config

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/djthorpe/gopi/v3"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type config struct {
	*flag.FlagSet
	args []string
}

type command struct {
	name string
}

///////////////////////////////////////////////////////////////////////////////
// NEW

func New(name string, args []string) gopi.Config {
	this := new(config)
	this.FlagSet = flag.NewFlagSet(name, flag.ContinueOnError)
	this.args = args
	return this
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *config) Parse() error {
	// Check for already parsed
	if this.FlagSet.Parsed() {
		return nil
	}
	// Perform parse
	if err := this.FlagSet.Parse(this.args); err != nil {
		return err
	}
	// Return success
	return nil
}

///////////////////////////////////////////////////////////////////////////////
// DEFINE FLAGS

func (this *config) FlagString(name, value, usage string) *string {
	return this.FlagSet.String(name, value, usage)
}

func (this *config) FlagBool(name string, value bool, usage string) *bool {
	return this.FlagSet.Bool(name, value, usage)
}

func (this *config) FlagUint(name string, value uint, usage string) *uint {
	return this.FlagSet.Uint(name, value, usage)
}

///////////////////////////////////////////////////////////////////////////////
// GET PROPERTIES

func (this *config) GetString(name string) string {
	if flag := this.FlagSet.Lookup(name); flag == nil {
		return ""
	} else {
		return flag.Value.String()
	}
}

func (this *config) GetBool(name string) bool {
	if flag := this.FlagSet.Lookup(name); flag == nil {
		return false
	} else if value_, err := strconv.ParseBool(flag.Value.String()); err != nil {
		return false
	} else {
		return value_
	}
}

func (this *config) GetUint(name string) uint {
	if flag := this.FlagSet.Lookup(name); flag == nil {
		return 0
	} else if value_, err := strconv.ParseUint(flag.Value.String(), 10, 32); err != nil {
		return 0
	} else {
		return uint(value_)
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *config) String() string {
	str := "<config"
	str += " name=" + strconv.Quote(this.FlagSet.Name())
	this.FlagSet.Visit(func(f *flag.Flag) {
		str += fmt.Sprintf(" %v=%q", f.Name, f.Value.String())
	})
	return str + ">"
}
