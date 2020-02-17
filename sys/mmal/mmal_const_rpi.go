//+build mmal

package mmal

import (
	"strings"
)

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: mmal
#include <interface/mmal/mmal.h>
*/
import "C"

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	MMAL_BOOL_FALSE = C.MMAL_BOOL_T(0)
	MMAL_BOOL_TRUE  = C.MMAL_BOOL_T(1)
)

const (
	MMAL_SUCCESS   MMAL_Status = iota
	MMAL_ENOMEM                // Out of memory
	MMAL_ENOSPC                // Out of resources (other than memory)
	MMAL_EINVAL                // Argument is invalid
	MMAL_ENOSYS                // Function not implemented
	MMAL_ENOENT                // No such file or directory
	MMAL_ENXIO                 // No such device or address
	MMAL_EIO                   // I/O error
	MMAL_ESPIPE                // Illegal seek
	MMAL_ECORRUPT              // Data is corrupt
	MMAL_ENOTREADY             // Component is not ready
	MMAL_ECONFIG               // Component is not configured
	MMAL_EISCONN               // Port is already connected
	MMAL_ENOTCONN              // Port is disconnected
	MMAL_EAGAIN                // Resource temporarily unavailable. Try again later
	MMAL_EFAULT                // Bad address
	MMAL_MIN       = MMAL_SUCCESS
	MMAL_MAX       = MMAL_EFAULT
)

const (
	MMAL_PORT_TYPE_UNKNOWN MMAL_PortType = iota
	MMAL_PORT_TYPE_CONTROL               // Control port
	MMAL_PORT_TYPE_INPUT                 // Input port
	MMAL_PORT_TYPE_OUTPUT                // Output port
	MMAL_PORT_TYPE_CLOCK                 // Clock port
	MMAL_PORT_TYPE_MAX     = MMAL_PORT_TYPE_CLOCK
	MMAL_PORT_TYPE_NONE    = MMAL_PORT_TYPE_UNKNOWN
)

const (
	MMAL_PORT_CAPABILITY_PASSTHROUGH                  MMAL_PortCapability = 0x01
	MMAL_PORT_CAPABILITY_ALLOCATION                   MMAL_PortCapability = 0x02
	MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE MMAL_PortCapability = 0x04
	MMAL_PORT_CAPABILITY_MAX                                              = MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE
	MMAL_PORT_CAPABILITY_MIN                                              = MMAL_PORT_CAPABILITY_PASSTHROUGH
)

const (
	MMAL_COMPONENT_DEFAULT_VIDEO_DECODER    = "vc.ril.video_decode"
	MMAL_COMPONENT_DEFAULT_VIDEO_ENCODER    = "vc.ril.video_encode"
	MMAL_COMPONENT_DEFAULT_VIDEO_RENDERER   = "vc.ril.video_render"
	MMAL_COMPONENT_DEFAULT_IMAGE_DECODER    = "vc.ril.image_decode"
	MMAL_COMPONENT_DEFAULT_IMAGE_ENCODER    = "vc.ril.image_encode"
	MMAL_COMPONENT_DEFAULT_CAMERA           = "vc.ril.camera"
	MMAL_COMPONENT_DEFAULT_VIDEO_CONVERTER  = "vc.video_convert"
	MMAL_COMPONENT_DEFAULT_SPLITTER         = "vc.splitter"
	MMAL_COMPONENT_DEFAULT_SCHEDULER        = "vc.scheduler"
	MMAL_COMPONENT_DEFAULT_VIDEO_INJECTER   = "vc.video_inject"
	MMAL_COMPONENT_DEFAULT_VIDEO_SPLITTER   = "vc.ril.video_splitter"
	MMAL_COMPONENT_DEFAULT_AUDIO_DECODER    = "none"
	MMAL_COMPONENT_DEFAULT_AUDIO_RENDERER   = "vc.ril.audio_render"
	MMAL_COMPONENT_DEFAULT_MIRACAST         = "vc.miracast"
	MMAL_COMPONENT_DEFAULT_CLOCK            = "vc.clock"
	MMAL_COMPONENT_DEFAULT_CAMERA_INFO      = "vc.camera_info"
	MMAL_COMPONENT_DEFAULT_CONTAINER_READER = "container_reader"
	MMAL_COMPONENT_DEFAULT_CONTAINER_WRITER = "container_writer"
)

