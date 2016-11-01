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
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
    #cgo LDFLAGS: -L/opt/vc/lib -lbcm_host
    #include "bcm_host.h"
	#include "vc_vchi_gencmd.h"

	int vc_gencmd_wrapper(char* response,int maxlen,const char* command) {
		return vc_gencmd(response,maxlen,command);
	}
*/
import "C"

////////////////////////////////////////////////////////////////////////////////

const (
	GENCMD_BUFFER_SIZE      = 1024
	GENCMD_COMMANDS         = "commands"
	GENCMD_MEASURE_TEMP     = "measure_temp"
	GENCMD_MEASURE_CLOCK    = "measure_clock arm core h264 isp v3d uart pwm emmc pixel vec hdmi dpi"
	GENCMD_MEASURE_VOLTS    = "measure_volts core sdram_c sdram_i sdram_p"
	GENCMD_CODEC_ENABLED    = "codec_enabled H264 MPG2 WVC1 MPG4 MJPG WMV9 VP8"
	GENCMD_MEMORY           = "get_mem arm gpu"
	GENCMD_OTPDUMP          = "otp_dump"
	GENCMD_OTPDUMP_SERIAL   = 28
	GENCMD_OTPDUMP_REVISION = 30
)

////////////////////////////////////////////////////////////////////////////////

var (
	REGEXP_COMMANDS *regexp.Regexp = regexp.MustCompile("commands=\"([^\"]+)\"")
	REGEXP_TEMP     *regexp.Regexp = regexp.MustCompile("temp=(\\d+\\.?\\d*)")
	REGEXP_CLOCK    *regexp.Regexp = regexp.MustCompile("frequency\\((\\d+)\\)=(\\d+)")
	REGEXP_VOLTAGE  *regexp.Regexp = regexp.MustCompile("volt=(\\d*\\.?\\d*)V")
	REGEXP_CODEC    *regexp.Regexp = regexp.MustCompile("(\\w+)=(enabled|disabled)")
	REGEXP_MEMORY   *regexp.Regexp = regexp.MustCompile("(\\w+)=(\\d+)M")
	REGEXP_OTP      *regexp.Regexp = regexp.MustCompile("(\\d\\d):([0123456789abcdefABCDEF]{8})")
)

////////////////////////////////////////////////////////////////////////////////

// See http://elinux.org/RPI_vcgencmd_usage for some example usage
func (this *RaspberryPi) VCGenCmd(command string) (string, error) {
	ccommand := C.CString(command)
	defer C.free(unsafe.Pointer(ccommand))
	cbuffer := make([]byte, GENCMD_BUFFER_SIZE)
	if int(C.vc_gencmd_wrapper((*C.char)(unsafe.Pointer(&cbuffer[0])), C.int(GENCMD_BUFFER_SIZE), (*C.char)(unsafe.Pointer(ccommand)))) != 0 {
		return "", ErrorGenCmd
	}
	return string(cbuffer), nil
}

////////////////////////////////////////////////////////////////////////////////

// Get a list of commands that can be executed
func (this *RaspberryPi) GetCommands() ([]string, error) {
	// retrieve value as text
	value, err := this.VCGenCmd(GENCMD_COMMANDS)
	if err != nil {
		return []string{}, err
	}

	// Find values within text
	match := REGEXP_COMMANDS.FindStringSubmatch(value)
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

// Get the core temperature in celcius
func (this *RaspberryPi) GetCoreTemperatureCelcius() (float64, error) {
	// retrieve value as text
	value, err := this.VCGenCmd(GENCMD_MEASURE_TEMP)
	if err != nil {
		return 0.0, err
	}

	// Find value within text
	match := REGEXP_TEMP.FindStringSubmatch(value)
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

// Return clock frequencies of various components
func (this *RaspberryPi) GetClockFrequencyHertz() (map[string]uint64, error) {
	// retrieve values as text
	command := strings.Split(GENCMD_MEASURE_CLOCK, " ")
	clocks := make(map[string]uint64, len(command))
	for _, name := range command[1:] {

		// Retrieve clock value
		value, err := this.VCGenCmd(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := REGEXP_CLOCK.FindStringSubmatch(value)
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

// Return voltage of various components
func (this *RaspberryPi) GetVolts() (map[string]float64, error) {
	// retrieve values as text
	command := strings.Split(GENCMD_MEASURE_VOLTS, " ")
	volts := make(map[string]float64, len(command))
	for _, name := range command[1:] {

		// Retrieve volt value
		value, err := this.VCGenCmd(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := REGEXP_VOLTAGE.FindStringSubmatch(value)
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

// Return set of codecs supported and/or not supported
func (this *RaspberryPi) GetCodecs() (map[string]bool, error) {
	// retrieve values as text
	command := strings.Split(GENCMD_CODEC_ENABLED, " ")
	codecs := make(map[string]bool, len(command))
	for _, name := range command[1:] {

		// Retrieve volt value
		value, err := this.VCGenCmd(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := REGEXP_CODEC.FindStringSubmatch(value)
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

// Return core and GPU memory sizes
func (this *RaspberryPi) GetMemoryMegabytes() (map[string]uint64, error) {
	// retrieve values as text
	command := strings.Split(GENCMD_MEMORY, " ")
	memories := make(map[string]uint64, len(command))
	for _, name := range command[1:] {

		// Retrieve memory value
		value, err := this.VCGenCmd(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := REGEXP_MEMORY.FindStringSubmatch(value)
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
