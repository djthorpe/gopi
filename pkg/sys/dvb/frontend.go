// +build dvb

package dvb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unsafe"
)

////////////////////////////////////////////////////////////////////////////////
// CGO INTERFACE

/*
#include <sys/ioctl.h>
#include <linux/dvb/frontend.h>
static int _FE_GET_INFO() { return FE_GET_INFO; }
static int _FE_DISEQC_RESET_OVERLOAD() { return FE_DISEQC_RESET_OVERLOAD; }
static int _FE_DISEQC_SEND_MASTER_CMD() { return FE_DISEQC_SEND_MASTER_CMD; }
static int _FE_DISEQC_RECV_SLAVE_REPLY() { return FE_DISEQC_RECV_SLAVE_REPLY; }
static int _FE_DISEQC_SEND_BURST() { return FE_DISEQC_SEND_BURST; }
static int _FE_SET_TONE() { return FE_SET_TONE; }
static int _FE_SET_VOLTAGE() { return FE_SET_VOLTAGE; }
static int _FE_ENABLE_HIGH_LNB_VOLTAGE() { return FE_ENABLE_HIGH_LNB_VOLTAGE; }
static int _FE_READ_STATUS() { return FE_READ_STATUS; }
static int _FE_READ_BER() { return FE_READ_BER; }
static int _FE_READ_SIGNAL_STRENGTH() { return FE_READ_SIGNAL_STRENGTH; }
static int _FE_READ_SNR() { return FE_READ_SNR; }
static int _FE_READ_UNCORRECTED_BLOCKS() { return FE_READ_UNCORRECTED_BLOCKS; }
static int _FE_SET_FRONTEND_TUNE_MODE() { return FE_SET_FRONTEND_TUNE_MODE; }
static int _FE_GET_EVENT() { return FE_GET_EVENT; }
static int _FE_DISHNETWORK_SEND_LEGACY_CMD() { return FE_DISHNETWORK_SEND_LEGACY_CMD; }
static int _FE_SET_PROPERTY() { return FE_SET_PROPERTY; }
static int _FE_GET_PROPERTY() { return FE_GET_PROPERTY; }
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	FEInfo           C.struct_dvb_frontend_info
	FECaps           uint64 // C.enum_fe_caps
	FEStatus         C.enum_fe_status
	FEKey            uint
	FEDeliverySystem C.enum_fe_delivery_system
	FEModulation     C.enum_fe_modulation
	FEInversion      C.enum_fe_spectral_inversion
	FECodeRate       C.enum_fe_code_rate
	FETransmitMode   C.enum_fe_transmit_mode
	FEGuardInterval  C.enum_fe_guard_interval
	FEHierarchy      C.enum_fe_hierarchy
	FEInterleaving   C.enum_fe_interleaving
	FEPropUint32     C.struct_dtv_property
)

type fePropUint32 struct {
	Key      FEKey
	reserved [3]uint32
	Data     uint32
}

type fePropEnum struct {
	Key      FEKey
	reserved [3]uint32
	Data     [32]uint8
	Len      uint32
}

type fePropStats struct {
	Key      FEKey
	reserved [3]uint32
	Len      uint8
	Data     [9 * 4]byte
}

type FEStats struct {
	Key    FEKey
	Values []FEStat
}

type FEStat struct {
	Scale uint8
	Value int64
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

var (
	FE_GET_INFO     = uintptr(C._FE_GET_INFO())
	FE_READ_STATUS  = uintptr(C._FE_READ_STATUS())
	FE_GET_PROPERTY = uintptr(C._FE_GET_PROPERTY())
	FE_SET_PROPERTY = uintptr(C._FE_SET_PROPERTY())
)

const (
	FE_IS_STUPID                  FECaps = C.FE_IS_STUPID
	FE_CAN_INVERSION_AUTO         FECaps = C.FE_CAN_INVERSION_AUTO
	FE_CAN_FEC_1_2                FECaps = C.FE_CAN_FEC_1_2
	FE_CAN_FEC_2_3                FECaps = C.FE_CAN_FEC_2_3
	FE_CAN_FEC_3_4                FECaps = C.FE_CAN_FEC_3_4
	FE_CAN_FEC_4_5                FECaps = C.FE_CAN_FEC_4_5
	FE_CAN_FEC_5_6                FECaps = C.FE_CAN_FEC_5_6
	FE_CAN_FEC_6_7                FECaps = C.FE_CAN_FEC_6_7
	FE_CAN_FEC_7_8                FECaps = C.FE_CAN_FEC_7_8
	FE_CAN_FEC_8_9                FECaps = C.FE_CAN_FEC_8_9
	FE_CAN_FEC_AUTO               FECaps = C.FE_CAN_FEC_AUTO
	FE_CAN_QPSK                   FECaps = C.FE_CAN_QPSK
	FE_CAN_QAM_16                 FECaps = C.FE_CAN_QAM_16
	FE_CAN_QAM_32                 FECaps = C.FE_CAN_QAM_32
	FE_CAN_QAM_64                 FECaps = C.FE_CAN_QAM_64
	FE_CAN_QAM_128                FECaps = C.FE_CAN_QAM_128
	FE_CAN_QAM_256                FECaps = C.FE_CAN_QAM_256
	FE_CAN_QAM_AUTO               FECaps = C.FE_CAN_QAM_AUTO
	FE_CAN_TRANSMISSION_MODE_AUTO FECaps = C.FE_CAN_TRANSMISSION_MODE_AUTO
	FE_CAN_BANDWIDTH_AUTO         FECaps = C.FE_CAN_BANDWIDTH_AUTO
	FE_CAN_GUARD_INTERVAL_AUTO    FECaps = C.FE_CAN_GUARD_INTERVAL_AUTO
	FE_CAN_HIERARCHY_AUTO         FECaps = C.FE_CAN_HIERARCHY_AUTO
	FE_CAN_8VSB                   FECaps = C.FE_CAN_8VSB
	FE_CAN_16VSB                  FECaps = C.FE_CAN_16VSB
	FE_HAS_EXTENDED_CAPS          FECaps = C.FE_HAS_EXTENDED_CAPS
	FE_CAN_MULTISTREAM            FECaps = C.FE_CAN_MULTISTREAM
	FE_CAN_TURBO_FEC              FECaps = C.FE_CAN_TURBO_FEC
	FE_CAN_2G_MODULATION          FECaps = C.FE_CAN_2G_MODULATION
	FE_NEEDS_BENDING              FECaps = C.FE_NEEDS_BENDING
	FE_CAN_RECOVER                FECaps = C.FE_CAN_RECOVER
	FE_CAN_MUTE_TS                FECaps = C.FE_CAN_MUTE_TS
	FE_CAPS_MIN                          = FE_CAN_INVERSION_AUTO
	FE_CAPS_MAX                          = FE_CAN_MUTE_TS
)

const (
	FE_NONE        FEStatus = C.FE_NONE
	FE_HAS_SIGNAL  FEStatus = C.FE_HAS_SIGNAL
	FE_HAS_CARRIER FEStatus = C.FE_HAS_CARRIER
	FE_HAS_VITERBI FEStatus = C.FE_HAS_VITERBI
	FE_HAS_SYNC    FEStatus = C.FE_HAS_SYNC
	FE_HAS_LOCK    FEStatus = C.FE_HAS_LOCK
	FE_TIMEDOUT    FEStatus = C.FE_TIMEDOUT
	FE_REINIT      FEStatus = C.FE_REINIT
	FE_STATUS_MIN           = FE_HAS_SIGNAL
	FE_STATUS_MAX           = FE_REINIT
)

const (
	DTV_UNDEFINED                      FEKey = C.DTV_UNDEFINED /* DVB Property Commands */
	DTV_TUNE                           FEKey = C.DTV_TUNE
	DTV_CLEAR                          FEKey = C.DTV_CLEAR
	DTV_FREQUENCY                      FEKey = C.DTV_FREQUENCY
	DTV_MODULATION                     FEKey = C.DTV_MODULATION
	DTV_BANDWIDTH_HZ                   FEKey = C.DTV_BANDWIDTH_HZ
	DTV_INVERSION                      FEKey = C.DTV_INVERSION
	DTV_DISEQC_MASTER                  FEKey = C.DTV_DISEQC_MASTER
	DTV_SYMBOL_RATE                    FEKey = C.DTV_SYMBOL_RATE
	DTV_INNER_FEC                      FEKey = C.DTV_INNER_FEC
	DTV_VOLTAGE                        FEKey = C.DTV_VOLTAGE
	DTV_TONE                           FEKey = C.DTV_TONE
	DTV_PILOT                          FEKey = C.DTV_PILOT
	DTV_ROLLOFF                        FEKey = C.DTV_ROLLOFF
	DTV_DISEQC_SLAVE_REPLY             FEKey = C.DTV_DISEQC_SLAVE_REPLY
	DTV_FE_CAPABILITY_COUNT            FEKey = C.DTV_FE_CAPABILITY_COUNT /* Basic enumeration set for querying unlimited capabilities */
	DTV_FE_CAPABILITY                  FEKey = C.DTV_FE_CAPABILITY
	DTV_DELIVERY_SYSTEM                FEKey = C.DTV_DELIVERY_SYSTEM
	DTV_ISDBT_PARTIAL_RECEPTION        FEKey = C.DTV_ISDBT_PARTIAL_RECEPTION /* ISDB-T and ISDB-Tsb */
	DTV_ISDBT_SOUND_BROADCASTING       FEKey = C.DTV_ISDBT_SOUND_BROADCASTING
	DTV_ISDBT_SB_SUBCHANNEL_ID         FEKey = C.DTV_ISDBT_SB_SUBCHANNEL_ID
	DTV_ISDBT_SB_SEGMENT_IDX           FEKey = C.DTV_ISDBT_SB_SEGMENT_IDX
	DTV_ISDBT_SB_SEGMENT_COUNT         FEKey = C.DTV_ISDBT_SB_SEGMENT_COUNT
	DTV_ISDBT_LAYERA_FEC               FEKey = C.DTV_ISDBT_LAYERA_FEC
	DTV_ISDBT_LAYERA_MODULATION        FEKey = C.DTV_ISDBT_LAYERA_MODULATION
	DTV_ISDBT_LAYERA_SEGMENT_COUNT     FEKey = C.DTV_ISDBT_LAYERA_SEGMENT_COUNT
	DTV_ISDBT_LAYERA_TIME_INTERLEAVING FEKey = C.DTV_ISDBT_LAYERA_TIME_INTERLEAVING
	DTV_ISDBT_LAYERB_FEC               FEKey = C.DTV_ISDBT_LAYERB_FEC
	DTV_ISDBT_LAYERB_MODULATION        FEKey = C.DTV_ISDBT_LAYERB_MODULATION
	DTV_ISDBT_LAYERB_SEGMENT_COUNT     FEKey = C.DTV_ISDBT_LAYERB_SEGMENT_COUNT
	DTV_ISDBT_LAYERB_TIME_INTERLEAVING FEKey = C.DTV_ISDBT_LAYERB_TIME_INTERLEAVING
	DTV_ISDBT_LAYERC_FEC               FEKey = C.DTV_ISDBT_LAYERC_FEC
	DTV_ISDBT_LAYERC_MODULATION        FEKey = C.DTV_ISDBT_LAYERC_MODULATION
	DTV_ISDBT_LAYERC_SEGMENT_COUNT     FEKey = C.DTV_ISDBT_LAYERC_SEGMENT_COUNT
	DTV_ISDBT_LAYERC_TIME_INTERLEAVING FEKey = C.DTV_ISDBT_LAYERC_TIME_INTERLEAVING
	DTV_API_VERSION                    FEKey = C.DTV_API_VERSION
	DTV_CODE_RATE_HP                   FEKey = C.DTV_CODE_RATE_HP
	DTV_CODE_RATE_LP                   FEKey = C.DTV_CODE_RATE_LP
	DTV_GUARD_INTERVAL                 FEKey = C.DTV_GUARD_INTERVAL
	DTV_TRANSMISSION_MODE              FEKey = C.DTV_TRANSMISSION_MODE
	DTV_HIERARCHY                      FEKey = C.DTV_HIERARCHY
	DTV_ISDBT_LAYER_ENABLED            FEKey = C.DTV_ISDBT_LAYER_ENABLED
	DTV_STREAM_ID                      FEKey = C.DTV_STREAM_ID
	DTV_ISDBS_TS_ID_LEGACY             FEKey = C.DTV_ISDBS_TS_ID_LEGACY
	DTV_DVBT2_PLP_ID_LEGACY            FEKey = C.DTV_DVBT2_PLP_ID_LEGACY
	DTV_ENUM_DELSYS                    FEKey = C.DTV_ENUM_DELSYS
	DTV_ATSCMH_FIC_VER                 FEKey = C.DTV_ATSCMH_FIC_VER /* ATSC-MH */
	DTV_ATSCMH_PARADE_ID               FEKey = C.DTV_ATSCMH_PARADE_ID
	DTV_ATSCMH_NOG                     FEKey = C.DTV_ATSCMH_NOG
	DTV_ATSCMH_TNOG                    FEKey = C.DTV_ATSCMH_TNOG
	DTV_ATSCMH_SGN                     FEKey = C.DTV_ATSCMH_SGN
	DTV_ATSCMH_PRC                     FEKey = C.DTV_ATSCMH_PRC
	DTV_ATSCMH_RS_FRAME_MODE           FEKey = C.DTV_ATSCMH_RS_FRAME_MODE
	DTV_ATSCMH_RS_FRAME_ENSEMBLE       FEKey = C.DTV_ATSCMH_RS_FRAME_ENSEMBLE
	DTV_ATSCMH_RS_CODE_MODE_PRI        FEKey = C.DTV_ATSCMH_RS_CODE_MODE_PRI
	DTV_ATSCMH_RS_CODE_MODE_SEC        FEKey = C.DTV_ATSCMH_RS_CODE_MODE_SEC
	DTV_ATSCMH_SCCC_BLOCK_MODE         FEKey = C.DTV_ATSCMH_SCCC_BLOCK_MODE
	DTV_ATSCMH_SCCC_CODE_MODE_A        FEKey = C.DTV_ATSCMH_SCCC_CODE_MODE_A
	DTV_ATSCMH_SCCC_CODE_MODE_B        FEKey = C.DTV_ATSCMH_SCCC_CODE_MODE_B
	DTV_ATSCMH_SCCC_CODE_MODE_C        FEKey = C.DTV_ATSCMH_SCCC_CODE_MODE_C
	DTV_ATSCMH_SCCC_CODE_MODE_D        FEKey = C.DTV_ATSCMH_SCCC_CODE_MODE_D
	DTV_INTERLEAVING                   FEKey = C.DTV_INTERLEAVING
	DTV_LNA                            FEKey = C.DTV_LNA
	DTV_STAT_SIGNAL_STRENGTH           FEKey = C.DTV_STAT_SIGNAL_STRENGTH /* Quality parameters */
	DTV_STAT_CNR                       FEKey = C.DTV_STAT_CNR
	DTV_STAT_PRE_ERROR_BIT_COUNT       FEKey = C.DTV_STAT_PRE_ERROR_BIT_COUNT
	DTV_STAT_PRE_TOTAL_BIT_COUNT       FEKey = C.DTV_STAT_PRE_TOTAL_BIT_COUNT
	DTV_STAT_POST_ERROR_BIT_COUNT      FEKey = C.DTV_STAT_POST_ERROR_BIT_COUNT
	DTV_STAT_POST_TOTAL_BIT_COUNT      FEKey = C.DTV_STAT_POST_TOTAL_BIT_COUNT
	DTV_STAT_ERROR_BLOCK_COUNT         FEKey = C.DTV_STAT_ERROR_BLOCK_COUNT
	DTV_STAT_TOTAL_BLOCK_COUNT         FEKey = C.DTV_STAT_TOTAL_BLOCK_COUNT
	DTV_SCRAMBLING_SEQUENCE_INDEX      FEKey = C.DTV_SCRAMBLING_SEQUENCE_INDEX /* Physical layer scrambling */
)

