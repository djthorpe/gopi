/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016
	All Rights Reserved

	For Licensing and Usage information, please see LICENSE.md
*/
package rpi /* import "github.com/djthorpe/gopi/device/rpi" */

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

import (
	gopi "github.com/djthorpe/gopi"
	util "github.com/djthorpe/gopi/util"
)

/*
    #cgo CFLAGS: -I/opt/vc/include -I/opt/vc/include/interface/vmcs_host -I/opt/vc/include/interface/vmcs_host/linux -I/opt/vc/include/interface/vcos/pthreads
    #cgo LDFLAGS: -L/opt/vc/lib -lbcm_host
    #include "bcm_host.h"
	#include "vc_vchi_gencmd.h"
	int vc_gencmd_wrap(char* response,int maxlen,const char* command) {
		return vc_gencmd(response,maxlen,command);
	}
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

// Configuration for the Raspberry Pi
type Hardware struct{}

type Device struct {
	log          *util.LoggerDevice // logger
	service      int                // service number
	serial       uint64
	revision     uint32
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	GENCMD_BUF_SIZE      = 1024
	GENCMD_SERVICE_NONE  = -1
	GENCMD_SERIAL_NONE   = 0
	GENCMD_REVISION_NONE = 0
)

// OTP (One Time Programmable) memory constants
const (
	GENCMD_OTP_DUMP          = "otp_dump"
	GENCMD_OTP_DUMP_SERIAL   = 28
	GENCMD_OTP_DUMP_REVISION = 30
	GENCMD_COMMANDS          = "commands"
	GENCMD_MEASURE_TEMP      = "measure_temp"
	GENCMD_MEASURE_CLOCK     = "measure_clock arm core h264 isp v3d uart pwm emmc pixel vec hdmi dpi"
	GENCMD_MEASURE_VOLTS     = "measure_volts core sdram_c sdram_i sdram_p"
	GENCMD_CODEC_ENABLED     = "codec_enabled H264 MPG2 WVC1 MPG4 MJPG WMV9 VP8"
	GENCMD_MEMORY            = "get_mem arm gpu"
)

////////////////////////////////////////////////////////////////////////////////
// GLOBAL VARIABLES

var (
	REGEXP_OTP_DUMP *regexp.Regexp = regexp.MustCompile("(\\d\\d):([0123456789abcdefABCDEF]{8})")
	REGEXP_TEMP     *regexp.Regexp = regexp.MustCompile("temp=(\\d+\\.?\\d*)")
	REGEXP_CLOCK    *regexp.Regexp = regexp.MustCompile("frequency\\((\\d+)\\)=(\\d+)")
	REGEXP_VOLTAGE  *regexp.Regexp = regexp.MustCompile("volt=(\\d*\\.?\\d*)V")
	REGEXP_CODEC    *regexp.Regexp = regexp.MustCompile("(\\w+)=(enabled|disabled)")
	REGEXP_MEMORY   *regexp.Regexp = regexp.MustCompile("(\\w+)=(\\d+)M")
	REGEXP_COMMANDS *regexp.Regexp = regexp.MustCompile("commands=\"([^\"]+)\"")
)

////////////////////////////////////////////////////////////////////////////////
// Open and close device

func (config Hardware) Open(log *util.LoggerDevice) (gopi.Driver, error) {
	log.Debug2("<rpi.Device>Open")
	if err := bcmHostInit(); err != nil {
		return nil, err
	}

	this := new(Device)
	this.log = log
	this.service = GENCMD_SERVICE_NONE
	this.serial = GENCMD_SERIAL_NONE
	this.revision = GENCMD_REVISION_NONE

	return this, nil
}

func (this *Device) Close() error {
	this.log.Debug2("<rpi.Device>Close")
	if this.service != GENCMD_SERVICE_NONE {
		if err := vcGencmdTerminate(); err != nil {
			bcmHostTerminate()
			return err
		}
	}
	if err := bcmHostTerminate(); err != nil {
		return err
	}
	return nil
}

func (this *Device) String() string {
	serial, _ := this.GetSerialNumber()
	revision, _ := this.GetRevision()
	model, pcb, _ := this.GetModel()
	processor, _ := this.GetProcessor()
	warranty_bit, _ := this.GetWarrantyBit()
	return fmt.Sprintf("<rpi.Device>{ serial_number=%08X revision=%04X model=%v pcb=%v processor=%v warranty_bit=%v }", serial, revision, model, pcb, processor, warranty_bit)
}

////////////////////////////////////////////////////////////////////////////////
// Get Device Information

func (this *Device) GetPeripheralAddress() uint32 {
	return bcmHostGetPeripheralAddress()
}

func (this *Device) GetPeripheralSize() uint32 {
	return bcmHostGetPeripheralSize()
}

// Return the 64-bit serial number for the device
func (this *Device) GetSerialNumber() (uint64, error) {
	// Return cached version
	if this.serial != GENCMD_SERIAL_NONE {
		return this.serial, nil
	}

	// Get embedded memory
	otp, err := this.GetOTP()
	if err != nil {
		return GENCMD_SERIAL_NONE, err
	}
	// Cache and return serial number
	this.serial = uint64(otp[GENCMD_OTP_DUMP_SERIAL])
	return this.serial, nil
}

// Return the 32-bit revision code for the device
func (this *Device) GetRevision() (uint32, error) {
	// Return cached version
	if this.revision != GENCMD_REVISION_NONE {
		return this.revision, nil
	}

	// Get embedded memory
	otp, err := this.GetOTP()
	if err != nil {
		return GENCMD_REVISION_NONE, err
	}

	// Cache and return revision number
	this.revision = uint32(otp[GENCMD_OTP_DUMP_REVISION])
	return this.revision, nil
}

// Return the size of a particular display
func (this *Device) GetDisplaySize(display uint16) (uint32, uint32) {
	return bcmGHostGetDisplaySize(display)
}

////////////////////////////////////////////////////////////////////////////////
// General Command Interface

// Execute a VideoCore "General Command" and return the results of
// that command. See http://elinux.org/RPI_vcgencmd_usage for some example
// usage
func (this *Device) GeneralCommand(command string) (string, error) {
	if this.service == GENCMD_SERVICE_NONE {
		var err error
		this.service, err = vcGencmdInit(this.log)
		if err != nil {
			this.service = GENCMD_SERVICE_NONE
			return "", err
		}
	}
	ccommand := C.CString(command)
	defer C.free(unsafe.Pointer(ccommand))
	cbuffer := make([]byte, GENCMD_BUF_SIZE)
	if int(C.vc_gencmd_wrap((*C.char)(unsafe.Pointer(&cbuffer[0])), C.int(GENCMD_BUF_SIZE), (*C.char)(unsafe.Pointer(ccommand)))) != 0 {
		return "", this.log.Error("General Command Error")
	}
	return string(cbuffer), nil
}

// Return OTP memory
func (this *Device) GetOTP() (map[byte]uint32, error) {
	// retrieve OTP
	value, err := this.GeneralCommand(GENCMD_OTP_DUMP)
	if err != nil {
		return nil, err
	}

	// find matches in the text
	matches := REGEXP_OTP_DUMP.FindAllStringSubmatch(value, -1)
	if len(matches) == 0 {
		return nil, this.log.Error("Bad Response from %v", GENCMD_OTP_DUMP)
	}
	otp := make(map[byte]uint32, len(matches))
	for _, match := range matches {
		if len(match) != 3 {
			return nil, this.log.Error("Bad Response from %v", GENCMD_OTP_DUMP)
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

// Get the core temperature in celcius
func (this *Device) GetCoreTemperatureCelcius() (float64, error) {
	// retrieve value as text
	value, err := this.GeneralCommand(GENCMD_MEASURE_TEMP)
	if err != nil {
		return 0.0, err
	}

	// Find value within text
	match := REGEXP_TEMP.FindStringSubmatch(value)
	if len(match) != 2 {
		return 0.0, this.log.Error("Bad Response from %v", GENCMD_MEASURE_TEMP)
	}

	// Convert to float64
	value2, err := strconv.ParseFloat(match[1], 64)
	if err != nil {
		return 0.0, err
	}

	// Return value as float64
	return value2, nil
}

// Get Commands
func (this *Device) GetCommands() ([]string, error) {
	// retrieve value as text
	value, err := this.GeneralCommand(GENCMD_COMMANDS)
	if err != nil {
		return []string{}, err
	}

	// Find values within text
	match := REGEXP_COMMANDS.FindStringSubmatch(value)
	if len(match) != 2 {
		return []string{}, this.log.Error("Bad Response from %v", GENCMD_COMMANDS)
	}

	// Split commands
	commands := strings.Split(match[1], ",")
	for i, command := range commands {
		commands[i] = strings.TrimSpace(command)
	}
	// return commands
	return commands, nil
}

// Return clock frequencies of various components
func (this *Device) GetClockFrequencyHertz() (map[string]uint64, error) {
	// retrieve values as text
	command := strings.Split(GENCMD_MEASURE_CLOCK, " ")
	clocks := make(map[string]uint64, len(command))
	for _, name := range command[1:] {

		// Retrieve clock value
		value, err := this.GeneralCommand(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := REGEXP_CLOCK.FindStringSubmatch(value)
		if len(match) != 3 {
			return nil, this.log.Error("Bad Response from %v", command[0])
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
func (this *Device) GetVolts() (map[string]float64, error) {
	// retrieve values as text
	command := strings.Split(GENCMD_MEASURE_VOLTS, " ")
	volts := make(map[string]float64, len(command))
	for _, name := range command[1:] {

		// Retrieve volt value
		value, err := this.GeneralCommand(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := REGEXP_VOLTAGE.FindStringSubmatch(value)
		if len(match) != 2 {
			return nil, this.log.Error("Bad Response from %v", command[0])
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
func (this *Device) GetCodecs() (map[string]bool, error) {
	// retrieve values as text
	command := strings.Split(GENCMD_CODEC_ENABLED, " ")
	codecs := make(map[string]bool, len(command))
	for _, name := range command[1:] {

		// Retrieve volt value
		value, err := this.GeneralCommand(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := REGEXP_CODEC.FindStringSubmatch(value)
		if len(match) != 3 {
			return nil, this.log.Error("Bad Response from %v", command[0])
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
func (this *Device) GetMemoryMegabytes() (map[string]uint64, error) {
	// retrieve values as text
	command := strings.Split(GENCMD_MEMORY, " ")
	memories := make(map[string]uint64, len(command))
	for _, name := range command[1:] {

		// Retrieve memory value
		value, err := this.GeneralCommand(command[0] + " " + name)
		if err != nil {
			return nil, err
		}

		// Find value within text
		match := REGEXP_MEMORY.FindStringSubmatch(value)
		if len(match) != 3 {
			return nil, this.log.Error("Bad Response from %v", command[0])
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

////////////////////////////////////////////////////////////////////////////////
// Private methods

func bcmHostInit() error {
	C.bcm_host_init()
	return nil
}

func bcmHostTerminate() error {
	C.bcm_host_deinit()
	return nil
}

func vcGencmdInit(log *util.LoggerDevice) (int, error) {
	service := int(C.vc_gencmd_init())
	if service < 0 {
		return -1, log.Error("vc_gencmd_init failed")
	}
	return service, nil
}

func vcGencmdTerminate() error {
	C.vc_gencmd_stop()
	return nil
}

func bcmHostGetPeripheralAddress() uint32 {
	return uint32(C.bcm_host_get_peripheral_address())
}

func bcmHostGetPeripheralSize() uint32 {
	return uint32(C.bcm_host_get_peripheral_size())
}

func bcmHostGetSDRAMAddress() uint32 {
	return uint32(C.bcm_host_get_sdram_address())
}

func bcmGHostGetDisplaySize(display uint16) (uint32, uint32) {
	var w, h uint32
	C.graphics_get_display_size((C.uint16_t)(display), (*C.uint32_t)(&w), (*C.uint32_t)(&h))
	return w, h
}
