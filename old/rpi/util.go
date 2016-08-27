/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	regexpTemperature = "temp=(\\d+\\.?\\d*)"
	regexpCommands    = "commands=\"([^\"]+)\""
	regexpClock       = "frequency\\((\\d+)\\)=(\\d+)"
	regexpVolt        = "volt=(\\d*\\.?\\d*)V"
	regexpCodec       = "(\\w+)=(enabled|disabled)"
	regexpMemory      = "(\\w+)=(\\d+)M"
	regexpOTP         = "(\\d\\d):([0123456789abcdefABCDEF]{8})"
)

var (
	rCommands    *regexp.Regexp
	rTemperature *regexp.Regexp
	rClock       *regexp.Regexp
	rVolt        *regexp.Regexp
	rCodec       *regexp.Regexp
	rMemory      *regexp.Regexp
	rOTP         *regexp.Regexp
)

////////////////////////////////////////////////////////////////////////////////

type State struct {
	revision uint32
	serial   uint64
}

////////////////////////////////////////////////////////////////////////////////

func New() *State {
	// create this object
	this := new(State)

	// initialize
	BCMHostInit()
	VCGenCmdInit()

	// Return this
	return this
}

func (this *State) Terminate() {
	VCGenCmdStop()
	BCMHostTerminate()
}

////////////////////////////////////////////////////////////////////////////////

func (this *State) GetCommands() ([]string, error) {
	// Compile regular expression
	if rCommands == nil {
		rCommands = regexp.MustCompile(regexpCommands)
	}

	// retrieve value as text
	value, err := VCGenCmd(GEMCMD_COMMANDS)
	if err != nil {
		return []string{}, err
	}

	// Find values within text
	match := rCommands.FindStringSubmatch(value)
	if len(match) != 2 {
		return []string{}, ErrorResponse
	}

	// Split commands
	commands := strings.Split(match[1], ",")
	for i, command := range commands {
		commands[i] = strings.TrimSpace(command)
	}
	// return commands
	return commands, nil
}

func (this *State) GetCoreTemperatureCelcius() (float64, error) {
	// Compile regular expression
	if rTemperature == nil {
		rTemperature = regexp.MustCompile(regexpTemperature)
	}

	// retrieve value as text
	value, err := VCGenCmd(GENCMD_MEASURE_TEMP)
	if err != nil {
		return 0.0, err
	}

	// Find value within text
	match := rTemperature.FindStringSubmatch(value)
	if len(match) != 2 {
		return 0.0, ErrorResponse
	}

	// Convert to float64
	value2, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0.0, err
	}

	// Return value as float64
	return value2, nil
}

func (this *State) GetClockFrequencyHertz() (map[string]uint64, error) {
	// Compile regular expression
	if rClock == nil {
		rClock = regexp.MustCompile(regexpClock)
	}

	// retrieve values as text
	command := strings.Split(GENCMD_MEASURE_CLOCK, " ")
	clocks := make(map[string]uint64, len(command))
	for _, name := range command[1:] {

		// Retrieve clock value
		value, err := VCGenCmd(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := rClock.FindStringSubmatch(value)
		if len(match) != 3 {
			return nil, ErrorResponse
		}

		// Convert to uint64
		value2, err := strconv.ParseUint(match[2], 10, 64)
		if err != nil {
			return nil, err
		}

		// Set value
		clocks[name] = value2
	}

	return clocks, nil
}

func (this *State) GetVolts() (map[string]float64, error) {
	// Compile regular expression
	if rVolt == nil {
		rVolt = regexp.MustCompile(regexpVolt)
	}

	// retrieve values as text
	command := strings.Split(GENCMD_MEASURE_VOLTS, " ")
	volts := make(map[string]float64, len(command))
	for _, name := range command[1:] {

		// Retrieve volt value
		value, err := VCGenCmd(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := rVolt.FindStringSubmatch(value)
		if len(match) != 2 {
			return nil, ErrorResponse
		}

		// Convert to uint64
		value2, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return nil, err
		}

		// Set value
		volts[name] = value2
	}

	return volts, nil
}

func (this *State) GetCodecs() (map[string]bool, error) {
	// Compile regular expression
	if rCodec == nil {
		rCodec = regexp.MustCompile(regexpCodec)
	}

	// retrieve values as text
	command := strings.Split(GENCMD_CODEC_ENABLED, " ")
	codecs := make(map[string]bool, len(command))
	for _, name := range command[1:] {

		// Retrieve volt value
		value, err := VCGenCmd(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := rCodec.FindStringSubmatch(value)
		if len(match) != 3 {
			return nil, ErrorResponse
		}

		// Convert to bool
		if match[2] == "enabled" {
			codecs[name] = true
		} else {
			codecs[name] = false
		}
	}

	return codecs, nil
}

func (this *State) GetMemoryMegabytes() (map[string]uint64, error) {
	// Compile regular expression
	if rMemory == nil {
		rMemory = regexp.MustCompile(regexpMemory)
	}

	// retrieve values as text
	command := strings.Split(GENCMD_MEMORY, " ")
	memories := make(map[string]uint64, len(command))
	for _, name := range command[1:] {

		// Retrieve memory value
		value, err := VCGenCmd(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := rMemory.FindStringSubmatch(value)
		if len(match) != 3 {
			return nil, ErrorResponse
		}

		// Convert to uint64
		value2, err := strconv.ParseUint(match[2], 10, 64)
		if err != nil {
			return nil, err
		}

		// Set value
		memories[name] = value2
	}

	return memories, nil
}

func (this *State) GetOTP() (map[byte]uint32, error) {
	// Compile regular expression
	if rOTP == nil {
		rOTP = regexp.MustCompile(regexpOTP)
	}
	// retrieve OTP
	value, err := VCGenCmd(GENCMD_OTPDUMP)
	if err != nil {
		return nil, err
	}

	// find matches in the text
	matches := rOTP.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return nil, ErrorResponse
	}
	otp := make(map[byte]uint32, len(matches))
	for _, match := range matches {
		if len(match) != 3 {
			return nil, ErrorResponse
		}
		index, err := strconv.ParseUint(match[1], 10, 8)
		if err != nil {
			return nil, err
		}
		value, err := strconv.ParseUint(match[2], 16, 32)
		if err != nil {
			return nil, err
		}
		otp[byte(index)] = uint32(value)
	}

	return otp, nil
}

func (this *State) GetSerial() (uint64, error) {
	// Return cached version
	if this.serial != 0 {
		return this.serial, nil
	}

	// Get embedded memory
	otp, err := this.GetOTP()
	if err != nil {
		return 0, err
	}
	// Cache and return serial number
	this.serial = uint64(otp[GENCMD_OTPDUMP_SERIAL])
	return this.serial, nil
}

func (this *State) GetRevision() (uint32, error) {
	// Return cached version
	if this.revision != 0 {
		return this.revision, nil
	}

	// Get embedded memory
	otp, err := this.GetOTP()
	if err != nil {
		return 0, err
	}

	// Cache and return revision number
	this.revision = uint32(otp[GENCMD_OTPDUMP_REVISION])
	return this.revision, nil
}