const (
	MMAL_PARAMETER_GROUP_COMMON   MMAL_ParameterType = (iota << 16)
	MMAL_PARAMETER_GROUP_CAMERA   MMAL_ParameterType = (iota << 16) // Camera-specific parameter ID group
	MMAL_PARAMETER_GROUP_VIDEO    MMAL_ParameterType = (iota << 16) // Video-specific parameter ID group
	MMAL_PARAMETER_GROUP_AUDIO    MMAL_ParameterType = (iota << 16) // Audio-specific parameter ID group
	MMAL_PARAMETER_GROUP_CLOCK    MMAL_ParameterType = (iota << 16) // Clock-specific parameter ID group
	MMAL_PARAMETER_GROUP_MIRACAST MMAL_ParameterType = (iota << 16) // Miracast-specific parameter ID group
	MMAL_PARAMETER_GROUP_MAX                         = MMAL_PARAMETER_GROUP_MIRACAST
	MMAL_PARAMETER_GROUP_MIN                         = MMAL_PARAMETER_GROUP_COMMON
)

const (
	// MMAL_PARAMETER_GROUP_COMMON
	_                                   MMAL_ParameterType = iota
	MMAL_PARAMETER_SUPPORTED_ENCODINGS                     // Takes a MMAL_PARAMETER_ENCODING_T
	MMAL_PARAMETER_URI                                     // Takes a MMAL_PARAMETER_URI_T
	MMAL_PARAMETER_CHANGE_EVENT_REQUEST                    // Takes a MMAL_PARAMETER_CHANGE_EVENT_REQUEST_T
	MMAL_PARAMETER_ZERO_COPY                               // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_BUFFER_REQUIREMENTS                     // Takes a MMAL_PARAMETER_BUFFER_REQUIREMENTS_T
	MMAL_PARAMETER_STATISTICS                              // Takes a MMAL_PARAMETER_STATISTICS_T
	MMAL_PARAMETER_CORE_STATISTICS                         // Takes a MMAL_PARAMETER_CORE_STATISTICS_T
	MMAL_PARAMETER_MEM_USAGE                               // Takes a MMAL_PARAMETER_MEM_USAGE_T
	MMAL_PARAMETER_BUFFER_FLAG_FILTER                      // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_SEEK                                    // Takes a MMAL_PARAMETER_SEEK_T
	MMAL_PARAMETER_POWERMON_ENABLE                         // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_LOGGING                                 // Takes a MMAL_PARAMETER_LOGGING_T
	MMAL_PARAMETER_SYSTEM_TIME                             // Takes a MMAL_PARAMETER_UINT64_T
	MMAL_PARAMETER_NO_IMAGE_PADDING                        // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_LOCKSTEP_ENABLE                         // Takes a MMAL_PARAMETER_BOOLEAN_T
)

