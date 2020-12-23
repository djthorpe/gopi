// +build drm

package drm

////////////////////////////////////////////////////////////////////////////////
// CGO

/*
#cgo pkg-config: libdrm
#include <xf86drm.h>
#include <xf86drmMode.h>
#include <libdrm/drm_fourcc.h>
*/
import "C"
import "strings"

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	ModeConnection uint
	ConnectorType  C.int
	ModeInfoType   uint
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ModeConnectionNone ModeConnection = iota
	ModeConnectionConnected
	ModeConnectionDisconnected
	ModeConnectionUnknown
)

const (
	DRM_MODE_FB_MODIFIERS = C.DRM_MODE_FB_MODIFIERS
)

const (
	DRM_FORMAT_MOD_LINEAR = C.DRM_FORMAT_MOD_LINEAR
)

const (
	DRM_CLIENT_CAP_STEREO_3D            = C.DRM_CLIENT_CAP_STEREO_3D
	DRM_CLIENT_CAP_UNIVERSAL_PLANES     = C.DRM_CLIENT_CAP_UNIVERSAL_PLANES
	DRM_CLIENT_CAP_ATOMIC               = C.DRM_CLIENT_CAP_ATOMIC
	DRM_CLIENT_CAP_ASPECT_RATIO         = C.DRM_CLIENT_CAP_ASPECT_RATIO
	DRM_CLIENT_CAP_WRITEBACK_CONNECTORS = C.DRM_CLIENT_CAP_WRITEBACK_CONNECTORS
)

const (
	DRM_MODE_PAGE_FLIP_EVENT           = C.DRM_MODE_PAGE_FLIP_EVENT
	DRM_MODE_PAGE_FLIP_ASYNC           = C.DRM_MODE_PAGE_FLIP_ASYNC
	DRM_MODE_PAGE_FLIP_TARGET_ABSOLUTE = C.DRM_MODE_PAGE_FLIP_TARGET_ABSOLUTE
	DRM_MODE_PAGE_FLIP_TARGET_RELATIVE = C.DRM_MODE_PAGE_FLIP_TARGET_RELATIVE
)

const (
	DRM_MODE_CONNECTOR_Unknown     ConnectorType = C.DRM_MODE_CONNECTOR_Unknown
	DRM_MODE_CONNECTOR_VGA         ConnectorType = C.DRM_MODE_CONNECTOR_VGA
	DRM_MODE_CONNECTOR_DVII        ConnectorType = C.DRM_MODE_CONNECTOR_DVII
	DRM_MODE_CONNECTOR_DVID        ConnectorType = C.DRM_MODE_CONNECTOR_DVID
	DRM_MODE_CONNECTOR_DVIA        ConnectorType = C.DRM_MODE_CONNECTOR_DVIA
	DRM_MODE_CONNECTOR_Composite   ConnectorType = C.DRM_MODE_CONNECTOR_Composite
	DRM_MODE_CONNECTOR_SVIDEO      ConnectorType = C.DRM_MODE_CONNECTOR_SVIDEO
	DRM_MODE_CONNECTOR_LVDS        ConnectorType = C.DRM_MODE_CONNECTOR_LVDS
	DRM_MODE_CONNECTOR_Component   ConnectorType = C.DRM_MODE_CONNECTOR_Component
	DRM_MODE_CONNECTOR_9PinDIN     ConnectorType = C.DRM_MODE_CONNECTOR_9PinDIN
	DRM_MODE_CONNECTOR_DisplayPort ConnectorType = C.DRM_MODE_CONNECTOR_DisplayPort
	DRM_MODE_CONNECTOR_HDMIA       ConnectorType = C.DRM_MODE_CONNECTOR_HDMIA
	DRM_MODE_CONNECTOR_HDMIB       ConnectorType = C.DRM_MODE_CONNECTOR_HDMIB
	DRM_MODE_CONNECTOR_TV          ConnectorType = C.DRM_MODE_CONNECTOR_TV
	DRM_MODE_CONNECTOR_eDP         ConnectorType = C.DRM_MODE_CONNECTOR_eDP
	DRM_MODE_CONNECTOR_VIRTUAL     ConnectorType = C.DRM_MODE_CONNECTOR_VIRTUAL
	DRM_MODE_CONNECTOR_DSI         ConnectorType = C.DRM_MODE_CONNECTOR_DSI
	DRM_MODE_CONNECTOR_DPI         ConnectorType = C.DRM_MODE_CONNECTOR_DPI
)

const (
	DRM_MODE_TYPE_BUILTIN ModeInfoType = (1 << iota)
	DRM_MODE_TYPE_CLOCK_C
	DRM_MODE_TYPE_CRTC_C
	DRM_MODE_TYPE_PREFERRED
	DRM_MODE_TYPE_DEFAULT
	DRM_MODE_TYPE_USERDEF
	DRM_MODE_TYPE_DRIVER
	DRM_MODE_TYPE_NONE ModeInfoType = 0
	DRM_MODE_TYPE_MIN               = DRM_MODE_TYPE_BUILTIN
	DRM_MODE_TYPE_MAX               = DRM_MODE_TYPE_DRIVER
)

