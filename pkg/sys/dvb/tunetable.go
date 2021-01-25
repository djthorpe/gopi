package dvb

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// CGO INTERFACE

/*
#include <sys/ioctl.h>
#include <linux/dvb/frontend.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type TuneParams struct {
	name  string
	flags scanflag

	DeliverySystem         FEDeliverySystem
	Frequency              uint32
	Modulation             FEModulation
	Bandwidth              uint32
	SymbolRate             uint32
	CodeRateLp, CodeRateHp FECodeRate
	InnerFEC               FECodeRate
	TransmitMode           FETransmitMode
	GuardInterval          FEGuardInterval
	Hierarchy              FEHierarchy
	Inversion              FEInversion
}

type scanflag int

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	scan_flag_none           scanflag = 0
	scan_flag_deliverysystem scanflag = (1 << iota)
	scan_flag_frequency
	scan_flag_modulation
	scan_flag_bandwidth
	scan_flag_symbolrate
	scan_flag_coderate
	scan_flag_innerfec
	scan_flag_transmitmode
	scan_flag_guardinterval
	scan_flag_hierarchy
	scan_flag_inversion
)

const (
	state_begin = iota
	state_comment
	state_channel
	state_kv
	state_error
)

var (
	reChannel  = regexp.MustCompile("^\\[(.*)\\]$")
	reKeyValue = regexp.MustCompile("^(\\w+)\\s*=\\s*(.*)$")
	rePidKey   = regexp.MustCompile("^PID_(\\w\\w+)$")
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

// NewTuneParams returns a new set of scan parameters
func NewTuneParams(name string) *TuneParams {
	this := new(TuneParams)
	this.name = name
	this.flags = scan_flag_none
	return this
}

// ReadTuneParamsTable returns tuning parameters from a file
// or handle
func ReadTuneParamsTable(r io.Reader) ([]*TuneParams, error) {
	reader := bufio.NewReader(r)
	pos := 1
	scans := []*TuneParams{}
	for {
		pos++
		if line, prefix, err := reader.ReadLine(); err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if prefix {
			return nil, fmt.Errorf("Line %v: Overflow", pos)
		} else if scans, err = parseScanline(scans, string(line)); err != nil {
			return nil, err
		}
	}
	return scans, nil
}

////////////////////////////////////////////////////////////////////////////////
// PROPERTIES

func (this *TuneParams) Name() string {
	return this.name
}

func (this *TuneParams) SetDeliverySystem(value FEDeliverySystem) {
	this.DeliverySystem = value
	this.flags |= scan_flag_deliverysystem
}

func (this *TuneParams) SetFrequency(value uint32) {
	this.Frequency = value
	this.flags |= scan_flag_frequency
}

func (this *TuneParams) SetModulation(value FEModulation) {
	this.Modulation = value
	this.flags |= scan_flag_modulation
}

func (this *TuneParams) SetBandwidth(value uint32) {
	this.Bandwidth = value
	this.flags |= scan_flag_bandwidth
}

func (this *TuneParams) SetSymbolrate(value uint32) {
	this.SymbolRate = value
	this.flags |= scan_flag_symbolrate
}

func (this *TuneParams) SetCoderate(lp, hp FECodeRate) {
	this.CodeRateLp = lp
	this.CodeRateHp = hp
	this.flags |= scan_flag_coderate
}

func (this *TuneParams) SetInnerFEC(value FECodeRate) {
	this.InnerFEC = value
	this.flags |= scan_flag_innerfec
}

func (this *TuneParams) SetTransmitMode(value FETransmitMode) {
	this.TransmitMode = value
	this.flags |= scan_flag_transmitmode
}

func (this *TuneParams) SetGuardInterval(value FEGuardInterval) {
	this.GuardInterval = value
	this.flags |= scan_flag_guardinterval
}

func (this *TuneParams) SetHierarchy(value FEHierarchy) {
	this.Hierarchy = value
	this.flags |= scan_flag_hierarchy
}

func (this *TuneParams) SetInversion(value FEInversion) {
	this.Inversion = value
	this.flags |= scan_flag_inversion
}

////////////////////////////////////////////////////////////////////////////////
// GET IOCTL PARAMETERS

func (this *TuneParams) params() []C.struct_dtv_property {
	params := []C.struct_dtv_property{}
	if this.flags&scan_flag_deliverysystem != 0 {
		params = append(params, propUint32(DTV_DELIVERY_SYSTEM, uint32(this.DeliverySystem)))
	}
	if this.flags&scan_flag_frequency != 0 {
		params = append(params, propUint32(DTV_FREQUENCY, uint32(this.Frequency)))
	}
	if this.flags&scan_flag_modulation != 0 {
		params = append(params, propUint32(DTV_MODULATION, uint32(this.Modulation)))
	}
	if this.flags&scan_flag_bandwidth != 0 {
		params = append(params, propUint32(DTV_BANDWIDTH_HZ, uint32(this.Bandwidth)))
	}
	if this.flags&scan_flag_symbolrate != 0 {
		params = append(params, propUint32(DTV_SYMBOL_RATE, uint32(this.SymbolRate)))
	}
	if this.flags&scan_flag_coderate != 0 {
		params = append(params, propUint32(DTV_CODE_RATE_LP, uint32(this.CodeRateLp)))
		params = append(params, propUint32(DTV_CODE_RATE_HP, uint32(this.CodeRateHp)))
	}
	if this.flags&scan_flag_innerfec != 0 {
		params = append(params, propUint32(DTV_INNER_FEC, uint32(this.InnerFEC)))
	}
	if this.flags&scan_flag_transmitmode != 0 {
		params = append(params, propUint32(DTV_TRANSMISSION_MODE, uint32(this.TransmitMode)))
	}
	if this.flags&scan_flag_guardinterval != 0 {
		params = append(params, propUint32(DTV_GUARD_INTERVAL, uint32(this.GuardInterval)))
	}
	if this.flags&scan_flag_hierarchy != 0 {
		params = append(params, propUint32(DTV_HIERARCHY, uint32(this.Hierarchy)))
	}
	if this.flags&scan_flag_inversion != 0 {
		params = append(params, propUint32(DTV_INVERSION, uint32(this.Inversion)))
	}
	return params
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *TuneParams) String() string {
	str := "<dvb.scan"
	str += " name=" + strconv.Quote(this.name)
	if this.flags&scan_flag_deliverysystem != 0 {
		str += " deliverysystem=" + fmt.Sprint(this.DeliverySystem)
	}
	if this.flags&scan_flag_frequency != 0 {
		str += " frequency=" + fmt.Sprint(this.Frequency)
	}
	if this.flags&scan_flag_modulation != 0 {
		str += " modulation=" + fmt.Sprint(this.Modulation)
	}
	if this.flags&scan_flag_bandwidth != 0 {
		str += " bandwidth=" + fmt.Sprint(this.Bandwidth)
	}
	if this.flags&scan_flag_symbolrate != 0 {
		str += " symbolrate=" + fmt.Sprint(this.SymbolRate)
	}
	if this.flags&scan_flag_coderate != 0 {
		str += " coderate_lp=" + fmt.Sprint(this.CodeRateLp)
		str += " coderate_hp=" + fmt.Sprint(this.CodeRateHp)
	}
	if this.flags&scan_flag_innerfec != 0 {
		str += " inner_fec=" + fmt.Sprint(this.InnerFEC)
	}
	if this.flags&scan_flag_transmitmode != 0 {
		str += " transmitmode=" + fmt.Sprint(this.TransmitMode)
	}
	if this.flags&scan_flag_guardinterval != 0 {
		str += " guardinterval=" + fmt.Sprint(this.GuardInterval)
	}
	if this.flags&scan_flag_hierarchy != 0 {
		str += " hierarchy=" + fmt.Sprint(this.Hierarchy)
	}
	if this.flags&scan_flag_inversion != 0 {
		str += " inversion=" + fmt.Sprint(this.Inversion)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func parseScanline(scans []*TuneParams, line string) ([]*TuneParams, error) {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, "#") {
		// Comment line
		return scans, nil
	} else if name := reChannel.FindStringSubmatch(line); len(name) > 0 {
		return append(scans, NewTuneParams(name[1])), nil
	} else if kv := reKeyValue.FindStringSubmatch(line); len(kv) > 0 && len(scans) > 0 {
		scan := scans[len(scans)-1]
		if err := parseScanKeyvalue(scan, kv[1], kv[2]); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("Parse error: %q", line)
	}

	// Return success
	return scans, nil
}

func parseScanKeyvalue(scan *TuneParams, key, value string) error {
	// Uppercase both key and value
	key = strings.ToUpper(key)
	value = strings.ToUpper(value)
	// Parse key
	switch key {
	case "DELIVERY_SYSTEM":
		if val, err := parseDeliverySystem(value); err != nil {
			return err
		} else {
			scan.SetDeliverySystem(val)
		}
	case "FREQUENCY":
		if val, err := strconv.ParseUint(value, 0, 32); err != nil {
			return err
		} else {
			scan.SetFrequency(uint32(val))
		}
	case "MODULATION":
		if val, err := parseModulation(value); err != nil {
			return err
		} else {
			scan.SetModulation(val)
		}
	case "BANDWIDTH_HZ", "BANDWIDTH":
		if val, err := strconv.ParseUint(value, 0, 32); err != nil {
			return err
		} else {
			scan.SetBandwidth(uint32(val))
		}
	case "SYMBOL_RATE":
		if val, err := strconv.ParseUint(value, 0, 32); err != nil {
			return err
		} else {
			scan.SetSymbolrate(uint32(val))
		}
	case "CODE_RATE_LP":
		if val, err := parseCoderate(value); err != nil {
			return err
		} else {
			scan.SetCoderate(val, scan.CodeRateHp)
		}
	case "CODE_RATE_HP":
		if val, err := parseCoderate(value); err != nil {
			return err
		} else {
			scan.SetCoderate(scan.CodeRateLp, val)
		}
	case "TRANSMISSION_MODE", "TRANSMIT_MODE":
		if val, err := parseTransmitmode(value); err != nil {
			return err
		} else {
			scan.SetTransmitMode(val)
		}
	case "INNER_FEC":
		if val, err := parseCoderate(value); err != nil {
			return err
		} else {
			scan.SetInnerFEC(val)
		}
	case "GUARD_INTERVAL":
		if val, err := parseGuardinterval(value); err != nil {
			return err
		} else {
			scan.SetGuardInterval(val)
		}
	case "HIERARCHY":
		if val, err := parseHierarchy(value); err != nil {
			return err
		} else {
			scan.SetHierarchy(val)
		}
	case "INVERSION":
		if val, err := parseInversion(value); err != nil {
			return err
		} else {
			scan.SetInversion(val)
		}
	case "STREAM_ID", "SERVICE_ID", "VIDEO_PID", "AUDIO_PID":
		// Ignore these parameters for now
	default:
		// Ignore PID value
		if rePidKey.MatchString(key) {
			return nil
		}
		return fmt.Errorf("Unsupported: %q", key)
	}

	// Return success
	return nil
}

func parseDeliverySystem(str string) (FEDeliverySystem, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, "/", "_"))
	for v := SYS_MIN; v <= SYS_MAX; v++ {
		if str == strings.TrimPrefix(fmt.Sprint(v), "SYS_") {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Unsupported: %q", str)
}

func parseModulation(str string) (FEModulation, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, "/", "_"))
	for v := MODULATION_MIN; v <= MODULATION_MAX; v++ {
		if str == fmt.Sprint(v) {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Unsupported: %q", str)
}

func parseCoderate(str string) (FECodeRate, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, "/", "_"))
	for v := CODERATE_MIN; v <= CODERATE_MAX; v++ {
		if str == strings.TrimPrefix(fmt.Sprint(v), "FEC_") {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Unsupported: %q", str)
}

func parseTransmitmode(str string) (FETransmitMode, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, "/", "_"))
	for v := TRANSMISSION_MODE_MIN; v <= TRANSMISSION_MODE_MAX; v++ {
		if str == strings.TrimPrefix(fmt.Sprint(v), "TRANSMISSION_MODE_") {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Unsupported: %q", str)
}

func parseGuardinterval(str string) (FEGuardInterval, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, "/", "_"))
	for v := GUARD_INTERVAL_MIN; v <= GUARD_INTERVAL_MAX; v++ {
		if str == strings.TrimPrefix(fmt.Sprint(v), "GUARD_INTERVAL_") {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Unsupported: %q", str)
}

func parseHierarchy(str string) (FEHierarchy, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, "/", "_"))
	for v := HIERARCHY_MIN; v <= HIERARCHY_MAX; v++ {
		if str == strings.TrimPrefix(fmt.Sprint(v), "HIERARCHY_") {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Unsupported: %q", str)
}

func parseInversion(str string) (FEInversion, error) {
	str = strings.ToUpper(strings.ReplaceAll(str, "/", "_"))
	for v := INVERSION_MIN; v <= INVERSION_MAX; v++ {
		if str == strings.TrimPrefix(fmt.Sprint(v), "INVERSION_") {
			return v, nil
		}
	}
	return 0, fmt.Errorf("Unsupported: %q", str)
}
