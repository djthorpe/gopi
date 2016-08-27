/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md

	vcgencmd.go is a command-line utility to print out all sorts of information
	from a Raspberry Pi using the vcgcmd interface. For example:

	vcgencmd temp
	vcgencmd clocks
	vcgencmd volts

    etc.
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
)

import (
	rpi "../device/rpi"
)

////////////////////////////////////////////////////////////////////////////

var (
	commandmap = map[string]func(*rpi.RaspberryPi) error{
		"all": allCommands,
	}
)

////////////////////////////////////////////////////////////////////////////////

func allCommands(pi *rpi.RaspberryPi) error {
	// return nil
	return nil
}

////////////////////////////////////////////////////////////////////////////////

func main() {

	// Open up the RaspberryPi interface
	rpi, err := rpi.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: ", err)
		os.Exit(-1)
	}
	defer rpi.Close()

	// Set flag usage, parse flags
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", path.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, " <command> can be one of the following: ")
		for k, _ := range commandmap {
			fmt.Fprintf(os.Stderr, "%s, ", k)
		}
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	// Read arguments, exit if no command
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(-1)
	}

	// Call command, report error
	if len(args) == 1 && commandmap[args[0]] != nil {
		if f := commandmap[args[0]]; f != nil {
			err = f(rpi)
		}
	} else {
		err = errors.New("Unknown command")
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}

}
