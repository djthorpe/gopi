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

var commandmap = map[string]func(*rpi.State) error{
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

func allCommand(pi *rpi.State) error {
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

	if err := modelCommand(pi); err != nil {
		return err
	}

	return nil
}

func tempCommand(pi *rpi.State) error {
	// print out temperature
	coretemp, err := pi.GetCoreTemperatureCelcius()
	if err != nil {
		return err
	}

	fmt.Printf("Temperature=%vºC\n", coretemp)
	return nil
}

func clocksCommand(pi *rpi.State) error {
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

func voltsCommand(pi *rpi.State) error {
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

func memoryCommand(pi *rpi.State) error {
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

func codecsCommand(pi *rpi.State) error {
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

func otpCommand(pi *rpi.State) error {
	// print out OTP memory
	otp, err := pi.GetOTP()
	if err != nil {
		return err
	}

	fmt.Println("OTP")
	for i, v := range otp {
		fmt.Printf("  %02d=%08X\n", i, v)
	}

	return nil
}

func serialCommand(pi *rpi.State) error {
	// print out Serial number
	serial, err := pi.GetSerial()
	if err != nil {
		return err
	}

	fmt.Printf("Serial=%016X\n", serial)

	return nil
}

func revisionCommand(pi *rpi.State) error {
	// print out Revision
	revision, err := pi.GetRevision()
	if err != nil {
		return err
	}

	fmt.Printf("Revision=%08X\n", revision)

	return nil
}

func modelCommand(pi *rpi.State) error {
	// print out Model information
	model, err := pi.GetModel()
	if err != nil {
		return err
	}

	fmt.Println("Model")
	fmt.Printf("  Revision=0x%08X\n", model.Revision)
	fmt.Printf("  Product=Respberry Pi %v (%v)\n", model.ProductString, model.Product)
	fmt.Printf("  Processor=%v (%v)\n", model.ProcessorString, model.Processor)
	fmt.Printf("  Memory=%vM\n", model.MemoryMB)
	fmt.Printf("  Manufacturer=%v (%v)\n", model.ManufacturerString, model.Manufacturer)
	fmt.Printf("  PCB Revision=%v\n", model.PCBRevision)
	fmt.Printf("  Peripheral Base=0x%08X\n", model.PeripheralBase)

	return nil
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

	var err error
	if len(args) == 1 && commandmap[args[0]] != nil {
		if f := commandmap[args[0]]; f != nil {
			err = f(pi)
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
