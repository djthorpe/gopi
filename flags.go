/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2017
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package gopi

import (
	"flag"
	"fmt"
	"time"
)

type Flags struct {
	flagset *flag.FlagSet
	flagmap map[string]bool
	name    string
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Create a new flags object
func NewFlags(name string) *Flags {
	this := new(Flags)
	this.flagset = flag.NewFlagSet(name, flag.ContinueOnError)
	this.flagmap = nil
	this.name = name
	return this
}

// Parse command line argumentsinto flags and pure arguments
func (this *Flags) Parse(args []string) error {

	// parse flags
	err := this.flagset.Parse(args)
	if err != nil {
		return err
	}

	// set hash of flags that were set
	this.flagmap = make(map[string]bool)
	this.flagset.Visit(func(f *flag.Flag) {
		this.flagmap[f.Name] = true
	})

	// return success
	return nil
}

// Parsed reports whether the command-line flags have been parsed
func (this *Flags) Parsed() bool {
	return this.flagset.Parsed()
}

// Name returns the name of the flagset (usually same as application)
func (this *Flags) Name() string {
	return this.name
}

// Args returns the command line arguments as an array which aren't flags
func (this *Flags) Args() []string {
	return this.flagset.Args()
}

// Flags returns the array of flags which were set on the command line
func (this *Flags) Flags() []string {
	if this.flagmap == nil {
		return []string{}
	}
	flags := make([]string, 0)
	for k := range this.flagmap {
		flags = append(flags, k)
	}
	return flags
}

// HasFlag returns a boolean indicating if a flag was set on the command line
func (this *Flags) HasFlag(name string) bool {
	if this.flagmap == nil {
		return false
	}
	_, exists := this.flagmap[name]
	return exists
}

// SetUsageFunc sets the usage function which prints
// usage information to stderr
func (this *Flags) SetUsageFunc(usage_func func(flags *Flags)) {
	this.flagset.Usage = func() {
		usage_func(this)
	}
}

// PrintUsage will call the usage function
func (this *Flags) PrintUsage() {
	this.flagset.Usage()
}

// PrintDefaults will output the flags to stderr
func (this *Flags) PrintDefaults() {
	this.flagset.PrintDefaults()
}

// String returns a human-readable form of the Flags object
func (this *Flags) String() string {
	return fmt.Sprintf("<app.Flags>{ parsed=%v name=%v flags=%v args=%v }", this.Parsed(), this.Name(), this.Flags(), this.Args())
}

////////////////////////////////////////////////////////////////////////////////
// DEFINE FLAGS

// FlagString defines string flag and return pointer to the flag value
func (this *Flags) FlagString(name string, value string, usage string) *string {
	if this.flagset == nil {
		return nil
	} else {
		return this.flagset.String(name, value, usage)
	}
}

// FlagBool defines a boolean flag and return pointer to the flag value
func (this *Flags) FlagBool(name string, value bool, usage string) *bool {
	if this.flagset == nil {
		return nil
	} else {
		return this.flagset.Bool(name, value, usage)
	}
}

// FlagDuration defines duration flag and return pointer to the flag value
func (this *Flags) FlagDuration(name string, value time.Duration, usage string) *time.Duration {
	if this.flagset == nil {
		return nil
	} else {
		return this.flagset.Duration(name, value, usage)
	}
}

// FlagInt defines integer flag and return pointer to the flag value
func (this *Flags) FlagInt(name string, value int, usage string) *int {
	if this.flagset == nil {
		return nil
	} else {
		return this.flagset.Int(name, value, usage)
	}
}

// FlagUint defines unsigned integer flag and return pointer to the flag value
func (this *Flags) FlagUint(name string, value uint, usage string) *uint {
	if this.flagset == nil {
		return nil
	} else {
		return this.flagset.Uint(name, value, usage)
	}
}

// FlagFloat64 defines float64 flag and return pointer to the flag value
func (this *Flags) FlagFloat64(name string, value float64, usage string) *float64 {
	if this.flagset == nil {
		return nil
	} else {
		return this.flagset.Float64(name, value, usage)
	}
}

////////////////////////////////////////////////////////////////////////////////
// GET FLAGS

// GetBool gets boolean value for a flag, and a boolean which indicates if the flag
// was set
func (this *Flags) GetBool(name string) (bool, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return false, false
	}
	return value.Value.(flag.Getter).Get().(bool), this.HasFlag(name)
}

// GetString gets string value for a flag, and a boolean which indicates if the flag
// was set
func (this *Flags) GetString(name string) (string, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return "", false
	}
	return value.Value.(flag.Getter).Get().(string), this.HasFlag(name)
}

// GetDuration gets duration value for a flag, and a boolean which indicates if the flag
// was set
func (this *Flags) GetDuration(name string) (time.Duration, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return time.Duration(0), false
	}
	return value.Value.(flag.Getter).Get().(time.Duration), this.HasFlag(name)
}

// GetInt gets integer value for a flag, and a boolean which indicates if the flag
// was set
func (this *Flags) GetInt(name string) (int, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0, false
	}
	return value.Value.(flag.Getter).Get().(int), this.HasFlag(name)
}

// GetUint gets unsigned integer value for a flag, and a boolean which indicates if
// the flag was set
func (this *Flags) GetUint(name string) (uint, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0, false
	}
	return value.Value.(flag.Getter).Get().(uint), this.HasFlag(name)
}

// GetUint16 gets unsigned integer value for a flag, and a boolean which indicates if
// the flag was set
func (this *Flags) GetUint16(name string) (uint16, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0, false
	}
	uint_value := value.Value.(flag.Getter).Get().(uint)
	return uint16(uint_value), this.HasFlag(name)
}

// GetFloat64 gets float64 value for a flag, and a boolean which indicates if
// the flag was set
func (this *Flags) GetFloat64(name string) (float64, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0.0, false
	}
	return value.Value.(flag.Getter).Get().(float64), this.HasFlag(name)
}

////////////////////////////////////////////////////////////////////////////////
// SET FLAGS

// Set a flag string value. The flag must have previously been configured
// using FlagXX method. Will return an error if the value couldn't be parsed
func (this *Flags) SetString(name, value string) error {
	if flag := this.flagset.Lookup(name); flag == nil {
		return fmt.Errorf("SetString: No such flag: %v", name)
	} else {
		return flag.Value.Set(value)
	}
}

// Set a flag uint value
func (this *Flags) SetUint(name string, value uint) error {
	return this.SetString(name, fmt.Sprint(value))
}

// Set a flag int value
func (this *Flags) SetInt(name string, value int) error {
	return this.SetString(name, fmt.Sprint(value))
}

// Set a flag bool value
func (this *Flags) SetBool(name string, value bool) error {
	return this.SetString(name, fmt.Sprint(value))
}

// Set a flag float64 value
func (this *Flags) SetFloat64(name string, value float64) error {
	return this.SetString(name, fmt.Sprint(value))
}

// Set a flag duration value
func (this *Flags) SetDuration(name string, value time.Duration) error {
	return this.SetString(name, fmt.Sprint(value))
}
