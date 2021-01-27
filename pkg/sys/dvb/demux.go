// +build dvb

package dvb

import (
	"fmt"
	"os"
	"strings"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO INTERFACE

/*
	#include <sys/ioctl.h>
	#include <linux/dvb/dmx.h>
	static int _DMX_START() { return DMX_START; }
	static int _DMX_STOP() { return DMX_STOP; }
	static int _DMX_SET_FILTER() { return DMX_SET_FILTER; }
	static int _DMX_SET_PES_FILTER() { return DMX_SET_PES_FILTER; }
	static int _DMX_SET_BUFFER_SIZE() { return DMX_SET_BUFFER_SIZE; }
	static int _DMX_GET_PES_PIDS() { return DMX_GET_PES_PIDS; }
	static int _DMX_GET_STC() { return DMX_GET_STC; }
	static int _DMX_ADD_PID() { return DMX_ADD_PID; }
	static int _DMX_REMOVE_PID() { return DMX_REMOVE_PID; }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	DMXInput         C.enum_dmx_input
	DMXOutput        C.enum_dmx_output
	DMXStreamType    C.enum_dmx_ts_pes
	DMXSectionFilter C.struct_dmx_sct_filter_params
	DMXStreamFilter  C.struct_dmx_pes_filter_params
	DMXFlag          uint32
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	DMX_OUT_DECODER     DMXOutput = C.DMX_OUT_DECODER     //  Streaming directly to decoder.
	DMX_OUT_TAP         DMXOutput = C.DMX_OUT_TAP         // Output going to a memory buffer
	DMX_OUT_TS_TAP      DMXOutput = C.DMX_OUT_TS_TAP      //  Output multiplexed into a new TS
	DMX_OUT_TSDEMUX_TAP DMXOutput = C.DMX_OUT_TSDEMUX_TAP // Like @DMX_OUT_TS_TAP but retrieved from the DMX device
)

const (
	DMX_IN_FRONTEND DMXInput = C.DMX_IN_FRONTEND // Input from a front-end device.
	DMX_IN_DVR      DMXInput = C.DMX_IN_DVR      // Input from the logical DVR device.
)

const (
	DMX_PES_AUDIO0    DMXStreamType = C.DMX_PES_AUDIO0    // first audio PID. Also referred as @DMX_PES_AUDIO.
	DMX_PES_VIDEO0    DMXStreamType = C.DMX_PES_VIDEO0    // first video PID. Also referred as @DMX_PES_VIDEO.
	DMX_PES_TELETEXT0 DMXStreamType = C.DMX_PES_TELETEXT0 // first teletext PID. Also referred as @DMX_PES_TELETEXT.
	DMX_PES_SUBTITLE0 DMXStreamType = C.DMX_PES_SUBTITLE0 // first subtitle PID. Also referred as @DMX_PES_SUBTITLE.
	DMX_PES_PCR0      DMXStreamType = C.DMX_PES_PCR0      // first Program Clock Reference PID.
	DMX_PES_AUDIO1    DMXStreamType = C.DMX_PES_AUDIO1    // second audio PID.
	DMX_PES_VIDEO1    DMXStreamType = C.DMX_PES_VIDEO1    // second video PID.
	DMX_PES_TELETEXT1 DMXStreamType = C.DMX_PES_TELETEXT1 // second teletext PID.
	DMX_PES_SUBTITLE1 DMXStreamType = C.DMX_PES_SUBTITLE1 // second subtitle PID.
	DMX_PES_PCR1      DMXStreamType = C.DMX_PES_PCR1      // second Program Clock Reference PID.
	DMX_PES_AUDIO2    DMXStreamType = C.DMX_PES_AUDIO2    // third audio PID.
	DMX_PES_VIDEO2    DMXStreamType = C.DMX_PES_VIDEO2    // third video PID.
	DMX_PES_TELETEXT2 DMXStreamType = C.DMX_PES_TELETEXT2 // third teletext PID.
	DMX_PES_SUBTITLE2 DMXStreamType = C.DMX_PES_SUBTITLE2 // third subtitle PID.
	DMX_PES_PCR2      DMXStreamType = C.DMX_PES_PCR2      // third Program Clock Reference PID.
	DMX_PES_AUDIO3    DMXStreamType = C.DMX_PES_AUDIO3    // fourth audio PID.
	DMX_PES_VIDEO3    DMXStreamType = C.DMX_PES_VIDEO3    // fourth video PID.
	DMX_PES_TELETEXT3 DMXStreamType = C.DMX_PES_TELETEXT3 // fourth teletext PID.
	DMX_PES_SUBTITLE3 DMXStreamType = C.DMX_PES_SUBTITLE3 // fourth subtitle PID.
	DMX_PES_PCR3      DMXStreamType = C.DMX_PES_PCR3      // fourth Program Clock Reference PID.
	DMX_PES_OTHER     DMXStreamType = C.DMX_PES_OTHER     // any other PID.
)

const (
	DMX_NONE            DMXFlag = 0
	DMX_CHECK_CRC       DMXFlag = C.DMX_CHECK_CRC
	DMX_ONESHOT         DMXFlag = C.DMX_ONESHOT
	DMX_IMMEDIATE_START DMXFlag = C.DMX_IMMEDIATE_START
	DMX_FLAGS_MIN               = DMX_CHECK_CRC
	DMX_FLAGS_MAX               = DMX_IMMEDIATE_START
)

var (
	DMX_START           = uintptr(C._DMX_START())
	DMX_STOP            = uintptr(C._DMX_STOP())
	DMX_SET_FILTER      = uintptr(C._DMX_SET_FILTER())
	DMX_SET_PES_FILTER  = uintptr(C._DMX_SET_PES_FILTER())
	DMX_SET_BUFFER_SIZE = uintptr(C._DMX_SET_BUFFER_SIZE())
	DMX_GET_PES_PIDS    = uintptr(C._DMX_GET_PES_PIDS())
	DMX_GET_STC         = uintptr(C._DMX_GET_STC())
	DMX_ADD_PID         = uintptr(C._DMX_ADD_PID())
	DMX_REMOVE_PID      = uintptr(C._DMX_REMOVE_PID())
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func DMXStart(fd uintptr) error {
	if err := dvb_ioctl(fd, DMX_START, unsafe.Pointer(nil)); err != 0 {
		return os.NewSyscallError("DMX_START", err)
	} else {
		return nil
	}
}

func DMXStop(fd uintptr) error {
	if err := dvb_ioctl(fd, DMX_STOP, unsafe.Pointer(nil)); err != 0 {
		return os.NewSyscallError("DMX_STOP", err)
	} else {
		return nil
	}
}

func DMXSetBufferSize(fd uintptr, size uint32) error {
	if err := dvb_ioctl(fd, DMX_SET_BUFFER_SIZE, unsafe.Pointer(uintptr(size))); err != 0 {
		return os.NewSyscallError("DMX_SET_BUFFER_SIZE", err)
	} else {
		return nil
	}
}

func DMXSetSectionFilter(fd uintptr, filter *DMXSectionFilter) error {
	if err := dvb_ioctl(fd, DMX_SET_FILTER, unsafe.Pointer(filter)); err != 0 {
		return os.NewSyscallError("DMX_SET_FILTER", err)
	} else {
		return nil
	}
}

func DMXSetStreamFilter(fd uintptr, filter *DMXStreamFilter) error {
	if err := dvb_ioctl(fd, DMX_SET_PES_FILTER, unsafe.Pointer(filter)); err != 0 {
		return os.NewSyscallError("DMX_SET_PES_FILTER", err)
	} else {
		return nil
	}
}

func DMXAddPid(fd uintptr, pid uint16) error {
	if err := dvb_ioctl(fd, DMX_ADD_PID, unsafe.Pointer(&pid)); err != 0 {
		return os.NewSyscallError("DMX_ADD_PID", err)
	} else {
		return nil
	}
}

func DMXRemovePid(fd uintptr, pid uint16) error {
	if err := dvb_ioctl(fd, DMX_REMOVE_PID, unsafe.Pointer(&pid)); err != 0 {
		return os.NewSyscallError("DMX_REMOVE_PID", err)
	} else {
		return nil
	}
}

func DMXGetStreamPids(fd uintptr) (map[DMXStreamType]uint16, error) {
	var pids [5]uint16
	if err := dvb_ioctl(fd, DMX_GET_PES_PIDS, unsafe.Pointer(&pids)); err != 0 {
		return nil, os.NewSyscallError("DMX_GET_PES_PIDS", err)
	}
	pidmap := make(map[DMXStreamType]uint16)
	for stream, pid := range pids {
		if pid != uint16(0xFFFF) {
			key := DMXStreamType(stream)
			pidmap[key] = pid
		}
	}
	return pidmap, nil
}

////////////////////////////////////////////////////////////////////////////////
// DMXSectionFilter

func NewSectionFilter(pid uint16, timeout uint32, flags DMXFlag) *DMXSectionFilter {
	filter := &DMXSectionFilter{
		pid:     C.__u16(pid),
		filter:  C.struct_dmx_filter{},
		timeout: C.__u32(timeout),
		flags:   C.__u32(flags),
	}
	return filter
}

func (f *DMXSectionFilter) Set(i int, tid, mask, mode uint8) {
	f.filter.filter[i] = C.__u8(tid)
	f.filter.mask[i] = C.__u8(mask)
	f.filter.mode[i] = C.__u8(mode)
}

func (f *DMXSectionFilter) String() string {
	str := "<dvb.dmx.sectionfilter"
	str += fmt.Sprintf(" pid=0x%04X", f.pid)
	str += fmt.Sprint(" filter=", f.filter)
	if f.timeout > 0 {
		str += fmt.Sprint(" timeout=", f.timeout)
	}
	if f.flags > 0 {
		str += fmt.Sprint(" flags=", DMXFlag(f.flags))
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// DMXStreamFilter

func NewStreamFilter(pid uint16, in DMXInput, out DMXOutput, stream DMXStreamType, flags DMXFlag) *DMXStreamFilter {
	filter := &DMXStreamFilter{
		pid:      C.__u16(pid),
		input:    C.enum_dmx_input(in),
		output:   C.enum_dmx_output(out),
		pes_type: C.enum_dmx_ts_pes(stream),
		flags:    C.__u32(flags),
	}
	return filter
}

func (f *DMXStreamFilter) String() string {
	str := "<dvb.dmx.streamfilter"
	str += fmt.Sprintf(" pid=0x%04X", f.pid)
	str += fmt.Sprint(" in=", DMXInput(f.input))
	str += fmt.Sprint(" out=", DMXOutput(f.output))
	str += fmt.Sprint(" stream=", DMXStreamType(f.pes_type))
	if f.flags > 0 {
		str += fmt.Sprint(" flags=", DMXFlag(f.flags))
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (f DMXOutput) String() string {
	switch f {
	case DMX_OUT_DECODER:
		return "DMX_OUT_DECODER"
	case DMX_OUT_TAP:
		return "DMX_OUT_TAP"
	case DMX_OUT_TS_TAP:
		return "DMX_OUT_TS_TAP"
	case DMX_OUT_TSDEMUX_TAP:
		return "DMX_OUT_TSDEMUX_TAP"
	default:
		return "[?? Invalid DMXOutput value]"
	}
}

func (f DMXInput) String() string {
	switch f {
	case DMX_IN_FRONTEND:
		return "DMX_IN_FRONTEND"
	case DMX_IN_DVR:
		return "DMX_IN_DVR"
	default:
		return "[?? Invalid DMXInput value]"
	}
}

func (f DMXStreamType) String() string {
	switch f {
	case DMX_PES_AUDIO0:
		return "DMX_PES_AUDIO0"
	case DMX_PES_VIDEO0:
		return "DMX_PES_VIDEO0"
	case DMX_PES_TELETEXT0:
		return "DMX_PES_TELETEXT0"
	case DMX_PES_SUBTITLE0:
		return "DMX_PES_SUBTITLE0"
	case DMX_PES_PCR0:
		return "DMX_PES_PCR0"
	case DMX_PES_AUDIO1:
		return "DMX_PES_AUDIO1"
	case DMX_PES_VIDEO1:
		return "DMX_PES_VIDEO1"
	case DMX_PES_TELETEXT1:
		return "DMX_PES_TELETEXT1"
	case DMX_PES_SUBTITLE1:
		return "DMX_PES_SUBTITLE1"
	case DMX_PES_PCR1:
		return "DMX_PES_PCR1"
	case DMX_PES_AUDIO2:
		return "DMX_PES_AUDIO2"
	case DMX_PES_VIDEO2:
		return "DMX_PES_VIDEO2"
	case DMX_PES_TELETEXT2:
		return "DMX_PES_TELETEXT2"
	case DMX_PES_SUBTITLE2:
		return "DMX_PES_SUBTITLE2"
	case DMX_PES_PCR2:
		return "DMX_PES_PCR2"
	case DMX_PES_AUDIO3:
		return "DMX_PES_AUDIO3"
	case DMX_PES_VIDEO3:
		return "DMX_PES_VIDEO3"
	case DMX_PES_TELETEXT3:
		return "DMX_PES_TELETEXT3"
	case DMX_PES_SUBTITLE3:
		return "DMX_PES_SUBTITLE3"
	case DMX_PES_PCR3:
		return "DMX_PES_PCR3"
	case DMX_PES_OTHER:
		return "DMX_PES_OTHER"
	default:
		return "[?? Invalid DMXStreamType value]"
	}
}

func (f DMXFlag) String() string {
	if f == DMX_NONE {
		return f.FlagString()
	}
	str := ""
	for v := DMX_FLAGS_MIN; v <= DMX_FLAGS_MAX; v <<= 1 {
		if f&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.Trim(str, "|")
}

func (f DMXFlag) FlagString() string {
	switch f {
	case DMX_NONE:
		return "DMX_NONE"
	case DMX_CHECK_CRC:
		return "DMX_CHECK_CRC"
	case DMX_ONESHOT:
		return "DMX_ONESHOT"
	case DMX_IMMEDIATE_START:
		return "DMX_IMMEDIATE_START"
	default:
		return "[?? Invalid DMXFlag value]"
	}
}
