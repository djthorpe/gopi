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
	"flag"
	"fmt"
	"os"
	"path"
	"strings"
)

import (
	rpi "../device/rpi"
)

////////////////////////////////////////////////////////////////////////////

var (
	commandmap = map[string]func(*rpi.RaspberryPi) error{
		"all":    allCommands,
		"temp":   tempCommand,
		"clocks": clocksCommand,
		"volts":  voltsCommand,
		"memory": memoryCommand,
		"codecs": codecsCommand,
		"otp":      otpCommand,
		"serial":   serialCommand,
		"revision": revisionCommand,
	}
)

////////////////////////////////////////////////////////////////////////////////

func allCommands(pi *rpi.RaspberryPi) error {

	if err := tempCommand(pi); err != nil {
		return err
	}

	if err := clocksCommand(pi); err != nil {
		return err
	}

	if err := voltsCommand(pi); err != nil {
		return err
	}

	if err := memoryCommand(pi); err != nil {
		return err
	}

	if err := codecsCommand(pi); err != nil {
		return err
	}

	if err := otpCommand(pi); err != nil {
		return err
	}	

	if err := serialCommand(pi); err != nil {
		return err
	}

	if err := revisionCommand(pi); err != nil {
		return err
	}

	// return success
	return nil
}

func tempCommand(pi *rpi.RaspberryPi) error {
	// print out temperature
	coretemp, err := pi.GetCoreTemperatureCelcius()
	if err != nil {
		return err
	}
	fmt.Printf("Temperature=%vºC\n", coretemp)

	// return success
	return nil
}

func clocksCommand(pi *rpi.RaspberryPi) error {
	// print out clocks
	clocks, err := pi.GetClockFrequencyHertz()
	if err != nil {
		return err
	}

	fmt.Println("Clock Frequency")
	for k, v := range clocks {
		fmt.Printf("  %v=%vMHz\n", k, (float64(v) / 1E6))
	}

	return nil
}

func voltsCommand(pi *rpi.RaspberryPi) error {
	// print out volts
	volts, err := pi.GetVolts()
	if err != nil {
		return err
	}

	fmt.Println("Volts")
	for k, v := range volts {
		fmt.Printf("  %v=%vV\n", k, v)
	}

	return nil
}

func memoryCommand(pi *rpi.RaspberryPi) error {
	// print out memory sizes
	memory, err := pi.GetMemoryMegabytes()
	if err != nil {
		return err
	}

	fmt.Println("Memory")
	for k, v := range memory {
		fmt.Printf("  %v=%vMB\n", k, v)
	}

	return nil
}


func codecsCommand(pi *rpi.RaspberryPi) error {
	// print out codecs
	codecs, err := pi.GetCodecs()
	if err != nil {
		return err
	}

	fmt.Println("Codecs")
	for k, v := range codecs {
		fmt.Printf("  %v=%v\n", k, v)
	}

	return nil
}

func otpCommand(pi *rpi.RaspberryPi) error {
	// print out OTP memory
	otp, err := pi.GetOTP()
	if err != nil {
		return err
	}

	fmt.Println("OTP")
	for i, v := range otp {
		fmt.Printf("  %02d=0x%08X\n", i, v)
	}

	return nil
}

func serialCommand(pi *rpi.RaspberryPi) error {
	// print out Serial number
	serial, err := pi.GetSerial()
	if err != nil {
		return err
	}

	fmt.Printf("Serial=0x%016X\n", serial)

	return nil
}

func revisionCommand(pi *rpi.RaspberryPi) error {
	// print out Revision
	revision, err := pi.GetRevision()
	if err != nil {
		return err
	}

	fmt.Printf("Revision=0x%08X\n", revision)

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
		vccommands, _ := rpi.GetCommands()
		for _, v := range vccommands {
			fmt.Fprintf(os.Stderr, "%s, ", v)
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
		var value string
		value, err := rpi.VCGenCmd(strings.Join(args, " "))
		if err == nil {
			fmt.Println(value)
		}
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}

}