const (
	// MMAL_PARAMETER_GROUP_VIDEO
	MMAL_PARAMETER_DISPLAYREGION                        MMAL_ParameterType = iota // Takes a MMAL_DISPLAYREGION_T
	MMAL_PARAMETER_SUPPORTED_PROFILES                                             // Takes a MMAL_PARAMETER_VIDEO_PROFILE_T
	MMAL_PARAMETER_PROFILE                                                        // Takes a MMAL_PARAMETER_VIDEO_PROFILE_T
	MMAL_PARAMETER_INTRAPERIOD                                                    // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_RATECONTROL                                                    // Takes a MMAL_PARAMETER_VIDEO_RATECONTROL_T
	MMAL_PARAMETER_NALUNITFORMAT                                                  // Takes a MMAL_PARAMETER_VIDEO_NALUNITFORMAT_T
	MMAL_PARAMETER_MINIMISE_FRAGMENTATION                                         // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_MB_ROWS_PER_SLICE                                              // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_LEVEL_EXTENSION                                          // Takes a MMAL_PARAMETER_VIDEO_LEVEL_EXTENSION_T
	MMAL_PARAMETER_VIDEO_EEDE_ENABLE                                              // Takes a MMAL_PARAMETER_VIDEO_EEDE_ENABLE_T
	MMAL_PARAMETER_VIDEO_EEDE_LOSSRATE                                            // Takes a MMAL_PARAMETER_VIDEO_EEDE_LOSSRATE_T
	MMAL_PARAMETER_VIDEO_REQUEST_I_FRAME                                          // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_INTRA_REFRESH                                            // Takes a MMAL_PARAMETER_VIDEO_INTRA_REFRESH_T
	MMAL_PARAMETER_VIDEO_IMMUTABLE_INPUT                                          // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_BIT_RATE                                                 // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_FRAME_RATE                                               // Takes a MMAL_PARAMETER_FRAME_RATE_T
	MMAL_PARAMETER_VIDEO_ENCODE_MIN_QUANT                                         // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ENCODE_MAX_QUANT                                         // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ENCODE_RC_MODEL                                          // Takes a MMAL_PARAMETER_VIDEO_ENCODE_RC_MODEL_T
	MMAL_PARAMETER_EXTRA_BUFFERS                                                  // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ALIGN_HORIZ                                              // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ALIGN_VERT                                               // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_DROPPABLE_PFRAMES                                        // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_INITIAL_QUANT                                     // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ENCODE_QP_P                                              // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ENCODE_RC_SLICE_DQUANT                                   // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ENCODE_FRAME_LIMIT_BITS                                  // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ENCODE_PEAK_RATE                                         // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_ENCODE_H264_DISABLE_CABAC                                // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_H264_LOW_LATENCY                                  // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_H264_AU_DELIMITERS                                // Takes a MMAL_PARAMETER_BOOLEAN_
	MMAL_PARAMETER_VIDEO_ENCODE_H264_DEBLOCK_IDC                                  // Takes a MMAL_PARAMETER_UINT32_
	MMAL_PARAMETER_VIDEO_ENCODE_H264_MB_INTRA_MODE                                // Takes a MMAL_PARAMETER_VIDEO_ENCODER_H264_MB_INTRA_MODES_T
	MMAL_PARAMETER_VIDEO_ENCODE_HEADER_ON_OPEN                                    // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_PRECODE_FOR_QP                                    // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_DRM_INIT_INFO                                            // Takes a MMAL_PARAMETER_VIDEO_DRM_INIT_INFO_T
	MMAL_PARAMETER_VIDEO_TIMESTAMP_FIFO                                           // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_DECODE_ERROR_CONCEALMENT                                 // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_DRM_PROTECT_BUFFER                                       // Takes a MMAL_PARAMETER_VIDEO_DRM_PROTECT_BUFFER_T
	MMAL_PARAMETER_VIDEO_DECODE_CONFIG_VD3                                        // Takes a MMAL_PARAMETER_BYTES_T
	MMAL_PARAMETER_VIDEO_ENCODE_H264_VCL_HRD_PARAMETERS                           // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_H264_LOW_DELAY_HRD_FLAG                           // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_INLINE_HEADER                                     // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_SEI_ENABLE                                        // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_INLINE_VECTORS                                    // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_RENDER_STATS                                             // Takes a MMAL_PARAMETER_VIDEO_RENDER_STATS_T
	MMAL_PARAMETER_VIDEO_INTERLACE_TYPE                                           // Takes a MMAL_PARAMETER_VIDEO_INTERLACE_TYPE_T
	MMAL_PARAMETER_VIDEO_INTERPOLATE_TIMESTAMPS                                   // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_ENCODE_SPS_TIMING                                        // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_MAX_NUM_CALLBACKS                                        // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_VIDEO_SOURCE_PATTERN                                           // Takes a MMAL_PARAMETER_SOURCE_PATTERN_T
	MMAL_PARAMETER_VIDEO_ENCODE_SEPARATE_NAL_BUFS                                 // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_VIDEO_DROPPABLE_PFRAME_LENGTH                                  // Takes a MMAL_PARAMETER_UINT32_T
)