const (
	SYS_UNDEFINED    FEDeliverySystem = C.SYS_UNDEFINED    // Undefined standard. Generally, indicates an error
	SYS_DVBC_ANNEX_A FEDeliverySystem = C.SYS_DVBC_ANNEX_A // Cable TV: DVB-C following ITU-T J.83 Annex A spec
	SYS_DVBC_ANNEX_B FEDeliverySystem = C.SYS_DVBC_ANNEX_B // Cable TV: DVB-C following ITU-T J.83 Annex B spec (ClearQAM)
	SYS_DVBC_ANNEX_C FEDeliverySystem = C.SYS_DVBC_ANNEX_C // Cable TV: DVB-C following ITU-T J.83 Annex C spec
	SYS_ISDBC        FEDeliverySystem = C.SYS_ISDBC        // Cable TV: ISDB-C (no drivers yet)
	SYS_DVBT         FEDeliverySystem = C.SYS_DVBT         // Terrestrial TV: DVB-T
	SYS_DVBT2        FEDeliverySystem = C.SYS_DVBT2        // Terrestrial TV: DVB-T2
	SYS_ISDBT        FEDeliverySystem = C.SYS_ISDBT        // Terrestrial TV: ISDB-T
	SYS_ATSC         FEDeliverySystem = C.SYS_ATSC         // Terrestrial TV: ATSC
	SYS_ATSCMH       FEDeliverySystem = C.SYS_ATSCMH       // Terrestrial TV (mobile): ATSC-M/H
	SYS_DTMB         FEDeliverySystem = C.SYS_DTMB         // Terrestrial TV: DTMB
	SYS_DVBS         FEDeliverySystem = C.SYS_DVBS         // Satellite TV: DVB-S
	SYS_DVBS2        FEDeliverySystem = C.SYS_DVBS2        // Satellite TV: DVB-S2
	SYS_TURBO        FEDeliverySystem = C.SYS_TURBO        // Satellite TV: DVB-S Turbo
	SYS_ISDBS        FEDeliverySystem = C.SYS_ISDBS        // Satellite TV: ISDB-S
	SYS_DAB          FEDeliverySystem = C.SYS_DAB          // Digital audio: DAB (not fully supported)
	SYS_DSS          FEDeliverySystem = C.SYS_DSS          // Satellite TV: DSS (not fully supported)
	SYS_CMMB         FEDeliverySystem = C.SYS_CMMB         // Terrestrial TV (mobile): CMMB (not fully supported)
	SYS_DVBH         FEDeliverySystem = C.SYS_DVBH         // Terrestrial TV (mobile): DVB-H (standard deprecated)
	SYS_MIN                           = SYS_DVBC_ANNEX_A
	SYS_MAX                           = SYS_DVBC_ANNEX_C
)

