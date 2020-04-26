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
	"io"
	"os"
	"strconv"
	"sync"
	"time"

	// Frameworks
	"github.com/djthorpe/gopi/v2"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type flagset struct {
	flags   *flag.FlagSet
	name    string
	flagMap map[gopi.FlagNS]map[string]*flag.Flag
	sync.Mutex
}

type stringValue struct {
	string
}

type boolValue struct {
	bool
}

type durationValue struct {
	time.Duration
}

type uintValue struct {
	uint
}

type intValue struct {
	int
}

type float64Value struct {
	float64
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// Newflags returns a new flagset object
func NewFlags(name string) gopi.Flags {
	return &flagset{
		flags: flag.NewFlagSet(name, flag.ContinueOnError),
		name:  name,
	}
}

// Parse command line arguments into flags and pure arguments
func (this *flagset) Parse(args []string) error {
	// set empty usage function
	devnull, err := os.Open(os.DevNull)
	if err != nil {
		return err
	}
	this.flags.SetOutput(devnull)
	defer devnull.Close()

	// set version and service flags
	this.setVersionFlags()
	this.setServiceFlags()

	// parse flags
	if err := this.flags.Parse(args); err == flag.ErrHelp {
		return gopi.ErrHelp
	} else if err != nil {
		return err
	}

	// set hash of flags that were set
	this.flags.VisitAll(func(f *flag.Flag) {
		this.setFlag(f.Name, gopi.FLAG_NS_DEFAULT, f)
	})

	// return success
	return nil
}

// Parsed reports whether the command-line flags have been parsed
func (this *flagset) Parsed() bool {
	return this.flags.Parsed()
}

// Name returns the name of the flagset (usually same as application)
func (this *flagset) Name() string {
	return this.name
}

// Args returns the command line arguments as an array which aren't flags
func (this *flagset) Args() []string {
	return this.flags.Args()
}

// Flags returns the command line arguments as an array
func (this *flagset) Flags(ns gopi.FlagNS) []*flag.Flag {
	if len(this.flagMap) == 0 {
		return nil
	} else if flagMap, exists := this.flagMap[ns]; exists == false {
		return nil
	} else {
		flags := make([]*flag.Flag, 0, len(flagMap))
		for _, v := range flagMap {
			flags = append(flags, v)
		}
		return flags
	}
}

// HasFlag returns true if a flag exists
func (this *flagset) HasFlag(name string, ns gopi.FlagNS) bool {
	if len(this.flagMap) == 0 {
		return false
	} else if flagMap, exists := this.flagMap[ns]; exists == false {
		return false
	} else if _, exists := flagMap[name]; exists {
		return true
	} else {
		return false
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *flagset) String() string {
	str := fmt.Sprintf("gopi.Flags parsed=%v name=%v args=%v ", this.Parsed(), strconv.Quote(this.Name()), this.Args())
	for ns, flagMap := range this.flagMap {
		str += fmt.Sprintf("%v={ ", ns)
		for k, v := range flagMap {
			str += k + "=" + strconv.Quote(v.Value.String()) + " "
		}
		str += "}"
	}
	return "<" + str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// DEFINE FLAGS

// FlagString defines string flag and return pointer to the flag value
func (this *flagset) FlagString(name, value, usage string) *string {
	return this.flags.String(name, value, usage)
}

// FlagBool defines a boolean flag and return pointer to the flag value
func (this *flagset) FlagBool(name string, value bool, usage string) *bool {
	return this.flags.Bool(name, value, usage)
}

// FlagDuration defines duration flag and return pointer to the flag value
func (this *flagset) FlagDuration(name string, value time.Duration, usage string) *time.Duration {
	return this.flags.Duration(name, value, usage)
}

// FlagInt defines integer flag and return pointer to the flag value
func (this *flagset) FlagInt(name string, value int, usage string) *int {
	return this.flags.Int(name, value, usage)
}

// FlagUint defines unsigned integer flag and return pointer to the flag value
func (this *flagset) FlagUint(name string, value uint, usage string) *uint {
	return this.flags.Uint(name, value, usage)
}

// FlagFloat64 defines float64 flag and return pointer to the flag value
func (this *flagset) FlagFloat64(name string, value float64, usage string) *float64 {
	return this.flags.Float64(name, value, usage)
}

////////////////////////////////////////////////////////////////////////////////
// GET FLAGS

// GetBool gets boolean value for a flag, or default value if not set
func (this *flagset) GetBool(name string, ns gopi.FlagNS) bool {
	if value := this.getFlag(name, gopi.FLAG_NS_DEFAULT); value == nil {
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
func (this *flagset) GetString(name string, ns gopi.FlagNS) string {
	if value := this.getFlag(name, gopi.FLAG_NS_DEFAULT); value == nil {
		return ""
	} else if getter, ok := value.Value.(flag.Getter); ok == false {
		return ""
	} else {
		return getter.String()
	}
}

// GetDuration gets duration value for a flag
func (this *flagset) GetDuration(name string, ns gopi.FlagNS) time.Duration {
	if value := this.getFlag(name, gopi.FLAG_NS_DEFAULT); value == nil {
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
func (this *flagset) GetInt(name string, ns gopi.FlagNS) int {
	if value := this.getFlag(name, gopi.FLAG_NS_DEFAULT); value == nil {
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
func (this *flagset) GetUint(name string, ns gopi.FlagNS) uint {
	if value := this.getFlag(name, gopi.FLAG_NS_DEFAULT); value == nil {
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
func (this *flagset) GetFloat64(name string, ns gopi.FlagNS) float64 {
	if value := this.getFlag(name, gopi.FLAG_NS_DEFAULT); value == nil {
		return 0
	} else if getter, ok := value.Value.(flag.Getter); ok == false {
		return 0
	} else if value_, ok := getter.Get().(float64); ok == false {
		return 0
	} else {
		return value_
	}
}

////////////////////////////////////////////////////////////////////////////////
// SET FLAGS

func (this *flagset) SetBool(name string, ns gopi.FlagNS, v bool) {
	this.setFlag(name, ns, &flag.Flag{
		Name:  name,
		Value: &boolValue{v},
	})
}

func (this *flagset) SetString(name string, ns gopi.FlagNS, v string) {
	this.setFlag(name, ns, &flag.Flag{
		Name:  name,
		Value: &stringValue{v},
	})
}

func (this *flagset) SetDuration(name string, ns gopi.FlagNS, v time.Duration) {
	this.setFlag(name, ns, &flag.Flag{
		Name:  name,
		Value: &durationValue{v},
	})
}

func (this *flagset) SetInt(name string, ns gopi.FlagNS, v int) {
	this.setFlag(name, ns, &flag.Flag{
		Name:  name,
		Value: &intValue{v},
	})
}

func (this *flagset) SetUint(name string, ns gopi.FlagNS, v uint) {
	this.setFlag(name, ns, &flag.Flag{
		Name:  name,
		Value: &uintValue{v},
	})
}

func (this *flagset) SetFloat64(name string, ns gopi.FlagNS, v float64) {
	this.setFlag(name, ns, &flag.Flag{
		Name:  name,
		Value: &float64Value{v},
	})
}

////////////////////////////////////////////////////////////////////////////////
// USAGE

func (this *flagset) Usage(io io.Writer) {
	this.flags.SetOutput(io)
	this.flags.PrintDefaults()
}

func (this *flagset) Version(io io.Writer) {
	if flagMap, exists := this.flagMap[gopi.FLAG_NS_VERSION]; exists {
		for k, v := range flagMap {
			fmt.Fprintf(io, "%-10s %v\n", k, v.Value)
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *flagset) setFlag(name string, ns gopi.FlagNS, value *flag.Flag) error {
	this.Lock()
	defer this.Unlock()

	if this.flagMap == nil {
		this.flagMap = make(map[gopi.FlagNS]map[string]*flag.Flag)
	}
	if _, exists := this.flagMap[ns]; exists == false {
		this.flagMap[ns] = make(map[string]*flag.Flag)
	}
	if flagMap, exists := this.flagMap[ns]; exists == false {
		return gopi.ErrInternalAppError
	} else {
		flagMap[name] = value
	}
	// Success
	return nil
}

func (this *flagset) getFlag(name string, ns gopi.FlagNS) *flag.Flag {
	this.Lock()
	defer this.Unlock()

	if this.flagMap == nil {
		return nil
	} else if flagMap, exists := this.flagMap[ns]; exists == false {
		return nil
	} else if flagValue, exists := flagMap[name]; exists == false {
		return nil
	} else {
		return flagValue
	}
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION OF VALUE

func (this *stringValue) Set(value string) error {
	this.string = value
	return nil
}

func (this *stringValue) String() string {
	return this.string
}

func (this *stringValue) Get() interface{} {
	return this.string
}

func (this *boolValue) Set(value string) error {
	if value_, err := strconv.ParseBool(value); err != nil {
		return err
	} else {
		this.bool = value_
		return nil
	}
}

func (this *boolValue) String() string {
	return fmt.Sprint(this.bool)
}

func (this *boolValue) Get() interface{} {
	return this.bool
}

func (this *durationValue) Set(value string) error {
	if value_, err := time.ParseDuration(value); err != nil {
		return err
	} else {
		this.Duration = value_
		return nil
	}
}

func (this *durationValue) String() string {
	return fmt.Sprint(this.Duration)
}

func (this *durationValue) Get() interface{} {
	return this.Duration
}

func (this *intValue) Set(value string) error {
	if value_, err := strconv.ParseInt(value, 10, 64); err != nil {
		return err
	} else {
		this.int = int(value_)
		return nil
	}
}

func (this *intValue) String() string {
	return fmt.Sprint(this.int)
}

func (this *intValue) Get() interface{} {
	return this.int
}

func (this *uintValue) Set(value string) error {
	if value_, err := strconv.ParseUint(value, 10, 64); err != nil {
		return err
	} else {
		this.uint = uint(value_)
		return nil
	}
}

func (this *uintValue) String() string {
	return fmt.Sprint(this.uint)
}

func (this *uintValue) Get() interface{} {
	return this.uint
}

func (this *float64Value) Set(value string) error {
	if value_, err := strconv.ParseFloat(value, 64); err != nil {
		return err
	} else {
		this.float64 = value_
		return nil
	}
}

func (this *float64Value) String() string {
	return fmt.Sprint(this.float64)
}

func (this *float64Value) Get() interface{} {
	return this.float64
}