const (
	// MMAL_PARAMETER_GROUP_CAMERA
	MMAL_PARAMETER_THUMBNAIL_CONFIGURATION     MMAL_ParameterType = iota // Takes a MMAL_PARAMETER_THUMBNAIL_CONFIG_T
	MMAL_PARAMETER_CAPTURE_QUALITY                                       // Unused?
	MMAL_PARAMETER_ROTATION                                              // Takes a MMAL_PARAMETER_INT32_T
	MMAL_PARAMETER_EXIF_DISABLE                                          // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_EXIF                                                  // Takes a MMAL_PARAMETER_EXIF_T
	MMAL_PARAMETER_AWB_MODE                                              // Takes a MMAL_PARAM_AWBMODE_T
	MMAL_PARAMETER_IMAGE_EFFECT                                          // Takes a MMAL_PARAMETER_IMAGEFX_T
	MMAL_PARAMETER_COLOUR_EFFECT                                         // Takes a MMAL_PARAMETER_COLOURFX_T
	MMAL_PARAMETER_FLICKER_AVOID                                         // Takes a MMAL_PARAMETER_FLICKERAVOID_T
	MMAL_PARAMETER_FLASH                                                 // Takes a MMAL_PARAMETER_FLASH_T
	MMAL_PARAMETER_REDEYE                                                // Takes a MMAL_PARAMETER_REDEYE_T
	MMAL_PARAMETER_FOCUS                                                 // Takes a MMAL_PARAMETER_FOCUS_T
	MMAL_PARAMETER_FOCAL_LENGTHS                                         // Unused?
	MMAL_PARAMETER_EXPOSURE_COMP                                         // Takes a MMAL_PARAMETER_INT32_T or MMAL_PARAMETER_RATIONAL_T
	MMAL_PARAMETER_ZOOM                                                  // Takes a MMAL_PARAMETER_SCALEFACTOR_T
	MMAL_PARAMETER_MIRROR                                                // Takes a MMAL_PARAMETER_MIRROR_T
	MMAL_PARAMETER_CAMERA_NUM                                            // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_CAPTURE                                               // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_EXPOSURE_MODE                                         // Takes a MMAL_PARAMETER_EXPOSUREMODE_T
	MMAL_PARAMETER_EXP_METERING_MODE                                     // Takes a MMAL_PARAMETER_EXPOSUREMETERINGMODE_T
	MMAL_PARAMETER_FOCUS_STATUS                                          // Takes a MMAL_PARAMETER_FOCUS_STATUS_T
	MMAL_PARAMETER_CAMERA_CONFIG                                         // Takes a MMAL_PARAMETER_CAMERA_CONFIG_T
	MMAL_PARAMETER_CAPTURE_STATUS                                        // Takes a MMAL_PARAMETER_CAPTURE_STATUS_T
	MMAL_PARAMETER_FACE_TRACK                                            // Takes a MMAL_PARAMETER_FACE_TRACK_T
	MMAL_PARAMETER_DRAW_BOX_FACES_AND_FOCUS                              // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_JPEG_Q_FACTOR                                         // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_FRAME_RATE                                            // Takes a MMAL_PARAMETER_FRAME_RATE_T
	MMAL_PARAMETER_USE_STC                                               // Takes a MMAL_PARAMETER_CAMERA_STC_MODE_T
	MMAL_PARAMETER_CAMERA_INFO                                           // Takes a MMAL_PARAMETER_CAMERA_INFO_T
	MMAL_PARAMETER_VIDEO_STABILISATION                                   // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_FACE_TRACK_RESULTS                                    // Takes a MMAL_PARAMETER_FACE_TRACK_RESULTS_T
	MMAL_PARAMETER_ENABLE_RAW_CAPTURE                                    // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_DPF_FILE                                              // Takes a MMAL_PARAMETER_URI_T
	MMAL_PARAMETER_ENABLE_DPF_FILE                                       // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_DPF_FAIL_IS_FATAL                                     // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_CAPTURE_MODE                                          // Takes a MMAL_PARAMETER_CAPTUREMODE_T
	MMAL_PARAMETER_FOCUS_REGIONS                                         // Takes a MMAL_PARAMETER_FOCUS_REGIONS_T
	MMAL_PARAMETER_INPUT_CROP                                            // Takes a MMAL_PARAMETER_INPUT_CROP_T
	MMAL_PARAMETER_SENSOR_INFORMATION                                    // Takes a MMAL_PARAMETER_SENSOR_INFORMATION_T
	MMAL_PARAMETER_FLASH_SELECT                                          // Takes a MMAL_PARAMETER_FLASH_SELECT_T
	MMAL_PARAMETER_FIELD_OF_VIEW                                         // Takes a MMAL_PARAMETER_FIELD_OF_VIEW_T
	MMAL_PARAMETER_HIGH_DYNAMIC_RANGE                                    // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_DYNAMIC_RANGE_COMPRESSION                             // Takes a MMAL_PARAMETER_DRC_T
	MMAL_PARAMETER_ALGORITHM_CONTROL                                     // Takes a MMAL_PARAMETER_ALGORITHM_CONTROL_T
	MMAL_PARAMETER_SHARPNESS                                             // Takes a MMAL_PARAMETER_RATIONAL_T
	MMAL_PARAMETER_CONTRAST                                              // Takes a MMAL_PARAMETER_RATIONAL_T
	MMAL_PARAMETER_BRIGHTNESS                                            // Takes a MMAL_PARAMETER_RATIONAL_T
	MMAL_PARAMETER_SATURATION                                            // Takes a MMAL_PARAMETER_RATIONAL_T
	MMAL_PARAMETER_ISO                                                   // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_ANTISHAKE                                             // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_IMAGE_EFFECT_PARAMETERS                               // Takes a MMAL_PARAMETER_IMAGEFX_PARAMETERS_T
	MMAL_PARAMETER_CAMERA_BURST_CAPTURE                                  // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_CAMERA_MIN_ISO                                        // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_CAMERA_USE_CASE                                       // Takes a MMAL_PARAMETER_CAMERA_USE_CASE_T
	MMAL_PARAMETER_CAPTURE_STATS_PASS                                    // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_CAMERA_CUSTOM_SENSOR_CONFIG                           // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_ENABLE_REGISTER_FILE                                  // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_REGISTER_FAIL_IS_FATAL                                // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_CONFIGFILE_REGISTERS                                  // Takes a MMAL_PARAMETER_CONFIGFILE_T
	MMAL_PARAMETER_CONFIGFILE_CHUNK_REGISTERS                            // Takes a MMAL_PARAMETER_CONFIGFILE_CHUNK_T
	MMAL_PARAMETER_JPEG_ATTACH_LOG                                       // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_ZERO_SHUTTER_LAG                                      // Takes a MMAL_PARAMETER_ZEROSHUTTERLAG_T
	MMAL_PARAMETER_FPS_RANGE                                             // Takes a MMAL_PARAMETER_FPS_RANGE_T
	MMAL_PARAMETER_CAPTURE_EXPOSURE_COMP                                 // Takes a MMAL_PARAMETER_INT32_T
	MMAL_PARAMETER_SW_SHARPEN_DISABLE                                    // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_FLASH_REQUIRED                                        // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_SW_SATURATION_DISABLE                                 // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_SHUTTER_SPEED                                         // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_CUSTOM_AWB_GAINS                                      // Takes a MMAL_PARAMETER_AWB_GAINS_T
	MMAL_PARAMETER_CAMERA_SETTINGS                                       // Takes a MMAL_PARAMETER_CAMERA_SETTINGS_T
	MMAL_PARAMETER_PRIVACY_INDICATOR                                     // Takes a MMAL_PARAMETER_PRIVACY_INDICATOR_T
	MMAL_PARAMETER_VIDEO_DENOISE                                         // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_STILLS_DENOISE                                        // Takes a MMAL_PARAMETER_BOOLEAN_T
	MMAL_PARAMETER_ANNOTATE                                              // Takes a MMAL_PARAMETER_CAMERA_ANNOTATE_T
	MMAL_PARAMETER_STEREOSCOPIC_MODE                                     // Takes a MMAL_PARAMETER_STEREOSCOPIC_MODE_T
	MMAL_PARAMETER_CAMERA_INTERFACE                                      // Takes a MMAL_PARAMETER_CAMERA_INTERFACE_T
	MMAL_PARAMETER_CAMERA_CLOCKING_MODE                                  // Takes a MMAL_PARAMETER_CAMERA_CLOCKING_MODE_T
	MMAL_PARAMETER_CAMERA_RX_CONFIG                                      // Takes a MMAL_PARAMETER_CAMERA_RX_CONFIG_T
	MMAL_PARAMETER_CAMERA_RX_TIMING                                      // Takes a MMAL_PARAMETER_CAMERA_RX_TIMING_T
	MMAL_PARAMETER_DPF_CONFIG                                            // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_JPEG_RESTART_INTERVAL                                 // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_CAMERA_ISP_BLOCK_OVERRIDE                             // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_LENS_SHADING_OVERRIDE                                 // Takes a MMAL_PARAMETER_LENS_SHADING_T
	MMAL_PARAMETER_BLACK_LEVEL                                           // Takes a MMAL_PARAMETER_UINT32_T
	MMAL_PARAMETER_RESIZE_PARAMS                                         // Takes a MMAL_PARAMETER_RESIZE_T
	MMAL_PARAMETER_CROP                                                  // Takes a MMAL_PARAMETER_CROP_T
	MMAL_PARAMETER_OUTPUT_SHIFT                                          // Takes a MMAL_PARAMETER_INT32_T
	MMAL_PARAMETER_CCM_SHIFT                                             // Takes a MMAL_PARAMETER_INT32_T
	MMAL_PARAMETER_CUSTOM_CCM                                            // Takes a MMAL_PARAMETER_CUSTOM_CCM_T
	MMAL_PARAMETER_ANALOG_GAIN                                           // Takes a MMAL_PARAMETER_RATIONAL_T
	MMAL_PARAMETER_DIGITAL_GAIN                                          // Takes a MMAL_PARAMETER_RATIONAL_T
)