const (
	QPSK           FEModulation = C.QPSK     // QPSK modulation
	QAM_16         FEModulation = C.QAM_16   // 16-QAM modulation
	QAM_32         FEModulation = C.QAM_32   // 32-QAM modulation
	QAM_64         FEModulation = C.QAM_64   // 64-QAM modulation
	QAM_128        FEModulation = C.QAM_128  // 128-QAM modulation
	QAM_256        FEModulation = C.QAM_256  // 256-QAM modulation
	QAM_AUTO       FEModulation = C.QAM_AUTO // Autodetect QAM modulation
	VSB_8          FEModulation = C.VSB_8    // 8-VSB modulation
	VSB_16         FEModulation = C.VSB_16   // 16-VSB modulation
	PSK_8          FEModulation = C.PSK_8    // 8-PSK modulation
	APSK_16        FEModulation = C.APSK_16  // 16-APSK modulation
	APSK_32        FEModulation = C.APSK_32  // 32-APSK modulation
	DQPSK          FEModulation = C.DQPSK    // DQPSK modulation
	QAM_4_NR       FEModulation = C.QAM_4_NR // 4-QAM-NR modulation
	MODULATION_MIN              = QPSK
	MODULATION_MAX              = QAM_4_NR
)

const (
	FEC_NONE     FECodeRate = C.FEC_NONE // No Forward Error Correction Code
	FEC_1_2      FECodeRate = C.FEC_1_2  // Forward Error Correction Code 1/2
	FEC_2_3      FECodeRate = C.FEC_2_3  // Forward Error Correction Code 2/3
	FEC_3_4      FECodeRate = C.FEC_3_4  // Forward Error Correction Code 3/4
	FEC_4_5      FECodeRate = C.FEC_4_5  // Forward Error Correction Code 4/5
	FEC_5_6      FECodeRate = C.FEC_5_6  // Forward Error Correction Code 5/6
	FEC_6_7      FECodeRate = C.FEC_6_7  // Forward Error Correction Code 6/7
	FEC_7_8      FECodeRate = C.FEC_7_8  // Forward Error Correction Code 7/8
	FEC_8_9      FECodeRate = C.FEC_8_9  // Forward Error Correction Code 8/9
	FEC_AUTO     FECodeRate = C.FEC_AUTO // Autodetect Error Correction Code
	FEC_3_5      FECodeRate = C.FEC_3_5  // Forward Error Correction Code 3/5
	FEC_9_10     FECodeRate = C.FEC_9_10 // Forward Error Correction Code 9/10
	FEC_2_5      FECodeRate = C.FEC_2_5  // Forward Error Correction Code 2/5
	CODERATE_MIN            = FEC_NONE
	CODERATE_MAX            = FEC_2_5
)

