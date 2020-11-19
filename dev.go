package gopi

////////////////////////////////////////////////////////////////////////////////
// ARGON ONE CASE

// Ref: https://github.com/Argon40Tech/Argon-ONE-i2c-Codes

type ArgonOnePowerMode uint

type ArgonOne interface {
	// Set fan duty cycle (0-100)
	SetFan(uint8) error

	// Set Power Mode
	SetPower(ArgonOnePowerMode) error
}

////////////////////////////////////////////////////////////////////////////////
// CONSTANTS

const (
	ARGONONE_POWER_DEFAULT ArgonOnePowerMode = iota
	ARGONONE_POWER_ALWAYSON
	ARGONONE_POWER_UART
)

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (v ArgonOnePowerMode) String() string {
	switch v {
	case ARGONONE_POWER_DEFAULT:
		return "ARGONONE_POWER_DEFAULT"
	case ARGONONE_POWER_ALWAYSON:
		return "ARGONONE_POWER_ALWAYSON"
	case ARGONONE_POWER_UART:
		return "ARGONONE_POWER_UART"
	default:
		return "[?? Invalid ArgonOnePowerMode value]"
	}
}
