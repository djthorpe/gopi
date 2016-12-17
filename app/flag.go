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
	"fmt"
	"time"
)

type Flags struct {
	flagset *flag.FlagSet
	flagmap map[string]bool
}

var (
	ErrHelp = flag.ErrHelp
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Create a new flags object
func NewFlags(name string) *Flags {
	this := new(Flags)
	this.flagset = flag.NewFlagSet(name, flag.ContinueOnError)
	this.flagmap = nil
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

// Return the command line arguments as an array which aren't flags
func (this *Flags) Args() []string {
	return this.flagset.Args()
}

// Return array of flags which were set on the command line
func (this *Flags) Flags() []string {
	if this.flagmap == nil {
		return []string{}
	}
	flags := make([]string, 0)
	for k, _ := range this.flagmap {
		flags = append(flags, k)
	}
	return flags
}

// Return boolean indicating if a flag was set on the command line
func (this *Flags) HasFlag(name string) bool {
	if this.flagmap == nil {
		return false
	}
	_, exists := this.flagmap[name]
	return exists
}

// Return human-readable form of the Flags object
func (this *Flags) String() string {
	return fmt.Sprintf("<app.Flags>{ flags=%v args=%v }", this.Flags(), this.Args())
}

////////////////////////////////////////////////////////////////////////////////
// SET FLAGS

// Define string flag and return pointer to the flag value
func (this *Flags) FlagString(name string, value string, usage string) *string {
	return this.flagset.String(name, value, usage)
}

// Define boolean flag and return pointer to the flag value
func (this *Flags) FlagBool(name string, value bool, usage string) *bool {
	return this.flagset.Bool(name, value, usage)
}

// Define duration flag and return pointer to the flag value
func (this *Flags) FlagDuration(name string, value time.Duration, usage string) *time.Duration {
	return this.flagset.Duration(name, value, usage)
}

// Define integer flag and return pointer to the flag value
func (this *Flags) FlagInt(name string, value int, usage string) *int {
	return this.flagset.Int(name, value, usage)
}

// Define unsigned integer flag and return pointer to the flag value
func (this *Flags) FlagUint(name string, value uint, usage string) *uint {
	return this.flagset.Uint(name, value, usage)
}

// Define float64 flag and return pointer to the flag value
func (this *Flags) FlagFloat64(name string, value float64, usage string) *float64 {
	return this.flagset.Float64(name, value, usage)
}

////////////////////////////////////////////////////////////////////////////////
// GET FLAGS

// Get boolean value for a flag, and a boolean which indicates if the flag
// was set
func (this *Flags) GetBool(name string) (bool, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return false, false
	}
	return value.Value.(flag.Getter).Get().(bool), this.HasFlag(name)
}

// Get string value for a flag, and a boolean which indicates if the flag
// was set
func (this *Flags) GetString(name string) (string, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return "", false
	}
	return value.Value.(flag.Getter).Get().(string), this.HasFlag(name)
}

// Get duration value for a flag, and a boolean which indicates if the flag
// was set
func (this *Flags) GetDuration(name string) (time.Duration, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return time.Duration(0), false
	}
	return value.Value.(flag.Getter).Get().(time.Duration), this.HasFlag(name)
}

// Get integer value for a flag, and a boolean which indicates if the flag
// was set
func (this *Flags) GetInt(name string) (int, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0, false
	}
	return value.Value.(flag.Getter).Get().(int), this.HasFlag(name)
}

// Get unsigned integer value for a flag, and a boolean which indicates if
// the flag was set
func (this *Flags) GetUint(name string) (uint, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0, false
	}
	return value.Value.(flag.Getter).Get().(uint), this.HasFlag(name)
}

// Get unsigned integer value for a flag, and a boolean which indicates if
// the flag was set
func (this *Flags) GetUint16(name string) (uint16, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0, false
	}
	uint_value := value.Value.(flag.Getter).Get().(uint)
	return uint16(uint_value), this.HasFlag(name)
}

// Get float64 value for a flag, and a boolean which indicates if
// the flag was set
func (this *Flags) GetFloat64(name string) (float64, bool) {
	value := this.flagset.Lookup(name)
	if value == nil {
		return 0.0, false
	}
	return value.Value.(flag.Getter).Get().(float64), this.HasFlag(name)
}
