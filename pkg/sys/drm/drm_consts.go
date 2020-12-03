// +build drm

package drm

////////////////////////////////////////////////////////////////////////////////
// TYPES

type (
	ModeConnection uint
)

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ModeConnectionNone ModeConnection = iota
	ModeConnectionConnected
	ModeConnectionDisconnected
	ModeConnectionUnknown
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (c ModeConnection) String() string {
	switch c {
	case 	ModeConnectionNone:
		return "ModeConnectionNone"
	case 	ModeConnectionConnected:
		return "ModeConnectionConnected"
	case 	ModeConnectionDisconnected:
		return "ModeConnectionDisconnected"
	case 	ModeConnectionUnknown:
		return "ModeConnectionUnknown"
	default:
		return "[?? Invalid ModeConnection value]"
	}
}