const (
	// MMAL_ES_TYPE_T
	MMAL_STREAM_TYPE_UNKNOWN    MMAL_StreamType = iota // Unknown elementary stream type
	MMAL_STREAM_TYPE_CONTROL                           // Elementary stream of control commands
	MMAL_STREAM_TYPE_AUDIO                             // Audio elementary stream
	MMAL_STREAM_TYPE_VIDEO                             //  Video elementary stream
	MMAL_STREAM_TYPE_SUBPICTURE                        // Sub-picture elementary stream (e.g. subtitles, overlays)
	MMAL_STREAM_TYPE_MIN        = MMAL_STREAM_TYPE_UNKNOWN
	MMAL_STREAM_TYPE_MAX        = MMAL_STREAM_TYPE_SUBPICTURE
)

const (
	// MMAL_StreamCompareFlags
	MMAL_STREAM_COMPARE_FLAG_TYPE               MMAL_StreamCompareFlags = 0x0001 // The type is different
	MMAL_STREAM_COMPARE_FLAG_ENCODING           MMAL_StreamCompareFlags = 0x0002 // The encoding is different
	MMAL_STREAM_COMPARE_FLAG_BITRATE            MMAL_StreamCompareFlags = 0x0004 // The bitrate is different
	MMAL_STREAM_COMPARE_FLAG_FLAGS              MMAL_StreamCompareFlags = 0x0008 // The flags are different
	MMAL_STREAM_COMPARE_FLAG_EXTRADATA          MMAL_StreamCompareFlags = 0x0010 // The extradata is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_RESOLUTION   MMAL_StreamCompareFlags = 0x0100 // The video resolution is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_CROPPING     MMAL_StreamCompareFlags = 0x0200 // The video cropping is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_FRAME_RATE   MMAL_StreamCompareFlags = 0x0400 // The video frame rate is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_ASPECT_RATIO MMAL_StreamCompareFlags = 0x0800 // The video aspect ratio is different
	MMAL_STREAM_COMPARE_FLAG_VIDEO_COLOR_SPACE  MMAL_StreamCompareFlags = 0x1000 // The video color space is different
	MMAL_STREAM_COMPARE_FLAG_MIN                                        = MMAL_STREAM_COMPARE_FLAG_TYPE
	MMAL_STREAM_COMPARE_FLAG_MAX                                        = MMAL_STREAM_COMPARE_FLAG_VIDEO_COLOR_SPACE
)

