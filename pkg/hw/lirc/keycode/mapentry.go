package keycode

import (
	"fmt"

	gopi "github.com/djthorpe/gopi/v3"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type mapentry struct {
	Key     gopi.KeyCode
	Device  gopi.InputDeviceType
	Code    uint32
	Comment string
	Index   int
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *mapentry) Equals(other *mapentry) bool {
	if other == nil {
		return false
	}
	if other.Key != this.Key {
		return false
	}
	if other.Code != this.Code {
		return false
	}
	if other.Device != this.Device {
		return false
	}
	if other.Comment != this.Comment {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *mapentry) String() string {
	str := "<keycode.map"
	if this.Key != gopi.KEYCODE_NONE {
		str += fmt.Sprintf(" key=%v", this.Key)
	}
	if this.Device != gopi.INPUT_DEVICE_NONE {
		str += " device=" + fmt.Sprint(this.Device)
	}
	if this.Code != 0 {
		str += " code=" + scancodeString(this.Code)
	}
	if this.Comment != "" {
		str += fmt.Sprintf(" comment=%q", this.Comment)
	}
	return str + ">"
}
