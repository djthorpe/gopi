package chromecast

import (
	"fmt"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Volume struct {
	Level float32 `json:"level,omitempty"`
	Muted bool    `json:"muted"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (v Volume) String() string {
	str := "<cast.volume"
	str += " level=" + fmt.Sprintf("%.2f", v.Level)
	str += " muted=" + fmt.Sprint(v.Muted)
	return str + ">"
}

func (v Volume) Equals(other Volume) bool {
	if v.Level != other.Level {
		return false
	}
	if v.Muted != other.Muted {
		return false
	}
	return true
}
