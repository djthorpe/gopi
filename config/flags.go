/*
  Go Language Raspberry Pi Interface
  (c) Copyright David Thorpe 2016-2020
  All Rights Reserved
  For Licensing and Usage information, please see LICENSE.md
*/

package config

import (
	"flag"
	"fmt"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Flags struct {
	flags   *flag.FlagSet
	flagMap map[string]*flag.Flag
	name    string
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// NewFlags returns a new flags object
func NewFlags(name string) *Flags {
	return &Flags{
		flag.NewFlagSet(name, flag.ContinueOnError),
		make(map[string]*flag.Flag),
		name,
	}
}

// Parse command line arguments into flags and pure arguments
func (this *Flags) Parse(args []string) error {
	// parse flags
	if err := this.flags.Parse(args); err == flag.ErrHelp {
		return gopi.ErrHelp
	} else if err != nil {
		return err
	}

	// set hash of flags that were set
	this.flags.Visit(func(f *flag.Flag) {
		this.flagMap[f.Name] = f
	})

	// return success
	return nil
}

// Parsed reports whether the command-line flags have been parsed
func (this *Flags) Parsed() bool {
	return this.flags.Parsed()
}

// Name returns the name of the flagset (usually same as application)
func (this *Flags) Name() string {
	return this.name
}

// Args returns the command line arguments as an array which aren't flags
func (this *Flags) Args() []string {
	return this.flags.Args()
}

// Flags returns the command line arguments as an array which aren't flags
func (this *Flags) Flags() []*flag.Flag {
	flags := make([]*flag.Flag, 0, len(this.flagMap))
	for _, v := range this.flagMap {
		flags = append(flags, v)
	}
	return flags
}

// Has returns true if a flag exists
func (this *Flags) Has(name string) bool {
	if _, exists := this.flagMap[name]; exists {
		return true
	} else {
		return false
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *Flags) String() string {
	return fmt.Sprintf("<gopi.Flags>{ parsed=%v name=%v flags=%v args=%v }", this.Parsed(), this.Name(), this.Flags(), this.Args())
}

////////////////////////////////////////////////////////////////////////////////
// DEFINE FLAGS

// FlagString defines string flag and return pointer to the flag value
func (this *Flags) FlagString(name, value, usage string) *string {
	return this.flags.String(name, value, usage)
}

// FlagBool defines a boolean flag and return pointer to the flag value
func (this *Flags) FlagBool(name string, value bool, usage string) *bool {
	return this.flags.Bool(name, value, usage)
}

// FlagDuration defines duration flag and return pointer to the flag value
func (this *Flags) FlagDuration(name string, value time.Duration, usage string) *time.Duration {
	return this.flags.Duration(name, value, usage)
}

// FlagInt defines integer flag and return pointer to the flag value
func (this *Flags) FlagInt(name string, value int, usage string) *int {
	return this.flags.Int(name, value, usage)
}

// FlagUint defines unsigned integer flag and return pointer to the flag value
func (this *Flags) FlagUint(name string, value uint, usage string) *uint {
	return this.flags.Uint(name, value, usage)
}

// FlagFloat64 defines float64 flag and return pointer to the flag value
func (this *Flags) FlagFloat64(name string, value float64, usage string) *float64 {
	return this.flags.Float64(name, value, usage)
}

////////////////////////////////////////////////////////////////////////////////
// GET FLAGS

// GetBool gets boolean value for a flag, or default value if not set
func (this *Flags) GetBool(name string) bool {
	if value := this.flags.Lookup(name); value == nil {
		return false
	} else if getter, ok := value.Value.(flag.Getter); ok == false {
		return false
	} else if value_, ok := getter.Get().(bool); ok == false {
		return false
	} else {
		return value_
	}
}

// GetString gets string value for a flag, or default value if not set
func (this *Flags) GetString(name string) string {
	if value := this.flags.Lookup(name); value == nil {
		return ""
	} else if getter, ok := value.Value.(flag.Getter); ok == false {
		return ""
	} else {
		return getter.String()
	}
}

// GetDuration gets duration value for a flag
func (this *Flags) GetDuration(name string) time.Duration {
	if value := this.flags.Lookup(name); value == nil {
		return 0
	} else if getter, ok := value.Value.(flag.Getter); ok == false {
		return 0
	} else if value_, ok := getter.Get().(time.Duration); ok == false {
		return 0
	} else {
		return value_
	}
}

// GetInt gets integer value for a flag
func (this *Flags) GetInt(name string) int {
	if value := this.flags.Lookup(name); value == nil {
		return 0
	} else if getter, ok := value.Value.(flag.Getter); ok == false {
		return 0
	} else if value_, ok := getter.Get().(int); ok == false {
		return 0
	} else {
		return value_
	}
}

// GetUint gets unsigned integer value for a flag
func (this *Flags) GetUint(name string) uint {
	if value := this.flags.Lookup(name); value == nil {
		return 0
	} else if getter, ok := value.Value.(flag.Getter); ok == false {
		return 0
	} else if value_, ok := getter.Get().(uint); ok == false {
		return 0
	} else {
		return value_
	}
}

// GetFloat64 gets float64 value for a flag
func (this *Flags) GetFloat64(name string) float64 {
	if value := this.flags.Lookup(name); value == nil {
		return 0
	} else if getter, ok := value.Value.(flag.Getter); ok == false {
		return 0
	} else if value_, ok := getter.Get().(float64); ok == false {
		return 0
	} else {
		return value_
	}
}