const (
	TRANSMISSION_MODE_AUTO  FETransmitMode = C.TRANSMISSION_MODE_AUTO  // Autodetect transmission mode
	TRANSMISSION_MODE_1K    FETransmitMode = C.TRANSMISSION_MODE_1K    // Transmission mode 1K
	TRANSMISSION_MODE_2K    FETransmitMode = C.TRANSMISSION_MODE_2K    // Transmission mode 2K
	TRANSMISSION_MODE_8K    FETransmitMode = C.TRANSMISSION_MODE_8K    // Transmission mode 8K
	TRANSMISSION_MODE_4K    FETransmitMode = C.TRANSMISSION_MODE_4K    // Transmission mode 4K
	TRANSMISSION_MODE_16K   FETransmitMode = C.TRANSMISSION_MODE_16K   // Transmission mode 16K
	TRANSMISSION_MODE_32K   FETransmitMode = C.TRANSMISSION_MODE_32K   // Transmission mode 32K
	TRANSMISSION_MODE_C1    FETransmitMode = C.TRANSMISSION_MODE_C1    // Single Carrier (C=1) transmission mode (DTMB only)
	TRANSMISSION_MODE_C3780 FETransmitMode = C.TRANSMISSION_MODE_C3780 // Multi Carrier (C=3780) transmission mode (DTMB only)
	TRANSMISSION_MODE_MIN                  = TRANSMISSION_MODE_2K
	TRANSMISSION_MODE_MAX                  = TRANSMISSION_MODE_C3780
)

const (
	GUARD_INTERVAL_AUTO   FEGuardInterval = C.GUARD_INTERVAL_AUTO   // Autodetect the guard interval
	GUARD_INTERVAL_1_128  FEGuardInterval = C.GUARD_INTERVAL_1_128  // Guard interval 1/128
	GUARD_INTERVAL_1_32   FEGuardInterval = C.GUARD_INTERVAL_1_32   // Guard interval 1/32
	GUARD_INTERVAL_1_16   FEGuardInterval = C.GUARD_INTERVAL_1_16   // Guard interval 1/16
	GUARD_INTERVAL_1_8    FEGuardInterval = C.GUARD_INTERVAL_1_8    // Guard interval 1/8
	GUARD_INTERVAL_1_4    FEGuardInterval = C.GUARD_INTERVAL_1_4    // Guard interval 1/4
	GUARD_INTERVAL_19_128 FEGuardInterval = C.GUARD_INTERVAL_19_128 // Guard interval 19/128
	GUARD_INTERVAL_19_256 FEGuardInterval = C.GUARD_INTERVAL_19_256 // Guard interval 19/256
	GUARD_INTERVAL_PN420  FEGuardInterval = C.GUARD_INTERVAL_PN420  // PN length 420 (1/4)
	GUARD_INTERVAL_PN595  FEGuardInterval = C.GUARD_INTERVAL_PN595  // PN length 595 (1/6)
	GUARD_INTERVAL_PN945  FEGuardInterval = C.GUARD_INTERVAL_PN945  // PN length 945 (1/9)
	GUARD_INTERVAL_MIN                    = GUARD_INTERVAL_1_32
	GUARD_INTERVAL_MAX                    = GUARD_INTERVAL_PN945
)

const (
	INVERSION_OFF  FEInversion = C.INVERSION_OFF  // Don't do spectral band inversion.
	INVERSION_ON   FEInversion = C.INVERSION_ON   // Do spectral band inversion.
	INVERSION_AUTO FEInversion = C.INVERSION_AUTO // Autodetect spectral band inversion.
	INVERSION_MIN              = INVERSION_OFF
	INVERSION_MAX              = INVERSION_AUTO
)

const (
	HIERARCHY_NONE FEHierarchy = C.HIERARCHY_NONE // No hierarchy
	HIERARCHY_AUTO FEHierarchy = C.HIERARCHY_AUTO // Autodetect hierarchy (if supported)
	HIERARCHY_1    FEHierarchy = C.HIERARCHY_1
	HIERARCHY_2    FEHierarchy = C.HIERARCHY_2
	HIERARCHY_4    FEHierarchy = C.HIERARCHY_4
	HIERARCHY_MIN              = HIERARCHY_NONE
	HIERARCHY_MAX              = HIERARCHY_AUTO
)

