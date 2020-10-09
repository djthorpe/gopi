package config

import (
	"context"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/djthorpe/gopi/v3"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type config struct {
	*flag.FlagSet
	args     []string
	commands []command
}

type command struct {
	name, usage string
	args        []string
	fn          gopi.CommandFunc
}

///////////////////////////////////////////////////////////////////////////////
// NEW

func New(name string, args []string) gopi.Config {
	this := new(config)
	this.FlagSet = flag.NewFlagSet(name, flag.ContinueOnError)
	this.FlagSet.Usage = this.usage
	this.args = args
	return this
}

func NewCommand(cmd command, args []string) gopi.Command {
	this := &cmd
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

func (this *config) Usage(name string) {
	if name == "" {
		this.usage()
	} else {
		fmt.Println("TODO: usage for", name)
	}
}

func (this *config) Command(name, usage string, fn gopi.CommandFunc) error {
	name = strings.ToLower(strings.TrimSpace(name))
	if name == "" {
		return gopi.ErrBadParameter.WithPrefix("name")
	} else if cmd := this.GetCommand([]string{name}); cmd != nil {
		return gopi.ErrBadParameter.WithPrefix("name")
	} else {
		this.commands = append(this.commands, command{
			name:  name,
			usage: usage,
			fn:    fn,
		})
	}

	// Return success
	return nil
}

func (this *config) GetCommand(args []string) gopi.Command {
	if len(this.commands) == 0 {
		return nil
	}
	if args == nil {
		args = this.FlagSet.Args()
	}
	if len(args) == 0 {
		return NewCommand(this.commands[0], args)
	}
	name := strings.ToLower(strings.TrimSpace(args[0]))
	for _, cmd := range this.commands {
		if cmd.name == name {
			return NewCommand(cmd, args[1:])
		}
	}

	// Command not found
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

func (this *config) String() string {
	str := "<config"
	str += " name=" + strconv.Quote(this.FlagSet.Name())
	this.FlagSet.Visit(func(f *flag.Flag) {
		str += fmt.Sprintf(" %v=%q", f.Name, f.Value.String())
	})
	return str + ">"
}

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

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *config) usage() {
	w := this.FlagSet.Output()
	name := this.FlagSet.Name()

	fmt.Fprintln(w, "Syntax:")
	fmt.Fprintf(w, "  %v (<flags>) <command> (<args>)\n", name)
	fmt.Fprintf(w, "  %v -help (<command>)\n", name)
	fmt.Fprintln(w, "\nCommands:")
	for _, cmd := range this.commands {
		fmt.Fprintf(w, "  %v %v\n  \t%v\n", cmd.name, "TODO", cmd.usage)
	}
	this.usageFlags("")
}

func (this *config) usageFlags(name string) {
	w := this.FlagSet.Output()
	if name == "" {
		fmt.Fprintf(w, "\nGlobal Flags:\n")
	} else {
		fmt.Fprintf(w, "\nFlags for %q:\n", name)
	}
	this.FlagSet.VisitAll(func(flag *flag.Flag) {
		arg, usage, def := flagUsage(flag)
		fmt.Fprintf(w, "  -%v %v\n  \t%v %v\n", flag.Name, arg, usage, def)
	})
}

// Get type and defaults to print in defaults
func flagUsage(f *flag.Flag) (string, string, string) {
	arg, usage := flag.UnquoteUsage(f)
	if f.DefValue == "" {
		return arg, usage, ""
	} else {
		if arg == "string" {
			return arg, usage, fmt.Sprintf("(default %q)", f.DefValue)
		} else {
			return arg, usage, fmt.Sprintf("(default %v)", f.DefValue)
		}
	}
}
