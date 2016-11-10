/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/

package app /* import "github.com/djthorpe/gopi/app" */

// import
import (
	"flag"
	"time"
)

type Flags struct {
	flagset *flag.FlagSet
}

var (
	ErrHelp = flag.ErrHelp
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func NewFlags(name string) *Flags {
	this := new(Flags)
	this.flagset = flag.NewFlagSet(name, flag.ContinueOnError)
	return this
}

func (this *Flags) Parse(args []string) error {
	return this.flagset.Parse(args)
}

func (this *Flags) Args() []string {
	return this.flagset.Args()
}

////////////////////////////////////////////////////////////////////////////////
// SET FLAGS

func (this *Flags) String(name string, value string, usage string) *string {
	return this.flagset.String(name,value,usage)
}

func (this *Flags) Bool(name string, value bool, usage string) *bool {
	return this.flagset.Bool(name,value,usage)
}

func (this *Flags) Duration(name string, value time.Duration, usage string) *time.Duration {
	return this.flagset.Duration(name,value,usage)
}

func (this *Flags) Int(name string, value int, usage string) *int {
	return this.flagset.Int(name,value,usage)
}

func (this *Flags) Uint(name string, value uint, usage string) *uint {
	return this.flagset.Uint(name,value,usage)
}

////////////////////////////////////////////////////////////////////////////////
// GET FLAGS

func (this *Flags) GetBool(name string) (bool, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return false, false
	}
	return value.Value.(flag.Getter).Get().(bool), true
}

func (this *Flags) GetString(name string) (string, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return "", false
	}
	return value.Value.(flag.Getter).Get().(string), true
}

func (this *Flags) GetDuration(name string) (time.Duration, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return time.Duration(0), false
	}
	return value.Value.(flag.Getter).Get().(time.Duration), true
}

func (this *Flags) GetInt(name string) (int, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0, false
	}
	return value.Value.(flag.Getter).Get().(int), true
}

func (this *Flags) GetUint(name string) (uint, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0, false
	}
	return value.Value.(flag.Getter).Get().(uint), true
}