const (
	SCALE_NONE     = 0
	SCALE_DECIBEL  = 1
	SCALE_RELATIVE = 2
	SCALE_COUNTER  = 3
)

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func FEGetInfo(fd uintptr) (FEInfo, error) {
	var info FEInfo
	if err := dvb_ioctl(fd, FE_GET_INFO, unsafe.Pointer(&info)); err != 0 {
		return info, os.NewSyscallError("FE_GET_INFO", err)
	} else {
		return info, nil
	}
}

func FEReadStatus(fd uintptr) (FEStatus, error) {
	var status FEStatus
	if err := dvb_ioctl(fd, FE_READ_STATUS, unsafe.Pointer(&status)); err != 0 {
		return status, os.NewSyscallError("FE_READ_STATUS", err)
	} else {
		return status, nil
	}
}

func FETune(fd uintptr, params *TuneParams) error {
	kv := []C.struct_dtv_property{}

	// Add clear
	kv = append(kv, propUint32(DTV_CLEAR, 0))
	kv = append(kv, params.params()...)
	kv = append(kv, propUint32(DTV_TUNE, 0))

	// Call ioctl
	properties := C.struct_dtv_properties{
		C.uint(len(kv)), (*C.struct_dtv_property)(&kv[0]),
	}
	if err := dvb_ioctl(fd, FE_SET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return os.NewSyscallError("FE_SET_PROPERTY", err)
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS: GET/SET PROPERTIES

func FEGetPropUint32(fd uintptr, key FEKey) (uint32, error) {
	property := propUint32(key, 0)
	properties := C.struct_dtv_properties{
		1, (*C.struct_dtv_property)(unsafe.Pointer(&property)),
	}
	if err := dvb_ioctl(fd, FE_GET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return 0, os.NewSyscallError("FE_GET_PROPERTY", err)
	} else {
		value := (*fePropUint32)(unsafe.Pointer(&property))
		return uint32(value.Data), nil
	}
}

func FESetPropUint32(fd uintptr, key FEKey, value uint32) error {
	property := propUint32(key, value)
	properties := C.struct_dtv_properties{
		1, (*C.struct_dtv_property)(unsafe.Pointer(&property)),
	}
	if err := dvb_ioctl(fd, FE_SET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return os.NewSyscallError("FE_SET_PROPERTY", err)
	} else {
		return nil
	}
}

func FEGetPropEnum(fd uintptr, key FEKey) ([]uint8, error) {
	property := propUint32(key, 0)
	properties := C.struct_dtv_properties{
		1, (*C.struct_dtv_property)(unsafe.Pointer(&property)),
	}
	if err := dvb_ioctl(fd, FE_GET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return nil, os.NewSyscallError("FE_GET_PROPERTY", err)
	} else {
		value := (*fePropEnum)(unsafe.Pointer(&property))
		return value.Data[0:value.Len], nil
	}
}

func FEGetVersion(fd uintptr) (uint, uint, error) {
	if version, err := FEGetPropUint32(fd, DTV_API_VERSION); err != nil {
		return 0, 0, err
	} else {
		major := version >> 8 & 0xFF
		minor := version & 0xFF
		return uint(major), uint(minor), nil
	}
}

func FEEnumDeliverySystems(fd uintptr) ([]FEDeliverySystem, error) {
	data, err := FEGetPropEnum(fd, DTV_ENUM_DELSYS)
	if err != nil {
		return nil, err
	}
	enum := make([]FEDeliverySystem, len(data))
	for i, value := range data {
		enum[i] = FEDeliverySystem(value)
	}
	return enum, nil
}

func FEGetPropStats(fd uintptr, keys ...FEKey) ([]FEStats, error) {
	kv := []C.struct_dtv_property{}
	for _, key := range keys {
		kv = append(kv, propUint32(key, 0))
	}

	// Call ioctl
	properties := C.struct_dtv_properties{
		C.uint(len(kv)), (*C.struct_dtv_property)(&kv[0]),
	}
	if err := dvb_ioctl(fd, FE_GET_PROPERTY, unsafe.Pointer(&properties)); err != 0 {
		return nil, os.NewSyscallError("FE_GET_PROPERTY", err)
	}

	// Convert into array of stats
	result := []FEStats{}
	for _, property := range kv {
		value := (*fePropStats)(unsafe.Pointer(&property))
		r := bytes.NewReader(value.Data[:])
		stat := FEStats{
			Key:    FEKey(property.cmd),
			Values: make([]FEStat, int(value.Len)),
		}
		for i := 0; i < int(value.Len); i++ {
			if err := binary.Read(r, binary.LittleEndian, &stat.Values[i]); err != nil {
				return nil, err
			}
		}
		result = append(result, stat)
	}

	// Return success
	return result, nil
}

////////////////////////////////////////////////////////////////////////////////
// FEInfo

func (i FEInfo) Name() string {
	return C.GoString(&i.name[0])
}

func (i FEInfo) FrequencyMin() uint32 {
	return uint32(i.frequency_min)
}

func (i FEInfo) FrequencyMax() uint32 {
	return uint32(i.frequency_max)
}

func (i FEInfo) FrequencyStepSize() uint32 {
	return uint32(i.frequency_stepsize)
}

func (i FEInfo) FrequencyTolerance() uint32 {
	return uint32(i.frequency_tolerance)
}

func (i FEInfo) SymbolrateMin() uint32 {
	return uint32(i.symbol_rate_min)
}

func (i FEInfo) SymbolrateMax() uint32 {
	return uint32(i.symbol_rate_max)
}

func (i FEInfo) SymbolrateTolerance() uint32 {
	return uint32(i.symbol_rate_max)
}

func (i FEInfo) Caps() FECaps {
	return FECaps(i.caps)
}

/////////////////////////////////////////////////////////////////////////////////
// FEStat

func (s FEStat) Decibel() float64 {
	return float64(s.Value) * 0.001
}

func (s FEStat) Relative() float64 {
	return float64(s.Value&0xFFFF) * 100 / float64(0xFFFF)
}

func (s FEStat) Counter() uint64 {
	return uint64(s.Value)
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (i FEInfo) String() string {
	str := "<dvb.fe.info"
	str += " name=" + strconv.Quote(i.Name())
	str += " caps=" + fmt.Sprint(i.Caps())
	str += " frequency=" + fmt.Sprintf("{ %v,%v }", i.FrequencyMin(), i.FrequencyMax())
	str += " symbolrate=" + fmt.Sprintf("{ %v,%v }", i.SymbolrateMin(), i.SymbolrateMax())
	return str + ">"
}

func (s FEStat) String() string {
	switch s.Scale {
	case SCALE_NONE:
		return "<nil>"
	case SCALE_COUNTER:
		return fmt.Sprint(s.Counter())
	case SCALE_DECIBEL:
		return fmt.Sprintf("%.1fdB", s.Decibel())
	case SCALE_RELATIVE:
		return fmt.Sprintf("%.1f%%", s.Relative())
	default:
		return "[?? Invalid FEStat value]"
	}
}

func (f FECaps) String() string {
	if f == FE_IS_STUPID {
		return f.FlagString()
	}
	str := ""
	for v := FE_CAPS_MIN; v <= FE_CAPS_MAX; v <<= 1 {
		if f&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.Trim(str, "|")
}

func (f FECaps) FlagString() string {
	switch f {
	case FE_IS_STUPID:
		return "FE_IS_STUPID"
	case FE_CAN_INVERSION_AUTO:
		return "FE_CAN_INVERSION_AUTO"
	case FE_CAN_FEC_1_2:
		return "FE_CAN_FEC_1_2"
	case FE_CAN_FEC_2_3:
		return "FE_CAN_FEC_2_3"
	case FE_CAN_FEC_3_4:
		return "FE_CAN_FEC_3_4"
	case FE_CAN_FEC_4_5:
		return "FE_CAN_FEC_4_5"
	case FE_CAN_FEC_5_6:
		return "FE_CAN_FEC_5_6"
	case FE_CAN_FEC_6_7:
		return "FE_CAN_FEC_6_7"
	case FE_CAN_FEC_7_8:
		return "FE_CAN_FEC_7_8"
	case FE_CAN_FEC_8_9:
		return "FE_CAN_FEC_8_9"
	case FE_CAN_FEC_AUTO:
		return "FE_CAN_FEC_AUTO"
	case FE_CAN_QPSK:
		return "FE_CAN_QPSK"
	case FE_CAN_QAM_16:
		return "FE_CAN_QAM_16"
	case FE_CAN_QAM_32:
		return "FE_CAN_QAM_32"
	case FE_CAN_QAM_64:
		return "FE_CAN_QAM_64"
	case FE_CAN_QAM_128:
		return "FE_CAN_QAM_128"
	case FE_CAN_QAM_256:
		return "FE_CAN_QAM_256"
	case FE_CAN_QAM_AUTO:
		return "FE_CAN_QAM_AUTO"
	case FE_CAN_TRANSMISSION_MODE_AUTO:
		return "FE_CAN_TRANSMISSION_MODE_AUTO"
	case FE_CAN_BANDWIDTH_AUTO:
		return "FE_CAN_BANDWIDTH_AUTO"
	case FE_CAN_GUARD_INTERVAL_AUTO:
		return "FE_CAN_GUARD_INTERVAL_AUTO"
	case FE_CAN_HIERARCHY_AUTO:
		return "FE_CAN_HIERARCHY_AUTO"
	case FE_CAN_8VSB:
		return "FE_CAN_8VSB"
	case FE_CAN_16VSB:
		return "FE_CAN_16VSB"
	case FE_HAS_EXTENDED_CAPS:
		return "FE_HAS_EXTENDED_CAPS"
	case FE_CAN_MULTISTREAM:
		return "FE_CAN_MULTISTREAM"
	case FE_CAN_TURBO_FEC:
		return "FE_CAN_TURBO_FEC"
	case FE_CAN_2G_MODULATION:
		return "FE_CAN_2G_MODULATION"
	case FE_NEEDS_BENDING:
		return "FE_NEEDS_BENDING"
	case FE_CAN_RECOVER:
		return "FE_CAN_RECOVER"
	case FE_CAN_MUTE_TS:
		return "FE_CAN_MUTE_TS"
	default:
		return "[?? Invalud FECaps value]"
	}
}

func (f FEStatus) String() string {
	if f == FE_NONE {
		return f.FlagString()
	}
	str := ""
	for v := FE_STATUS_MIN; v <= FE_STATUS_MAX; v <<= 1 {
		if f&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.Trim(str, "|")
}

func (f FEStatus) FlagString() string {
	switch f {
	case FE_NONE:
		return "FE_NONE"
	case FE_HAS_SIGNAL:
		return "FE_HAS_SIGNAL"
	case FE_HAS_CARRIER:
		return "FE_HAS_CARRIER"
	case FE_HAS_VITERBI:
		return "FE_HAS_VITERBI"
	case FE_HAS_SYNC:
		return "FE_HAS_SYNC"
	case FE_HAS_LOCK:
		return "FE_HAS_LOCK"
	case FE_TIMEDOUT:
		return "FE_TIMEDOUT"
	case FE_REINIT:
		return "FE_REINIT"
	default:
		return "[?? Invalid FEStatus value]"
	}
}

func (f FEDeliverySystem) String() string {
	switch f {
	case SYS_UNDEFINED:
		return "SYS_UNDEFINED"
	case SYS_DVBC_ANNEX_A:
		return "SYS_DVBC_ANNEX_A"
	case SYS_DVBC_ANNEX_B:
		return "SYS_DVBC_ANNEX_B"
	case SYS_DVBC_ANNEX_C:
		return "SYS_DVBC_ANNEX_C"
	case SYS_ISDBC:
		return "SYS_ISDBC"
	case SYS_DVBT:
		return "SYS_DVBT"
	case SYS_DVBT2:
		return "SYS_DVBT2"
	case SYS_ISDBT:
		return "SYS_ISDBT"
	case SYS_ATSC:
		return "SYS_ATSC"
	case SYS_ATSCMH:
		return "SYS_ATSCMH"
	case SYS_DTMB:
		return "SYS_DTMB"
	case SYS_DVBS:
		return "SYS_DVBS"
	case SYS_DVBS2:
		return "SYS_DVBS2"
	case SYS_TURBO:
		return "SYS_TURBO"
	case SYS_ISDBS:
		return "SYS_ISDBS"
	case SYS_DAB:
		return "SYS_DAB"
	case SYS_DSS:
		return "SYS_DSS"
	case SYS_CMMB:
		return "SYS_CMMB"
	case SYS_DVBH:
		return "SYS_DVBH"
	default:
		return "[?? Invalid FEDeliverySystem value]"
	}
}

func (f FEModulation) String() string {
	switch f {
	case QPSK:
		return "QPSK"
	case QAM_16:
		return "QAM_16"
	case QAM_32:
		return "QAM_32"
	case QAM_64:
		return "QAM_64"
	case QAM_128:
		return "QAM_128"
	case QAM_256:
		return "QAM_256"
	case QAM_AUTO:
		return "QAM_AUTO"
	case VSB_8:
		return "VSB_8"
	case VSB_16:
		return "VSB_16"
	case PSK_8:
		return "PSK_8"
	case APSK_16:
		return "APSK_16"
	case APSK_32:
		return "APSK_32"
	case DQPSK:
		return "DQPSK"
	case QAM_4_NR:
		return "QAM_4_NR"
	default:
		return "[?? Invalid FEModulation value]"
	}
}

func (f FECodeRate) String() string {
	switch f {
	case FEC_NONE:
		return "FEC_NONE"
	case FEC_1_2:
		return "FEC_1_2"
	case FEC_2_3:
		return "FEC_2_3"
	case FEC_3_4:
		return "FEC_3_4"
	case FEC_4_5:
		return "FEC_4_5"
	case FEC_5_6:
		return "FEC_5_6"
	case FEC_6_7:
		return "FEC_6_7"
	case FEC_7_8:
		return "FEC_7_8"
	case FEC_8_9:
		return "FEC_8_9"
	case FEC_AUTO:
		return "FEC_AUTO"
	case FEC_3_5:
		return "FEC_3_5"
	case FEC_9_10:
		return "FEC_9_10"
	case FEC_2_5:
		return "FEC_2_5"
	default:
		return "[?? Invalid FECodeRate value]"
	}
}

func (f FETransmitMode) String() string {
	switch f {
	case TRANSMISSION_MODE_AUTO:
		return "TRANSMISSION_MODE_AUTO"
	case TRANSMISSION_MODE_1K:
		return "TRANSMISSION_MODE_1K"
	case TRANSMISSION_MODE_2K:
		return "TRANSMISSION_MODE_2K"
	case TRANSMISSION_MODE_8K:
		return "TRANSMISSION_MODE_8K"
	case TRANSMISSION_MODE_4K:
		return "TRANSMISSION_MODE_4K"
	case TRANSMISSION_MODE_16K:
		return "TRANSMISSION_MODE_16K"
	case TRANSMISSION_MODE_32K:
		return "TRANSMISSION_MODE_32K"
	case TRANSMISSION_MODE_C1:
		return "TRANSMISSION_MODE_C1"
	case TRANSMISSION_MODE_C3780:
		return "TRANSMISSION_MODE_C3780"
	default:
		return "[?? Invalid FETransmitMode value]"
	}
}

func (f FEGuardInterval) String() string {
	switch f {
	case GUARD_INTERVAL_AUTO:
		return "GUARD_INTERVAL_AUTO"
	case GUARD_INTERVAL_1_128:
		return "GUARD_INTERVAL_1_128"
	case GUARD_INTERVAL_1_32:
		return "GUARD_INTERVAL_1_32"
	case GUARD_INTERVAL_1_16:
		return "GUARD_INTERVAL_1_16"
	case GUARD_INTERVAL_1_8:
		return "GUARD_INTERVAL_1_8"
	case GUARD_INTERVAL_1_4:
		return "GUARD_INTERVAL_1_4"
	case GUARD_INTERVAL_19_128:
		return "GUARD_INTERVAL_19_128"
	case GUARD_INTERVAL_19_256:
		return "GUARD_INTERVAL_19_256"
	case GUARD_INTERVAL_PN420:
		return "GUARD_INTERVAL_PN420"
	case GUARD_INTERVAL_PN595:
		return "GUARD_INTERVAL_PN595"
	case GUARD_INTERVAL_PN945:
		return "GUARD_INTERVAL_PN945"
	default:
		return "[?? Invalid FEGuardInterval value]"
	}
}

func (f FEInversion) String() string {
	switch f {
	case INVERSION_OFF:
		return "INVERSION_OFF"
	case INVERSION_ON:
		return "INVERSION_ON"
	case INVERSION_AUTO:
		return "INVERSION_AUTO"
	default:
		return "[?? Invalid FEInversion value]"
	}
}

func (f FEHierarchy) String() string {
	switch f {
	case HIERARCHY_NONE:
		return "HIERARCHY_NONE"
	case HIERARCHY_AUTO:
		return "HIERARCHY_AUTO"
	case HIERARCHY_1:
		return "HIERARCHY_1"
	case HIERARCHY_2:
		return "HIERARCHY_2"
	case HIERARCHY_4:
		return "HIERARCHY_4"
	default:
		return "[?? Invalid FEHierarchy value]"
	}
}

func (k FEKey) String() string {
	switch k {
	case DTV_UNDEFINED:
		return "DTV_UNDEFINED"
	case DTV_TUNE:
		return "DTV_TUNE"
	case DTV_CLEAR:
		return "DTV_CLEAR"
	case DTV_FREQUENCY:
		return "DTV_FREQUENCY"
	case DTV_MODULATION:
		return "DTV_MODULATION"
	case DTV_BANDWIDTH_HZ:
		return "DTV_BANDWIDTH_HZ"
	case DTV_INVERSION:
		return "DTV_INVERSION"
	case DTV_DISEQC_MASTER:
		return "DTV_DISEQC_MASTER"
	case DTV_SYMBOL_RATE:
		return "DTV_SYMBOL_RATE"
	case DTV_INNER_FEC:
		return "DTV_INNER_FEC"
	case DTV_VOLTAGE:
		return "DTV_VOLTAGE"
	case DTV_TONE:
		return "DTV_TONE"
	case DTV_PILOT:
		return "DTV_PILOT"
	case DTV_ROLLOFF:
		return "DTV_ROLLOFF"
	case DTV_DISEQC_SLAVE_REPLY:
		return "DTV_DISEQC_SLAVE_REPLY"
	case DTV_FE_CAPABILITY_COUNT:
		return "DTV_FE_CAPABILITY_COUNT"
	case DTV_FE_CAPABILITY:
		return "DTV_FE_CAPABILITY"
	case DTV_DELIVERY_SYSTEM:
		return "DTV_DELIVERY_SYSTEM"
	case DTV_ISDBT_PARTIAL_RECEPTION:
		return "DTV_ISDBT_PARTIAL_RECEPTION"
	case DTV_ISDBT_SOUND_BROADCASTING:
		return "DTV_ISDBT_SOUND_BROADCASTING"
	case DTV_ISDBT_SB_SUBCHANNEL_ID:
		return "DTV_ISDBT_SB_SUBCHANNEL_ID"
	case DTV_ISDBT_SB_SEGMENT_IDX:
		return "DTV_ISDBT_SB_SEGMENT_IDX"
	case DTV_ISDBT_SB_SEGMENT_COUNT:
		return "DTV_ISDBT_SB_SEGMENT_COUNT"
	case DTV_ISDBT_LAYERA_FEC:
		return "DTV_ISDBT_LAYERA_FEC"
	case DTV_ISDBT_LAYERA_MODULATION:
		return "DTV_ISDBT_LAYERA_MODULATION"
	case DTV_ISDBT_LAYERA_SEGMENT_COUNT:
		return "DTV_ISDBT_LAYERA_SEGMENT_COUNT"
	case DTV_ISDBT_LAYERA_TIME_INTERLEAVING:
		return "DTV_ISDBT_LAYERA_TIME_INTERLEAVING"
	case DTV_ISDBT_LAYERB_FEC:
		return "DTV_ISDBT_LAYERB_FEC"
	case DTV_ISDBT_LAYERB_MODULATION:
		return "DTV_ISDBT_LAYERB_MODULATION"
	case DTV_ISDBT_LAYERB_SEGMENT_COUNT:
		return "DTV_ISDBT_LAYERB_SEGMENT_COUNT"
	case DTV_ISDBT_LAYERB_TIME_INTERLEAVING:
		return "DTV_ISDBT_LAYERB_TIME_INTERLEAVING"
	case DTV_ISDBT_LAYERC_FEC:
		return "DTV_ISDBT_LAYERC_FEC"
	case DTV_ISDBT_LAYERC_MODULATION:
		return "DTV_ISDBT_LAYERC_MODULATION"
	case DTV_ISDBT_LAYERC_SEGMENT_COUNT:
		return "DTV_ISDBT_LAYERC_SEGMENT_COUNT"
	case DTV_ISDBT_LAYERC_TIME_INTERLEAVING:
		return "DTV_ISDBT_LAYERC_TIME_INTERLEAVING"
	case DTV_API_VERSION:
		return "DTV_API_VERSION"
	case DTV_CODE_RATE_HP:
		return "DTV_CODE_RATE_HP"
	case DTV_CODE_RATE_LP:
		return "DTV_CODE_RATE_LP"
	case DTV_GUARD_INTERVAL:
		return "DTV_GUARD_INTERVAL"
	case DTV_TRANSMISSION_MODE:
		return "DTV_TRANSMISSION_MODE"
	case DTV_HIERARCHY:
		return "DTV_HIERARCHY"
	case DTV_ISDBT_LAYER_ENABLED:
		return "DTV_ISDBT_LAYER_ENABLED"
	case DTV_STREAM_ID:
		return "DTV_STREAM_ID"
	case DTV_ENUM_DELSYS:
		return "DTV_ENUM_DELSYS"
	case DTV_ATSCMH_FIC_VER:
		return "DTV_ATSCMH_FIC_VER"
	case DTV_ATSCMH_PARADE_ID:
		return "DTV_ATSCMH_PARADE_ID"
	case DTV_ATSCMH_NOG:
		return "DTV_ATSCMH_NOG"
	case DTV_ATSCMH_TNOG:
		return "DTV_ATSCMH_TNOG"
	case DTV_ATSCMH_SGN:
		return "DTV_ATSCMH_SGN"
	case DTV_ATSCMH_PRC:
		return "DTV_ATSCMH_PRC"
	case DTV_ATSCMH_RS_FRAME_MODE:
		return "DTV_ATSCMH_RS_FRAME_MODE"
	case DTV_ATSCMH_RS_FRAME_ENSEMBLE:
		return "DTV_ATSCMH_RS_FRAME_ENSEMBLE"
	case DTV_ATSCMH_RS_CODE_MODE_PRI:
		return "DTV_ATSCMH_RS_CODE_MODE_PRI"
	case DTV_ATSCMH_RS_CODE_MODE_SEC:
		return "DTV_ATSCMH_RS_CODE_MODE_SEC"
	case DTV_ATSCMH_SCCC_BLOCK_MODE:
		return "DTV_ATSCMH_SCCC_BLOCK_MODE"
	case DTV_ATSCMH_SCCC_CODE_MODE_A:
		return "DTV_ATSCMH_SCCC_CODE_MODE_A"
	case DTV_ATSCMH_SCCC_CODE_MODE_B:
		return "DTV_ATSCMH_SCCC_CODE_MODE_B"
	case DTV_ATSCMH_SCCC_CODE_MODE_C:
		return "DTV_ATSCMH_SCCC_CODE_MODE_C"
	case DTV_ATSCMH_SCCC_CODE_MODE_D:
		return "DTV_ATSCMH_SCCC_CODE_MODE_D"
	case DTV_INTERLEAVING:
		return "DTV_INTERLEAVING"
	case DTV_LNA:
		return "DTV_LNA"
	case DTV_STAT_SIGNAL_STRENGTH:
		return "DTV_STAT_SIGNAL_STRENGTH"
	case DTV_STAT_CNR:
		return "DTV_STAT_CNR"
	case DTV_STAT_PRE_ERROR_BIT_COUNT:
		return "DTV_STAT_PRE_ERROR_BIT_COUNT"
	case DTV_STAT_PRE_TOTAL_BIT_COUNT:
		return "DTV_STAT_PRE_TOTAL_BIT_COUNT"
	case DTV_STAT_POST_ERROR_BIT_COUNT:
		return "DTV_STAT_POST_ERROR_BIT_COUNT"
	case DTV_STAT_POST_TOTAL_BIT_COUNT:
		return "DTV_STAT_POST_TOTAL_BIT_COUNT"
	case DTV_STAT_ERROR_BLOCK_COUNT:
		return "DTV_STAT_ERROR_BLOCK_COUNT"
	case DTV_STAT_TOTAL_BLOCK_COUNT:
		return "DTV_STAT_TOTAL_BLOCK_COUNT"
	case DTV_SCRAMBLING_SEQUENCE_INDEX:
		return "DTV_SCRAMBLING_SEQUENCE_INDEX"
	default:
		return "[?? Invalid FEKey value]"
	}
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func propUint32(cmd FEKey, value uint32) C.struct_dtv_property {
	v := C.struct_dtv_property{
		cmd: C.uint(cmd),
	}
	v_ := (*fePropUint32)(unsafe.Pointer(&v))
	v_.Data = value
	return v
}
