package argonone

import (
	"context"
	"fmt"
	"os"
	"time"

	gopi "github.com/djthorpe/gopi/v3"

	// Units
	_ "github.com/djthorpe/gopi/v3/pkg/hw/i2c"
	_ "github.com/djthorpe/gopi/v3/pkg/hw/platform"
	_ "github.com/djthorpe/gopi/v3/pkg/log"
	_ "github.com/djthorpe/gopi/v3/pkg/metrics"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type argonone struct {
	gopi.Unit
	gopi.I2C
	gopi.Platform
	gopi.Logger
	gopi.Metrics

	bus         gopi.I2CBus
	slave       uint8
	tempzone    string
	measurement string
	fan         *fanValue
}

////////////////////////////////////////////////////////////////////////////////
// GLOBALS

const (
	fanMin, fanMax = 0x00, 0x64       // Minimum and maximum fan duty cycle values
	fanDelay       = time.Second * 30 // Delay for this number of seconds to set new fan value

	// The period for measuring CPU temperature
	measureDelta = 5 * time.Second
)

var (
	// Default fan configuration, if celcius is greater or equal to
	// value then return fan value or else return zero
	fanConfig = fanConfigArr{
		{55.0, 10},
		{60.0, 50},
		{65.0, 100},
	}
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func (this *argonone) Define(cfg gopi.Config) error {
	cfg.FlagUint("i2c.bus", 1, "I2C Bus")
	cfg.FlagUint("i2c.slave", 0x1A, "I2C Slave")
	cfg.FlagString("argonone.zone", "", "Temperature Zone name")
	cfg.FlagString("argonone.measurement", "argonone", "Measurement name")
	return nil
}

func (this *argonone) New(cfg gopi.Config) error {
	// Check I2C
	if bus, slave := this.i2cBusSlave(cfg); slave == 0 {
		return fmt.Errorf("Missing I2C interface")
	} else if detected, err := this.I2C.DetectSlave(bus, slave); err != nil {
		return err
	} else if detected == false {
		return fmt.Errorf("Missing I2C slave (slave 0x%02X)", slave)
	} else if err := this.I2C.SetSlave(bus, slave); err != nil {
		return err
	} else {
		this.bus = bus
		this.slave = slave
	}

	// Set Hysteresis to prevent flapping
	this.fan = NewFanValue(fanDelay)

	// Check platform
	if this.Platform == nil {
		return fmt.Errorf("Missing Platform interface")
	} else {
		this.tempzone = cfg.GetString("argonone.zone")
	}

	// Define measurement
	if this.Metrics == nil {
		return fmt.Errorf("Missing Metrics interface")
	} else if measurement := cfg.GetString("argonone.measurement"); measurement != "" {
		if host, err := os.Hostname(); err != nil {
			return err
		} else if tags := this.Metrics.NewFields(fmt.Sprintf("host=%q", host)); tags == nil {
			return gopi.ErrInternalAppError
		} else if m, err := this.Metrics.NewMeasurement(measurement, "celcius float32 fan uint8", tags...); err != nil {
			return err
		} else {
			this.Debug("measurement=", m)
			this.measurement = m.Name()
		}
	}

	// Return success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// RUN

func (this *argonone) Run(ctx context.Context) error {
	timer := time.NewTimer(time.Nanosecond)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			if celcius := this.getTemperature(); celcius != 0 {
				if err := this.setFanForTemperature(celcius); err != nil {
					this.Print("SetFanForTemperature: ", err)
				}
			}
			timer.Reset(measureDelta)
		case <-ctx.Done():
			return nil
		}
	}
}

////////////////////////////////////////////////////////////////////////////////
// PUBLIC METHODS

func (this *argonone) SetFan(value uint8) error {
	if value < fanMin || value > fanMax {
		return gopi.ErrBadParameter.WithPrefix("SetFan")
	} else if n, err := this.I2C.Write(this.bus, []byte{value}); err != nil {
		return err
	} else if n != 1 {
		return gopi.ErrUnexpectedResponse
	} else {
		return nil
	}
}

func (this *argonone) SetPower(mode gopi.ArgonOnePowerMode) error {
	var value uint8
	switch mode {
	case gopi.ARGONONE_POWER_DEFAULT:
		value = 0xFD
	case gopi.ARGONONE_POWER_ALWAYSON:
		value = 0xFE
	case gopi.ARGONONE_POWER_UART:
		value = 0xFF
	default:
		return gopi.ErrBadParameter.WithPrefix("SetPower")
	}
	if n, err := this.I2C.Write(this.bus, []byte{value}); err != nil {
		return err
	} else if n != 1 {
		return gopi.ErrUnexpectedResponse
	} else {
		return nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *argonone) String() string {
	str := "<argonone"
	str += " bus=" + fmt.Sprint(this.bus)
	if this.slave != 0 {
		str += " slave=" + fmt.Sprint(this.slave)
	}
	return str + ">"
}

////////////////////////////////////////////////////////////////////////////////
// PRIVATE METHODS

func (this *argonone) i2cBusSlave(cfg gopi.Config) (gopi.I2CBus, uint8) {
	bus := gopi.I2CBus(cfg.GetUint("i2c.bus"))
	slave := uint8(cfg.GetUint("i2c.slave"))
	if this.I2C == nil || this.i2cHasBus(bus) == false || slave == 0 {
		return bus, 0
	} else {
		return bus, slave
	}
}

func (this *argonone) i2cHasBus(bus gopi.I2CBus) bool {
	for _, d := range this.I2C.Devices() {
		if bus == d {
			return true
		}
	}
	return false
}

func (this *argonone) getTemperature() float32 {
	measurements := this.Platform.TemperatureZones()
	if len(measurements) == 0 {
		return 0
	}

	// Iterate through measurements, matching the zone
	for k, v := range measurements {
		if this.tempzone == "" || this.tempzone == k {
			return v
		}
	}

	// No temperature found
	return 0
}

func (this *argonone) setFanForTemperature(celcius float32) error {
	// Obtain fan value for temperature
	fan := fanConfig.fanForTemperature(celcius)

	// Report measurement
	if this.measurement != "" {
		if err := this.Metrics.Emit(this.measurement, celcius, fan); err != nil {
			return err
		}
	}

	// Determine if we should set the fan
	if this.fan.Set(fan) {
		this.Debugf("Setting fan => %d%%", fan)
		if err := this.SetFan(fan); err != nil {
			return err
		}
	}

	// Return success
	return nil
}
