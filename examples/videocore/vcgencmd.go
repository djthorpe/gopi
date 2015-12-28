/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md

	vcgencmd.go is a command-line utility to print out all sorts of information
	from a Raspberry Pi using the VCGenCmd interface. For example:

	vcgencmd temp
	vcgencmd clocks
	vcgencmd volts

    etc.
*/
package main

import (
	"flag"
	"fmt"
	"github.com/djthorpe/gopi/rpi"
	"os"
	"path"
	"strings"
)

////////////////////////////////////////////////////////////////////////////

var commandmap = map[string]func(*rpi.State){
	"all":      allCommand,
	"temp":     tempCommand,
	"clocks":   clocksCommand,
	"volts":    voltsCommand,
	"memory":   memoryCommand,
	"codecs":   codecsCommand,
	"otp":      otpCommand,
	"serial":   serialCommand,
	"revision": revisionCommand,
	"model":    modelCommand,
}

////////////////////////////////////////////////////////////////////////////

func allCommand(pi *rpi.State) {
	tempCommand(pi)
	clocksCommand(pi)
	voltsCommand(pi)
	memoryCommand(pi)
	codecsCommand(pi)
	otpCommand(pi)
	serialCommand(pi)
	revisionCommand(pi)
	modelCommand(pi)
}

func tempCommand(pi *rpi.State) {
	// print out temperature
	coretemp, err := pi.GetCoreTemperatureCelcius()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Printf("Temperature=%vºC\n", coretemp)
}

func clocksCommand(pi *rpi.State) {
	// print out clocks
	clocks, err := pi.GetClockFrequencyHertz()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("Clock Frequency")
	for k, v := range clocks {
		fmt.Printf("  %v=%vMHz\n", k, (float64(v) / 1E6))
	}
}

func voltsCommand(pi *rpi.State) {
	// print out volts
	volts, err := pi.GetVolts()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("Volts")
	for k, v := range volts {
		fmt.Printf("  %v=%vV\n", k, v)
	}
}

func memoryCommand(pi *rpi.State) {
	// print out memory sizes
	memory, err := pi.GetMemoryMegabytes()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("Memory")
	for k, v := range memory {
		fmt.Printf("  %v=%vMB\n", k, v)
	}
}

func codecsCommand(pi *rpi.State) {
	// print out codecs
	codecs, err := pi.GetCodecs()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("Codecs")
	for k, v := range codecs {
		fmt.Printf("  %v=%v\n", k, v)
	}
}

func otpCommand(pi *rpi.State) {
	// print out OTP memory
	otp, err := pi.GetOTP()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Println("OTP")
	for i, v := range otp {
		fmt.Printf("  %02d=%08X\n", i, v)
	}
}

func serialCommand(pi *rpi.State) {
	// print out Serial number
	serial, err := pi.GetSerial()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Printf("Serial=%016X\n", serial)
}

func revisionCommand(pi *rpi.State) {
	// print out Revision
	revision, err := pi.GetRevision()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Printf("Revision=%08X\n", revision)
}

func modelCommand(pi *rpi.State) {
	// print out Revision
	model, err := pi.GetModel()
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	fmt.Printf("Model=%+v\n", model)
}

////////////////////////////////////////////////////////////////////////////

func main() {

	pi := rpi.New()
	defer pi.Terminate()

	////////////////////////////////////////////////////////////////////////////

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s <command>\n", path.Base(os.Args[0]))

		fmt.Fprintf(os.Stderr, " <command> can be one of the following: ")
		for k, _ := range commandmap {
			fmt.Fprintf(os.Stderr, "%s, ", k)
		}
		vccommands, _ := pi.GetCommands()
		for _, v := range vccommands {
			fmt.Fprintf(os.Stderr, "%s, ", v)
		}
		fmt.Fprintf(os.Stderr, "\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	////////////////////////////////////////////////////////////////////////////

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	////////////////////////////////////////////////////////////////////////////

	var done bool

	// attempt to run a custom command
	if len(args) == 1 {
		if f := commandmap[args[0]]; f != nil {
			f(pi)
			done = true
		}
	}

	// if custom command not run, use VCGenCmd
	if done == false {
		fmt.Println(rpi.VCGenCmd(strings.Join(args, " ")))
	}

}
