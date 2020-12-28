package config

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/djthorpe/gopi/v3"
	"github.com/hashicorp/go-multierror"
)

///////////////////////////////////////////////////////////////////////////////
// TYPES

type config struct {
	*flag.FlagSet
	args     []string
	commands *command
	flags    map[string][]string
}

///////////////////////////////////////////////////////////////////////////////
// NEW

func New(name string, args []string) gopi.Config {
	this := new(config)
	this.FlagSet = flag.NewFlagSet(name, flag.ContinueOnError)
	this.FlagSet.Usage = func() {
		if cmd, _ := this.GetCommand(nil); cmd == nil {
			this.usageAll()
		} else {
			this.usageOne(cmd)
		}
	}
	this.args = args
	this.flags = make(map[string][]string)
	this.commands = NewCommand(name, "", "", args, nil)
	return this
}

///////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *config) Version() gopi.Version {
	return NewVersion(this.FlagSet.Name())
}

func (this *config) Parse() error {
	// Check for already parsed
	if this.FlagSet.Parsed() {
		return nil
	}

	// Perform parse and translate error
	if err := this.FlagSet.Parse(this.args); err == flag.ErrHelp {
		return gopi.ErrHelp
	} else if err != nil {
		return err
	}

	// Return success
	return nil
}

// Usage prints out the command usage instructions for all
// commands or for a specific command
func (this *config) Usage(name string) {
	name = strings.ToLower(strings.TrimSpace(name))
	parts := strings.Fields(name)
	if cmd, _ := this.GetCommand(parts); cmd != nil {
		this.usageOne(cmd)
	} else {
		this.usageAll()
	}
}

// Command defines a command and associated function call
func (this *config) Command(name, usage string, fn gopi.CommandFunc) error {
	name = strings.ToLower(strings.TrimSpace(name))
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return gopi.ErrBadParameter.WithPrefix("name")
	}
	cmd := this.commands
	for i, part := range parts {
		last := (i == len(parts)-1)
		if last {
			// Check for duplicate
			if getCommand(part, cmd.commands) != nil {
				return gopi.ErrDuplicateEntry.WithPrefix("Command ", name)
			}
			// Append commmand
			cmd.commands = append(cmd.commands, NewCommand(part, usage, "", nil, fn))
		}
		if cmd = getCommand(part, cmd.commands); cmd == nil {
			return gopi.ErrNotFound.WithPrefix("Command ", part)
		}
	}
	// Return success
	return nil
}

// GetCommand returns command with arguments or nil if
// a command was not registered which matches the signature,
// or if the wrong flags have been set for the command.
// When the arguments parameter is nil, the arguments from
// the FlagSet are used.
func (this *config) GetCommand(args []string) (gopi.Command, error) {
	if len(this.commands.commands) == 0 {
		return nil, gopi.ErrNotFound.WithPrefix("GetCommand")
	}
	if args == nil {
		args = this.FlagSet.Args()
	}
	// Iterate through commands
	cmd := this.commands
	i := 0
	for i = range args {
		child := getCommand(args[i], cmd.commands)
		if child == nil {
			i--
			break
		} else {
			cmd = child
		}
	}

	// Special case where root command matched
	if cmd == this.commands {
		return nil, nil
	}

	// Create a new command from the existing one, setting arguments
	cmd = NewCommand(strings.Join(args[:i+1], " "), cmd.usage, cmd.syntax, args[i+1:], cmd.fn)

	// Iterate through set flags
	var result error
	this.FlagSet.Visit(func(f *flag.Flag) {
		if this.flagIsGlobal(f) {
			cmd.flags[f.Name] = f
		} else if this.flagIsLocal(f, cmd.Name()) {
			cmd.flags[f.Name] = f
		} else {
			result = multierror.Append(result, fmt.Errorf("%w: Flag set but not defined: -%v", gopi.ErrBadParameter, f.Name))
		}
	})

	// Success
	return cmd, multierror.Flatten(result)
}

///////////////////////////////////////////////////////////////////////////////
// DEFINE FLAGS

func (this *config) FlagString(name, value, usage string, cmds ...string) *string {
	this.flags[name] = cmds
	return this.FlagSet.String(name, value, usage)
}

func (this *config) FlagBool(name string, value bool, usage string, cmds ...string) *bool {
	this.flags[name] = cmds
	return this.FlagSet.Bool(name, value, usage)
}

func (this *config) FlagUint(name string, value uint, usage string, cmds ...string) *uint {
	this.flags[name] = cmds
	return this.FlagSet.Uint(name, value, usage)
}

func (this *config) FlagInt(name string, value int, usage string, cmds ...string) *int {
	this.flags[name] = cmds
	return this.FlagSet.Int(name, value, usage)
}