const (
	DRM_PLANE_TYPE_OVERLAY = C.DRM_PLANE_TYPE_OVERLAY
	DRM_PLANE_TYPE_PRIMARY = C.DRM_PLANE_TYPE_PRIMARY
	DRM_PLANE_TYPE_CURSOR  = C.DRM_PLANE_TYPE_CURSOR
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c ModeConnection) String() string {
	switch c {
	case ModeConnectionNone:
		return "ModeConnectionNone"
	case ModeConnectionConnected:
		return "ModeConnectionConnected"
	case ModeConnectionDisconnected:
		return "ModeConnectionDisconnected"
	case ModeConnectionUnknown:
		return "ModeConnectionUnknown"
	default:
		return "[?? Invalid ModeConnection value]"
	}
}

func (c ConnectorType) String() string {
	switch c {
	case DRM_MODE_CONNECTOR_Unknown:
		return "DRM_MODE_CONNECTOR_Unknown"
	case DRM_MODE_CONNECTOR_VGA:
		return "DRM_MODE_CONNECTOR_VGA"
	case DRM_MODE_CONNECTOR_DVII:
		return "DRM_MODE_CONNECTOR_DVII"
	case DRM_MODE_CONNECTOR_DVID:
		return "DRM_MODE_CONNECTOR_DVID"
	case DRM_MODE_CONNECTOR_DVIA:
		return "DRM_MODE_CONNECTOR_DVIA"
	case DRM_MODE_CONNECTOR_Composite:
		return "DRM_MODE_CONNECTOR_Composite"
	case DRM_MODE_CONNECTOR_SVIDEO:
		return "DRM_MODE_CONNECTOR_SVIDEO"
	case DRM_MODE_CONNECTOR_LVDS:
		return "DRM_MODE_CONNECTOR_LVDS"
	case DRM_MODE_CONNECTOR_Component:
		return "DRM_MODE_CONNECTOR_Component"
	case DRM_MODE_CONNECTOR_9PinDIN:
		return "DRM_MODE_CONNECTOR_9PinDIN"
	case DRM_MODE_CONNECTOR_DisplayPort:
		return "DRM_MODE_CONNECTOR_DisplayPort"
	case DRM_MODE_CONNECTOR_HDMIA:
		return "DRM_MODE_CONNECTOR_HDMIA"
	case DRM_MODE_CONNECTOR_HDMIB:
		return "DRM_MODE_CONNECTOR_HDMIB"
	case DRM_MODE_CONNECTOR_TV:
		return "DRM_MODE_CONNECTOR_TV"
	case DRM_MODE_CONNECTOR_eDP:
		return "DRM_MODE_CONNECTOR_eDP"
	case DRM_MODE_CONNECTOR_VIRTUAL:
		return "DRM_MODE_CONNECTOR_VIRTUAL"
	case DRM_MODE_CONNECTOR_DSI:
		return "DRM_MODE_CONNECTOR_DSI"
	case DRM_MODE_CONNECTOR_DPI:
		return "DRM_MODE_CONNECTOR_DPI"
	default:
		return "[?? Unknown ConnectorType value]"
	}
}

func (t ModeInfoType) String() string {
	if t == DRM_MODE_TYPE_NONE {
		return t.FlagString()
	}
	str := ""
	for v := DRM_MODE_TYPE_MIN; v <= DRM_MODE_TYPE_MAX; v <<= 1 {
		if t&v == v {
			str += v.FlagString() + "|"
		}
	}
	return strings.TrimSuffix(str, "|")
}

func (t ModeInfoType) FlagString() string {
	switch t {
	case DRM_MODE_TYPE_NONE:
		return "DRM_MODE_TYPE_NONE"
	case DRM_MODE_TYPE_BUILTIN:
		return "DRM_MODE_TYPE_BUILTIN"
	case DRM_MODE_TYPE_CLOCK_C:
		return "DRM_MODE_TYPE_CLOCK_C"
	case DRM_MODE_TYPE_CRTC_C:
		return "DRM_MODE_TYPE_CRTC_C"
	case DRM_MODE_TYPE_PREFERRED:
		return "DRM_MODE_TYPE_PREFERRED"
	case DRM_MODE_TYPE_DEFAULT:
		return "DRM_MODE_TYPE_DEFAULT"
	case DRM_MODE_TYPE_USERDEF:
		return "DRM_MODE_TYPE_USERDEF"
	case DRM_MODE_TYPE_DRIVER:
		return "DRM_MODE_TYPE_DRIVER"
	default:
		return "[?? Invalid ModeInfoType value]"
	}
}