const (
	MMAL_DISPLAY_SET_NONE        C.uint32_t = 0x0000
	MMAL_DISPLAY_SET_NUM         C.uint32_t = 0x0001
	MMAL_DISPLAY_SET_FULLSCREEN  C.uint32_t = 0x0002
	MMAL_DISPLAY_SET_TRANSFORM   C.uint32_t = 0x0004
	MMAL_DISPLAY_SET_DEST_RECT   C.uint32_t = 0x0008
	MMAL_DISPLAY_SET_SRC_RECT    C.uint32_t = 0x0010
	MMAL_DISPLAY_SET_MODE        C.uint32_t = 0x0020
	MMAL_DISPLAY_SET_PIXEL       C.uint32_t = 0x0040
	MMAL_DISPLAY_SET_NOASPECT    C.uint32_t = 0x0080
	MMAL_DISPLAY_SET_LAYER       C.uint32_t = 0x0100
	MMAL_DISPLAY_SET_COPYPROTECT C.uint32_t = 0x0200
	MMAL_DISPLAY_SET_ALPHA       C.uint32_t = 0x0400
)

const (
	MMAL_PARAMETER_SEEK_FLAG_PRECISE = 0x01
	MMAL_PARAMETER_SEEK_FLAG_FORWARD = 0x02
)

