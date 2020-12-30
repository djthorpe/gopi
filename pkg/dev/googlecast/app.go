package googlecast

import (
	"strconv"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type App struct {
	AppId        string `json:"appId"`
	DisplayName  string `json:"displayName"`
	IsIdleScreen bool   `json:"isIdleScreen"`
	SessionId    string `json:"sessionId"`
	StatusText   string `json:"statusText"`
	TransportId  string `json:"transportId"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (a App) Equals(other App) bool {
	if a.AppId != other.AppId {
		return false
	}
	if a.DisplayName != other.DisplayName {
		return false
	}
	if a.IsIdleScreen != other.IsIdleScreen {
		return false
	}
	if a.SessionId != other.SessionId {
		return false
	}
	if a.StatusText != other.StatusText {
		return false
	}
	if a.TransportId != other.TransportId {
		return false
	}
	return true
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (a App) String() string {
	str := "<cast.app"
	str += " id=" + strconv.Quote(a.AppId)
	if a.DisplayName != "" {
		str += " name=" + strconv.Quote(a.DisplayName)
	}
	if a.StatusText != "" {
		str += " statusText=" + strconv.Quote(a.StatusText)
	}
	if a.TransportId != "" {
		str += " transportId=" + strconv.Quote(a.TransportId)
	}
	if a.SessionId != "" {
		str += " sessionId=" + strconv.Quote(a.SessionId)
	}
	if a.IsIdleScreen {
		str += " isIdleScreen=true"
	}
	return str + ">"
}