func (this *config) FlagDuration(name string, value time.Duration, usage string, cmds ...string) *time.Duration {
	this.flags[name] = cmds
	return this.FlagSet.Duration(name, value, usage)
}

func (this *config) FlagFloat(name string, value float64, usage string, cmds ...string) *float64 {
	this.flags[name] = cmds
	return this.FlagSet.Float64(name, value, usage)
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
	} else if value_, err := strconv.ParseUint(flag.Value.String(), 0, 64); err != nil {
		return 0
	} else {
		return uint(value_)
	}
}

func (this *config) GetFloat(name string) float64 {
	if flag := this.FlagSet.Lookup(name); flag == nil {
		return 0
	} else if value_, err := strconv.ParseFloat(flag.Value.String(), 64); err != nil {
		return 0
	} else {
		return float64(value_)
	}
}

func (this *config) GetInt(name string) int {
	if flag := this.FlagSet.Lookup(name); flag == nil {
		return 0
	} else if value_, err := strconv.ParseInt(flag.Value.String(), 0, 64); err != nil {
		return 0
	} else {
		return int(value_)
	}
}

func (this *config) GetDuration(name string) time.Duration {
	if flag := this.FlagSet.Lookup(name); flag == nil {
		return 0
	} else if value_, err := time.ParseDuration(flag.Value.String()); err != nil {
		return 0
	} else {
		return value_
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

///////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *config) usageAll() {
	w := this.FlagSet.Output()
	name := this.FlagSet.Name()

	fmt.Fprintln(w, "Syntax:")
	if len(this.commands.commands) == 0 {
		fmt.Fprintf(w, "  %v (<flags>) (<args>)\n", name)
		fmt.Fprintf(w, "  %v -help\n", name)
	} else {
		fmt.Fprintf(w, "  %v (<flags>) <command> (<args>)\n", name)
		fmt.Fprintf(w, "  %v -help (<command>)\n", name)
		fmt.Fprintln(w, "\nCommands:")
	}
	usageCommand(w, this.commands, nil)
	this.usageFlags("")
}

func usageCommand(w io.Writer, cmd *command, path []string) {
	if path != nil {
		fmt.Fprintf(w, "  %v %v\n  \t%v\n", strings.Join(path, " "), cmd.syntax, cmd.usage)
	}
	for _, cmd := range cmd.commands {
		usageCommand(w, cmd, append(path, cmd.name))
	}
}

func (this *config) usageOne(cmd gopi.Command) {
	w := this.Output()
	name := this.FlagSet.Name()
	syntax, description := cmd.Usage()

	fmt.Fprintln(w, description)
	fmt.Fprintln(w, "\nSyntax:")
	fmt.Fprintf(w, "  %v (<flags>) %v %v\n", name, cmd.Name(), syntax)
	this.usageFlags(cmd.Name())
	this.usageFlags("")
}

func (this *config) usageFlags(name string) {
	w := this.FlagSet.Output()

	// Gather the relevant flags
	flags := []*flag.Flag{}
	this.FlagSet.VisitAll(func(flag *flag.Flag) {
		if name == "" && this.flagIsGlobal(flag) {
			flags = append(flags, flag)
		} else if name != "" && this.flagIsLocal(flag, name) {
			flags = append(flags, flag)
		}
	})

	// Don't output if no flags
	if len(flags) == 0 {
		return
	}

	// Print header
	if name == "" {
		fmt.Fprintf(w, "\nGlobal Flags:\n")
	} else {
		fmt.Fprintf(w, "\nFlags for %q:\n", name)
	}

	// Print flags
	for _, flag := range flags {
		arg, usage, def := flagUsage(flag)
		fmt.Fprintf(w, "  -%v %v\n  \t%v %v\n", flag.Name, arg, usage, def)
	}
}

func (this *config) flagIsGlobal(f *flag.Flag) bool {
	if cmds, exists := this.flags[f.Name]; exists == false {
		return false
	} else {
		return len(cmds) == 0
	}
}

func (this *config) flagIsLocal(f *flag.Flag, name string) bool {
	if cmds, exists := this.flags[f.Name]; exists {
		for _, cmd := range cmds {
			if name == cmd {
				return true
			}
		}
	}
	return false
}

// Get type and defaults to print in defaults
func flagUsage(f *flag.Flag) (string, string, string) {
	arg, usage := flag.UnquoteUsage(f)
	if f.DefValue == "" {
		return arg, usage, ""
	} else {
		if arg == "string" {
			return arg, usage, fmt.Sprintf("(default %q)", f.DefValue)
		} else if arg == "" { // Boolean value
			return arg, usage, ""
		} else {
			return arg, usage, fmt.Sprintf("(default %v)", f.DefValue)
		}
	}
}

// Return command from name
func getCommand(name string, cmds []*command) *command {
	for _, cmd := range cmds {
		if cmd.name == name {
			return cmd
		}
	}
	return nil
}