const (
	// Up to 20 supported encodings in array
	MMAL_PARAMETER_SUPPORTED_ENCODINGS_ARRAY_SIZE = 20 * C.sizeof_uint32_t
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (s MMAL_Status) Error() string {
	switch s {
	case MMAL_SUCCESS:
		return "MMAL_SUCCESS"
	case MMAL_ENOMEM:
		return "MMAL_ENOMEM: Out of memory"
	case MMAL_ENOSPC:
		return "MMAL_ENOSPC: Out of resources (other than memory)"
	case MMAL_EINVAL:
		return "MMAL_EINVAL: Argument is invalid"
	case MMAL_ENOSYS:
		return "MMAL_ENOSYS: Function not implemented"
	case MMAL_ENOENT:
		return "MMAL_ENOENT: No such file or directory"
	case MMAL_ENXIO:
		return "MMAL_ENXIO: No such device or address"
	case MMAL_EIO:
		return "MMAL_EIO: I/O error"
	case MMAL_ESPIPE:
		return "MMAL_ESPIPE: Illegal seek"
	case MMAL_ECORRUPT:
		return "MMAL_ECORRUPT: Data is corrupt"
	case MMAL_ENOTREADY:
		return "MMAL_ENOTREADY: Component is not ready"
	case MMAL_ECONFIG:
		return "MMAL_ECONFIG: Component is not configured"
	case MMAL_EISCONN:
		return "MMAL_EISCONN: Port is already connected"
	case MMAL_ENOTCONN:
		return "MMAL_ENOTCONN: Port is disconnected"
	case MMAL_EAGAIN:
		return "MMAL_EAGAIN: Resource temporarily unavailable. Try again later"
	case MMAL_EFAULT:
		return "MMAL_EFAULT: Bad address"
	default:
		return "[?? Invalid MMAL_StatusType value]"
	}
}

func (p MMAL_PortType) String() string {
	switch p {
	case MMAL_PORT_TYPE_UNKNOWN:
		return "MMAL_PORT_TYPE_UNKNOWN"
	case MMAL_PORT_TYPE_CONTROL:
		return "MMAL_PORT_TYPE_CONTROL"
	case MMAL_PORT_TYPE_INPUT:
		return "MMAL_PORT_TYPE_INPUT"
	case MMAL_PORT_TYPE_OUTPUT:
		return "MMAL_PORT_TYPE_OUTPUT"
	case MMAL_PORT_TYPE_CLOCK:
		return "MMAL_PORT_TYPE_CLOCK"
	default:
		return "[?? Invalid MMAL_PortType value]"
	}
}

func (c MMAL_PortCapability) String() string {
	parts := ""
	for flag := MMAL_PORT_CAPABILITY_MIN; flag <= MMAL_PORT_CAPABILITY_MAX; flag <<= 1 {
		if c&flag == 0 {
			continue
		}
		switch flag {
		case MMAL_PORT_CAPABILITY_PASSTHROUGH:
			parts += "|" + "MMAL_PORT_CAPABILITY_PASSTHROUGH"
		case MMAL_PORT_CAPABILITY_ALLOCATION:
			parts += "|" + "MMAL_PORT_CAPABILITY_ALLOCATION"
		case MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE:
			parts += "|" + "MMAL_PORT_CAPABILITY_SUPPORTS_EVENT_FORMAT_CHANGE"
		default:
			parts += "|" + "[?? Invalid MMAL_PortCapability value]"
		}
	}
	return strings.Trim(parts, "|")
}

func (s MMAL_StreamType) String() string {
	switch s {
	case MMAL_STREAM_TYPE_UNKNOWN:
		return "MMAL_STREAM_TYPE_UNKNOWN"
	case MMAL_STREAM_TYPE_CONTROL:
		return "MMAL_STREAM_TYPE_CONTROL"
	case MMAL_STREAM_TYPE_AUDIO:
		return "MMAL_STREAM_TYPE_AUDIO"
	case MMAL_STREAM_TYPE_VIDEO:
		return "MMAL_STREAM_TYPE_VIDEO"
	case MMAL_STREAM_TYPE_SUBPICTURE:
		return "MMAL_STREAM_TYPE_SUBPICTURE"
	default:
		return "[?? Invalid MMAL_StreamType value]"
	}
}

func (f MMAL_StreamCompareFlags) String() string {
	parts := ""
	for flag := MMAL_STREAM_COMPARE_FLAG_MIN; flag <= MMAL_STREAM_COMPARE_FLAG_MAX; flag <<= 1 {
		if f&flag == 0 {
			continue
		}
		switch flag {
		case MMAL_STREAM_COMPARE_FLAG_TYPE:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_TYPE"
		case MMAL_STREAM_COMPARE_FLAG_ENCODING:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_ENCODING"
		case MMAL_STREAM_COMPARE_FLAG_BITRATE:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_BITRATE"
		case MMAL_STREAM_COMPARE_FLAG_FLAGS:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_FLAGS"
		case MMAL_STREAM_COMPARE_FLAG_EXTRADATA:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_EXTRADATA"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_RESOLUTION:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_RESOLUTION"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_CROPPING:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_CROPPING"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_FRAME_RATE:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_FRAME_RATE"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_ASPECT_RATIO:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_ASPECT_RATIO"
		case MMAL_STREAM_COMPARE_FLAG_VIDEO_COLOR_SPACE:
			parts += "|" + "MMAL_STREAM_COMPARE_FLAG_VIDEO_COLOR_SPACE"
		default:
			parts += "|" + "[?? Invalid MMAL_StreamCompareFlags value]"
		}
	}
	return strings.Trim(parts, "|")
}
