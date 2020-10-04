package config

import (
	"flag"
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

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY
/*
func (this *config) String() string {
	str := "<config"
	str += " name=" + strconv.Quote(this.FlagSet.Name())
	this.FlagSet.Visit(func(f *flag.Flag) {
		str += fmt.Sprintf(" %v=%q", f.Name, f.Value.String())
	})
	return str + ">"
}
*/
